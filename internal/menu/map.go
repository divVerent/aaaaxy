// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package menu

import (
	"fmt"
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/font"
	"github.com/divVerent/aaaaaa/internal/image"
	"github.com/divVerent/aaaaaa/internal/input"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
	"github.com/divVerent/aaaaaa/internal/player_state"
)

type MapScreen struct {
	Menu       *Menu
	Level      *level.Level
	CurrentCP  string
	SortedLocs []string

	cpSprite          *ebiten.Image
	cpSelectedSprite  *ebiten.Image
	deadEndSprite     *ebiten.Image
	cpCheckmarkSprite *ebiten.Image
	whiteImage        *ebiten.Image
}

// TODO: parametrize.
const (
	firstCP = "leap_of_faith"
)

func (s *MapScreen) Init(m *Menu) error {
	s.Menu = m
	s.CurrentCP = s.Menu.World.PlayerState.LastCheckpoint()
	if s.CurrentCP == "" {
		// Have no checkpoint yet - start the game right away.
		return s.Menu.SwitchToGame()
	}
	var err error
	s.cpSprite, err = image.Load("sprites", "checkpoint.png")
	if err != nil {
		return err
	}
	s.cpSelectedSprite, err = image.Load("sprites", "checkpoint_selected.png")
	if err != nil {
		return err
	}
	s.deadEndSprite, err = image.Load("sprites", "dead_end.png")
	if err != nil {
		return err
	}
	s.cpCheckmarkSprite, err = image.Load("sprites", "checkpoint_checkmark.png")
	if err != nil {
		return err
	}
	s.whiteImage = ebiten.NewImage(1, 1)
	s.whiteImage.Fill(color.Gray{255})

	s.SortedLocs = nil
	for name := range s.Menu.World.Level.Checkpoints {
		s.SortedLocs = append(s.SortedLocs, name)
	}
	// Note: we do not care for the actual order, just that it does not change between frames.
	sort.Strings(s.SortedLocs)

	return nil
}

func (s *MapScreen) moveBy(d m.Delta) {
	loc := s.Menu.World.Level.CheckpointLocations
	cpLoc := loc.Locs[s.CurrentCP]
	edge, found := cpLoc.NextByDir[d]
	if !found {
		return
	}
	edgeSeen := s.Menu.World.PlayerState.CheckpointsWalked(s.CurrentCP, edge.Other)
	reverseSeen := s.Menu.World.PlayerState.CheckpointsWalked(edge.Other, s.CurrentCP)
	if !edgeSeen && !reverseSeen {
		// Don't know this yet :)
		return
	}
	if s.Menu.World.Level.Checkpoints[edge.Other].Properties["dead_end"] == "true" {
		// A dead end!
		return
	}
	s.CurrentCP = edge.Other
	s.Menu.MoveSound(nil)
}

func (s *MapScreen) exit() error {
	if s.CurrentCP != firstCP && s.Menu.World.PlayerState.CheckpointSeen(firstCP) != player_state.NotSeen {
		s.CurrentCP = firstCP
		return s.Menu.MoveSound(nil)
	}
	return s.Menu.ActivateSound(s.Menu.SwitchToScreen(&MainScreen{}))
}

func (s *MapScreen) Update() error {
	if input.Exit.JustHit {
		return s.exit()
	}
	if input.Left.JustHit {
		s.moveBy(m.West())
	}
	if input.Right.JustHit {
		s.moveBy(m.East())
	}
	if input.Up.JustHit {
		s.moveBy(m.North())
	}
	if input.Down.JustHit {
		s.moveBy(m.South())
	}
	if input.Jump.JustHit || input.Action.JustHit {
		return s.Menu.ActivateSound(s.Menu.SwitchToCheckpoint(s.CurrentCP))
	}
	return nil
}

func (s *MapScreen) Draw(screen *ebiten.Image) {
	h := engine.GameHeight
	w := engine.GameWidth
	x := w / 2
	fgs := color.NRGBA{R: 255, G: 255, B: 75, A: 255}
	bgs := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	lineColor := color.NRGBA{R: 170, G: 170, B: 170, A: 255}
	darkLineColor := color.NRGBA{R: 75, G: 75, B: 75, A: 255}
	font.MenuBig.Draw(screen, "Pick-a-Path", m.Pos{X: x, Y: h / 8}, true, fgs, bgs)
	cpText := s.Menu.World.Level.Checkpoints[s.CurrentCP].Properties["text"]
	seen, total := s.Menu.World.PlayerState.TnihSignsSeen(s.CurrentCP)
	if total > 0 {
		cpText += fmt.Sprintf(" (%d/%d)", seen, total)
	}
	font.Menu.Draw(screen, cpText, m.Pos{X: x, Y: 7 * h / 8}, true, fgs, bgs)

	// Draw all known checkpoints.
	loc := s.Menu.World.Level.CheckpointLocations
	mapWidth := w
	mapHeight := h / 2
	if mapWidth*loc.Rect.Size.DY > mapHeight*loc.Rect.Size.DX {
		mapWidth = mapHeight * loc.Rect.Size.DX / loc.Rect.Size.DY
	} else {
		mapHeight = mapWidth * loc.Rect.Size.DY / loc.Rect.Size.DX
	}
	mapRect := m.Rect{
		Origin: m.Pos{X: (w - mapWidth) / 2, Y: (h - mapHeight) / 2},
		Size:   m.Delta{DX: mapWidth, DY: mapHeight},
	}
	// ebitenutil.DrawRect(screen, float64(mapRect.Origin.X), float64(mapRect.Origin.Y), float64(mapRect.Size.DX), float64(mapRect.Size.DY), bgs)
	// First draw all edges.
	for _, cpName := range s.SortedLocs {
		cpLoc := loc.Locs[cpName]
		if s.Menu.World.PlayerState.CheckpointSeen(cpName) == player_state.NotSeen {
			continue
		}
		pos := cpLoc.MapPos.FromRectToRect(loc.Rect, mapRect)
		for _, dir := range level.AllCheckpointDirs {
			edge, found := cpLoc.NextByDir[dir]
			if !found || !edge.Forward || edge.Optional {
				continue
			}
			otherName := edge.Other
			edgeSeen := s.Menu.World.PlayerState.CheckpointsWalked(cpName, otherName)
			closePos := pos.Add(dir.Mul(7))
			options := &ebiten.DrawTrianglesOptions{
				CompositeMode: ebiten.CompositeModeSourceOver,
				Filter:        ebiten.FilterNearest,
			}
			geoM := &ebiten.GeoM{}
			geoM.Scale(0, 0)
			if edgeSeen {
				otherLoc := loc.Locs[otherName]
				otherPos := otherLoc.MapPos.FromRectToRect(loc.Rect, mapRect)
				farPos := otherPos.Sub(dir.Mul(7))
				engine.DrawPolyLine(screen, 6.0, []m.Pos{pos, closePos, farPos, otherPos}, s.whiteImage, darkLineColor, geoM, options)
				engine.DrawPolyLine(screen, 3.0, []m.Pos{pos, closePos, farPos, otherPos}, s.whiteImage, lineColor, geoM, options)
			} else {
				engine.DrawPolyLine(screen, 3.0, []m.Pos{pos, closePos}, s.whiteImage, darkLineColor, geoM, options)
			}
		}
	}
	// Then draw the CPs.
	for _, cpName := range s.SortedLocs {
		cpLoc := loc.Locs[cpName]
		if s.Menu.World.PlayerState.CheckpointSeen(cpName) == player_state.NotSeen {
			continue
		}
		sprite := s.cpSprite
		if cpName == s.CurrentCP {
			sprite = s.cpSelectedSprite
		} else if s.Menu.World.Level.Checkpoints[cpName].Properties["dead_end"] == "true" {
			sprite = s.deadEndSprite
		}
		pos := cpLoc.MapPos.FromRectToRect(loc.Rect, mapRect)
		opts := ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeSourceOver,
			Filter:        ebiten.FilterNearest,
		}
		opts.GeoM.Translate(float64(pos.X-7), float64(pos.Y-7))
		screen.DrawImage(sprite, &opts)
	}
	// Finally the checkmarks.
	for _, cpName := range s.SortedLocs {
		cpLoc := loc.Locs[cpName]
		if s.Menu.World.PlayerState.CheckpointSeen(cpName) == player_state.NotSeen {
			continue
		}
		if seen, total := s.Menu.World.PlayerState.TnihSignsSeen(cpName); seen < total {
			continue
		}
		sprite := s.cpCheckmarkSprite
		pos := cpLoc.MapPos.FromRectToRect(loc.Rect, mapRect)
		opts := ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeSourceOver,
			Filter:        ebiten.FilterNearest,
		}
		opts.GeoM.Translate(float64(pos.X-7), float64(pos.Y-7))
		screen.DrawImage(sprite, &opts)
	}
}

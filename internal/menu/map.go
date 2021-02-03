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
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/font"
	"github.com/divVerent/aaaaaa/internal/image"
	"github.com/divVerent/aaaaaa/internal/input"
	m "github.com/divVerent/aaaaaa/internal/math"
)

type MapScreen struct {
	Menu      *Menu
	Level     *engine.Level
	CurrentCP string

	cpSprite         *ebiten.Image
	cpSelectedSprite *ebiten.Image
	deadEndSprite    *ebiten.Image
	whiteImage       *ebiten.Image
}

func (s *MapScreen) Init(m *Menu) error {
	s.Menu = m
	s.CurrentCP = s.Menu.World.Level.Player.PersistentState["last_checkpoint"]
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
	s.whiteImage = ebiten.NewImage(1, 1)
	s.whiteImage.Fill(color.Gray{255})
	return nil
}

func (s *MapScreen) moveBy(d m.Delta) {
	loc := s.Menu.World.Level.CheckpointLocations
	cpLoc := loc.Locs[s.CurrentCP]
	edge, found := cpLoc.NextByDir[d]
	if !found {
		return
	}
	edgeSeen := s.Menu.World.Level.Player.PersistentState["checkpoints_walked."+s.CurrentCP+"."+edge.Other] != ""
	reverseSeen := s.Menu.World.Level.Player.PersistentState["checkpoints_walked."+edge.Other+"."+s.CurrentCP] != ""
	if !edgeSeen && !reverseSeen {
		// Don't know this yet :)
		return
	}
	if s.Menu.World.Level.Checkpoints[edge.Other].Properties["dead_end"] == "true" {
		// A dead end!
		return
	}
	s.CurrentCP = edge.Other
}

func (s *MapScreen) Update() error {
	if input.Exit.JustHit {
		return s.Menu.SwitchToScreen(&MainScreen{})
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
		return s.Menu.SwitchToCheckpoint(s.CurrentCP)
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
	font.MenuBig.Draw(screen, "Pick-a-Path", m.Pos{X: x, Y: h / 7}, true, fgs, bgs)

	// Draw all known checkpoints.
	loc := s.Menu.World.Level.CheckpointLocations
	mapSize := 2 * h / 3
	mapWidth := mapSize
	mapHeight := mapSize
	if loc.Rect.Size.DX > loc.Rect.Size.DY {
		mapHeight = mapWidth * loc.Rect.Size.DY / loc.Rect.Size.DX
	} else {
		mapWidth = mapHeight * loc.Rect.Size.DX / loc.Rect.Size.DY
	}
	mapRect := m.Rect{
		Origin: m.Pos{X: (w - mapWidth) / 2, Y: ((h - h/4) - mapHeight) / 2},
		Size:   m.Delta{DX: mapWidth, DY: mapHeight},
	}
	// First draw all edges.
	for cpName, cpLoc := range loc.Locs {
		cpSeen := s.Menu.World.Level.Player.PersistentState["checkpoint_seen."+cpName] != ""
		if !cpSeen {
			continue
		}
		pos := cpLoc.MapPos.FromRectToRect(loc.Rect, mapRect)
		for dir, edge := range cpLoc.NextByDir {
			if !edge.Forward {
				continue
			}
			otherName := edge.Other
			edgeSeen := s.Menu.World.Level.Player.PersistentState["checkpoints_walked."+cpName+"."+otherName] != ""
			closePos := pos.Add(dir.Mul(7))
			options := &ebiten.DrawTrianglesOptions{
				CompositeMode: ebiten.CompositeModeSourceOver,
				Filter:        ebiten.FilterNearest,
				Address:       ebiten.AddressUnsafe,
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
	for cpName, cpLoc := range loc.Locs {
		cpSeen := s.Menu.World.Level.Player.PersistentState["checkpoint_seen."+cpName] != ""
		if !cpSeen {
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
}

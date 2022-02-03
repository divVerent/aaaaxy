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
	"math/rand"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/fun"
	"github.com/divVerent/aaaaxy/internal/image"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/level"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/playerstate"
)

type MapScreen struct {
	Controller  *Controller
	Level       *level.Level
	CurrentCP   string
	SortedLocs  []string
	SortedEdges map[string][]level.CheckpointEdge
	CPPos       map[string]m.Pos
	MapRect     m.Rect

	cpSprite                *ebiten.Image
	cpSelectedSprite        *ebiten.Image
	cpFlippedSprite         *ebiten.Image
	cpFlippedSelectedSprite *ebiten.Image
	deadEndSprite           *ebiten.Image
	cpCheckmarkSprite       *ebiten.Image
	whiteImage              *ebiten.Image
}

// TODO: parametrize.
const (
	firstCP = "leap_of_faith"

	edgeFarAttachDistance = 7
	edgeThickness         = 3
	mouseDistance         = 7
)

func (s *MapScreen) Init(c *Controller) error {
	s.Controller = c
	s.CurrentCP = s.Controller.World.PlayerState.LastCheckpoint()
	if s.CurrentCP == "" {
		// Have no checkpoint yet - start the game right away.
		return s.Controller.SwitchToGame()
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
	s.cpFlippedSprite, err = image.Load("sprites", "checkpoint_flipped.png")
	if err != nil {
		return err
	}
	s.cpFlippedSelectedSprite, err = image.Load("sprites", "checkpoint_flipped_selected.png")
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
	for name := range s.Controller.World.Level.Checkpoints {
		if name == "" {
			continue
		}
		s.SortedLocs = append(s.SortedLocs, name)
	}
	// Note: we do not care for the actual order, just that it does not change between frames.
	sort.Strings(s.SortedLocs)
	// Now also yield a deterministic edge order.
	s.SortedEdges = make(map[string][]level.CheckpointEdge, len(s.SortedLocs))
	loc := s.Controller.World.Level.CheckpointLocations
	for _, cpName := range s.SortedLocs {
		cpLoc := loc.Locs[cpName]
		edges := make([]level.CheckpointEdge, 0, len(cpLoc.NextByDir)+len(cpLoc.NextDeadEnds))
		for _, edge := range cpLoc.NextByDir {
			edges = append(edges, edge)
		}
		edges = append(edges, cpLoc.NextDeadEnds...)
		sort.Slice(edges, func(i, j int) bool {
			return edges[i].Other < edges[j].Other
		})
		s.SortedEdges[cpName] = edges
	}

	mapWidth := engine.GameWidth
	mapHeight := engine.GameHeight / 2
	if mapWidth*loc.Rect.Size.DY > mapHeight*loc.Rect.Size.DX {
		mapWidth = mapHeight * loc.Rect.Size.DX / loc.Rect.Size.DY
	} else {
		mapHeight = mapWidth * loc.Rect.Size.DY / loc.Rect.Size.DX
	}
	s.MapRect = m.Rect{
		Origin: m.Pos{X: (engine.GameWidth - mapWidth) / 2, Y: (engine.GameHeight - mapHeight) / 2},
		Size:   m.Delta{DX: mapWidth, DY: mapHeight},
	}

	s.CPPos = make(map[string]m.Pos, len(s.SortedLocs))
	for _, cpName := range s.SortedLocs {
		cpLoc := loc.Locs[cpName]
		pos := cpLoc.MapPos.FromRectToRect(loc.Rect, s.MapRect)
		s.CPPos[cpName] = pos
	}

	return nil
}

func (s *MapScreen) moveBy(d m.Delta) {
	loc := s.Controller.World.Level.CheckpointLocations
	cpLoc := loc.Locs[s.CurrentCP]
	edge, found := cpLoc.NextByDir[d]
	if !found {
		return
	}
	edgeSeen := s.Controller.World.PlayerState.CheckpointsWalked(s.CurrentCP, edge.Other)
	if !edgeSeen {
		// Don't know this yet :)
		return
	}
	if s.Controller.World.Level.Checkpoints[edge.Other].Properties["dead_end"] == "true" {
		// A dead end!
		return
	}
	s.CurrentCP = edge.Other
	s.Controller.MoveSound(nil)
}

func (s *MapScreen) moveTo(pos m.Pos) bool {
	hitCP := ""
	for _, cpName := range s.SortedLocs {
		if s.Controller.World.PlayerState.CheckpointSeen(cpName) == playerstate.NotSeen {
			// Don't know this yet :)
			continue
		}
		if s.Controller.World.Level.Checkpoints[cpName].Properties["dead_end"] == "true" {
			// A dead end!
			continue
		}
		cpPos := s.CPPos[cpName]
		d := cpPos.Delta(pos).Length2()
		if d < mouseDistance*mouseDistance {
			if hitCP != "" {
				return false // Not unique.
			}
			hitCP = cpName
		}
		cpName = cpName
	}
	if hitCP == "" {
		return false // Nothing hit.
	}
	if hitCP == s.CurrentCP {
		// No change.
		return true
	}
	s.CurrentCP = hitCP
	s.Controller.MoveSound(nil)
	return true
}

func (s *MapScreen) exit() error {
	if s.CurrentCP != firstCP && s.Controller.World.PlayerState.CheckpointSeen(firstCP) != playerstate.NotSeen {
		s.CurrentCP = firstCP
		return s.Controller.MoveSound(nil)
	}
	return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&MainScreen{}))
}

func (s *MapScreen) Update() error {
	mousePos, mouseState := input.Mouse()
	clicked := false
	if mouseState != input.NoMouse {
		clicked = s.moveTo(mousePos)
	}
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
	if input.Jump.JustHit || input.Action.JustHit || (clicked && mouseState == input.ClickingMouse) {
		return s.Controller.ActivateSound(s.Controller.SwitchToCheckpoint(s.CurrentCP))
	}
	return nil
}

func (s *MapScreen) Draw(screen *ebiten.Image) {
	h := engine.GameHeight
	w := engine.GameWidth
	x := w / 2
	fgs := color.NRGBA{R: 255, G: 255, B: 85, A: 255}
	bgs := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	takenRouteColor := color.NRGBA{R: 170, G: 170, B: 170, A: 255}
	selectedRouteColor := color.NRGBA{R: 255, G: 255, B: 85, A: 255}
	unseenPathToSeenCPColor := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	unseenPathToUnseenCPColor := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	font.MenuBig.Draw(screen, "Pick-a-Path", m.Pos{X: x, Y: h / 8}, true, fgs, bgs)
	cpText := fun.FormatText(&s.Controller.World.PlayerState, s.Controller.World.Level.Checkpoints[s.CurrentCP].Properties["text"])
	seen, total := s.Controller.World.PlayerState.TnihSignsSeen(s.CurrentCP)
	if total > 0 {
		cpText += fmt.Sprintf(" (%d/%d)", seen, total)
	}
	font.Menu.Draw(screen, cpText, m.Pos{X: x, Y: 7 * h / 8}, true, fgs, bgs)

	// Draw all known checkpoints.
	opts := ebiten.DrawImageOptions{
		CompositeMode: ebiten.CompositeModeSourceOver,
		Filter:        ebiten.FilterNearest,
	}
	opts.GeoM.Scale(float64(s.MapRect.Size.DX+24), float64(s.MapRect.Size.DY+24))
	opts.GeoM.Translate(float64(s.MapRect.Origin.X-12), float64(s.MapRect.Origin.Y-12))
	opts.ColorM.Scale(1.0/3.0, 1.0/3.0, 1.0/3.0, 2.0/3.0)
	screen.DrawImage(s.whiteImage, &opts)
	// First draw all edges.
	cpPos := make(map[string]m.Pos, len(s.SortedLocs))
	for cpName, pos := range s.CPPos {
		if s.Controller.World.Level.Checkpoints[cpName].Properties["hub"] == "true" {
			pos.X += rand.Intn(3) - 1
			pos.Y += rand.Intn(3) - 1
		}
		cpPos[cpName] = pos
	}
	for z := 0; z < 3; z++ {
		for _, cpName := range s.SortedLocs {
			if s.Controller.World.PlayerState.CheckpointSeen(cpName) == playerstate.NotSeen {
				continue
			}
			pos := cpPos[cpName]
			for _, edge := range s.SortedEdges[cpName] {
				// We only draw forward non-optional edges; all others are for keyboard navigation only.
				if !edge.Forward || edge.Optional {
					continue
				}
				otherName := edge.Other
				edgeSeen := s.Controller.World.PlayerState.CheckpointsWalked(cpName, otherName)
				// Unseen edges leading to a secret are only drawn if the game has already been completed fully (Any% is not enough).
				if !edgeSeen && s.Controller.World.Level.Checkpoints[otherName].Properties["secret"] == "true" {
					if !s.Controller.World.PlayerState.SpeedrunCategories().ContainAll(playerstate.AnyPercentSpeedrun | playerstate.AllCheckpointsSpeedrun) {
						continue
					}
					// Even if we draw them, make edges pointing at secrets flicker.
					if rand.Intn(2) == 0 {
						continue
					}
				}
				endPos := cpPos[otherName]
				color := takenRouteColor
				switch z {
				case 0:
					// Selected edge is drawn first so it is clear what overlaps it.
					if !edgeSeen || !(cpName == s.CurrentCP || otherName == s.CurrentCP) {
						continue
					}
					color = selectedRouteColor
				case 1:
					// Normal edges are drawn next.
					if !edgeSeen || (cpName == s.CurrentCP || otherName == s.CurrentCP) {
						continue
					}
				case 2:
					// Missing edges are drawn last so one can always see them.
					if edgeSeen {
						continue
					}
					if s.Controller.World.PlayerState.CheckpointSeen(otherName) == playerstate.NotSeen {
						color = unseenPathToUnseenCPColor
					} else {
						color = unseenPathToSeenCPColor
					}
					endPos = pos.Add(endPos.Delta(pos).WithMaxLengthFixed(m.NewFixed(edgeFarAttachDistance)))
				}
				options := &ebiten.DrawTrianglesOptions{
					CompositeMode: ebiten.CompositeModeSourceOver,
					Filter:        ebiten.FilterNearest,
				}
				geoM := &ebiten.GeoM{}
				geoM.Scale(0, 0)
				engine.DrawPolyLine(screen, edgeThickness, []m.Pos{pos, endPos}, s.whiteImage, color, geoM, options)
			}
		}
	}
	// Then draw the CPs.
	for _, cpName := range s.SortedLocs {
		var sprite *ebiten.Image
		switch s.Controller.World.PlayerState.CheckpointSeen(cpName) {
		case playerstate.NotSeen:
			continue
		case playerstate.SeenNormal:
			if cpName == s.CurrentCP {
				sprite = s.cpSelectedSprite
			} else {
				sprite = s.cpSprite
			}
		case playerstate.SeenFlipped:
			if cpName == s.CurrentCP {
				sprite = s.cpFlippedSelectedSprite
			} else {
				sprite = s.cpFlippedSprite
			}
		}
		if s.Controller.World.Level.Checkpoints[cpName].Properties["dead_end"] == "true" {
			sprite = s.deadEndSprite
		}
		pos := cpPos[cpName]
		opts := ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeSourceOver,
			Filter:        ebiten.FilterNearest,
		}
		opts.GeoM.Translate(float64(pos.X-7), float64(pos.Y-7))
		screen.DrawImage(sprite, &opts)
	}
	// Finally the checkmarks.
	for _, cpName := range s.SortedLocs {
		if s.Controller.World.PlayerState.CheckpointSeen(cpName) == playerstate.NotSeen {
			continue
		}
		if seen, total := s.Controller.World.PlayerState.TnihSignsSeen(cpName); seen < total {
			continue
		}
		sprite := s.cpCheckmarkSprite
		pos := cpPos[cpName]
		opts := ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeSourceOver,
			Filter:        ebiten.FilterNearest,
		}
		opts.GeoM.Translate(float64(pos.X-7), float64(pos.Y-7))
		screen.DrawImage(sprite, &opts)
	}
}

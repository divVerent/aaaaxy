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
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/font"
	"github.com/divVerent/aaaaaa/internal/input"
	m "github.com/divVerent/aaaaaa/internal/math"
)

type MapScreen struct {
	Menu      *Menu
	Level     *engine.Level
	CurrentCP string
}

func (s *MapScreen) Init(m *Menu) error {
	s.Menu = m
	s.CurrentCP = s.Menu.World.Level.Player.PersistentState["last_checkpoint"]
	if s.CurrentCP == "" {
		// Have no checkpoint yet - start the game right away.
		return s.Menu.SwitchToGame()
	}
	return nil
}

func (s *MapScreen) moveBy(d m.Delta) {
	loc := s.Menu.World.Level.CheckpointLocations
	cpLoc := loc.Locs[s.CurrentCP]
	edge, found := cpLoc.NextByDir[d]
	if !found {
		return
	}
	otherSeen := s.Menu.World.Level.Player.PersistentState["checkpoint_seen."+edge.Other] != ""
	if !otherSeen {
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
	fgs := color.NRGBA{R: 255, G: 255, B: 85, A: 255}
	bgs := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	font.MenuBig.Draw(screen, "Pick-a-Path", m.Pos{X: x, Y: h / 8}, true, fgs, bgs)

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
			otherSeen := s.Menu.World.Level.Player.PersistentState["checkpoint_seen."+otherName] != ""
			closePos := pos.Add(dir.Mul(5))
			if otherSeen {
				otherLoc := loc.Locs[otherName]
				otherPos := otherLoc.MapPos.FromRectToRect(loc.Rect, mapRect)
				farPos := otherPos.Sub(dir.Mul(5))
				ebitenutil.DrawLine(screen, float64(pos.X), float64(pos.Y), float64(closePos.X), float64(closePos.Y), fgs)
				ebitenutil.DrawLine(screen, float64(closePos.X), float64(closePos.Y), float64(farPos.X), float64(farPos.Y), fgs)
				ebitenutil.DrawLine(screen, float64(farPos.X), float64(farPos.Y), float64(otherPos.X), float64(otherPos.Y), fgs)
			} else {
				ebitenutil.DrawLine(screen, float64(pos.X), float64(pos.Y), float64(closePos.X), float64(closePos.Y), bgs)
			}
		}
	}
	// Then draw the CPs.
	for cpName, cpLoc := range loc.Locs {
		cpSeen := s.Menu.World.Level.Player.PersistentState["checkpoint_seen."+cpName] != ""
		if !cpSeen {
			continue
		}
		cpSprite := "checkpoint"
		if s.Menu.World.Level.Checkpoints[cpName].Properties["dead_end"] == "true" {
			cpSprite = "dead_end"
		}
		pos := cpLoc.MapPos.FromRectToRect(loc.Rect, mapRect)
		if cpName == s.CurrentCP {
			cpSprite += "_selected"
		}
		// TODO Draw sprite!
		if cpName == s.CurrentCP {
			ebitenutil.DrawRect(screen, float64(pos.X-4), float64(pos.Y-4), 8, 8, fgs)
		} else {
			ebitenutil.DrawRect(screen, float64(pos.X-2), float64(pos.Y-2), 4, 4, fgs)
		}
	}
}

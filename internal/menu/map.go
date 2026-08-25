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
	"cmp"
	"fmt"
	"image/color"
	"maps"
	"math/rand"
	"slices"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/fun"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/m"
	"github.com/divVerent/aaaaxy/internal/palette"
	"github.com/divVerent/aaaaxy/internal/picture"
	"github.com/divVerent/aaaaxy/internal/playerstate"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

var (
	debugAntiAlias = flag.Bool("debug_anti_alias", true, "allow anti aliasing")
)

type MapScreen struct {
	Controller  *Controller
	Level       *level.Level
	FirstCP     string
	CurrentCP   string
	SortedLocs  []string
	SortedEdges map[string][]level.CheckpointEdge
	CPPos       map[string]m.Pos
	MapRect     m.Rect
	WalkFrame   int

	cpSprite                *ebiten.Image
	cpSelectedSprite        *ebiten.Image
	cpFlippedSprite         *ebiten.Image
	cpFlippedSelectedSprite *ebiten.Image
	deadEndSprite           *ebiten.Image
	cpCheckmarkSprite       *ebiten.Image

	nameHovered bool
}

// TODO: parametrize.
const (
	edgeFarAttachDistance = 9
	edgeThickness         = 3
	mouseDistance         = 16
	walkSpeed             = 0.2
	mapBorder             = 9
)

func (s *MapScreen) Init(c *Controller) error {
	s.Controller = c
	s.CurrentCP = s.Controller.World.PlayerState.LastCheckpoint()
	if s.CurrentCP == "" {
		// Have no checkpoint yet - start the game right away.
		return s.Controller.SwitchToGame()
	}
	var err error
	s.cpSprite, err = picture.Load("sprites", "checkpoint.png")
	if err != nil {
		return err
	}
	s.cpSelectedSprite, err = picture.Load("sprites", "checkpoint_selected.png")
	if err != nil {
		return err
	}
	s.cpFlippedSprite, err = picture.Load("sprites", "checkpoint_flipped.png")
	if err != nil {
		return err
	}
	s.cpFlippedSelectedSprite, err = picture.Load("sprites", "checkpoint_flipped_selected.png")
	if err != nil {
		return err
	}
	s.deadEndSprite, err = picture.Load("sprites", "dead_end.png")
	if err != nil {
		return err
	}
	s.cpCheckmarkSprite, err = picture.Load("sprites", "checkpoint_checkmark.png")
	if err != nil {
		return err
	}

	var parseErr error

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
		if propmap.ValueOrP(s.Controller.World.Level.Checkpoints[cpName].Properties, "first", false, &parseErr) {
			if s.FirstCP != "" {
				return fmt.Errorf("more than one first checkpoint is not allowed: got %v and %v, want only one", s.FirstCP, cpName)
			}
			s.FirstCP = cpName
		}
	}

	mapWidth := engine.GameWidth
	mapHeight := 3 * engine.GameHeight / 4
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

	return parseErr
}

func (s *MapScreen) moveBy(d m.Delta) {
	loc := s.Controller.World.Level.CheckpointLocations
	cpLoc := loc.Locs[s.CurrentCP]
	edge, found := cpLoc.NextByDir[d]
	if !found {
		return
	}
	edgeSeen := s.Controller.World.PlayerState.CheckpointsWalked(s.CurrentCP, edge.Other)
	otherSeen := s.Controller.World.PlayerState.CheckpointSeen(edge.Other) != playerstate.NotSeen
	if !edgeSeen && !otherSeen {
		// Don't know this yet :)
		return
	}
	if propmap.ValueOrP(s.Controller.World.Level.Checkpoints[edge.Other].Properties, "dead_end", false, nil) {
		// A dead end!
		return
	}
	s.CurrentCP = edge.Other
	s.Controller.MoveSound(nil)
}

func (s *MapScreen) moveTo(pos m.Pos) bool {
	hitCP := ""
	hitD := int64(mouseDistance * mouseDistance)
	for _, cpName := range s.SortedLocs {
		if s.Controller.World.PlayerState.CheckpointSeen(cpName) == playerstate.NotSeen {
			// Don't know this yet :)
			continue
		}
		if propmap.ValueOrP(s.Controller.World.Level.Checkpoints[cpName].Properties, "dead_end", false, nil) {
			// A dead end!
			continue
		}
		cpPos := s.CPPos[cpName]
		d := cpPos.Delta(pos).Length2()
		if d < hitD {
			hitCP = cpName
			hitD = d
		}
	}
	if hitCP == "" && pos.Y > s.MapRect.OppositeCorner().Y+mapBorder {
		// TODO: this may be an off-by-one error?
		s.nameHovered = true
		return true
	}
	if hitCP == "" {
		s.nameHovered = false
		return false // Nothing hit.
	}
	s.nameHovered = true
	if hitCP == s.CurrentCP {
		// No change.
		return true
	}
	s.CurrentCP = hitCP
	s.Controller.MoveSound(nil)
	return true
}

func (s *MapScreen) exit() error {
	if s.CurrentCP != s.FirstCP && s.Controller.World.PlayerState.CheckpointSeen(s.FirstCP) != playerstate.NotSeen {
		s.CurrentCP = s.FirstCP
		return s.Controller.MoveSound(nil)
	}
	return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&MainScreen{}))
}

func (s *MapScreen) Update() error {
	s.WalkFrame++
	mousePos, mouseState := input.Mouse()
	clicked := false
	if mouseState == input.NoMouse {
		s.nameHovered = false
	} else {
		clicked = s.moveTo(mousePos)
	}
	if input.Exit.JustHit || (!clicked && mouseState == input.ClickingMouse) {
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
	fgn := palette.EGA(palette.LightGrey, 255)
	bgn := palette.EGA(palette.DarkGrey, 255)
	fgs := palette.EGA(palette.Yellow, 255)
	bgs := palette.EGA(palette.Black, 255)
	takenRouteColor := palette.EGA(palette.LightGrey, 255)
	selectedRouteColor := palette.EGA(palette.Yellow, 255)
	unseenPathToSeenCPColor := palette.EGA(palette.White, 255)
	unseenPathToUnseenCPColor := palette.EGA(palette.Black, 255)
	unseenPathBlinkColor := palette.EGA(palette.DarkGrey, 255)
	font.ByName["MenuBig"].Draw(screen, locale.G.Get("Pick-a-Path"), m.Pos{X: x, Y: h / 12}, font.Center, fgs, bgs)
	cpText := fun.FormatText(&s.Controller.World.PlayerState, propmap.ValueP(s.Controller.World.Level.Checkpoints[s.CurrentCP].Properties, "text", "", nil))
	seen, total := s.Controller.World.PlayerState.TnihSignsSeen(s.CurrentCP)
	if total > 0 {
		cpText = locale.G.Get("%s (%d/%d)", cpText, seen, total)
	}
	fg, bg := fgn, bgn
	if s.nameHovered {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, cpText, m.Pos{X: x, Y: 11*h/12 + 12}, font.Center, fg, bg)

	// Draw all known checkpoints.
	opts := ebiten.DrawImageOptions{
		Blend:  ebiten.BlendSourceOver,
		Filter: ebiten.FilterNearest,
	}
	opts.GeoM.Scale(float64(s.MapRect.Size.DX+2*mapBorder), float64(s.MapRect.Size.DY+2*mapBorder))
	opts.GeoM.Translate(float64(s.MapRect.Origin.X-mapBorder), float64(s.MapRect.Origin.Y-mapBorder))
	opts.ColorScale.Scale(2.0/9.0, 2.0/9.0, 2.0/9.0, 2.0/3.0) // Color: #555555 at 2/3 alpha.
	screen.DrawImage(s.Controller.WhiteImage, &opts)

	// First draw all edges.
	cpPos := make(map[string]m.Pos, len(s.SortedLocs))
	for cpName, pos := range s.CPPos {
		if propmap.ValueOrP(s.Controller.World.Level.Checkpoints[cpName].Properties, "hub", false, nil) {
			pos.X += rand.Intn(3) - 1
			pos.Y += rand.Intn(3) - 1
		}
		cpPos[cpName] = pos
	}
	focusMissing := s.Controller.World.PlayerState.SpeedrunCategories().ContainAll(playerstate.AllCheckpointsSpeedrun)
	revealSecrets := s.Controller.World.PlayerState.SpeedrunCategories().ContainAll(playerstate.AnyPercentSpeedrun | playerstate.AllCheckpointsSpeedrun)
	focusSecrets := s.Controller.World.PlayerState.SpeedrunCategories().ContainAll(playerstate.AnyPercentSpeedrun | playerstate.AllCheckpointsSpeedrun | playerstate.AllPathsSpeedrun)

	paths := map[color.NRGBA]*vector.Path{}
	pathFor := func(col color.NRGBA) *vector.Path {
		path := paths[col]
		if path == nil {
			path = &vector.Path{}
			paths[col] = path
		}
		return path
	}

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
			isSecret := propmap.ValueOrP(s.Controller.World.Level.Checkpoints[otherName].Properties, "secret", false, nil)
			// Unseen edges leading to a secret are only drawn if the game has already been completed fully (Any% is not enough).
			if !edgeSeen && isSecret && !revealSecrets {
				continue
			}
			startPos := pos
			endPos := cpPos[otherName]
			var startPos2, endPos2 m.Pos
			col := takenRouteColor

			if edgeSeen {
				if cpName == s.CurrentCP || otherName == s.CurrentCP {
					col = selectedRouteColor
				} else {
					col = takenRouteColor
				}
			} else {
				otherSeen := s.Controller.World.PlayerState.CheckpointSeen(otherName) != playerstate.NotSeen
				if cpName == s.CurrentCP || otherName == s.CurrentCP {
					col = selectedRouteColor
				} else if otherSeen {
					col = unseenPathToSeenCPColor
				} else {
					col = unseenPathToUnseenCPColor
				}

				if focusMissing && rand.Intn(2) == 0 {
					if focusSecrets {
						// Once all CPs are hit and paths are set, blink secret paths.
						if isSecret {
							col = unseenPathBlinkColor
						}
					} else if revealSecrets {
						// Once all CPs are hit, show secret paths and blink missing paths.
						if !isSecret {
							col = unseenPathBlinkColor
						}
					} else {
						// By default, hide secret paths and blink missing paths.
						col = unseenPathBlinkColor
					}
				}

				dp := endPos.Delta(startPos)
				section := m.NewFixed(edgeFarAttachDistance)
				length := dp.LengthFixed()
				if otherSeen {
					if length < section {
						// Leave endPos as is. We would make it longer.
					} else {
						// Animate missing paths when the other side is seen to indicate direction.
						a := m.NewFixed(s.WalkFrame).Mul(m.NewFixedFloat64(walkSpeed)).Mod(length)
						b := (a + section).Mod(length)
						if a < b {
							startPos, endPos = pos.Add(dp.WithLengthFixed(a)), pos.Add(dp.WithLengthFixed(b))
						} else {
							startPos2, endPos2 = pos.Add(dp.WithLengthFixed(a)), endPos
							startPos, endPos = pos, pos.Add(dp.WithLengthFixed(b))
						}
					}
				} else {
					// Don't reveal actual CP location.
					endPos = pos.Add(dp.WithLengthFixed(section))
				}
			}

			// Adding 0.5, as vector module thinks in pixel corners,
			// and CP coordinates are pixel centers.
			path := pathFor(col)
			path.MoveTo(float32(startPos.X)+0.5, float32(startPos.Y)+0.5)
			path.LineTo(float32(endPos.X)+0.5, float32(endPos.Y)+0.5)
			if startPos2 != endPos2 {
				path.MoveTo(float32(startPos2.X)+0.5, float32(startPos2.Y)+0.5)
				path.LineTo(float32(endPos2.X)+0.5, float32(endPos2.Y)+0.5)
			}
		}
	}

	for _, col := range slices.SortedFunc(maps.Keys(paths), func(a, b color.NRGBA) int {
		// Reverse brightness comparison.
		return -cmp.Or(
			cmp.Compare(a.A, b.A),
			cmp.Compare(a.G, b.G),
			cmp.Compare(a.R, b.R),
			cmp.Compare(a.B, b.B),
		)
	}) {
		path := paths[col]
		strokeOptions := vector.StrokeOptions{
			Width: edgeThickness,
		}
		drawPathOptions := vector.DrawPathOptions{
			AntiAlias: *debugAntiAlias,
		}
		drawPathOptions.ColorScale.ScaleWithColor(col)
		vector.StrokePath(screen, path, &strokeOptions, &drawPathOptions)
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
		if propmap.ValueOrP(s.Controller.World.Level.Checkpoints[cpName].Properties, "dead_end", false, nil) {
			sprite = s.deadEndSprite
		}
		pos := cpPos[cpName]
		opts := ebiten.DrawImageOptions{
			Blend:  ebiten.BlendSourceOver,
			Filter: ebiten.FilterNearest,
		}
		opts.GeoM.Translate(float64(pos.X-7), float64(pos.Y-7))
		if propmap.ValueOrP(s.Controller.World.Level.Checkpoints[cpName].Properties, "final", false, nil) {
			c := rand.Float32() * 2.0
			opts.ColorScale.Scale(c, c, c, 1.0)
		}
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
			Blend:  ebiten.BlendSourceOver,
			Filter: ebiten.FilterNearest,
		}
		opts.GeoM.Translate(float64(pos.X-7), float64(pos.Y-7))
		screen.DrawImage(sprite, &opts)
	}
}

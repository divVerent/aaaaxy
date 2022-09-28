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

package misc

import (
	"fmt"
	"image/color"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/constants"
	"github.com/divVerent/aaaaxy/internal/level"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

// SpriteBase is a base class for sprites.
// To instantiate it, just set the entity image, then forward to this.
type SpriteBase struct {
	ZDefault int
}

func (s *SpriteBase) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	var parseErr error
	w.SetSolid(e, propmap.ValueOrP(sp.Properties, "solid", false, &parseErr))
	w.SetOpaque(e, propmap.ValueOrP(sp.Properties, "opaque", false, &parseErr))
	if t := propmap.ValueOrP(sp.Properties, "player_solid", propmap.TriState{}, &parseErr); t.Active {
		w.MutateContentsBool(e, level.PlayerSolidContents, t.Value)
	}
	if t := propmap.ValueOrP(sp.Properties, "object_solid", propmap.TriState{}, &parseErr); t.Active {
		w.MutateContentsBool(e, level.ObjectSolidContents, t.Value)
	}
	e.Alpha = propmap.ValueOrP(sp.Properties, "alpha", 1.0, &parseErr)
	mapBlackTo := propmap.ValueOrP(sp.Properties, "map_black_to", color.NRGBA{R: 0, G: 0, B: 0, A: 0}, &parseErr)
	e.ColorAdd[0] = float64(mapBlackTo.R) / 255.0
	e.ColorAdd[1] = float64(mapBlackTo.G) / 255.0
	e.ColorAdd[2] = float64(mapBlackTo.B) / 255.0
	e.ColorAdd[3] = float64(mapBlackTo.A) / 255.0
	mapWhiteTo := propmap.ValueOrP(sp.Properties, "map_white_to", color.NRGBA{R: 255, G: 255, B: 255, A: 255}, &parseErr)
	e.ColorMod[0] = float64(mapWhiteTo.R)/255.0 - e.ColorAdd[0]
	e.ColorMod[1] = float64(mapWhiteTo.G)/255.0 - e.ColorAdd[1]
	e.ColorMod[2] = float64(mapWhiteTo.B)/255.0 - e.ColorAdd[2]
	e.ColorMod[3] = float64(mapWhiteTo.A)/255.0 - e.ColorAdd[3]
	z := propmap.ValueOrP(sp.Properties, "z_index", s.ZDefault, &parseErr)
	if z != s.ZDefault && (z < constants.MinSpriteZ || z > constants.MaxSpriteZ) {
		return fmt.Errorf("z index out of range: got %v, want %v..%v", z, constants.MinSpriteZ, constants.MaxSpriteZ)
	}
	w.SetZIndex(e, z)
	if propmap.ValueOrP(sp.Properties, "no_transform", false, &parseErr) {
		// Undo transform of orientation by tile.
		e.Orientation = sp.Orientation
	}
	if e.Transform.Determinant() < 0 {
		// e.Orientation: in-editor transform. Applied first.
		// Normally the formula is e.Transform.Inverse().Concat(e.Orientation).
		// This flips the view on the _image_ X axis.
		flip := propmap.StringOr(sp.Properties, "no_flip", "")
		switch flip {
		case "x":
			e.Orientation = e.Orientation.Concat(m.FlipX())
		case "y":
			e.Orientation = e.Orientation.Concat(m.FlipY())
		case "", "false":
			// Nothing to do.
		default:
			return fmt.Errorf("invalid no_flip value: got %v, want one of empty, x, y, false", flip)
		}
	}

	// Field contains orientation OF THE PLAYER to make it easier in the map editor.
	// So it is actually a transform as far as this code is concerned.
	requiredTransforms := propmap.ValueOrP(sp.Properties, "required_orientation", m.Orientations{}, &parseErr)
	if len(requiredTransforms) != 0 {
		show := false
		for _, requiredTransform := range requiredTransforms {
			if e.Transform == requiredTransform {
				show = true
			} else if e.Transform == requiredTransform.Concat(m.FlipX()) {
				show = true
			}
		}
		if !show {
			// Hide.
			e.Alpha = 0.0
			w.MutateContentsBool(e, level.AllContents, false)
		}
	}

	return parseErr
}

// The other methods to reduce code duplication in implementors.

func (s *SpriteBase) Despawn() {}

func (s *SpriteBase) Update() {}

func (s *SpriteBase) Touch(other *engine.Entity) {}

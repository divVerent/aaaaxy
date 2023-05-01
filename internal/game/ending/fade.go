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

package ending

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2/colorm"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

// FadeTarget fades the screen out.
type FadeTarget struct {
	World *engine.World

	Frames int
	Frame  int
	State  bool

	ColorM colorm.ColorM
}

func (f *FadeTarget) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	f.World = w

	var parseErr error

	durationTime := propmap.ValueOrP(sp.Properties, "duration", time.Second, &parseErr)
	f.Frames = int((durationTime*engine.GameTPS + (time.Second / 2)) / time.Second)
	if f.Frames < 1 {
		f.Frames = 1
	}

	// We want a color matrix that maps A1 to A'1, B1 to B'1, C1 to C'1, D1 to D'1.
	// So we build two color matrices - fromM maps a pentahedron to 0 A0 B0 C0 D0, toM maps the same pentahedron to 0 A' B' C' D'.
	// Then toM * fromM^-1 will be what we need.
	var fromM, toM colorm.ColorM

	c := propmap.ValueP(sp.Properties, "from_color_a", color.NRGBA{}, &parseErr)
	fromM.SetElement(0, 0, float64(c.R)/255.0)
	fromM.SetElement(1, 0, float64(c.G)/255.0)
	fromM.SetElement(2, 0, float64(c.B)/255.0)
	fromM.SetElement(3, 0, 1.0)
	c = propmap.ValueP(sp.Properties, "from_color_b", color.NRGBA{}, &parseErr)
	fromM.SetElement(0, 1, float64(c.R)/255.0)
	fromM.SetElement(1, 1, float64(c.G)/255.0)
	fromM.SetElement(2, 1, float64(c.B)/255.0)
	fromM.SetElement(3, 1, 1.0)
	c = propmap.ValueP(sp.Properties, "from_color_c", color.NRGBA{}, &parseErr)
	fromM.SetElement(0, 2, float64(c.R)/255.0)
	fromM.SetElement(1, 2, float64(c.G)/255.0)
	fromM.SetElement(2, 2, float64(c.B)/255.0)
	fromM.SetElement(3, 2, 1.0)
	c = propmap.ValueP(sp.Properties, "from_color_d", color.NRGBA{}, &parseErr)
	fromM.SetElement(0, 3, float64(c.R)/255.0)
	fromM.SetElement(1, 3, float64(c.G)/255.0)
	fromM.SetElement(2, 3, float64(c.B)/255.0)
	fromM.SetElement(3, 3, 1.0)
	// In addition, add another row to keep the alpha channel invariant.
	fromM.SetElement(0, 4, float64(c.R)/255.0)
	fromM.SetElement(1, 4, float64(c.G)/255.0)
	fromM.SetElement(2, 4, float64(c.B)/255.0)
	fromM.SetElement(3, 4, 0.0)

	c = propmap.ValueP(sp.Properties, "to_color_a", color.NRGBA{}, &parseErr)
	toM.SetElement(0, 0, float64(c.R)/255.0)
	toM.SetElement(1, 0, float64(c.G)/255.0)
	toM.SetElement(2, 0, float64(c.B)/255.0)
	toM.SetElement(3, 0, 1.0)
	c = propmap.ValueP(sp.Properties, "to_color_b", color.NRGBA{}, &parseErr)
	toM.SetElement(0, 1, float64(c.R)/255.0)
	toM.SetElement(1, 1, float64(c.G)/255.0)
	toM.SetElement(2, 1, float64(c.B)/255.0)
	toM.SetElement(3, 1, 1.0)
	c = propmap.ValueP(sp.Properties, "to_color_c", color.NRGBA{}, &parseErr)
	toM.SetElement(0, 2, float64(c.R)/255.0)
	toM.SetElement(1, 2, float64(c.G)/255.0)
	toM.SetElement(2, 2, float64(c.B)/255.0)
	toM.SetElement(3, 2, 1.0)
	c = propmap.ValueP(sp.Properties, "to_color_d", color.NRGBA{}, &parseErr)
	toM.SetElement(0, 3, float64(c.R)/255.0)
	toM.SetElement(1, 3, float64(c.G)/255.0)
	toM.SetElement(2, 3, float64(c.B)/255.0)
	toM.SetElement(3, 3, 1.0)
	// In addition, add another row to keep the alpha channel invariant.
	toM.SetElement(0, 4, float64(c.R)/255.0)
	toM.SetElement(1, 4, float64(c.G)/255.0)
	toM.SetElement(2, 4, float64(c.B)/255.0)
	toM.SetElement(3, 4, 0.0)

	// Do not try inverting if parsing failed.
	if parseErr != nil {
		return parseErr
	}

	f.ColorM = fromM
	f.ColorM.Invert()
	f.ColorM.Concat(toM)

	return nil
}

func (f *FadeTarget) Despawn() {}

func (f *FadeTarget) Update() {
	if f.Frame <= 0 {
		return
	}
	f.Frame--

	factor := 1.0 - float64(f.Frame)/float64(f.Frames) // Is 1.0 in the last execution.

	// Linearly interpolate the matrix.
	var colorM colorm.ColorM
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			identity := 0.0
			if i == j {
				identity = 1.0
			}
			colorM.SetElement(i, j, f.ColorM.Element(i, j)*factor+identity*(1.0-factor))
		}
	}
	f.World.GlobalColorM.Concat(colorM)

	if f.Frame == 0 {
		// Keep showing this effect when at the end.
		f.Frame = 1
	}
}

func (f *FadeTarget) SetState(originator, predecessor *engine.Entity, state bool) {
	if state == f.State {
		return
	}
	f.State = state
	if state {
		f.Frame = f.Frames
	} else {
		f.Frame = 0
	}
}

func (f *FadeTarget) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&FadeTarget{})
}

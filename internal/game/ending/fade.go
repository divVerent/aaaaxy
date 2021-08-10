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
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/level"
)

// FadeTarget fades the screen out.
type FadeTarget struct {
	World *engine.World

	Frames int
	Frame  int
	State  bool

	Base   [3]float64
	Normal [3]float64
}

func (f *FadeTarget) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	f.World = w

	durationString := sp.Properties["duration"]
	durationTime, err := time.ParseDuration(durationString)
	if err != nil {
		return fmt.Errorf("could not parse duration time: %v", durationString)
	}
	f.Frames = int((durationTime*engine.GameTPS + (time.Second / 2)) / time.Second)
	if f.Frames < 1 {
		f.Frames = 1
	}

	var cA, cB, cC [3]float64

	var r, g, b, a int
	colorString := sp.Properties["invariant_color_a"]
	if _, err := fmt.Sscanf(colorString, "#%02x%02x%02x%02x", &a, &r, &g, &b); err != nil {
		return fmt.Errorf("could not decode color %q: %v", colorString, err)
	}
	cA[0] = float64(r) / 255.0
	cA[1] = float64(g) / 255.0
	cA[2] = float64(b) / 255.0
	colorString = sp.Properties["invariant_color_b"]
	if _, err := fmt.Sscanf(colorString, "#%02x%02x%02x%02x", &a, &r, &g, &b); err != nil {
		return fmt.Errorf("could not decode color %q: %v", colorString, err)
	}
	cB[0] = float64(r) / 255.0
	cB[1] = float64(g) / 255.0
	cB[2] = float64(b) / 255.0
	colorString = sp.Properties["invariant_color_c"]
	if _, err := fmt.Sscanf(colorString, "#%02x%02x%02x%02x", &a, &r, &g, &b); err != nil {
		return fmt.Errorf("could not decode color %q: %v", colorString, err)
	}
	cC[0] = float64(r) / 255.0
	cC[1] = float64(g) / 255.0
	cC[2] = float64(b) / 255.0

	f.Base[0] = cA[0]
	f.Base[1] = cA[1]
	f.Base[2] = cA[2]

	dB := [3]float64{cB[0] - cA[0], cB[1] - cA[1], cB[2] - cA[2]}
	dC := [3]float64{cC[0] - cA[0], cC[1] - cA[1], cC[2] - cA[2]}
	var n [3]float64
	n[0] = dB[1]*dC[2] - dB[2]*dC[1]
	n[1] = dB[2]*dC[0] - dB[0]*dC[2]
	n[2] = dB[0]*dC[1] - dB[1]*dC[0]
	l := math.Sqrt(n[0]*n[0] + n[1]*n[1] + n[2]*n[2])
	f.Normal[0] = n[0] / l
	f.Normal[1] = n[1] / l
	f.Normal[2] = n[2] / l

	return nil
}

func (f *FadeTarget) Despawn() {}

func (f *FadeTarget) Update() {
	if f.Frame <= 0 {
		return
	}
	f.Frame--
	// Fade AWAY from triangle.
	// normalFactor := float64(f.Frames)/(float64(f.Frame)+0.5) - 1.0 // Avoid division by zero.
	// Fade TO triangle.
	normalFactor := float64(f.Frame)/float64(f.Frames) - 1.0

	var colorM ebiten.ColorM
	delta := f.Base
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			isDiag := 0.0
			if i == j {
				isDiag = 1.0
			}
			e := isDiag + normalFactor*f.Normal[i]*f.Normal[j]
			colorM.SetElement(i, j, e)
			delta[i] -= e * f.Base[j]
		}
	}
	colorM.Translate(delta[0], delta[1], delta[2], 0.0)

	f.World.GlobalColorM = colorM
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

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
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/level"
)

// FadeTarget fades the screen out.
type FadeTarget struct {
	World *engine.World

	Color1, Color2, Color3 [3]float64
}

func (f *FadeTarget) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	f.World = w
	// Note: duration, color_1, color_2, color_3.
	// Precompute from this a base color and the vector to cancel (cross product of color diffs).
	return nil
}

func (f *FadeTarget) Despawn() {}

func (f *FadeTarget) Update() {
	a := 0.5 // Fraction of time passed.
	f := 1.0 / (1.0 - a)

	// Find color matrix so that:
	// - f = 0 maps Color1, Color2, Color3 to themselves, and all other colors to their plane.
	// - f = 1 is identity.
	// - Rest follows linearly.
	// So, form of matrix is:
	// M(x-A) + A
	// where
	// n = normalize((B-A) cross (C-A))
	//
	// Mx = x - (1-f)*n*dot(n, x)
	// M = id - (1-f)*n*n^T
	// T = A - MA

	// Once we have this, we can make the global ColorM an entity property, collect all entity ColorMs in a prepass, and then apply them to all object and world rendering.
}

func (f *FadeTarget) SetState(originator, predecessor *engine.Entity, state bool) {
	// TODO implement.
}

func (f *FadeTarget) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&FadeTarget{})
}

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
}

func (f *FadeTarget) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	f.World = w
	// Note: duration, map_to_black, map_to_white.
	return nil
}

func (f *FadeTarget) Despawn() {}

func (f *FadeTarget) Update() {}

func (f *FadeTarget) SetState(originator, predecessor *engine.Entity, state bool) {
	// TODO implement.
}

func (f *FadeTarget) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&FadeTarget{})
}

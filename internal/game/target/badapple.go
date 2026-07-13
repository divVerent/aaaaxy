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

package target

import (
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/level"
)

// BadAppleTarget prints the given text to console when activated.
// Setting state to ON saves the current text, setting state to OFF dumps it.
type BadAppleTarget struct {
	World *engine.World
}

func (b *BadAppleTarget) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	b.World = w
	return nil
}

func (b *BadAppleTarget) Despawn() {}

func (b *BadAppleTarget) Update() {}

func (b *BadAppleTarget) Touch(other *engine.Entity) {}

func (b *BadAppleTarget) SetState(originator, predecessor *engine.Entity, state bool) {
	if state {
		b.World.PlayerState.StartBadApple()
	}
}

func init() {
	engine.RegisterEntityType(&BadAppleTarget{})
}

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
	"github.com/divVerent/aaaaxy/internal/game/interfaces"
	"github.com/divVerent/aaaaxy/internal/level"
)

// Goal makes the player move towards it when activated.
type Goal struct {
	World  *engine.World
	Entity *engine.Entity
}

func (g *Goal) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	g.World = w
	g.Entity = e
	return nil
}

func (g *Goal) Despawn() {
	g.SetState(nil, nil, false)
}

func (g *Goal) Update() {}

func (g *Goal) SetState(originator, predecessor *engine.Entity, state bool) {
	if state {
		g.World.Player.Impl.(interfaces.SetGoaler).SetGoal(g.Entity)
	} else {
		g.World.Player.Impl.(interfaces.SetGoaler).SetGoal(nil)
	}
}

func (g *Goal) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&Goal{})
}

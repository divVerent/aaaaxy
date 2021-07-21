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

// StopTimerTarget makes the player move towards it when activated.
type StopTimerTarget struct {
	World *engine.World
}

func (s *StopTimerTarget) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	s.World = w
	return nil
}

func (s *StopTimerTarget) Despawn() {}

func (s *StopTimerTarget) Update() {}

func (s *StopTimerTarget) SetState(originator, predecessor *engine.Entity, state bool) {
	s.World.TimerStopped = state
}

func (s *StopTimerTarget) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&StopTimerTarget{})
}

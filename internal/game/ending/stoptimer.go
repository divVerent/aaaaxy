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
	"github.com/divVerent/aaaaxy/internal/fun"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/log"
)

// StopTimerTarget makes the player move towards it when activated.
type StopTimerTarget struct {
	World *engine.World
}

func (s *StopTimerTarget) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	s.World = w
	return nil
}

func (s *StopTimerTarget) Despawn() {}

func (s *StopTimerTarget) Update() {}

func (s *StopTimerTarget) SetState(originator, predecessor *engine.Entity, state bool) {
	if !state {
		if s.World.TimerStopped {
			log.Fatalf("Sorry, can't restart the timer again.")
		}
		return
	}
	if s.World.TimerStopped {
		return
	}
	s.World.TimerStopped = true
	s.World.PlayerState.SubFrame() // The ending frame doesn't count.

	s.World.PlayerState.SetWon()
	s.World.PlayerState.KickBadApple()

	err := s.World.Save()
	if err != nil {
		log.Errorf("could not save game: %v", err)
	}

	log.Infof("%v", fun.FormatText(&s.World.PlayerState,
		"your time: {{GameTime}}; your speedrun categories: {{SpeedrunCategories}}; try next: {{SpeedrunTryNext}}."))
}

func (s *StopTimerTarget) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&StopTimerTarget{})
}

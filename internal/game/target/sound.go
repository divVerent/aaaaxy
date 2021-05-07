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
	"fmt"
	"time"

	"github.com/divVerent/aaaaaa/internal/audiowrap"
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/mixins"
	"github.com/divVerent/aaaaaa/internal/level"
	"github.com/divVerent/aaaaaa/internal/sound"
)

// SoundTarget just changes the music track to the given one.
type SoundTarget struct {
	World  *engine.World
	Entity *engine.Entity

	Sound  *sound.Sound
	Player *audiowrap.Player

	Target mixins.TargetSelection
	State  bool
}

func (s *SoundTarget) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	s.World = w
	s.Entity = e
	var err error
	s.Sound, err = sound.Load(sp.Properties["sound"])
	if err != nil {
		return fmt.Errorf("could not load sound: %v", err)
	}
	s.Target = mixins.ParseTarget(sp.Properties["target"])
	s.State = sp.Properties["state"] != "false"
	return nil
}

func (s *SoundTarget) Despawn() {
	if s.Player != nil {
		s.Player.Close()
	}
}

func (s *SoundTarget) Update() {
	if s.Player != nil && !s.Player.IsPlaying() {
		s.Player.Close()
		s.Player = nil
		mixins.SetStateOfTarget(s.World, s.Entity, s.Target, !s.State)
	}
}

func (s *SoundTarget) Touch(other *engine.Entity) {}

func (s *SoundTarget) SetState(by *engine.Entity, state bool) {
	if state {
		if s.Player != nil {
			if s.Player.Current() < 100*time.Millisecond {
				return
			}
			s.Player.Close()
		}
		s.Player = s.Sound.Play()
		// Do we need this redirection?
		mixins.SetStateOfTarget(s.World, s.Entity, s.Target, s.State)
	} else {
		if s.Player != nil {
			s.Player.Close()
			s.Player = nil
			mixins.SetStateOfTarget(s.World, s.Entity, s.Target, !s.State)
		}
	}
}

func init() {
	engine.RegisterEntityType(&SoundTarget{})
}

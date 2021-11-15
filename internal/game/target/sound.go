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

	"github.com/divVerent/aaaaxy/internal/audiowrap"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/sound"
)

// SoundTarget just changes the music track to the given one.
type SoundTarget struct {
	World  *engine.World
	Entity *engine.Entity

	Sound  *sound.Sound
	Player *audiowrap.Player

	Target      mixins.TargetSelection
	StopWhenOff bool
	State       bool
	Originator  *engine.Entity

	Active bool
	Frames int
	Frame  int
}

func (s *SoundTarget) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	s.World = w
	s.Entity = e
	var err error
	s.Sound, err = sound.Load(sp.Properties["sound"])
	if err != nil {
		return fmt.Errorf("could not load sound: %v", err)
	}
	s.StopWhenOff = sp.Properties["stop_when_off"] == "true" // default false
	s.Target = mixins.ParseTarget(sp.Properties["target"])
	s.State = sp.Properties["state"] != "false"
	s.Frames = -1
	durationString := sp.Properties["duration"]
	var soundTime time.Duration
	if durationString == "" {
		soundTime = s.Sound.Duration()
		if len(s.Target) != 0 {
			return fmt.Errorf("a sound with target must have a duration - please set %v's duration to about %v", sp.Properties["sound"], soundTime)
		}
	} else {
		soundTime, err = time.ParseDuration(durationString)
		if err != nil {
			return fmt.Errorf("could not parse sound duration: %v", durationString)
		}
	}
	s.Frames = int((soundTime*engine.GameTPS + (time.Second / 2)) / time.Second)
	return nil
}

func (s *SoundTarget) Despawn() {
	if s.Player != nil {
		s.Player.Close()
	}
}

func (s *SoundTarget) Update() {
	// Game logic.
	if s.Frame > 0 {
		s.Frame--
		if s.Frame == 0 {
			mixins.SetStateOfTarget(s.World, s.Originator, s.Entity, s.Target, !s.State)
		}
	}
}

func (s *SoundTarget) Touch(other *engine.Entity) {}

func (s *SoundTarget) SetState(originator, predecessor *engine.Entity, state bool) {
	if state {
		// Game logic.
		if s.Active {
			return
		}
		s.Active = true
		s.Frame = s.Frames
		s.Originator = originator
		mixins.SetStateOfTarget(s.World, originator, s.Entity, s.Target, s.State)

		// Sound logic.
		if s.Player != nil {
			s.Player.Close()
		}
		s.Player = s.Sound.Play()
	} else {
		// Game logic.
		if !s.Active {
			return
		}
		s.Active = false
		s.Frame = 0
		mixins.SetStateOfTarget(s.World, originator, s.Entity, s.Target, !s.State)

		// Sound logic.
		if s.Player != nil && s.StopWhenOff {
			s.Player.Close()
			s.Player = nil
		}
	}
}

func init() {
	engine.RegisterEntityType(&SoundTarget{})
}

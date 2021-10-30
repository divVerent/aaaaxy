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

package misc

import (
	"fmt"
	"time"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/level"
	m "github.com/divVerent/aaaaxy/internal/math"
)

// MovingAnimation is a simple entity type that moves in a specified direction.
// Optionally despawns when hitting solid.
type MovingAnimation struct {
	Animation
	mixins.Moving
	mixins.Fadable
	mixins.NonSolidTouchable

	World  *engine.World
	Entity *engine.Entity

	Alpha          float64
	FadeOnTouch    bool
	RespawnOnTouch bool
	StopOnTouch    bool

	FramesToMove int
	FramesToFade int
}

func (s *MovingAnimation) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	s.World = w
	s.Entity = e
	err := s.Animation.Spawn(w, sp, e)
	if err != nil {
		return err
	}
	s.Alpha = e.Alpha
	s.Moving.Init(w, sp, e, level.ObjectSolidContents, s.handleTouch)
	err = s.Fadable.Init(w, sp, e)
	if err != nil {
		return err
	}
	s.NonSolidTouchable.Init(w, e)
	s.FadeOnTouch = sp.Properties["fade_on_touch"] == "true"
	s.RespawnOnTouch = sp.Properties["respawn_on_touch"] == "true"
	s.StopOnTouch = sp.Properties["stop_on_touch"] == "true"
	timeToMoveString := sp.Properties["time_to_move"]
	if timeToMoveString != "" {
		timeToMove, err := time.ParseDuration(timeToMoveString)
		if err != nil {
			return fmt.Errorf("could not parse time to move: %v", timeToMoveString)
		}
		s.FramesToMove = int((timeToMove*engine.GameTPS + (time.Second / 2)) / time.Second)
	}
	timeToFadeString := sp.Properties["time_to_fade"]
	if timeToFadeString != "" {
		timeToFade, err := time.ParseDuration(timeToFadeString)
		if err != nil {
			return fmt.Errorf("could not parse time to fade: %v", timeToFadeString)
		}
		s.FramesToFade = int((timeToFade*engine.GameTPS + (time.Second / 2)) / time.Second)
	}

	return nil
}

func (s *MovingAnimation) Update() {
	if s.FramesToMove > 0 {
		s.FramesToMove--
	} else {
		s.Moving.Update()
	}
	if s.FramesToFade > 0 {
		s.FramesToFade--
		if s.FramesToFade == 0 {
			s.SetState(s.Entity, s.Entity, s.Invert)
		}
	}
	s.Animation.Update()
	s.Fadable.Update()
	s.NonSolidTouchable.Update()
}

func (s *MovingAnimation) Touch(other *engine.Entity) {
	if other != nil && (other.Contents()&level.ObjectSolidContents == 0) {
		// Exclude some "fake hits" by NonSolidTouchable as that one does not care for contents (trace does).
		return
	}
	if other == s.World.Player {
		if s.RespawnOnTouch {
			s.World.RespawnPlayer(s.World.PlayerState.LastCheckpoint(), false)
		}
	} else {
		if s.FadeOnTouch {
			s.SetState(other, s.Entity, s.Invert)
		}
		if s.StopOnTouch {
			s.Velocity = m.Delta{}
		}
	}
}

func (s *MovingAnimation) handleTouch(trace engine.TraceResult) {
	s.World.TouchEvent(s.Entity, trace.HitEntity)
}

func init() {
	engine.RegisterEntityType(&MovingAnimation{})
}

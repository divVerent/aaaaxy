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
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/mixins"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
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

	Alpha float64

	FadeOnTouch    bool
	RespawnOnTouch bool
	StopOnTouch    bool
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
	return nil
}

func (s *MovingAnimation) Update() {
	s.Moving.Update()
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
			s.World.RespawnPlayer(s.World.PlayerState.LastCheckpoint())
		}
	} else {
		if s.FadeOnTouch {
			s.SetState(s.Entity, s.Invert)
		}
	}
	if s.StopOnTouch {
		s.Velocity = m.Delta{}
	}
}

func (s *MovingAnimation) handleTouch(trace engine.TraceResult) {
	s.Touch(trace.HitEntity)
}

func init() {
	engine.RegisterEntityType(&MovingAnimation{})
}

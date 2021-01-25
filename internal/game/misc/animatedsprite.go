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
)

const (
	fadeFrames     = 16
	solidThreshold = 12
)

// AnimatedSprite is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
// Can be toggled from outside.
type AnimatedSprite struct {
	Sprite
	Entity *engine.Entity

	Alpha  float64
	Solid  bool
	Opaque bool

	AnimDir   int
	AnimFrame int
}

func (s *AnimatedSprite) Spawn(w *engine.World, sp *engine.Spawnable, e *engine.Entity) error {
	err := s.Sprite.Spawn(w, sp, e)
	if err != nil {
		return nil
	}

	s.Entity = e

	// Collect the sprite info.
	s.Alpha = s.Entity.Alpha
	s.Solid = s.Entity.Solid
	s.Opaque = s.Entity.Opaque

	// Load the initial state.
	initialState := sp.Properties["initial_state"] != "false" // Defaults to true.
	s.SetState(initialState)

	// Skip the animation on initial load.
	if initialState {
		s.AnimFrame = fadeFrames
	} else {
		s.AnimFrame = 0
	}
	s.Update()

	return nil
}

func (s *AnimatedSprite) SetState(state bool) {
	if state {
		s.AnimDir = 1
	} else {
		s.AnimDir = -1
	}
}

func (s *AnimatedSprite) Update() {
	s.Sprite.Update()

	s.AnimFrame += s.AnimDir

	if s.AnimFrame <= 0 {
		s.Entity.Alpha = 0
		s.AnimFrame = 0
	} else if s.AnimFrame >= fadeFrames {
		s.Entity.Alpha = s.Alpha
		s.AnimFrame = fadeFrames
	} else {
		alpha := float64(s.AnimFrame) / float64(fadeFrames)
		s.Entity.Alpha = alpha * s.Alpha
	}

	if s.AnimFrame >= solidThreshold {
		s.Entity.Solid = s.Solid
		s.Entity.Opaque = s.Opaque
	} else {
		s.Entity.Solid = false
		s.Entity.Opaque = false
	}
}

func init() {
	engine.RegisterEntityType(&AnimatedSprite{})
}

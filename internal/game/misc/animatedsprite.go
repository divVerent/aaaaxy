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

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/mixins"
	"github.com/divVerent/aaaaaa/internal/level"
)

const (
	defaultFadeFrames = 16
	solidThreshold    = 12
)

// AnimatedSprite is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
// Can be toggled from outside.
type AnimatedSprite struct {
	Sprite
	mixins.Settable
	World  *engine.World
	Entity *engine.Entity

	Alpha      float64
	Solid      bool
	Opaque     bool
	FadeFrames int

	AnimDir   int
	AnimFrame int
}

func (s *AnimatedSprite) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	err := s.Sprite.Spawn(w, sp, e)
	if err != nil {
		return nil
	}
	s.Settable.Init(sp)

	s.World = w
	s.Entity = e

	// Collect the sprite info.
	s.Alpha = s.Entity.Alpha
	s.Solid = s.Entity.Solid()
	s.Opaque = s.Entity.Opaque()

	fadeString := sp.Properties["fade_time"]
	if fadeString != "" {
		animTime, err := time.ParseDuration(fadeString)
		if err != nil {
			return fmt.Errorf("could not parse fade time: %v", fadeString)
		}
		s.FadeFrames = int((animTime*engine.GameTPS + (time.Second / 2)) / time.Second)
		if s.FadeFrames < 1 {
			s.FadeFrames = 1
		}
	} else {
		s.FadeFrames = defaultFadeFrames
	}

	// Skip the animation on initial load.
	if s.Settable.State {
		s.AnimFrame = s.FadeFrames
	} else {
		s.AnimFrame = 0
	}
	s.Update()

	return nil
}

func (s *AnimatedSprite) Update() {
	s.Sprite.Update()

	if s.Settable.State {
		s.AnimFrame++
	} else {
		s.AnimFrame--
	}

	if s.AnimFrame <= 0 {
		s.Entity.Alpha = 0
		s.AnimFrame = 0
	} else if s.AnimFrame >= s.FadeFrames {
		s.Entity.Alpha = s.Alpha
		s.AnimFrame = s.FadeFrames
	} else {
		alpha := float64(s.AnimFrame) / float64(s.FadeFrames)
		s.Entity.Alpha = alpha * s.Alpha
	}

	if s.AnimFrame >= solidThreshold {
		s.World.SetSolid(s.Entity, s.Solid)
		s.World.SetOpaque(s.Entity, s.Opaque)
	} else {
		s.World.SetSolid(s.Entity, false)
		s.World.SetOpaque(s.Entity, false)
	}
}

func init() {
	engine.RegisterEntityType(&AnimatedSprite{})
}

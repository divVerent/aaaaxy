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
	solidThreshold    = 8
	opaqueThreshold   = 16
)

// SwitchableSprite is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
// Can be toggled from outside.
type SwitchableSprite struct {
	Sprite
	mixins.Settable
	World  *engine.World
	Entity *engine.Entity

	Alpha      float64
	Contents   level.Contents
	FadeFrames int

	AnimDir   int
	AnimFrame int
}

func (s *SwitchableSprite) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	s.Settable.Init(sp)

	// Default sprites for easy mapping and coherence.
	if sp.Properties["image"] == "" {
		if s.Settable.State {
			// This block shows by default, and thus whenever it is "off".
			sp.Properties["image"] = "switchblock_off.png"
		} else {
			// This block shows only when switched on.
			sp.Properties["image"] = "switchblock_on.png"
		}
		if sp.Properties["solid"] == "" {
			sp.Properties["solid"] = "true"
		}
		if sp.Properties["no_transform"] == "" {
			sp.Properties["no_transform"] = "true"
		}
	}

	err := s.Sprite.Spawn(w, sp, e)
	if err != nil {
		return nil
	}

	s.World = w
	s.Entity = e

	// Collect the sprite info.
	s.Alpha = s.Entity.Alpha
	s.Contents = s.Entity.Contents()

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

func (s *SwitchableSprite) Update() {
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

	s.World.MutateContentsBool(s.Entity, s.Contents&level.SolidContents, s.AnimFrame >= solidThreshold)
	s.World.MutateContentsBool(s.Entity, s.Contents&level.OpaqueContents, s.AnimFrame >= opaqueThreshold)
}

func init() {
	engine.RegisterEntityType(&SwitchableSprite{})
}

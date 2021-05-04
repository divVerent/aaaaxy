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
)

// SwitchableSprite is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
// Can be toggled from outside.
type SwitchableSprite struct {
	Sprite
	mixins.Fadable
}

func (s *SwitchableSprite) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	// Default sprites for easy mapping and coherence.
	if sp.Properties["image"] == "" {
		if sp.Properties["initial_state"] != "false" { // Default true.
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
		return err
	}
	s.Fadable.Init(w, sp, e)
	return nil
}

func (s *SwitchableSprite) Update() {
	s.Sprite.Update()
	s.Fadable.Update()
}

func init() {
	engine.RegisterEntityType(&SwitchableSprite{})
}

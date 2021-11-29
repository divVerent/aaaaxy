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
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/level"
)

// ExitButton shows as Esc key or Start button depending on input device.
type ExitButton struct {
	SwitchableSprite
}

func (s *ExitButton) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	// Just a normal sprite.
	if sp.Properties["fade_time"] == "" {
		sp.Properties["fade_time"] = "10s"
	}
	if sp.Properties["no_flip"] == "" {
		sp.Properties["no_flip"] = "x"
	}
	switch input.ExitButton() {
	default: // case input.Escape:
		sp.Properties["image"] = "esc.png"
	case input.Backspace:
		sp.Properties["image"] = "backspace.png"
	case input.Start:
		sp.Properties["image"] = "start.png"
	}
	s.SwitchableSprite.Spawn(w, sp, e)
	return nil
}

func init() {
	engine.RegisterEntityType(&ExitButton{})
}

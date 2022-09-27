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
	"strings"
	"time"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

// ExitButton shows as Esc key or Start button depending on input device.
type ExitButton struct {
	SwitchableSprite
}

func (s *ExitButton) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	// HACK: adjust some defaults.
	propmap.SetDefault(sp.Properties, "fade_time", 10*time.Second)
	propmap.SetDefault(sp.Properties, "no_flip", "x")
	switch input.ExitButton() {
	default: // case input.Escape:
		propmap.SetDefault(sp.Properties, "image", "esc.png")
	case input.Backspace, input.Back:
		propmap.SetDefault(sp.Properties, "image", "backspace.png")
	case input.Start:
		propmap.SetDefault(sp.Properties, "image", "start.png")
	}
	s.SwitchableSprite.Spawn(w, sp, e)

	// Can turn off based on player abilities :)
	if abilities := propmap.StringOr(sp.Properties, "abilities", ""); abilities != "" {
		haveAll := true
		for _, a := range strings.Split(abilities, " ") {
			if !w.PlayerState.HasAbility(a) {
				haveAll = false
				break
			}
		}
		if haveAll {
			// Hide.
			e.Alpha = 0.0
			w.MutateContentsBool(e, level.AllContents, false)
			s.Fadable.Alpha = 0.0
		}
	}

	return nil
}

func init() {
	engine.RegisterEntityType(&ExitButton{})
}

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

package mixins

import (
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/interfaces"
)

// NonSolidTouchable implements the Touch handler for nonsolid entities.
// Overrides Update(), so if additional handling is desired, calls must be chained.
type NonSolidTouchable struct {
	World           *engine.World
	Entity          *engine.Entity
	NotifyUntouched bool
}

func (t *NonSolidTouchable) Init(w *engine.World, e *engine.Entity) error {
	t.World = w
	t.Entity = e
	return nil
}

func (t *NonSolidTouchable) Update() {
	// NOTE: These Touch events are NOT symmetric like all others! The other entity is NOT notified that we touched it.
	touched := false
	t.World.ForEachEntity(func(e *engine.Entity) {
		if e == t.Entity {
			return
		}
		// It has to be something that can move.
		if _, ok := e.Impl.(interfaces.Velocityer); !ok {
			return
		}
		// Should we filter stronger? Like, only triggers?
		delta := t.Entity.Rect.Delta(e.Rect)
		if delta.IsZero() {
			t.Entity.Impl.Touch(e)
			touched = true
		}
	})
	if !touched && t.NotifyUntouched {
		t.Entity.Impl.Touch(nil)
	}
}

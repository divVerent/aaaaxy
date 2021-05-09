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

package trigger

import (
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/interfaces"
	"github.com/divVerent/aaaaaa/internal/game/mixins"
	"github.com/divVerent/aaaaaa/internal/level"
)

// VVVVVV enables/disables gravity flipping when jumping.
type VVVVVV struct {
	mixins.NonSolidTouchable

	State bool
	Text  string
}

func (v *VVVVVV) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	v.NonSolidTouchable.Init(w, e)
	v.State = s.Properties["state"] != "false" // Default true.
	v.Text = s.Properties["text"]
	return nil
}

func (v *VVVVVV) Despawn() {}

func (v *VVVVVV) Update() {
	v.NonSolidTouchable.Update()
}

func (v *VVVVVV) Touch(other *engine.Entity) {
	if other != v.World.Player {
		return
	}
	v.World.Player.Impl.(interfaces.VVVVVVer).SetVVVVVV(v.State, v.Text)
}

func init() {
	engine.RegisterEntityType(&VVVVVV{})
}

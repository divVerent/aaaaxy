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
	"github.com/divVerent/aaaaaa/internal/game/mixins"
	"github.com/divVerent/aaaaaa/internal/game/target"
	"github.com/divVerent/aaaaaa/internal/level"
)

// SwitchMusic just changes the music track to the given one.
type SwitchMusic struct {
	mixins.NonSolidTouchable
	target.SwitchMusicTarget
}

func (s *SwitchMusic) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	s.NonSolidTouchable.Init(w, e)
	return s.SwitchMusicTarget.Spawn(w, sp, e)
}

func (s *SwitchMusic) Despawn() {}

func (s *SwitchMusic) Update() {
	s.NonSolidTouchable.Update()
	s.SwitchMusicTarget.Update()
}

func (s *SwitchMusic) Touch(other *engine.Entity) {
	if other != s.World.Player {
		return
	}
	s.SetState(true)
}

func init() {
	engine.RegisterEntityType(&SwitchMusic{})
}

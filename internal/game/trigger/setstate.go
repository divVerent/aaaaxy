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

// SetState overrides the boolean state of a warpzone or entity.
type SetState struct {
	mixins.NonSolidTouchable
	target.SetStateTarget

	SendUntouch    bool
	SendEveryFrame bool

	Touching bool
	Touched  bool
}

func (s *SetState) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	s.NonSolidTouchable.Init(w, e)
	err := s.SetStateTarget.Spawn(w, sp, e)
	if err != nil {
		return err
	}
	s.SendUntouch = sp.Properties["send_untouch"] == "true"
	s.SendEveryFrame = sp.Properties["send_every_frame"] == "true"
	return nil
}

func (s *SetState) Despawn() {}

func (s *SetState) Update() {
	s.NonSolidTouchable.Update()
	s.SetStateTarget.Update()
	if s.Touched && !s.Touching && s.SendUntouch {
		s.SetState(false)
	}
	s.Touching, s.Touched = false, s.Touching
}

func (s *SetState) Touch(other *engine.Entity) {
	if s.SendEveryFrame || (!s.Touching && !s.Touched) {
		s.SetState(true)
	}
	s.Touching = true
}

func init() {
	engine.RegisterEntityType(&SetState{})
}

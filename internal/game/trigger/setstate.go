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
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/game/target"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/m"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

// SetState overrides the boolean state of a warpzone or entity.
type SetState struct {
	World  *engine.World
	Entity *engine.Entity
	mixins.NonSolidTouchable
	target.SetStateTarget

	SendUntouch bool
	SendOnce    bool
	PlayerOnly  bool

	Touching   bool
	Touched    bool
	State      bool
	Originator *engine.Entity
}

func (s *SetState) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	s.World = w
	s.Entity = e
	s.NonSolidTouchable.Init(w, e)
	err := s.SetStateTarget.Spawn(w, sp, e)
	if err != nil {
		return err
	}
	var parseErr error
	s.SendUntouch = propmap.ValueOrP(sp.Properties, "send_untouch", false, &parseErr)
	s.SendOnce = propmap.ValueOrP(sp.Properties, "send_once", false, &parseErr)
	s.PlayerOnly = propmap.ValueOrP(sp.Properties, "player_only", false, &parseErr)

	requiredTransforms := propmap.ValueOrP(sp.Properties, "required_orientation", m.Orientations{}, &parseErr)
	if len(requiredTransforms) != 0 {
		show := false
		for _, requiredTransform := range requiredTransforms {
			if e.Transform == requiredTransform {
				show = true
			} else if e.Transform == requiredTransform.Concat(m.FlipX()) {
				show = true
			}
		}
		if !show {
			// Disable.
			s.Target = nil
		}
	}

	return parseErr
}

func (s *SetState) Despawn() {}

func (s *SetState) Update() {
	s.NonSolidTouchable.Update()
	s.SetStateTarget.Update()
	if s.Touched && !s.Touching && s.SendUntouch {
		s.State = false
		s.SetState(s.Originator, s.Entity, false)
	}
	s.Touching, s.Touched = false, s.Touching
}

func (s *SetState) Touch(other *engine.Entity) {
	if s.PlayerOnly && other != s.World.Player {
		return
	}
	if !s.SendOnce || (!s.Touching && !s.Touched) {
		s.State = true
		s.SetState(other, s.Entity, true)
		s.Originator = other
	}
	s.Touching = true
}

func init() {
	engine.RegisterEntityType(&SetState{})
}

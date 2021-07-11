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
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/level"
)

// SwitchableText is a simple entity type that renders text.
// Can be toggled from outside.
type SwitchableText struct {
	Text
	mixins.Fadable
}

func (s *SwitchableText) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	err := s.Text.Spawn(w, sp, e)
	if err != nil {
		return err
	}
	s.Fadable.Init(w, sp, e)
	return nil
}

func (s *SwitchableText) Update() {
	s.Text.Update()
	s.Fadable.Update()
}

func init() {
	engine.RegisterEntityType(&SwitchableText{})
}

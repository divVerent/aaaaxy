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

// MovableSprite is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
// Can be toggled from outside and will move back/forth.
type MovableSprite struct {
	Sprite
	mixins.Movable
}

func (s *MovableSprite) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	err := s.Sprite.Spawn(w, sp, e)
	if err != nil {
		return err
	}
	contents := level.ObjectSolidContents
	if sp.Properties["hit_opaque"] == "true" {
		contents = level.OpaqueContents
	}
	s.Movable.Init(w, sp, e, contents)
	return nil
}

func (s *MovableSprite) Update() {
	s.Sprite.Update()
	s.Movable.Update()
}

func init() {
	engine.RegisterEntityType(&MovableSprite{})
}

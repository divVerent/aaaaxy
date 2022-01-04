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
	"github.com/divVerent/aaaaxy/internal/game/constants"
	"github.com/divVerent/aaaaxy/internal/level"
)

// CoverSprite is a Sprite with a high Z index. Use seldomly!
type CoverSprite struct {
	Sprite
}

func (s *CoverSprite) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	s.Sprite.ZDefault = constants.CoverSpriteZ
	return s.Sprite.Spawn(w, sp, e)
}

func init() {
	engine.RegisterEntityType(&CoverSprite{})
}

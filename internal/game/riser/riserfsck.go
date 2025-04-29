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

package riser

import (
	"fmt"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/picture"
	"github.com/divVerent/aaaaxy/internal/sound"
)

// RiserFsck is an object that kills Risers when touching them.
// Note that if the original spawnpoint of the riser is visible, the riser will then respawn there.
type RiserFsck struct {
	mixins.NonSolidTouchable
	World *engine.World

	Sound *sound.Sound
}

func (r *RiserFsck) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	r.World = w

	err := r.NonSolidTouchable.Init(w, e)
	if err != nil {
		return fmt.Errorf("could not initialize nonsolidtouchbale: %w", err)
	}

	e.Image, err = picture.Load("sprites", "riserfsck.png")
	if err != nil {
		return fmt.Errorf("could not load riserfsck sprite: %w", err)
	}

	r.Sound, err = sound.Load("riserfsck.ogg")
	if err != nil {
		return fmt.Errorf("could not load riserfsck sound: %w", err)
	}

	return nil
}

func (r *RiserFsck) Despawn() {}

func (r *RiserFsck) Touch(other *engine.Entity) {
	if _, ok := other.Impl.(*Riser); !ok {
		return
	}
	if other.Detached() {
		return
	}
	r.World.Detach(other)
	r.Sound.Play()
}

func init() {
	engine.RegisterEntityType(&RiserFsck{})
}

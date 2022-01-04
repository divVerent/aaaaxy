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

package checkpoint

import (
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/level"
)

// Checkpoint remembers that it was hit and allows spawning from there again. Also displays a text.
type Checkpoint struct {
	mixins.NonSolidTouchable
	CheckpointTarget
}

func (c *Checkpoint) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	c.NonSolidTouchable.Init(w, e)
	return c.CheckpointTarget.Spawn(w, sp, e)
}

func (c *Checkpoint) Update() {
	c.NonSolidTouchable.Update()
	c.CheckpointTarget.Update()
}

func (c *Checkpoint) Touch(other *engine.Entity) {
	if other != c.CheckpointTarget.World.Player {
		return
	}
	c.SetState(other, c.CheckpointTarget.Entity, true)
}

func init() {
	engine.RegisterEntityType(&Checkpoint{})
}

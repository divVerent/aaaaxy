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
	"github.com/divVerent/aaaaaa/internal/level"
)

// JumpPad, when hit by the player, sends the player on path to set destination.
// Note that sadly, JumpPads are rarely ever useful in rooms that can be used in multiple orientations.
// May want to introduce required orientation like with checkpoints.
// Or could require player to hit jumppad from above.
type SwitchableJumpPad struct {
	mixins.Settable
	JumpPad
}

func (j *SwitchableJumpPad) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	j.Settable.Init(s)
	return j.JumpPad.Spawn(w, s, e)
}

func (j *SwitchableJumpPad) Touch(other *engine.Entity) {
	if !j.Settable.State {
		return
	}
	j.JumpPad.Touch(other)
}

func init() {
	engine.RegisterEntityType(&SwitchableJumpPad{})
}

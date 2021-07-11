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
	"fmt"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/interfaces"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/image"
	"github.com/divVerent/aaaaxy/internal/level"
	m "github.com/divVerent/aaaaxy/internal/math"
)

// VVVVVV enables/disables gravity flipping when jumping through.
type VVVVVV struct {
	mixins.NonSolidTouchable

	State        bool
	OnGroundVec  m.Delta
	ResetGravity bool
}

func (v *VVVVVV) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	v.NonSolidTouchable.Init(w, e)
	var err error
	e.Image, err = image.Load("sprites", "v.png")
	if err != nil {
		return fmt.Errorf("could not load vvvvvv image: %v", err)
	}
	v.ResetGravity = s.Properties["reset_gravity"] != "false" // Default true.
	if onGroundVecStr := s.Properties["gravity_direction"]; onGroundVecStr != "" {
		_, err := fmt.Sscanf(onGroundVecStr, "%d %d", &v.OnGroundVec.DX, &v.OnGroundVec.DY)
		if err != nil {
			return fmt.Errorf("invalid gravity_direction: %v", err)
		}
	}
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
	side := other.Rect.Center().Delta(v.Entity.Rect.Center()).Dot(v.Entity.Orientation.Right) > 0
	down := v.OnGroundVec
	if !side && v.ResetGravity {
		down = m.Delta{DX: 0, DY: 1}
	}
	v.World.Player.Impl.(interfaces.VVVVVVer).SetVVVVVV(side, down)
}

func init() {
	engine.RegisterEntityType(&VVVVVV{})
}

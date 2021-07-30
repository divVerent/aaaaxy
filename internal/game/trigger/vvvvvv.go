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

	State                bool
	VVVVVVOnGroundVec    m.Delta
	NormalOnGroundVec    m.Delta
	VVVVVVVelocityFactor float64
	NormalVelocityFactor float64
}

func (v *VVVVVV) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	v.NonSolidTouchable.Init(w, e)
	var err error
	e.Image, err = image.Load("sprites", "v.png")
	if err != nil {
		return fmt.Errorf("could not load vvvvvv image: %v", err)
	}
	if onGroundVecStr := s.Properties["gravity_direction"]; onGroundVecStr != "" {
		_, err := fmt.Sscanf(onGroundVecStr, "%d %d", &v.NormalOnGroundVec.DX, &v.NormalOnGroundVec.DY)
		if err != nil {
			return fmt.Errorf("invalid gravity_direction: %v", err)
		}
	}
	if onGroundVecStr := s.Properties["vvvvvv_gravity_direction"]; onGroundVecStr != "" {
		_, err := fmt.Sscanf(onGroundVecStr, "%d %d", &v.VVVVVVOnGroundVec.DX, &v.VVVVVVOnGroundVec.DY)
		if err != nil {
			return fmt.Errorf("invalid vvvvvv_gravity_direction: %v", err)
		}
	}
	v.NormalVelocityFactor = 1.0
	if factorStr := s.Properties["velocity_factor"]; factorStr != "" {
		_, err := fmt.Sscanf(factorStr, "%f", &v.NormalVelocityFactor)
		if err != nil {
			return fmt.Errorf("invalid velocity_factor: %v", err)
		}
	}
	v.VVVVVVVelocityFactor = 1.0
	if factorStr := s.Properties["vvvvvv_velocity_factor"]; factorStr != "" {
		_, err := fmt.Sscanf(factorStr, "%f", &v.VVVVVVVelocityFactor)
		if err != nil {
			return fmt.Errorf("invalid vvvvvv_velocity_factor: %v", err)
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
	vel := other.Impl.(interfaces.Velocityer).ReadVelocity()
	velSide := vel.Dot(v.Entity.Orientation.Right) > 0
	if side != velSide {
		return
	}
	down := v.NormalOnGroundVec
	factor := v.NormalVelocityFactor
	if side {
		down = v.VVVVVVOnGroundVec
		factor = v.VVVVVVVelocityFactor
	}
	other.Impl.(interfaces.VVVVVVer).SetVVVVVV(side, down, factor)
}

func init() {
	engine.RegisterEntityType(&VVVVVV{})
}

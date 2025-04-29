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
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/picture"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

// VVVVVV enables/disables gravity flipping when jumping through.
type VVVVVV struct {
	mixins.NonSolidTouchable

	State                bool
	VVVVVVGravityFlip    bool
	NormalGravityFlip    bool
	VVVVVVVelocityFactor float64
	NormalVelocityFactor float64
}

func (v *VVVVVV) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	v.NonSolidTouchable.Init(w, e)
	var err error
	e.Image, err = picture.Load("sprites", "v.png")
	if err != nil {
		return fmt.Errorf("could not load vvvvvv image: %w", err)
	}
	var parseErr error
	v.NormalGravityFlip = propmap.ValueOrP(sp.Properties, "gravity_flip", false, &parseErr)        // default false
	v.VVVVVVGravityFlip = propmap.ValueOrP(sp.Properties, "vvvvvv_gravity_flip", false, &parseErr) // default false
	v.NormalVelocityFactor = propmap.ValueOrP(sp.Properties, "velocity_factor", 1.0, &parseErr)
	v.VVVVVVVelocityFactor = propmap.ValueOrP(sp.Properties, "vvvvvv_velocity_factor", 1.0, &parseErr)
	return parseErr
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
	flip := v.NormalGravityFlip
	factor := v.NormalVelocityFactor
	if side {
		flip = v.VVVVVVGravityFlip
		factor = v.VVVVVVVelocityFactor
	}
	onGroundVec := other.Impl.(interfaces.Physics).ReadOnGroundVec()
	if flip {
		onGroundVec = onGroundVec.Mul(-1)
	}
	other.Impl.(interfaces.VVVVVVer).SetVVVVVV(side, onGroundVec, factor)
}

func init() {
	engine.RegisterEntityType(&VVVVVV{})
}

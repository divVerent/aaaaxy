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

package mixins

import (
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/constants"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/m"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

// Movable a mixin to make an object move back/forth when toggled.
// Must be initialized _after_ alpha and contents are set by the entity.
// Is shown in the "off" state in the editor.
// It moves using physics using a set acceleration value.
type Movable struct {
	Settable
	Physics
	World  *engine.World
	Entity *engine.Entity

	Acceleration m.Fixed
	From, To     m.Pos

	AnimDir int
}

func (v *Movable) Init(w *engine.World, sp *level.SpawnableProps, e *engine.Entity, contents level.Contents) error {
	v.Settable.Init(sp)

	v.World = w
	v.Entity = e

	var parseErr error
	accel := propmap.ValueOrP(sp.Properties, "acceleration", 0.0, &parseErr)
	if accel != 0 {
		v.Acceleration = m.NewFixedFloat64(accel * constants.SubPixelScale / engine.GameTPS / engine.GameTPS)
	} else {
		v.Acceleration = m.NewFixed(constants.Gravity)
	}

	delta := propmap.ValueP(sp.Properties, "delta", m.Delta{}, &parseErr)
	v.From = e.Rect.Origin
	v.To = e.Rect.Origin.Add(e.Transform.Inverse().Apply(delta))

	// No animation on initial load.
	if v.Settable.State {
		v.Entity.Rect.Origin = v.To
	}

	v.Physics.Init(w, e, contents, func(trace engine.TraceResult) {})

	return parseErr
}

func (v *Movable) Update() {
	// Compute new velocity.
	var target m.Pos
	if v.Settable.State {
		target = v.To
	} else {
		target = v.From
	}
	deltaSub := target.Delta(v.Entity.Rect.Origin).Mul(constants.SubPixelScale)
	deltaSub = deltaSub.Add(m.Delta{DX: constants.SubPixelScale / 2, DY: constants.SubPixelScale / 2}).Sub(v.SubPixel)
	if deltaSub.IsZero() {
		v.Physics.Velocity = m.Delta{}
	} else {
		curSpeed := v.Physics.Velocity.LengthFixed()
		wantSpeed := curSpeed + v.Acceleration
		v.Physics.Velocity = deltaSub.WithMaxLengthFixed(wantSpeed)
	}

	// Move.
	v.Physics.Update()
	// Note: this object does not get pushed by other ground.
	v.Physics.GroundEntity = nil
}

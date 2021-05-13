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
	"fmt"
	"math"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/constants"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
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

	Acceleration float64
	From, To     m.Pos

	AnimDir int
}

func (v *Movable) Init(w *engine.World, s *level.Spawnable, e *engine.Entity, contents level.Contents) error {
	v.Settable.Init(s)

	v.World = w
	v.Entity = e

	accelString := s.Properties["acceleration"]
	if accelString != "" {
		var accel float64
		_, err := fmt.Sscanf(accelString, "%v", &accel)
		if err != nil {
			return fmt.Errorf("failed to parse acceleration %q: %v", accelString, err)
		}
		v.Acceleration = accel * constants.SubPixelScale / engine.GameTPS / engine.GameTPS
	} else {
		v.Acceleration = constants.Gravity
	}

	var delta m.Delta
	_, err := fmt.Sscanf(s.Properties["delta"], "%d %d", &delta.DX, &delta.DY)
	if err != nil {
		return fmt.Errorf("failed to parse delta: %v", err)
	}
	v.From = e.Rect.Origin
	v.To = e.Rect.Origin.Add(e.Transform.Inverse().Apply(delta))

	// No animation on initial load.
	if v.Settable.State {
		v.Entity.Rect.Origin = v.To
	}

	v.Physics.Init(w, e, contents, func(trace engine.TraceResult) {})

	return nil
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
	deltaLen := math.Sqrt(float64(deltaSub.Length2()))
	if deltaLen == 0 {
		v.Physics.Velocity = m.Delta{}
	} else {
		curSpeed := math.Sqrt(float64(v.Physics.Velocity.Length2()))
		wantSpeed := curSpeed + v.Acceleration
		maxSpeed := deltaLen
		if wantSpeed > maxSpeed {
			wantSpeed = maxSpeed
		}
		v.Physics.Velocity = deltaSub.MulFloat(wantSpeed / deltaLen)
	}

	// Move.
	v.Physics.Update()
}

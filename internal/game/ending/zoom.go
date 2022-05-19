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

package ending

import (
	"fmt"
	"time"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/level"
	m "github.com/divVerent/aaaaxy/internal/math"
)

// ZoomTarget zooms the screen out.
type ZoomTarget struct {
	World *engine.World

	State  bool
	Frames int
	Frame  int
}

func (z *ZoomTarget) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	z.World = w

	durationString := sp.Properties["duration"]
	durationTime, err := time.ParseDuration(durationString)
	if err != nil {
		return fmt.Errorf("could not parse duration time: %s", durationString)
	}
	z.Frames = int((durationTime*engine.GameTPS + (time.Second / 2)) / time.Second)
	if z.Frames < 1 {
		z.Frames = 1
	}

	return nil
}

func (z *ZoomTarget) Despawn() {}

func (z *ZoomTarget) Update() {
	if z.Frame > 0 {
		z.Frame--
		z.World.MaxVisiblePixels = int(m.Delta{DX: engine.GameWidth, DY: engine.GameHeight}.Length() * 0.5 * float64(z.Frame) / float64(z.Frames))
	}
}

func (z *ZoomTarget) SetState(originator, predecessor *engine.Entity, state bool) {
	if state == z.State {
		return
	}
	z.State = state
	if state {
		z.Frame = z.Frames
	} else {
		z.Frame = 0
	}
}

func (z *ZoomTarget) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&ZoomTarget{})
}

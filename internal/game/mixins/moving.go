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

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/constants"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
)

// Moving is a physics object with initial velocity. That's all.
type Moving struct {
	Physics
	TouchedSomething bool
}

func (v *Moving) Init(w *engine.World, s *level.Spawnable, e *engine.Entity, contents level.Contents, handleTouch func(engine.TraceResult)) error {
	v.Physics.Init(w, e, contents, handleTouch)
	if str := s.Properties["velocity"]; str != "" {
		var dx, dy float64
		if _, err := fmt.Sscanf(str, "%f %f", &dx, &dy); err != nil {
			return fmt.Errorf("Failed to parse velocity %q: %v", str, err)
		}
		v.Physics.Velocity = e.Transform.Inverse().Apply(m.Delta{
			DX: m.Rint(dx * constants.SubPixelScale / engine.GameTPS),
			DY: m.Rint(dy * constants.SubPixelScale / engine.GameTPS),
		})
	}
	return nil
}

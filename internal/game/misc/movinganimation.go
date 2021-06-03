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

package misc

import (
	"fmt"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/constants"
	"github.com/divVerent/aaaaaa/internal/game/mixins"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
)

const (
	FadeFrames = 16
)

// MovingAnimation is a simple entity type that moves in a specified direction.
// Optionally despawns when hitting solid.
type MovingAnimation struct {
	Animation
	mixins.Physics

	Alpha float64

	DespawnOnTouch bool // TODO implement.
	Despawning     bool
	FadeFrame      int
}

func (s *MovingAnimation) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	err := s.Animation.Spawn(w, sp, e)
	if err != nil {
		return err
	}
	s.Alpha = e.Alpha
	s.Physics.Init(w, e, level.ObjectSolidContents, s.handleTouch)
	if str := sp.Properties["velocity"]; str != "" {
		var dx, dy float64
		if _, err := fmt.Sscanf(str, "%f %f", &dx, &dy); err != nil {
			return fmt.Errorf("Failed to parse velocity %q: %v", str, err)
		}
		s.Physics.Velocity = m.Delta{
			DX: m.Rint(dx * constants.SubPixelScale / engine.GameTPS),
			DY: m.Rint(dy * constants.SubPixelScale / engine.GameTPS),
		}
	}
	return nil
}

func (s *MovingAnimation) Update() {
	s.Physics.Update()
	s.Animation.Update()

	if s.Despawning {
		if s.FadeFrame > 0 {
			s.FadeFrame--
		}
		if s.FadeFrame == 0 {
			s.World.Despawn(s.Physics.Entity)
		}
	} else {
		if s.FadeFrame < FadeFrames {
			s.FadeFrame++
		}
	}
	s.Physics.Entity.Alpha = s.Alpha * float64(s.FadeFrame) / float64(FadeFrames)
}

func (s *MovingAnimation) handleTouch(trace engine.TraceResult) {
	if s.DespawnOnTouch && trace.HitEntity != s.World.Player {
		s.Despawning = true
	}
}

func init() {
	engine.RegisterEntityType(&MovingAnimation{})
}

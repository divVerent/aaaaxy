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

package target

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/mixins"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
)

// SequenceTarget sends a given string to a SequenceCollector when triggered.
type SequenceTarget struct {
	World  *engine.World
	Entity *engine.Entity

	Target   string
	Sequence string
}

func (s *SequenceTarget) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	s.World = w
	s.Entity = e
	s.Target = sp.Properties["target"]
	s.Sequence = sp.Properties["sequence"]
	return nil
}

func (s *SequenceTarget) Despawn() {}

func (s *SequenceTarget) Update() {}

func (s *SequenceTarget) SetState(state bool) {
	if !state {
		return
	}
	for _, ent := range s.World.FindName(s.Target) {
		collector, ok := ent.Impl.(*SequenceCollector)
		if !ok {
			log.Printf("Target of SequenceTarget is not a SequenceCollector: %T, name: %v", ent, s.Target)
		}
		collector.Append(s.Sequence)
	}
}

func (s *SequenceTarget) Touch(other *engine.Entity) {}

func (s *SequenceTarget) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

// SequenceCollector waits for receiving strings, and sends a trigger event when the correct string was received.
type SequenceCollector struct {
	World  *engine.World
	Entity *engine.Entity

	Sequence string
	Target   mixins.TargetSelection
	State    bool

	Current string
}

func (s *SequenceCollector) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	s.World = w
	s.Entity = e
	s.Sequence = sp.Properties["sequence"]
	s.Target = mixins.ParseTarget(sp.Properties["target"])
	s.State = sp.Properties["state"] != "false"
	return nil
}

func (s *SequenceCollector) Despawn() {}

func (s *SequenceCollector) Update() {}

func (s *SequenceCollector) Append(str string) {
	matched := s.Current == s.Sequence
	s.Current += str
	if len(s.Current) > len(s.Sequence) {
		s.Current = s.Current[len(s.Current)-len(s.Sequence):]
	}
	matches := s.Current == s.Sequence
	if matches != matched {
		mixins.SetStateOfTarget(s.World, s.Entity, s.Target, s.State == matches)
	}
}

func (s *SequenceCollector) Touch(other *engine.Entity) {}

func (s *SequenceCollector) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&SequenceTarget{})
	engine.RegisterEntityType(&SequenceCollector{})
}

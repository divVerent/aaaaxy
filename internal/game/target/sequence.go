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
	"github.com/divVerent/aaaaxy/internal/log"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/level"
)

// SequenceTarget sends a given string to a SequenceCollector when triggered.
type SequenceTarget struct {
	World  *engine.World
	Entity *engine.Entity

	Target   string
	Sequence string

	State bool
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

func (s *SequenceTarget) SetState(originator, predecessor *engine.Entity, state bool) {
	// Only respond to state transitions.
	if state == s.State {
		return
	}
	s.State = state
	// Only respond to switching on.
	if !state {
		return
	}
	for _, ent := range s.World.FindName(s.Target) {
		collector, ok := ent.Impl.(*SequenceCollector)
		if !ok {
			log.Errorf("Target of SequenceTarget is not a SequenceCollector: %T, name: %v", ent, s.Target)
		}
		collector.Append(originator, s.Sequence)
	}
}

func (s *SequenceTarget) Touch(other *engine.Entity) {}

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

func (s *SequenceCollector) Append(originator *engine.Entity, str string) {
	matched := s.Current == s.Sequence
	s.Current += str
	if len(s.Current) > len(s.Sequence) {
		s.Current = s.Current[len(s.Current)-len(s.Sequence):]
	}
	matches := s.Current == s.Sequence
	if matches && !matched {
		mixins.SetStateOfTarget(s.World, originator, s.Entity, s.Target, true)
		// TODO(divVerent): Maybe also add a send_untouch-like feature to send "off" events too?
	}
}

func (s *SequenceCollector) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&SequenceTarget{})
	engine.RegisterEntityType(&SequenceCollector{})
}

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
	"fmt"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/level"
)

// LogicalGate sends a signal along if ANY incoming target triggers.
type LogicalGate struct {
	World  *engine.World
	Entity *engine.Entity

	Target        mixins.TargetSelection
	Invert        bool
	CountRequired int
	IgnoreOff     bool

	IncomingState map[engine.EntityIncarnation]struct{}
	State         bool
	Originator    *engine.Entity
}

func (g *LogicalGate) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	g.World = w
	g.Entity = e
	g.Target = mixins.ParseTarget(sp.Properties["target"])
	g.CountRequired = 1                                 // An "or" gate by default.
	g.Invert = sp.Properties["invert"] == "true"        // false by default.
	g.IgnoreOff = sp.Properties["ignore_off"] == "true" // false by default.
	if countStr := sp.Properties["count_required"]; countStr != "" {
		_, err := fmt.Sscanf(countStr, "%d", &g.CountRequired)
		if err != nil {
			return fmt.Errorf("failed to parse count_required: got %q, want integer: %v", countStr, err)
		}
	}
	g.IncomingState = map[engine.EntityIncarnation]struct{}{}
	return nil
}

func (g *LogicalGate) Despawn() {}

func (g *LogicalGate) Update() {
	for ent := range g.IncomingState {
		if !g.World.EntityIsAlive(ent) {
			delete(g.IncomingState, ent)
		}
	}
	g.MaybeSendEvent(true)
}

func (g *LogicalGate) Touch(other *engine.Entity) {}

func (g *LogicalGate) SetState(originator, predecessor *engine.Entity, state bool) {
	if state {
		g.IncomingState[predecessor.Incarnation] = struct{}{}
	} else if !g.IgnoreOff {
		delete(g.IncomingState, predecessor.Incarnation)
	}
	g.Originator = originator
}

func (g *LogicalGate) MaybeSendEvent(sendEveryFrame bool) {
	newState := len(g.IncomingState) >= g.CountRequired
	if newState == g.State && !(sendEveryFrame && newState) {
		return
	}
	g.State = newState
	mixins.SetStateOfTarget(g.World, g.Originator, g.Entity, g.Target, newState != g.Invert)
}

func init() {
	engine.RegisterEntityType(&LogicalGate{})
}

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
	"time"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/level"
)

type delayEvent struct {
	FramesLeft int
	Originator *engine.Entity
	State      bool
}

// DelayTarget delays all state changes by the given amount.
type DelayTarget struct {
	World  *engine.World
	Entity *engine.Entity

	DelayFrames int

	Events []delayEvent

	Target     mixins.TargetSelection
	Originator *engine.Entity
}

func (d *DelayTarget) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	d.World = w
	d.Entity = e

	delayString := sp.Properties["delay"]
	delayTime, err := time.ParseDuration(delayString)
	if err != nil {
		return fmt.Errorf("could not parse delay time: %v", delayString)
	}
	d.DelayFrames = int((delayTime*engine.GameTPS + (time.Second / 2)) / time.Second)
	if d.DelayFrames < 1 {
		d.DelayFrames = 1
	}

	d.Target = mixins.ParseTarget(sp.Properties["target"])
	return nil
}

func (d *DelayTarget) Despawn() {}

func (d *DelayTarget) Update() {
	for i := range d.Events {
		ev := &d.Events[i]
		if ev.FramesLeft <= 0 {
			continue
		}
		ev.FramesLeft--
		if ev.FramesLeft == 0 {
			mixins.SetStateOfTarget(d.World, ev.Originator, d.Entity, d.Target, ev.State)
		}
	}
}

func (d *DelayTarget) Touch(other *engine.Entity) {}

func (d *DelayTarget) SetState(originator, predecessor *engine.Entity, state bool) {
	newEvent := delayEvent{
		FramesLeft: d.DelayFrames,
		Originator: originator,
		State:      state,
	}
	for i := range d.Events {
		ev := &d.Events[i]
		if ev.FramesLeft <= 0 {
			*ev = newEvent
			return
		}
	}
	d.Events = append(d.Events, newEvent)
}

func init() {
	engine.RegisterEntityType(&DelayTarget{})
}

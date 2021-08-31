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
	"strings"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/log"
)

// Settable implements the SetState handler for settable entities.
type Settable struct {
	State  bool
	Invert bool
}

// SetState changes the state of the entity.
func (s *Settable) SetState(originator, predecessor *engine.Entity, state bool) {
	s.State = state != s.Invert
}

// Init initializes the initial state of the entity.
func (s *Settable) Init(sp *level.Spawnable) error {
	s.Invert = sp.Properties["invert"] == "true" // Default false.
	s.State = s.Invert
	return nil
}

// stateSetter is an entity that contains this mixin.
type stateSetter interface {
	SetState(originator, predecessor *engine.Entity, state bool)
}

// SetStateOfEntity sets the state of an entity, if available.
// Returns whether the setting was successful.
func SetStateOfEntity(originator, predecessor *engine.Entity, of *engine.Entity, state bool) bool {
	setter, ok := of.Impl.(stateSetter)
	if !ok {
		return false
	}
	setter.SetState(originator, predecessor, state)
	return true
}

type TargetSelection []string

func ParseTarget(target string) TargetSelection {
	return TargetSelection(strings.Split(target, " "))
}

// SetStateOfTarget toggles the state of all entities of the given target name to the given state.
// Includes WarpZones too.
// Excludes the given entity (should be the caller).
func SetStateOfTarget(w *engine.World, originator, predecessor *engine.Entity, targets TargetSelection, state bool) {
	for _, target := range targets {
		if target == "" {
			continue
		}
		thisState := state
		if target[0] == '!' {
			thisState = !state
			target = target[1:]
		}
		if target == "" {
			continue
		}
		if target[0] == '=' {
			target = target[1:]
			var closest *engine.Entity
			for _, ent := range w.FindName(target) {
				if ent == predecessor {
					log.Fatalf("Should this even happen? An entity targeting itself? Name is %v.", target)
					continue
				}
				if closest == nil || closest.Rect.Delta(w.Player.Rect).Norm1() > ent.Rect.Delta(w.Player.Rect).Norm1() {
					closest = ent
				}
			}
			if closest != nil {
				if !SetStateOfEntity(originator, predecessor, closest, state) {
					log.Errorf("Tried to set state of a non-supporting entity: %T, name: %v", closest, target)
				}
			}
		} else {
			w.SetWarpZoneState(target, thisState)
			for _, ent := range w.FindName(target) {
				if ent == predecessor {
					log.Fatalf("Should this even happen? An entity targeting itself? Name is %v.", target)
					continue
				}
				if !SetStateOfEntity(originator, predecessor, ent, thisState) {
					log.Errorf("Tried to set state of a non-supporting entity: %T, name: %v", ent, target)
				}
			}
		}
	}
}

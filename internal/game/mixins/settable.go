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
	"log"
	"strings"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/level"
)

// Settable implements the SetState handler for settable entities.
type Settable struct {
	State bool
}

// SetState changes the state of the entity.
func (s *Settable) SetState(by *engine.Entity, state bool) {
	s.State = state
}

// Init initializes the initial state of the entity.
func (s *Settable) Init(sp *level.Spawnable) error {
	s.State = sp.Properties["initial_state"] != "false" // Default true.
	return nil
}

// stateSetter is an entity that contains this mixin.
type stateSetter interface {
	SetState(by *engine.Entity, state bool)
}

// SetStateOfEntity sets the state of an entity, if available.
// Returns whether the setting was successful.
func SetStateOfEntity(by *engine.Entity, of *engine.Entity, state bool) bool {
	setter, ok := of.Impl.(stateSetter)
	if !ok {
		return false
	}
	setter.SetState(by, state)
	return true
}

type TargetSelection []string

func ParseTarget(target string) TargetSelection {
	return TargetSelection(strings.Split(target, " "))
}

// SetStateOfTarget toggles the state of all entities of the given target name to the given state.
// Includes WarpZones too.
// Excludes the given entity (should be the caller).
func SetStateOfTarget(w *engine.World, by *engine.Entity, targets TargetSelection, state bool) {
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
				if ent == by {
					continue
				}
				if closest == nil || closest.Rect.Delta(w.Player.Rect).Norm1() > ent.Rect.Delta(w.Player.Rect).Norm1() {
					closest = ent
				}
			}
			if closest != nil {
				if !SetStateOfEntity(by, closest, state) {
					log.Printf("Tried to set state of a non-supporting entity: %T, name: %v", closest, target)
				}
			}
		} else {
			w.SetWarpZoneState(target, thisState)
			for _, ent := range w.FindName(target) {
				if ent == by {
					continue
				}
				if !SetStateOfEntity(by, ent, thisState) {
					log.Printf("Tried to set state of a non-supporting entity: %T, name: %v", ent, target)
				}
			}
		}
	}
}

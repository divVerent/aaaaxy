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

package level

import (
	m "github.com/divVerent/aaaaxy/internal/math"
)

// PersistentState is how entities retain values across loading/unloading and in
// save games.
type PersistentState map[string]string

// A Spawnable is a blueprint to create an Entity.
type Spawnable struct {
	ID EntityID

	// Type.
	EntityType string

	// Location.
	LevelPos    m.Pos
	RectInTile  m.Rect
	Orientation m.Orientation

	// Other properties.
	Properties map[string]string

	// Persistent entity state, if any, shall be kept in this map.
	PersistentState map[string]string `hash:"-"`
}

func (sp *Spawnable) Clone() *Spawnable {
	// First make a shallow copy.
	outSp := new(Spawnable)
	*outSp = *sp
	// Now "deepend" all that's needed.
	outSp.PersistentState = make(map[string]string, len(sp.PersistentState))
	for k, v := range sp.PersistentState {
		outSp.PersistentState[k] = v
	}
	return outSp
}

// EntityID represents an unique ID of an entity.
type EntityID int

// InvalidEntityID is an ID that cannot be used in Tiled.
// Tiled's first entity has ID 1.
const InvalidEntityID EntityID = 0

// IsValid returns whether an EntityID is valid for an actual entity.
func (e EntityID) IsValid() bool {
	return e != InvalidEntityID
}

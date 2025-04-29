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
	"github.com/divVerent/aaaaxy/internal/m"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

// PersistentState is how entities retain values across loading/unloading and in
// save games.
type PersistentState = propmap.Map

// A SpawnableProps is a blueprint to create an Entity but without its position.
type SpawnableProps struct {
	// Type.
	EntityType string

	// Orientation.
	Orientation m.Orientation

	// Other properties.
	Properties propmap.Map

	// Persistent entity state, if any, shall be kept in this map.
	PersistentState PersistentState `hash:"-"`

	// SpawnTilesGrowth is how much extra pixels around the entity to consider
	// for spawning.
	SpawnTilesGrowth m.Delta
}

// A Spawnable is a blueprint to create an Entity in a level.
type Spawnable struct {
	SpawnableProps

	// The ID of the entity in the map.
	ID EntityID

	// Location.
	LevelPos   m.Pos
	RectInTile m.Rect
}

func (sp *Spawnable) Clone() *Spawnable {
	// First make a shallow copy.
	outSp := new(Spawnable)
	*outSp = *sp
	// Now "deepend" all that's needed.
	outSp.PersistentState = propmap.New()
	propmap.ForEach(sp.PersistentState, func(k, v string) error {
		propmap.Set(outSp.PersistentState, k, v)
		return nil
	})
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

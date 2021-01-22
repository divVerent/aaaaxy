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

package engine

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"

	m "github.com/divVerent/aaaaaa/internal/math"
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

// An Entity is an object that exists in the game.
type Entity struct {
	Incarnation    EntityIncarnation
	visibilityMark uint

	// Info needed for gameplay.
	Solid     bool
	Opaque    bool
	Rect      m.Rect
	Transform m.Orientation // Possibly needed for initialization.

	// Info needed for rendering.
	Orientation  m.Orientation
	Image        *ebiten.Image
	RenderOffset m.Delta
	ResizeImage  bool // Conceptually incompatible with RenderOffset.
	Alpha        float64
	ZIndex       int

	// Entity's own state.
	Impl EntityImpl
}

type (
	// EntityID represents an unique ID of an entity.
	EntityID int
	// EntityIncarnation represents a specific incarnation of an entity. Entities spawn more than once if their tile is seen more than once.
	EntityIncarnation struct {
		ID      EntityID
		TilePos m.Pos
	}
)

type EntityImpl interface {
	// Spawn initializes the entity based on a Spawnable.
	// Receiver will be a zero struct of the entity type.
	// Will usually remember a reference to the World and Entity.
	// ID, Pos, Size and Orientation of the entity will be preset but may be changed.
	Spawn(w *World, s *Spawnable, e *Entity) error

	// Despawn notifies the entity that it will be deleted.
	Despawn()

	// Update asks the entity to update its state.
	Update()

	// Touch notifies the entity that it was hit by another entity moving.
	Touch(other *Entity)

	// Draw the entity's overlay. Useful for entities that are more than just a sprite.
	DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta)
}

// entityTypes is a helper map to know how to spawn an entity.
var entityTypes = map[string]EntityImpl{}

// RegisterEntityType adds an entity type to the spawn system.
// To be called from init() functions of entity implementations.
func RegisterEntityType(t EntityImpl) {
	typeName := reflect.TypeOf(t).Elem().Name()
	if entityTypes[typeName] != nil {
		log.Panicf("Duplicate entity type %q", typeName)
	}
	entityTypes[typeName] = t
	log.Printf("Registered entity type %q", typeName)
}

// Spawn turns a Spawnable into an Entity.
func (s *Spawnable) Spawn(w *World, tilePos m.Pos, t *Tile) (*Entity, error) {
	tInv := t.Transform.Inverse()
	originTilePos := tilePos.Add(tInv.Apply(s.LevelPos.Delta(t.LevelPos)))
	incarnation := EntityIncarnation{
		ID:      s.ID,
		TilePos: originTilePos,
	}
	if e := w.Entities[incarnation]; e != nil {
		return e, nil
	}
	eTmpl := entityTypes[s.EntityType]
	if eTmpl == nil {
		return nil, fmt.Errorf("unknown entity type %q", s.EntityType)
	}
	eImplVal := reflect.New(reflect.TypeOf(eTmpl).Elem())
	eImplVal.Elem().Set(reflect.ValueOf(eTmpl).Elem())
	eImpl := eImplVal.Interface().(EntityImpl)
	e := &Entity{
		Incarnation: incarnation,
		Impl:        eImpl,
		Transform:   t.Transform,
	}
	pivot2InTile := m.Pos{X: TileSize, Y: TileSize}
	e.Rect = tInv.ApplyToRect2(pivot2InTile, s.RectInTile)
	e.Rect.Origin = originTilePos.Mul(TileSize).Add(e.Rect.Origin.Delta(m.Pos{}))
	e.Orientation = tInv.Concat(s.Orientation)
	e.Alpha = 1.0
	err := eImpl.Spawn(w, s, e)
	if err != nil {
		return nil, err
	}
	w.Entities[incarnation] = e
	return e, nil
}

// PlayerEntityImpl defines some additional methods player entities must have.
type PlayerEntityImpl interface {
	EntityImpl

	// EyePos is the location of the eye in world coordinates.
	EyePos() m.Pos

	// LookPos is the desired screen center position.
	LookPos() m.Pos

	// Respawned() notifies the entity that the world respawned it.
	Respawned()
}

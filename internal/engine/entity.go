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

	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
)

// An Entity is an object that exists in the game.
type Entity struct {
	Incarnation    EntityIncarnation
	visibilityMark uint

	// Info needed for gameplay.
	solid     bool
	opaque    bool
	Rect      m.Rect
	Transform m.Orientation // Possibly needed for initialization.
	name      string        // Possibly searched for.

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

// EntityIncarnation represents a specific incarnation of an entity. Entities spawn more than once if their tile is seen more than once.
type EntityIncarnation struct {
	ID      level.EntityID
	TilePos m.Pos
}

type EntityImpl interface {
	// Spawn initializes the entity based on a Spawnable.
	// Receiver will be a zero struct of the entity type.
	// Will usually remember a reference to the World and Entity.
	// ID, Pos, Size and Orientation of the entity will be preset but may be changed.
	Spawn(w *World, s *level.Spawnable, e *Entity) error

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
func (w *World) Spawn(s *level.Spawnable, tilePos m.Pos, t *level.Tile) (*Entity, error) {
	tInv := t.Transform.Inverse()
	originTilePos := tilePos.Add(tInv.Apply(s.LevelPos.Delta(t.LevelPos)))
	incarnation := EntityIncarnation{
		ID:      s.ID,
		TilePos: originTilePos,
	}
	if e := w.entities[incarnation]; e != nil {
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
		Transform:   t.Transform,
		name:        s.Properties["name"],
		Impl:        eImpl,
	}
	pivot2InTile := m.Pos{X: level.TileSize, Y: level.TileSize}
	e.Rect = tInv.ApplyToRect2(pivot2InTile, s.RectInTile)
	e.Rect.Origin = originTilePos.Mul(level.TileSize).Add(e.Rect.Origin.Delta(m.Pos{}))
	e.Orientation = tInv.Concat(s.Orientation)
	e.Alpha = 1.0
	w.link(e)
	err := eImpl.Spawn(w, s, e)
	if err != nil {
		w.unlink(e)
		return nil, err
	}
	return e, nil
}

// SetSolid makes an entity solid (or not).
func (w *World) SetSolid(e *Entity, solid bool) {
	w.unlink(e)
	e.solid = solid
	w.link(e)
}

// SetOpaque makes an entity opaque (or not).
func (w *World) SetOpaque(e *Entity, opaque bool) {
	w.unlink(e)
	e.opaque = opaque
	w.link(e)
}

func (e *Entity) Solid() bool {
	return e.solid
}

func (e *Entity) Opaque() bool {
	return e.opaque
}

func (e *Entity) Name() string {
	return e.name
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

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
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/m"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

// An Entity is an object that exists in the game.
type Entity struct {
	// Info needed for management.
	Incarnation      EntityIncarnation
	SpawnTilesGrowth m.Delta

	// Info needed for gameplay.
	contents     level.Contents
	Rect         m.Rect
	BorderPixels int           // Border applied to ALL sides. Used for entity tracing only.
	Transform    m.Orientation // Possibly needed for initialization.
	name         string        // Possibly searched for.
	RequireTiles bool          // Entity requires tiles to be loaded.

	// Info needed for rendering.
	Orientation  m.Orientation
	Image        *ebiten.Image
	RenderOffset m.Delta
	ResizeImage  bool // Conceptually incompatible with RenderOffset.
	Alpha        float64
	ColorAdd     [4]float64
	ColorMod     [4]float64
	zIndex       int

	// Intrusive list state.
	indexInListPlusOne [numLists]int

	// Entity's own state.
	Impl EntityImpl
}

// EntityIncarnation represents a specific incarnation of an entity. Entities spawn more than once if their tile is seen more than once.
type EntityIncarnation struct {
	ID      level.EntityID
	TilePos m.Pos
}

func (e EntityIncarnation) IsValid() bool {
	return e.ID.IsValid()
}

type EntityImpl interface {
	// Spawn initializes the entity based on a Spawnable.
	// Receiver will be a zero struct of the entity type.
	// Will usually remember a reference to the World and Entity.
	// ID, Pos, Size and Orientation of the entity will be preset but may be changed.
	Spawn(w *World, sp *level.SpawnableProps, e *Entity) error

	// Despawn notifies the entity that it will be deleted.
	Despawn()

	// Update asks the entity to update its state.
	Update()

	// Touch notifies the entity that it was hit by another entity moving.
	Touch(other *Entity)
}

// Some entities fulfill PrecacheImpl. These will get precached.
type Precacher interface {
	// Precache gets called during level loading to preload anything the entity may need.
	Precache(sp *level.Spawnable) error
}

// Some entities get a pre-despawn notification.
type PreDespawner interface {
	// PreDespawn gets called when the game enters the menu. After this, either Despawn or Update will happen eventually.
	PreDespawn()
}

// entityTypes is a helper map to know how to spawn an entity.
var entityTypes = map[string]EntityImpl{}

// RegisterEntityType adds an entity type to the spawn system.
// To be called from init() functions of entity implementations.
func RegisterEntityType(t EntityImpl) {
	typeName := reflect.TypeOf(t).Elem().Name()
	if entityTypes[typeName] != nil {
		log.Fatalf("duplicate entity type: %v", typeName)
	}
	entityTypes[typeName] = t
	log.Debugf("registered entity type %q", typeName)
}

// Precache all entities.
func precacheEntities(lvl *level.Level) error {
	var err error
	precached := map[level.EntityID]struct{}{}
	lvl.ForEachTile(func(pos m.Pos, t *level.LevelTile) {
		for _, sp := range t.Tile.Spawnables {
			if err != nil {
				break
			}
			if _, found := precached[sp.ID]; found {
				continue
			}
			precached[sp.ID] = struct{}{}
			eTmpl := entityTypes[sp.EntityType]
			if eTmpl == nil {
				err = fmt.Errorf("unknown entity type %q", sp.EntityType)
				break
			}
			if precacher, ok := eTmpl.(Precacher); ok {
				err = precacher.Precache(sp)
				if err != nil {
					err = fmt.Errorf("failed to precache entity %v: %w", sp, err)
				}
			}
		}
	})
	return err
}

// spawnAt spawns a given entity at a given location.
// Still need to provide a Spawnable. Transform should be the transform of the origin entity or tile.
func (w *World) spawnAt(sp *level.SpawnableProps, rect m.Rect, transform, tInv m.Orientation, incarnation EntityIncarnation) (*Entity, error) {
	eTmpl := entityTypes[sp.EntityType]
	if eTmpl == nil {
		return nil, fmt.Errorf("unknown entity type %q", sp.EntityType)
	}
	eImplVal := reflect.New(reflect.TypeOf(eTmpl).Elem())
	eImplVal.Elem().Set(reflect.ValueOf(eTmpl).Elem())
	eImpl := eImplVal.Interface().(EntityImpl)
	e := &Entity{
		Incarnation:      incarnation,
		Transform:        transform,
		name:             propmap.StringOr(sp.Properties, "name", ""),
		Impl:             eImpl,
		Rect:             rect,
		Orientation:      tInv.Concat(sp.Orientation),
		SpawnTilesGrowth: sp.SpawnTilesGrowth,
	}
	e.Alpha = 1.0
	e.ColorMod[0] = 1.0
	e.ColorMod[1] = 1.0
	e.ColorMod[2] = 1.0
	e.ColorMod[3] = 1.0
	w.link(e)
	err := eImpl.Spawn(w, sp, e)
	if err != nil {
		w.unlink(e)
		return nil, err
	}
	return e, nil
}

// Spawn turns a Spawnable into an Entity.
func (w *World) Spawn(sp *level.Spawnable, tilePos m.Pos, t *level.Tile) (*Entity, error) {
	tInv := t.Transform.Inverse()
	originTilePos := tilePos.Add(tInv.Apply(sp.LevelPos.Delta(t.LevelPos)))
	incarnation := EntityIncarnation{
		ID:      sp.ID,
		TilePos: originTilePos,
	}
	if _, found := w.incarnations[incarnation]; found {
		return nil, nil
	}
	if propmap.Has(sp.Properties, "only_if_can_switch_level") {
		restrict, err := propmap.ValueOr(sp.Properties, "only_if_can_switch_level", false)
		if err != nil {
			return nil, err
		}
		if restrict {
			if !CanSwitchLevel() {
				return nil, nil
			}
		}
	}
	pivot2InTile := m.Pos{X: level.TileSize, Y: level.TileSize}
	rect := tInv.ApplyToRect2(pivot2InTile, sp.RectInTile)
	rect.Origin = originTilePos.Mul(level.TileSize).Add(rect.Origin.Delta(m.Pos{}))
	return w.spawnAt(&sp.SpawnableProps, rect, t.Transform, tInv, incarnation)
}

// SpawnDetached spawns a detached new entity.
// Note that it will inherit the spawning entity's transform; sp.Orientation should be set to the same as the transform if one wants to undo that.
func (w *World) SpawnDetached(sp *level.SpawnableProps, rect m.Rect, orientation m.Orientation, by *Entity) (*Entity, error) {
	return w.spawnAt(sp, rect, by.Transform, by.Transform.Inverse(), EntityIncarnation{
		ID:      level.InvalidEntityID,
		TilePos: by.Incarnation.TilePos,
	})
}

func (w *World) Despawn(e *Entity) {
	e.Impl.Despawn()
	w.unlink(e)
}

// MutateContents mutates an entity's contents.
func (w *World) MutateContents(e *Entity, mask, set level.Contents) {
	if e.contents&mask == set {
		return
	}
	w.unlink(e)
	e.contents &= ^mask
	e.contents |= set
	w.link(e)
}

// MutateContentsBool mutates an entity's contents.
func (w *World) MutateContentsBool(e *Entity, mask level.Contents, set bool) {
	if set {
		w.MutateContents(e, mask, mask)
	} else {
		w.MutateContents(e, mask, 0)
	}
}

// SetSolid makes an entity solid (or not).
func (w *World) SetSolid(e *Entity, solid bool) {
	w.MutateContentsBool(e, level.SolidContents, solid)
}

// SetOpaque makes an entity opaque (or not).
func (w *World) SetOpaque(e *Entity, opaque bool) {
	w.MutateContentsBool(e, level.OpaqueContents, opaque)
}

// SetZIndex sets an entity's Z index.
func (w *World) SetZIndex(e *Entity, index int) {
	if e.zIndex == index {
		return
	}
	w.unlink(e)
	e.zIndex = index
	w.link(e)
}

// Detach detaches an entity from its spawn origin.
// If the spawn origin is onscreen, this will respawn the entity from there this frame.
func (w *World) Detach(e *Entity) {
	if !e.Incarnation.IsValid() {
		return
	}
	w.unlink(e)
	e.Incarnation.ID = level.InvalidEntityID
	w.link(e)
}

func (e *Entity) Detached() bool {
	return e.Incarnation.ID == level.InvalidEntityID
}

func (e *Entity) ZIndex() int {
	return e.zIndex
}

func (e *Entity) Contents() level.Contents {
	return e.contents
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

	// DebugPos64 returns the position and velocity for debug purposes, in subpixels.
	DebugPos64() (x int64, y int64, vx int64, vy int64)

	// Respawned() notifies the entity that the world respawned it.
	Respawned()
}

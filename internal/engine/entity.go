package engine

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"

	m "github.com/divVerent/aaaaaa/internal/math"
)

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
	// Ensures that entity state can be fully serialized/deserialized.
	PersistentState map[string]string
}

// An Entity is an object that exists in the game.
type Entity struct {
	ID             EntityID
	VisibilityMark uint

	// Info needed for gameplay.
	Solid  bool
	Opaque bool
	Rect   m.Rect

	// Info needed for rendering.
	Orientation m.Orientation
	Image       *ebiten.Image

	// Entity's own state.
	Impl EntityImpl
}

// EntityID represents an unique ID of an entity.
type EntityID int

type EntityImpl interface {
	// Spawn initializes the entity based on a Spawnable.
	// Will usually remember a reference to the World and Entity.
	// ID, Pos, Size and Orientation of the entity will be preset but may be changed.
	Spawn(w *World, s *Spawnable, e *Entity) error

	// Despawn notifies the entity that it will be deleted.
	Despawn()

	// Update asks the entity to update its state.
	Update()

	// Touch notifies the entity that it was hit by another entity moving.
	Touch(other *Entity)
}

// entityTypes is a helper map to know how to spawn an entity.
var entityTypes = map[string]EntityImpl{}

// RegisterEntityType adds an entity type to the spawn system.
// To be called from init() functions of entity implementations.
func RegisterEntityType(t EntityImpl) {
	typeName := reflect.TypeOf(t).Name()
	if entityTypes[typeName] != nil {
		log.Panicf("Duplicate entity type %q", typeName)
	}
	entityTypes[typeName] = t
	log.Printf("Registered entity type %q", typeName)
}

// Spawn turns a Spawnable into an Entity.
func (s *Spawnable) Spawn(w *World, tilePos m.Pos, t *Tile) (*Entity, error) {
	eImpl := entityTypes[s.EntityType]
	if eImpl == nil {
		return nil, fmt.Errorf("unknown entity type %q", s.EntityType)
	}
	e := &Entity{
		ID:   s.ID,
		Impl: eImpl,
	}
	tInv := t.Transform.Inverse()
	pivot2InTile := m.Pos{X: TileSize - 1, Y: TileSize - 1}
	e.Rect = tInv.ApplyToRect2(pivot2InTile, s.RectInTile)
	e.Rect.Origin = tilePos.Mul(TileSize).Add(e.Rect.Origin.Delta(m.Pos{}))
	e.Orientation = tInv.Concat(s.Orientation)
	err := eImpl.Spawn(w, s, e)
	if err != nil {
		return nil, err
	}
	// TODO w.LinkEntity(tilePos, e) - shall load in all tiles of e.
	return e, nil
}

package aaaaaa

import (
	"github.com/hajimehoshi/ebiten/v2"

	m "github.com/divVerent/aaaaaa/internal/math"
)

// An Entity is an object that exists in the game.
type Entity struct {
	ID             EntityID
	VisibilityMark uint

	// Info needed for gameplay.
	Solid  bool
	Opaque bool
	Pos    m.Pos
	Size   m.Delta

	// Info needed for rendering.
	Orientation m.Orientation
	Image       *ebiten.Image

	// Entity's own state.
	Impl EntityImpl
}

// EntityID represents an unique ID of an entity.
type EntityID int

type EntityImpl interface {
	// Update asks the entity to update its state.
	Update()

	// Touch notifies the entity that it was hit by another entity moving.
	Touch(other *Entity)
}

package aaaaaa

import (
	"github.com/hajimehoshi/ebiten/v2"

	m "github.com/divVerent/aaaaaa/internal/math"
)

// A Tile is a single game tile.
type Tile struct {
	// Info needed for gameplay.
	Solid      bool
	Opaque     bool
	Spawnables []*Spawnable // NOTE: not adjusted for transform!

	// Info needed for loading more tiles.
	LevelPos       m.Pos
	Transform      m.Orientation
	VisibilityMark uint

	// Info needed for rendering.
	Orientation m.Orientation
	Image       *ebiten.Image

	// Debug info.
	LoadedFromNeighbor m.Pos
}

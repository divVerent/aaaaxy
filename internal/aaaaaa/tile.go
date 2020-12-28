package aaaaaa

import (
	m "github.com/divVerent/aaaaaa/internal/math"
	"github.com/hajimehoshi/ebiten/v2"
)

// A Tile is a single game tile.
type Tile struct {
	// Info needed for gameplay.
	Solid      bool
	Spawnables []*Spawnable // NOTE: not adjusted for transform!

	// Info needed for loading more tiles.
	LevelPos       m.Pos
	Transform      m.Orientation
	VisibilityMark uint

	// Info needed for rendering.
	Orientation m.Orientation
	Image       *ebiten.Image
}

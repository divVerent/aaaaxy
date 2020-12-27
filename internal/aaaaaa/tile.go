package aaaaaa

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// A Tile is a single game tile.
type Tile struct {
	// Info needed for gameplay.
	Solid bool

	// Info needed for loading more tiles.
	LevelPos  Pos
	Transform Orientation

	// Info needed for rendering.
	Orientation Orientation
	Image       *ebiten.Image
}

package aaaaaa

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// A Tile is a single game tile.
type Tile struct {
	// Info needed for gameplay.
	Solid bool

	// Info needed for loading more tiles.
	levelPos levelPos

	// Info needed for rendering.
	orientation Orientation
	image       *ebiten.Image
}

package aaaaaa

// A Tile is a single game tile.
type Tile struct {
	// Info needed for gameplay.
	Solid bool

	// Info needed for loading more tiles.
	levelPos tilePos

	// Info needed for rendering.
	orientation Orientation
}

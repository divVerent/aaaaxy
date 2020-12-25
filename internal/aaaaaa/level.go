package aaaaaa

// level is a parsed form of a loaded level.
type level struct {
	tiles     map[levelPos]Tile
	entities  map[levelPos][]spawnable
	warpzones map[levelPos]warpzone
}

// levelPos is a position in the level.
type levelPos struct {
	c, r int
}

// levelTile is a single tile in the level.
type levelTile struct {
	tile       Tile
	spawnables []*spawnable
	warp       *warpzone
}

// warpzone represents a warp tile. Whenever anything enters this tile, it gets
// moved to "to" and the direction transformed by "transform". For the game to
// work, every warpzone must be paired with an exact opposite elsewhere. This
// is ensured at load time.
type warpzone struct {
	to        levelPos
	transform Orientation
}

type spawnable struct {
	// Entity ID. Used to decide what needs spawning. Unique within a level.
	id EntityID
}

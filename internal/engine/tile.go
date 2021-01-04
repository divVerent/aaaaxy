package engine

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
	visibilityMark uint

	// Info needed for rendering.
	Orientation m.Orientation
	Image       *ebiten.Image

	// If provided, these are used instead of image for "nicer" rotation (e.g. for shadow effects).
	// Because Orientation is also set, looking these up is tricky; we want things to show up as in the editor but potentially rotated.
	// We know:
	// - Transform * Orientation = orientationInEditor
	// - If we pick tile I and render at orientation O, we actually render at full orientation O * I.
	// - BUT lighting direction orientation is just O.
	// - we want O = orientationInEditor.
	// - Solve: Orientation = orientationInEditor * I
	// - Orientation = (Transform * Orientation) * I
	// - O = Transform Orientation
	// - I = O^-1 Orientation
	ImageByOrientation map[m.Orientation]*ebiten.Image

	// Debug info.
	LoadedFromNeighbor m.Pos
}

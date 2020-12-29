package aaaaaa

import (
	m "github.com/divVerent/aaaaaa/internal/math"
)

const (
	// GameWidth is the width of the displayed game area.
	GameWidth = 640
	// GameHeight is the height of the displayed game area.
	GameHeight = 360
	// GameTPS is the game ticks per second.
	GameTPS = 60

	// TileSize is the size of each tile graphics.
	TileSize = 16
	// SubPixels is the number of subpixels ("physics pixels") per tile pixel.
	// Note that only physics entities (like player) actually track their subpixels; subpixels do not actually exist elsewhere.
	SubPixels = 16
	// SweepStep is the distance between visibility traces in pixels. Lower means worse performance.
	SweepStep = 4
	// NumSweepTraces is the number of sweep operations we need.
	NumSweepTraces = 2 * (GameWidth + GameHeight) / SweepStep
	// ExpandSize is the amount of pixels to expand the visible area by.
	ExpandSize = 8
	// BlurSize is the amount of pixels to blur the visible area by.
	BlurSize = 4
	// ExpandTiles is the number of tiles beyond tiles hit by a trace that may need to be displayed.
	// As map design may need to take this into account, try to keep it at 1.
	ExpandTiles = (ExpandSize + BlurSize + SweepStep + TileSize - 1) / TileSize

	// PlayerWidth is the width of the player.
	PlayerWidth = TileSize - 2
	// PlayerHeight is the height of the player.
	PlayerHeight = 2*TileSize - 2
)

//ExpandStep is a single expansion step.
type ExpandStep struct {
	from, to m.Delta
}

var (
	// ExpandSteps is the list of steps to walk from each marked tile to expand.
	ExpandSteps = []ExpandStep{
		// First expansion tile.
		{m.Delta{}, m.Delta{1, 0}},
		{m.Delta{}, m.Delta{0, -1}},
		{m.Delta{}, m.Delta{-1, 0}},
		{m.Delta{}, m.Delta{0, 1}},
		{m.Delta{}, m.Delta{1, -1}},
		{m.Delta{}, m.Delta{-1, -1}},
		{m.Delta{}, m.Delta{-1, 1}},
		{m.Delta{}, m.Delta{1, 1}},
	}
)

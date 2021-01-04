package engine

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
	// sweepStep is the distance between visibility traces in pixels. Lower means worse performance.
	sweepStep = 4
	// numSweepTraces is the number of sweep operations we need.
	numSweepTraces = 2 * (GameWidth + GameHeight) / sweepStep
	// expandSize is the amount of pixels to expand the visible area by.
	expandSize = 6
	// blurSize is the amount of pixels to blur the visible area by.
	blurSize = 6
	// expandTiles is the number of tiles beyond tiles hit by a trace that may need to be displayed.
	// As map design may need to take this into account, try to keep it at 1.
	expandTiles = (expandSize + blurSize + sweepStep + TileSize - 1) / TileSize

	// MinEntitySize is the smallest allowed entity size.
	MinEntitySize = 8

	// frameBlurSize is how much the previous frame is to be blurred.
	frameBlurSize = 2
	// frameDarkenAlpha is how much the previous frame is to be darkened.
	frameDarkenAlpha = 0.98

	// How much to scroll towards focus point each frame.
	scrollPerFrame = 0.05
	// Minimum distance from screen edge when scrolling.
	scrollMinDistance = 2 * TileSize
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

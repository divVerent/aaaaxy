// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package engine

import (
	"github.com/divVerent/aaaaxy/internal/level"
	m "github.com/divVerent/aaaaxy/internal/math"
)

const (
	// GameWidth is the width of the displayed game area.
	GameWidth = 640
	// GameHeight is the height of the displayed game area.
	GameHeight = 360
	// GameTPS is the game ticks per second.
	GameTPS = 60

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
	expandTiles = (expandSize + blurSize + sweepStep + level.TileSize - 1) / level.TileSize

	// MinEntitySize is the smallest allowed entity size.
	MinEntitySize = 8

	// frameBlurSize is how much the previous frame is to be blurred.
	frameBlurSize = 1
	// frameDarkenAlpha is how much the previous frame is to be darkened relatively.
	frameDarkenAlpha = 0.99
	// frameDarkenAlpha is how much the previous frame is to be darkened absolutely.
	frameDarkenAmount = 1.0 / 255.0

	// How much to scroll towards focus point each frame.
	scrollPerFrame = 0.1
	// Minimum distance from screen edge when scrolling.
	scrollMinDistance = 2 * level.TileSize

	// Fully "fade in" in one second.
	pixelsPerSpawnFrame = (GameWidth / 2) / 60

	// borderWindowWidth is the maximum amount of pixels loaded outside the screen.
	// Must be at least the largest entity width plus two tiles to cover for misalignment.
	borderWindowWidth = 1264 + 2*level.TileSize
	// borderWindowHeight is the maximum amount of pixels loaded outside the screen.
	// Must be at least the largest entity height plus two tiles to cover for misalignment.
	borderWindowHeight = 6480 + 2*level.TileSize

	// tileWindowWidth is the maximum known width in tiles.
	tileWindowWidth = (GameWidth+2*borderWindowWidth+level.TileSize-2)/level.TileSize + 1
	// tileWindowHeight is the maximum known width in tiles.
	tileWindowHeight = (GameHeight+2*borderWindowHeight+level.TileSize-2)/level.TileSize + 1
)

//expandStep is a single expansion step.
type expandStep struct {
	from, from2, from3, to m.Delta
}

var (
	// expandSteps is the list of steps to walk from each marked tile to expand.
	expandSteps = []expandStep{
		// First expansion tile.
		{m.Delta{}, m.Delta{}, m.Delta{}, m.Delta{DX: 1, DY: 0}},
		{m.Delta{}, m.Delta{}, m.Delta{}, m.Delta{DX: 0, DY: -1}},
		{m.Delta{}, m.Delta{}, m.Delta{}, m.Delta{DX: -1, DY: 0}},
		{m.Delta{}, m.Delta{}, m.Delta{}, m.Delta{DX: 0, DY: 1}},
		{m.Delta{DX: 1, DY: 0}, m.Delta{DX: 0, DY: -1}, m.Delta{}, m.Delta{DX: 1, DY: -1}},
		{m.Delta{DX: -1, DY: 0}, m.Delta{DX: 0, DY: -1}, m.Delta{}, m.Delta{DX: -1, DY: -1}},
		{m.Delta{DX: -1, DY: 0}, m.Delta{DX: 0, DY: 1}, m.Delta{}, m.Delta{DX: -1, DY: 1}},
		{m.Delta{DX: 1, DY: 0}, m.Delta{DX: 0, DY: 1}, m.Delta{}, m.Delta{DX: 1, DY: 1}},
	}
)

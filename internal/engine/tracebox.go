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
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/m"
)

func (l *normalizedLine) traceBoxTiles(w *World, o TraceOptions, enlarge m.Delta, result *TraceResult) {
	if o.PathOut != nil {
		*o.PathOut = append(*o.PathOut, l.Origin.Div(level.TileSize))
	}
	// Find the corner in direction of the trace.
	var adjustment m.Delta
	if l.XDir > 0 {
		adjustment.DX = enlarge.DX
	}
	if l.YDir > 0 {
		adjustment.DY = enlarge.DY
	}

	// Trace that corner path against the tilemap.
	prevOrigin := l.Origin
	prevTarget := l.Target
	l.Origin = l.Origin.Add(adjustment)
	l.Target = l.Target.Add(adjustment)
	l.walkTiles(func(prevTile, nextTile m.Pos, delta m.Delta, prevPixelAdj, nextPixelAdj m.Pos) error {
		// First, unadjust.
		prevPixel := prevPixelAdj.Sub(adjustment)
		nextPixel := nextPixelAdj.Sub(adjustment)
		// Check the newly hit tiles.
		if nextPixel.X != prevPixel.X {
			// X move.
			// Check all newly hit tiles in Y range.
			top := m.Div(nextPixel.Y, level.TileSize)
			bottom := m.Div(nextPixel.Y+enlarge.DY, level.TileSize)
			for y := top; y <= bottom; y++ {
				tilePos := m.Pos{X: nextTile.X, Y: y}
				var tile *level.Tile
				if o.LoadTiles {
					tile = w.LoadTile(tilePos.Sub(delta), tilePos, delta)
				} else {
					tile = w.Tile(tilePos)
				}
				if tile == nil {
					result.EndPos = prevPixel
					result.HitDelta = delta
					// result.HitFogOfWar = true
					return errTraceDone
				}
				if o.Contents&tile.Contents != 0 {
					result.EndPos = prevPixel
					result.HitDelta = delta
					// result.HitTilePos = nextTile
					// result.HitTile = tile
					return errTraceDone
				}
			}
		} else {
			// Y move.
			// Check all newly hit tiles in X range.
			left := m.Div(nextPixel.X, level.TileSize)
			right := m.Div(nextPixel.X+enlarge.DX, level.TileSize)
			for x := left; x <= right; x++ {
				tilePos := m.Pos{X: x, Y: nextTile.Y}
				var tile *level.Tile
				if o.LoadTiles {
					tile = w.LoadTile(tilePos.Sub(delta), tilePos, delta)
				} else {
					tile = w.Tile(tilePos)
				}
				if tile == nil {
					result.EndPos = prevPixel
					result.HitDelta = delta
					// result.HitFogOfWar = true
					return errTraceDone
				}
				if o.Contents&tile.Contents != 0 {
					result.EndPos = prevPixel
					result.HitDelta = delta
					// result.HitTilePos = nextTile
					// result.HitTile = tile
					return errTraceDone
				}
			}
		}
		return nil
	})
	l.Origin = prevOrigin
	l.Target = prevTarget
}

func traceBox(w *World, from m.Rect, to m.Pos, o TraceOptions) TraceResult {
	if o.Contents == level.NoContents {
		log.Fatalf("do not know what to stop at - need to specify Contents in every trace")
	}
	result := TraceResult{
		EndPos: to,
		// HitTile:     nil,
		// HitEntities: nil,
		// HitFogOfWar: false,
	}

	if from.Origin == to {
		// Empty trace? Nothign we can hit.
		return result
	}

	l := normalizeLine(from.Origin, to)
	// As from != to, we know NumSteps > 0.

	enlarge := from.Size.Sub(m.Delta{DX: 1, DY: 1})
	maxBorder := 0
	if o.ForEnt != nil {
		maxBorder = o.ForEnt.BorderPixels
	}

	if !o.NoTiles {
		l.traceBoxTiles(w, o, enlarge, &result)
	}

	if !o.NoEntities {
		l.traceEntities(w, o, enlarge, maxBorder, &result)
	}

	return result
}

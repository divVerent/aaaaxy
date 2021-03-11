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
	"math"

	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
)

func (l *normalizedLine) traceBoxTiles(w *World, o TraceOptions, enlarge m.Delta, result *TraceResult) {
	result.EndPos = l.Origin
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
	// TODO: do this more efficiently than copying the line. We maybe can
	// even reuse the same adjustment for entity tracing? Or just edit the
	// l object and reset when done?
	ll := *l
	ll.Origin = ll.Origin.Add(adjustment)
	ll.Target = ll.Target.Add(adjustment)
	err := ll.walkTiles(func(prevTile, nextTile m.Pos, delta m.Delta, prevPixelAdj, nextPixelAdj m.Pos) error {
		// First, unadjust.
		prevPixel := prevPixelAdj.Sub(adjustment)
		nextPixel := nextPixelAdj.Sub(adjustment)
		// Record the EndPos as prevPixel was sure fine.
		result.EndPos = prevPixel
		// Check the newly hit tiles.
		if nextPixel.X != prevPixel.X {
			// X move.
			// Check all newly hit tiles in Y range.
			// TODO: One of these divisions is redundant. Worth optimizing?
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
					// result.HitFogOfWar = true
					return traceDoneErr
				}
				if o.Mode == HitSolid && tile.Solid || o.Mode == HitOpaque && tile.Opaque {
					// result.HitTilePos = nextTile
					// result.HitTile = tile
					return traceDoneErr
				}
			}
		} else {
			// Y move.
			// Check all newly hit tiles in X range.
			// TODO: One of these divisions is redundant. Worth optimizing?
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
					// result.HitFogOfWar = true
					return traceDoneErr
				}
				if o.Mode == HitSolid && tile.Solid || o.Mode == HitOpaque && tile.Opaque {
					// result.HitTilePos = nextTile
					// result.HitTile = tile
					return traceDoneErr
				}
			}
		}
		return nil
	})
	if err != traceDoneErr {
		result.EndPos = l.Target
	}
	result.Score = TraceScore{
		TraceDistance:  result.EndPos.Delta(l.Origin).Norm1(),
		EntityDistance: math.MaxInt32, // Not an entity.
	}
}

func traceBox(w *World, from m.Rect, to m.Pos, o TraceOptions) TraceResult {
	result := TraceResult{
		EndPos: to,
		// HitTile:     nil,
		// HitEntity:   nil,
		// HitFogOfWar: false,
	}

	if from.Origin == to {
		// Empty trace? Nothign we can hit.
		return result
	}

	l := normalizeLine(from.Origin, to)
	// As from != to, we know NumSteps > 0.

	enlarge := from.Size.Sub(m.Delta{DX: 1, DY: 1})

	if !o.NoTiles {
		l.traceBoxTiles(w, o, enlarge, &result)
	}

	if !o.NoEntities {
		l.traceEntities(w, o, enlarge, &result)
	}

	return result
}

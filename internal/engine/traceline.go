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

// traceLineTiles cuts the given trace by hits against the tilemap.
// l must have been initialized to finish at the current EndPos.
func (l *normalizedLine) traceLineTiles(w *World, o TraceOptions, result *TraceResult) {
	result.EndPos = l.Origin
	if o.PathOut != nil {
		*o.PathOut = append(*o.PathOut, l.Origin.Div(level.TileSize))
	}
	err := l.walkTiles(func(prevTile, nextTile m.Pos, delta m.Delta, prevPixel, nextPixel m.Pos) error {
		// Record the EndPos as the prevPixel was sure fine.
		result.EndPos = prevPixel
		// Check the newly hit tile(s).
		var tile *level.Tile
		if o.LoadTiles {
			tile = w.LoadTile(prevTile, nextTile, delta)
		} else {
			tile = w.Tile(nextTile)
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
		if o.PathOut != nil {
			*o.PathOut = append(*o.PathOut, nextTile)
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

// traceLine moves from from to to and yields info about where this hit solid etc.
func traceLine(w *World, from, to m.Pos, o TraceOptions) TraceResult {
	result := TraceResult{
		EndPos: to,
		// HitTile:     nil,
		// HitEntity:   nil,
		// HitFogOfWar: false,
	}

	if o.PathOut != nil {
		*o.PathOut = (*o.PathOut)[:0]
	}

	if from == to {
		// Empty trace? Nothign we can hit.
		return result
	}

	l := normalizeLine(from, to)
	// As from != to, we know NumSteps > 0.

	if !o.NoTiles {
		l.traceLineTiles(w, o, &result)
	}

	if !o.NoEntities {
		l.traceEntities(w, o, m.Delta{}, &result)
	}

	return result
}

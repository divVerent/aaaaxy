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

// traceLineTiles cuts the given trace by hits against the tilemap.
// l must have been initialized to finish at the current EndPos.
func (l *normalizedLine) traceLineTiles(w *World, o TraceOptions, result *TraceResult) {
	if o.PathOut != nil {
		*o.PathOut = append(*o.PathOut, l.Origin.Div(level.TileSize))
	}
	l.walkTiles(func(prevTile, nextTile m.Pos, delta m.Delta, prevPixel, nextPixel m.Pos) error {
		// Check the newly hit tile(s).
		var tile *level.Tile
		if o.LoadTiles {
			tile = w.LoadTile(prevTile, nextTile, delta)
		} else {
			tile = w.Tile(nextTile)
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
		if o.PathOut != nil {
			*o.PathOut = append(*o.PathOut, nextTile)
		}
		return nil
	})
}

// traceLine moves from from to to and yields info about where this hit solid etc.
func traceLine(w *World, from, to m.Pos, o TraceOptions) TraceResult {
	if o.Contents == level.NoContents {
		log.Fatalf("do not know what to stop at - need to specify Contents in every trace")
	}
	result := TraceResult{
		EndPos: to,
		// HitTile:     nil,
		// HitEntities: nil,
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
		l.traceEntities(w, o, m.Delta{}, 0, &result)
	}

	return result
}

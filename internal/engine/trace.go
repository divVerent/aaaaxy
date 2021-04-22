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
	"errors"

	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
)

type TraceOptions struct {
	// Contents is the OR'd set of contents to stop at (whether we want to do a visibility or collision trace).
	Contents level.Contents
	// If NoTiles is set, we ignore hits against tiles.
	NoTiles bool
	// If NoEntities is set, we ignore hits against entities.
	NoEntities bool
	// IgnoreEnt is the entity that shall be ignored when tracing.
	IgnoreEnt *Entity
	// ForEnt is the entity on whose behalf the trace is done.
	ForEnt *Entity
	// If LoadTiles is set, not yet known tiles will be loaded in by the trace operation.
	// Otherwise hitting a not-yet-loaded tile will end the trace.
	// Only valid on line traces.
	LoadTiles bool
	// If set, the trace path will be collected into this array. Provided here to reduce memory allocation.
	PathOut *[]m.Pos
}

// TraceResult returns the status of a trace operation.
type TraceResult struct {
	// EndPos is the pixel the trace ended on (the last nonsolid pixel).
	EndPos m.Pos
	// HitDelta is the one-pixel delta that hit the obstacle.
	HitDelta m.Delta
	// // HitTilePos is the position of the tile that stopped the trace, if any (in this case, HitTile will also be set).
	// HitTilePos m.Pos
	// // HitTile is the tile that stopped the trace, if any.
	// HitTile *level.Tile
	// HitEntity is the entity that stopped the trace, if any.
	HitEntity *Entity
	// // HitFogOfWar is set if the trace ended by hitting an unloaded tile.
	// HitFogOfWar bool
	// Score is a number used to decide which of multiple traces to keep.
	// Typically related to the trace distance and which entity was hit if any.
	Score TraceScore
}

// TraceScore is a scoring value of a trace.
type TraceScore struct {
	// TraceDistance is the length of the trace.
	TraceDistance int
	// EntityZ is the Z index of the entity hit.
	EntityZ int
	// EntityDistance is the distance between entity centers of the traces.
	// This is used as a tie breaker.
	EntityDistance int
}

// Less returns whether this score is smaller than the other.
func (s TraceScore) Less(o TraceScore) bool {
	// Prefer lower TraceDistance.
	if s.TraceDistance < o.TraceDistance {
		return true
	}
	if s.TraceDistance > o.TraceDistance {
		return false
	}
	// Prefer higher EntityZ.
	if s.EntityZ < o.EntityZ {
		return false
	}
	if s.EntityZ > o.EntityZ {
		return true
	}
	// Prefer lower EntityDistance.
	return s.EntityDistance < o.EntityDistance
}

// A normalizedLine represents a line to trace on.
// NOTE: Pixel (i, j) is on line IF:
// - Assuming:
//   - i0 := MAX(0, 2*i-1)
//   - i1: MIN(2*i+1, NumSteps)
//   - j0 := (Height*i0+NumSteps)/(2*NumSteps)
//   - j1 := (Height*i1+NumSteps)/(2*NumSteps)
// - Then: j in {j0, j1}.
type normalizedLine struct {
	Origin   m.Pos
	Target   m.Pos
	NumSteps int
	Height   int
	XDir     int
	YDir     int
	ScanX    bool
}

func normalizeLine(from, to m.Pos) normalizedLine {
	delta := to.Delta(from)
	absDelta := delta
	l := normalizedLine{
		Origin: from,
		Target: to,
	}
	l.XDir = 1
	if absDelta.DX < 0 {
		absDelta.DX = -absDelta.DX
		l.XDir = -1
	}
	l.YDir = 1
	if absDelta.DY < 0 {
		absDelta.DY = -absDelta.DY
		l.YDir = -1
	}
	l.ScanX = true
	l.NumSteps = absDelta.DX
	l.Height = absDelta.DY
	if l.NumSteps < l.Height {
		l.NumSteps, l.Height = l.Height, l.NumSteps
		l.ScanX = false
	}
	return l
}

func (l *normalizedLine) fromPos(p m.Pos) (int, int) {
	if l.ScanX {
		return (p.X - l.Origin.X) * l.XDir, (p.Y - l.Origin.Y) * l.YDir
	} else {
		return (p.Y - l.Origin.Y) * l.YDir, (p.X - l.Origin.X) * l.XDir
	}
}

func (l *normalizedLine) fromRect(r m.Rect, enlarge m.Delta) (int, int, int, int) {
	i0, j0 := l.fromPos(r.Origin.Sub(enlarge))
	i1, j1 := l.fromPos(r.OppositeCorner())
	if i0 > i1 {
		i0, i1 = i1, i0
	}
	if j0 > j1 {
		j0, j1 = j1, j0
	}
	return i0, j0, i1, j1
}

func (l *normalizedLine) toPos(i, j int) m.Pos {
	if l.ScanX {
		return m.Pos{X: l.Origin.X + l.XDir*i, Y: l.Origin.Y + l.YDir*j}
	} else {
		return m.Pos{X: l.Origin.X + l.XDir*j, Y: l.Origin.Y + l.YDir*i}
	}
}

func (l *normalizedLine) toDelta(u, v int) m.Delta {
	if l.ScanX {
		return m.Delta{DX: l.XDir * u, DY: l.YDir * v}
	} else {
		return m.Delta{DX: l.XDir * v, DY: l.YDir * u}
	}
}

// walkTiles yields all tile intersections on the line from start to end of the line.
func (l *normalizedLine) walkTiles(check func(prevTile, nextTile m.Pos, delta m.Delta, prevPixel, nextPixel m.Pos) error) error {
	// Algorithm idea:
	// - INIT: calculate iMod, jMod, scanI, scanJ.
	// - SEARCH:
	//   - Find nextI > scanI so that i % level.TileSize == iMod.
	//     - Actually can compute once, then just add level.TileSize.
	//   - Find nextJ > scanJ so that j % level.TileSize == jMod.
	//     - Actually can just conditionally add level.TileSize whenever we hit new tile.
	//   - Compute nextJI from nextJ like i00 below.
	//   - If nextI < nextJI:
	//     - Set nextJ = f(nextI) like j00.
	//     - Yield (nextI-1, nextJ) as endpos in current tile.
	//     - Set scanI, scanJ = nextI, nextJ.
	//   - If nextI == nextJI:
	//     - Set nextJ = f(nextI) like j00. Actually == nextJ.
	//     - Yield (nextI-1, nextJ) as endpos in current tile.
	//     - Yield (nextI, nextJ) as endpos in next tile.
	//     - Set scanI, scanJ = nextI, nextJ
	//   - If nextI > nextJI:
	//     - Yield (nextJI, nextJ-1) as endpos in current tile.
	//     - Set scanI, scanJ = nextJI, nextJ.
	// nextI, nextJ are the next i or j values that cross a tile border.
	tile := l.Origin.Div(level.TileSize)
	var nextI, nextJ int
	var iDelta, jDelta m.Delta
	if l.ScanX {
		if l.XDir > 0 {
			nextI = 1 + m.Mod(-l.Origin.X-1, level.TileSize)
		} else {
			nextI = 1 + m.Mod(l.Origin.X, level.TileSize)
		}
		if l.YDir > 0 {
			nextJ = 1 + m.Mod(-l.Origin.Y-1, level.TileSize)
		} else {
			nextJ = 1 + m.Mod(l.Origin.Y, level.TileSize)
		}
		iDelta = m.Delta{DX: l.XDir, DY: 0}
		jDelta = m.Delta{DX: 0, DY: l.YDir}
	} else {
		if l.YDir > 0 {
			nextI = 1 + m.Mod(-l.Origin.Y-1, level.TileSize)
		} else {
			nextI = 1 + m.Mod(l.Origin.Y, level.TileSize)
		}
		if l.XDir > 0 {
			nextJ = 1 + m.Mod(-l.Origin.X-1, level.TileSize)
		} else {
			nextJ = 1 + m.Mod(l.Origin.X, level.TileSize)
		}
		iDelta = m.Delta{DX: 0, DY: l.YDir}
		jDelta = m.Delta{DX: l.XDir, DY: 0}
	}
	if l.Height == 0 {
		// Special handling for x-only traces.
		for {
			if nextI > l.NumSteps {
				return nil
			}
			nextTile := tile.Add(iDelta)
			if err := check(tile, nextTile, iDelta, l.toPos(nextI-1, 0), l.toPos(nextI, 0)); err != nil {
				return err
			}
			tile = nextTile
			nextI += level.TileSize
		}
	}
	for {
		// Compute the i for nextJ. It is the SMALLEST i of the potential group.
		nextJI := (l.NumSteps*(2*nextJ-1) + l.Height - 1) / (2 * l.Height) // Same as i00 below.
		if nextJI < nextI {
			if nextJ > l.Height {
				return nil
			}
			nextTile := tile.Add(jDelta)
			if err := check(tile, nextTile, jDelta, l.toPos(nextJI, nextJ-1), l.toPos(nextJI, nextJ)); err != nil {
				return err
			}
			tile = nextTile
			nextJ += level.TileSize
		} else if nextJI > nextI {
			if nextI > l.NumSteps {
				return nil
			}
			nextTile := tile.Add(iDelta)
			// Compute the j for nextI. It is the SMALLEST j of the potential group.
			nextIJ := (l.Height*(2*nextI-1) + l.NumSteps) / (2 * l.NumSteps) // Same as j00 below.
			if err := check(tile, nextTile, iDelta, l.toPos(nextI-1, nextIJ), l.toPos(nextI, nextIJ)); err != nil {
				return err
			}
			tile = nextTile
			nextI += level.TileSize
		} else { // nextJI == nextI
			// We cross both boundaries.
			// By our line drawing algorithm, we always walk i first.
			if nextI > l.NumSteps {
				return nil
			}
			nextTile := tile.Add(iDelta)
			if err := check(tile, nextTile, iDelta, l.toPos(nextI-1, nextJ-1), l.toPos(nextI, nextJ-1)); err != nil {
				return err
			}
			tile = nextTile
			if nextJ > l.Height {
				return nil
			}
			nextTile = tile.Add(jDelta)
			if err := check(tile, nextTile, jDelta, l.toPos(nextI, nextJ-1), l.toPos(nextI, nextJ)); err != nil {
				return err
			}
			tile = nextTile
			nextI += level.TileSize
			nextJ += level.TileSize
		}
	}
}

var traceDoneErr = errors.New("traceDone")

// traceEntity returns whether the line from from to to hits the entity, as well as the last coordinate not hitting yet.
func (l *normalizedLine) traceEntity(ent *Entity, enlarge m.Delta) (bool, m.Pos, m.Delta) {
	i0, j0, i1, j1 := l.fromRect(ent.Rect, enlarge)
	if hit, i, j, u, v := traceLineBox(l.NumSteps, l.Height, i0, j0, i1, j1); hit {
		return true, l.toPos(i, j), l.toDelta(u, v)
	}
	// Not hit.
	return false, m.Pos{}, m.Delta{}
}

// traceEntities clips the given trace against all entities.
// l must have been initialized to hit the current EndPos anywhere on its path.
func (l *normalizedLine) traceEntities(w *World, o TraceOptions, enlarge m.Delta, result *TraceResult) {
	// Clip the trace to first entity hit.
	ents := w.FindContents(o.Contents)
	for _, ent := range ents {
		if ent == o.IgnoreEnt {
			continue
		}
		if hit, endPos, delta := l.traceEntity(ent, enlarge); hit {
			score := TraceScore{
				TraceDistance: endPos.Delta(l.Origin).Norm1(),
			}
			if o.ForEnt != nil {
				score.EntityDistance = ent.Rect.Center().Delta(o.ForEnt.Rect.Center()).Norm1()
				score.EntityZ = ent.ZIndex()
			}
			if score.Less(result.Score) {
				result.EndPos = endPos
				result.HitDelta = delta
				result.HitEntity = ent
				result.Score = score
			}
		}
	}
	if result.HitEntity != nil {
		endTile := result.EndPos.Div(level.TileSize)
		if o.PathOut != nil {
			for i, pos := range *o.PathOut {
				if pos == endTile {
					*o.PathOut = (*o.PathOut)[:(i + 1)]
				}
			}
		}
		// result.HitTilePos = m.Pos{}
		// result.HitTile = nil
		// result.HitFogOfWar = false
	}
}

// traceLineBox checks if from..to intersects with box, and if so, returns the pixel right before the intersection.
// i, j must be positive and i > j. The box is described by i0, j0, i1, j1 such that i0 <= i1 and j0 <= j1.
func traceLineBox(i, j, i0, j0, i1, j1 int) (bool, int, int, int, int) {
	// Is the box even hittable?
	if j < j0 || j1 < 0 {
		return false, 0, 0, 0, 0
	}
	if i < i0 || i1 < 0 {
		return false, 0, 0, 0, 0
	}
	if i0 <= 0 && j0 <= 0 {
		// We already overlap. Consider this a non-hit so we can get out of solid.
		return false, 0, 0, 0, 0
	}

	// TODO: some of the cases here may be redundant.
	// Do we really need four hit cases?
	// Maybe it's enough to just have the j0<=j00<=j1 and the i0<=i00<=i1 case?

	// Formula is: y(x) = x * j / i.
	// Pixels hit by x are thus: round(y(x-0.5)), round(y(x+0.5)).
	// Do we hit at i0?
	// Note that we can only ever hit two pixels at once because i > j >= 0.
	i200 := 2*i0 - 1
	if i200 < 0 {
		i200 = 0
	}
	i201 := 2*i0 + 1
	if i201 > 2*i {
		i201 = 2 * i
	}
	j00 := (j*i200 + i) / (2 * i)
	j01 := (j*i201 + i) / (2 * i)
	// Better to make this a range?
	if j00 >= j0 && j00 <= j1 {
		// Return the last pixel before hit.
		return true, i0 - 1, j00, 1, 0
	}
	if j01 >= j0 && j01 <= j1 {
		// Return the last pixel before.
		return true, i0, j01 - 1, 0, 1
	}

	if j == 0 {
		// i movement only.
		// But we already checked all we need.
		return false, 0, 0, 0, 0
	}

	// Do we hit at j0?
	// We need the first x so that round(y(x+0.5)) = j0 and the last x so that round(y(x-0.5)) = j0.
	// If that range intersects with [i0, i1] AND our valid range for i, we have a hit.
	j200 := 2*j0 - 1
	if j200 < 0 {
		j200 = 0
	}
	j201 := 2*j0 + 1
	if j201 > 2*j {
		j201 = 2 * j
	}
	i00 := (i*j200 + j - 1) / (2 * j) // Fulfills "translating to j01 yields j0" and is min.
	i01 := (i*j201 + j - 1) / (2 * j) // Fulfills "translating to j00 yields j0" and is max.

	// Compare ranges.
	if i00 >= i0 && i00 <= i1 {
		return true, i00, j0 - 1, 0, 1
	}
	iHit := i0
	if iHit < i00+1 {
		iHit = i00 + 1
	}
	if iHit <= i01 && iHit <= i1 {
		return true, iHit - 1, j0, 1, 0
	}

	// No hit.
	return false, 0, 0, 0, 0
}

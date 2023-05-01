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
	"sort"

	"github.com/divVerent/aaaaxy/internal/level"
	m "github.com/divVerent/aaaaxy/internal/math"
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
	// HitEntities are all the entities that stopped the trace simultaneously, if any.
	// They are sorted in decreasing order of closeness to the player; be aware that some code will only consider the first member.
	HitEntities []*Entity
	// // HitFogOfWar is set if the trace ended by hitting an unloaded tile.
	// HitFogOfWar bool
}

// TraceScore is a scoring value of a trace.
type traceScore struct {
	// TraceDistance is the length of the trace.
	traceDistance int
	// EntityZ is the Z index of the entity hit.
	entityZ int
	// EntityDistance is the distance between entity centers of the traces.
	// This is used as a tie breaker.
	entityDistance int
}

// CompareCoarse returns <0 if s < 0, >0 if s > 0, 0 otherwise.
func (s traceScore) CompareCoarse(o traceScore) int {
	// Prefer lower TraceDistance.
	d := s.traceDistance - o.traceDistance
	if d != 0 {
		return d
	}
	// Prefer higher EntityZ.
	return o.entityZ - s.entityZ
}

// CompareFine returns <0 if s < 0, >0 if s > 0, 0 otherwise, assuming CompareCoarse was 0.
func (s traceScore) CompareFine(o traceScore) int {
	// Prefer lower EntityDistance.
	return s.entityDistance - o.entityDistance
}

// A normalizedLine represents a line to trace on.
// NOTE: Pixel (i, j) is on line IF:
// - Assuming:
//   - i0 := MAX(0, 2*i-1)
//   - i1: MIN(2*i+1, NumSteps)
//   - j0 := (Height*i0+NumSteps)/(2*NumSteps)
//   - j1 := (Height*i1+NumSteps)/(2*NumSteps)
//
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

var errTraceDone = errors.New("traceDone")

// traceEntity returns whether the line from from to to hits the entity, as well as the last coordinate not hitting yet.
func (l *normalizedLine) traceEntity(ent *Entity, enlarge m.Delta, maxBorder int) (bool, m.Pos, m.Delta) {
	border := ent.BorderPixels
	if border > maxBorder {
		border = maxBorder
	}
	i0, j0, i1, j1 := l.fromRect(ent.Rect, enlarge)
	i0 -= border
	j0 -= border
	i1 += border
	j1 += border
	if hit, i, j, u, v := traceLineBox(l.NumSteps, l.Height, i0, j0, i1, j1); hit {
		return true, l.toPos(i, j), l.toDelta(u, v)
	}
	// Not hit.
	return false, m.Pos{}, m.Delta{}
}

type traceHit struct {
	endPos    m.Pos
	hitDelta  m.Delta
	hitEntity *Entity
	score     traceScore
}

// traceEntities clips the given trace against all entities.
// l must have been initialized to hit the current EndPos anywhere on its path.
func (l *normalizedLine) traceEntities(w *World, o TraceOptions, enlarge m.Delta, maxBorder int, result *TraceResult) {
	worldDist := result.EndPos.Delta(l.Origin).Norm1()

	// Clip the trace to first entity hit.
	ents := w.FindContents(o.Contents)

	var hits []traceHit

	for _, ent := range ents {
		if ent == o.IgnoreEnt {
			continue
		}
		if hit, endPos, delta := l.traceEntity(ent, enlarge, maxBorder); hit {
			dist := endPos.Delta(l.Origin).Norm1()
			if dist > worldDist {
				continue
			}
			score := traceScore{
				traceDistance: dist,
				entityZ:       ent.ZIndex(),
			}
			if o.ForEnt != nil {
				score.entityDistance = ent.Rect.Center().Delta(o.ForEnt.Rect.Center()).Norm1()
			}
			if len(hits) != 0 {
				cmp := score.CompareCoarse(hits[0].score)
				if cmp > 0 {
					continue
				}
				if cmp < 0 {
					hits = hits[:0]
				}
			}
			hits = append(hits, traceHit{
				endPos:    endPos,
				hitDelta:  delta,
				hitEntity: ent,
				score:     score,
			})
		}
	}

	if len(hits) == 0 {
		return
	}

	// Move the closest hit to the start.
	// Yes, this may be more expensive, but it makes the game usually more deterministic regarding touch event ordering.
	sort.SliceStable(hits, func(i, j int) bool {
		return hits[i].score.CompareFine(hits[j].score) < 0
	})

	// Return all trace hits.
	result.HitEntities = make([]*Entity, len(hits))
	for i, hit := range hits {
		result.HitEntities[i] = hit.hitEntity
	}

	// Return the closest hit properties.
	result.EndPos = hits[0].endPos
	result.HitDelta = hits[0].hitDelta

	// Return the end tile.
	endTile := result.EndPos.Div(level.TileSize)
	if o.PathOut != nil {
		for i, pos := range *o.PathOut {
			if pos == endTile {
				*o.PathOut = (*o.PathOut)[:(i + 1)]
			}
		}
	}

	// Fields that no longer exist:
	// result.HitTilePos = m.Pos{}
	// result.HitTile = nil
	// result.HitFogOfWar = false
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

	// Formula is: y(x) = x * j / i.
	// Pixels hit by x are thus: round(y(x-0.5)), round(y(x+0.5)).
	// Do we hit at i0?
	// Note that we can only ever hit two pixels at once because i > j >= 0.
	i200 := 2*i0 - 1
	if i200 < 0 {
		i200 = 0
	}
	j00 := (j*i200 + i) / (2 * i)
	// If the collision happens in i direction, it must be when entering a column, not when leaving it.
	// Note that the collision may still happen in the same column but in j direction.
	if j00 >= j0 && j00 <= j1 {
		// Return the last pixel before hit.
		return true, i0 - 1, j00, 1, 0
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
	i00 := (i*j200 + j - 1) / (2 * j) // Fulfills "translating to j01 yields j0" and is min.

	// Compare ranges.
	// If the collision happens in j direction, it must be when entering a column.
	// If the collision happens later in i direction, the code above should have already caught it.
	if i00 >= i0 && i00 <= i1 {
		return true, i00, j0 - 1, 0, 1
	}

	// No hit.
	return false, 0, 0, 0, 0
}

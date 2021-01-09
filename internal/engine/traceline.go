package engine

import (
	"errors"
	"log"

	m "github.com/divVerent/aaaaaa/internal/math"
)

// TraceMode indicates what kind of tiles/objects we want to hit.
type TraceMode int

const (
	// HitSolid indicates we want to hit solid (non passable) tiles.
	HitSolid TraceMode = iota

	// HitOpaque indicates we want to hit opaque (non see through) tiles.
	HitOpaque
)

type TraceOptions struct {
	// Mode is the TraceMode to trace by (whether we want to do a visibility or collision trace).
	Mode TraceMode
	// If NoTiles is set, we ignore hits against tiles.
	NoTiles bool
	// If NoEntities is set, we ignore hits against entities.
	NoEntities bool
	// IgnoreEnt is the entity that shall be ignored when tracing.
	IgnoreEnt *Entity
	// If LoadTiles is set, not yet known tiles will be loaded in by the trace operation.
	// Otherwise hitting a not-yet-loaded tile will end the trace.
	// Only valid on line traces.
	LoadTiles bool
}

// TraceResult returns the status of a trace operation.
type TraceResult struct {
	// EndPos is the pixel the trace ended on (the last nonsolid pixel).
	EndPos m.Pos
	// Path is the set of tiles touched, not including what stopped the trace.
	// Only set by line traces.
	Path []m.Pos
	// hitTilePos is the position of the tile that stopped the trace, if any.
	HitTilePos *m.Pos
	// HitTile is the tile that stopped the trace, if any.
	HitTile *Tile
	// HitEntity is the entity that stopped the trace, if any.
	HitEntity *Entity
	// HitFogOfWar is set if the trace ended by hitting an unloaded tile.
	HitFogOfWar bool
}

type normalizedLine struct {
	Origin   m.Pos
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

func (l *normalizedLine) fromRect(r m.Rect) (int, int, int, int) {
	i0, j0 := l.fromPos(r.Origin)
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

// traceEntity returns whether the line from from to to hits the entity, as well as the last coordinate not hitting yet.
func traceEntity(from, to m.Pos, ent *Entity) (bool, m.Pos) {
	l := normalizeLine(from, to)
	if l.NumSteps == 0 {
		// Start point is end point. Nothing to do.
		return false, m.Pos{}
	}
	i0, j0, i1, j1 := l.fromRect(ent.Rect)
	if hit, i, j := traceLineBox(l.NumSteps, l.Height, i0, j0, i1, j1); hit {
		return true, l.toPos(i, j)
	}
	// Not hit.
	return false, m.Pos{}
}

func (l *normalizedLine) nextTile(i int) int {
	// Locate the next end-of-tile position.
	// We are at end of tile if either (potentially swapped x/y):
	// - l.Origin.X + l.XDir*i == l.XDir>0 ? TileSize-1 : 0
	//   - Can precompute and then always advance by exactly TileSize.
	// - l.Origin.Y + l.YDir*j0 == l.YDir>0 ? TileSize-1 : 0 AND j1 != j0
	//   - where j0 = (l.Height*(2*i-1) + l.NumSteps) / twiceSteps
	//   - where j1 = (l.Height*(2*i+1) + l.NumSteps) / twiceSteps
	//   - Can we rule out "shooting ahead" one pixel another way? Sure.
	//   - Can we precompute?
	return 0 // TODO.
}

// walkLine walks on pixels from from to to, calling the check() function on every pixel hit.
// Any two adjacent positions hit are exactly 1 pixel apart (no diagonal steps).
func walkLine(from, to m.Pos, check func(pixel m.Pos) error) error {
	l := normalizeLine(from, to)
	if l.NumSteps == 0 {
		// Start point is end point. Nothing to do.
		return check(from)
	}
	twiceSteps := 2 * l.NumSteps
	prevPixel, prevTile := from, from.Div(TileSize)
	for i := 0; i <= l.NumSteps; i++ {
		i0 := 2*i - 1
		if i0 < 0 {
			i0 = 0
		}
		i1 := 2*i + 1
		if i1 > 2*l.NumSteps {
			i1 = 2 * l.NumSteps
		}
		j0 := (l.Height*i0 + l.NumSteps) / twiceSteps
		j1 := (l.Height*i1 + l.NumSteps) / twiceSteps
		for j := j0; j <= j1; j++ {
			pixel := l.toPos(i, j)
			// Only call the callback if we hit the end of a tile, or the end of the trace.
			// Should speed up tracing SUBSTANTIALLY by saving lots of callback invocations.
			tile := pixel.Div(TileSize)
			if tile != prevTile {
				err := check(prevPixel)
				if err != nil {
					return err
				}
				prevTile = tile
			}
			prevPixel = pixel
		}
	}
	return check(to)
}

// traceLine moves from from to to and yields info about where this hit solid etc.
func traceLine(w *World, from, to m.Pos, o TraceOptions) TraceResult {
	// TODO write an optimized implementation. We do the naive one first.

	result := TraceResult{
		EndPos:      to,
		Path:        nil,
		HitTilePos:  nil,
		HitTile:     nil,
		HitEntity:   nil,
		HitFogOfWar: false,
	}

	if !o.NoTiles {
		result.EndPos = from
		var prevTilePos m.Pos
		havePrevTile := false
		doneErr := errors.New("done")
		walkLine(from, to, func(pixel m.Pos) error {
			tilePos := pixel.Div(TileSize)
			if tilePos == prevTilePos && havePrevTile {
				result.EndPos = pixel
				return nil
			}
			if o.LoadTiles && havePrevTile {
				w.LoadTile(prevTilePos, tilePos.Delta(prevTilePos))
			}
			tile := w.Tiles[tilePos]
			if tile == nil {
				if !havePrevTile {
					log.Panicf("Traced from nonexistent tile %v (loaded: %v)", tilePos, w.Tiles)
				}
				result.HitFogOfWar = true
				return doneErr
			}
			if o.Mode == HitSolid && tile.Solid || o.Mode == HitOpaque && tile.Opaque {
				result.HitTilePos = &tilePos
				result.HitTile = tile
				return doneErr
			}
			result.Path = append(result.Path, tilePos)
			havePrevTile = true
			prevTilePos = tilePos
			result.EndPos = pixel
			return nil
		})
	}

	if !o.NoEntities {
		// Clip the trace to first entity hit.
		var closestEnt *Entity
		var closestEndPos m.Pos
		closestDistance := result.EndPos.Delta(from).Norm1()
		for _, ent := range w.Entities {
			if ent == o.IgnoreEnt {
				continue
			}
			if o.Mode == HitSolid && !ent.Solid || o.Mode == HitOpaque && !ent.Opaque {
				continue
			}
			if hit, endPos := traceEntity(from, to, ent); hit {
				distance := endPos.Delta(from).Norm1()
				if distance < closestDistance {
					closestEnt, closestEndPos, closestDistance = ent, endPos, distance
				}
			}
		}
		if closestEnt != nil {
			result.EndPos = closestEndPos
			endTile := closestEndPos.Div(TileSize)
			for i, pos := range result.Path {
				if pos == endTile {
					result.Path = result.Path[:(i + 1)]
				}
			}
			result.HitTilePos = nil
			result.HitTile = nil
			result.HitEntity = closestEnt
			result.HitFogOfWar = false
		}
	}

	return result
}

// traceLineBox checks if from..to intersects with box, and if so, returns the pixel right before the intersection.
// i, j must be positive and i > j. The box is described by i0, j0, i1, j1 such that i0 <= i1 and j0 <= j1.
func traceLineBox(i, j, i0, j0, i1, j1 int) (bool, int, int) {
	// Is the box even hittable?
	if j < j0 || j1 < 0 {
		return false, 0, 0
	}
	if i < i0 || i1 < 0 {
		return false, 0, 0
	}
	if i0 <= 0 && j0 <= 0 {
		// We already overlap. Consider this a non-hit so we can get out of solid.
		return false, 0, 0
	}

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
		return true, i0 - 1, j00
	}
	if j01 >= j0 && j01 <= j1 {
		// Return the last pixel before.
		return true, i0, j01 - 1
	}

	if j == 0 {
		// i movement only.
		// But we already checked all we need.
		return false, 0, 0
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
		return true, i00, j0 - 1
	}
	iHit := i0
	if iHit < i00+1 {
		iHit = i00 + 1
	}
	if iHit <= i01 && iHit <= i1 {
		return true, iHit - 1, j0
	}

	// No hit.
	return false, 0, 0
}

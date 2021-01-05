package engine

import (
	"errors"

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
	// ForEnt is the entity that shall be ignored when tracing.
	ForEnt *Entity
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

// walkLine walks on pixels from from to to, calling the check() function on every pixel hit.
// Any two adjacent positions hit are exactly 1 pixel apart (no diagonal steps).
func walkLine(from, to m.Pos, check func(pixel m.Pos) error) error {
	delta := to.Delta(from)
	absDelta := delta
	xDir := 1
	if absDelta.DX < 0 {
		absDelta.DX = -absDelta.DX
		xDir = -1
	}
	yDir := 1
	if absDelta.DY < 0 {
		absDelta.DY = -absDelta.DY
		yDir = -1
	}
	scanX := true
	numSteps := absDelta.DX
	height := absDelta.DY
	if numSteps < height {
		numSteps, height = height, numSteps
		scanX = false
	}
	if numSteps == 0 {
		// Start point is end point. Nothing to do.
		return check(from)
	}
	twiceSteps := 2 * numSteps
	prevPixel, prevTile := from, from.Div(TileSize)
	for i := 0; i <= numSteps; i++ {
		i0 := 2*i - 1
		if i0 < 0 {
			i0 = 0
		}
		i1 := 2*i + 1
		if i1 > 2*numSteps {
			i1 = 2 * numSteps
		}
		j0 := (height*i0 + numSteps) / twiceSteps
		j1 := (height*i1 + numSteps) / twiceSteps
		for j := j0; j <= j1; j++ {
			var pixel m.Pos
			if scanX {
				pixel = m.Pos{X: from.X + xDir*i, Y: from.Y + yDir*j}
			} else {
				pixel = m.Pos{X: from.X + xDir*j, Y: from.Y + yDir*i}
			}
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
			if tilePos == prevTilePos {
				result.EndPos = pixel
				return nil
			}
			if o.LoadTiles && havePrevTile {
				w.LoadTile(prevTilePos, tilePos.Delta(prevTilePos))
			}
			tile := w.Tiles[tilePos]
			if tile == nil {
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
		for i, pos := range result.Path {
			i = i
			pos = pos
			// if pos is on an entity, clip here
		}
	}

	return result
}

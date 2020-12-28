package aaaaaa

import (
	"log"

	m "github.com/divVerent/aaaaaa/internal/math"

	"github.com/hajimehoshi/ebiten/v2"
)

// World represents the current game state including its entities.
type World struct {
	// tiles are all tiles currently loaded.
	Tiles map[m.Pos]*Tile
	// entities are all entities currently loaded.
	Entities map[EntityID]*Entity
	// Player is the ID of the player entity.
	PlayerID EntityID
	// scrollPos is the current screen scrolling position.
	ScrollPos m.Pos
	// scrollTarget is where we want to scroll to.
	ScrollTarget m.Pos
	// scrollSpeed is the speed of scrolling to ScrollTarget, or 0 if not aiming for a target.
	ScrollSpeed int
	// level is the current tilemap (universal covering with warpzones).
	Level *Level
	// SpawnMark is the current mark value to detect visible tiles/objects.
	SpawnMark uint
}

func NewWorld() *World {
	// Load map.
	level, err := LoadLevel("map")
	if err != nil {
		log.Panicf("Could not load level: %v", err)
	}
	w := World{
		Level: level,
	}

	// Create player entity.
	w.PlayerID = w.Level.Player.ID
	// TODO actually spawn the player properly.
	w.Entities[w.PlayerID] = &Entity{
		ID:   w.Level.Player.ID,
		Pos:  w.Level.Player.LevelPos.Scale(TileSize, 1).Add(w.Level.Player.PosInTile),
		Size: w.Level.Player.Size,
	}

	// Load in the tiles the player is standing on.
	tile := w.Level.Tiles[w.Level.Player.LevelPos].Tile
	w.Tiles[w.Level.Player.LevelPos] = &tile
	w.LoadTilesForBox(w.Entities[w.PlayerID].Pos, w.Entities[w.PlayerID].Size, w.Level.Player.LevelPos)

	return &w
}

func (w *World) traceLineAndMark(p m.Pos, d m.Delta) TraceResult {
	result := w.TraceLine(p, d, TraceOptions{
		LoadTiles: true,
	})
	for _, tilePos := range result.Path {
		w.Tiles[tilePos].SpawnMark = w.SpawnMark
	}
	return result
}

func (w *World) Update() error {
	// TODO Let all entities move/act. Fetch player position.

	// Update ScrollPos based on player position and scroll target.
	player := w.Entities[w.PlayerID]
	w.ScrollPos = player.Pos

	// Unmark all tiles and entities (just bump mark index).
	w.SpawnMark++

	// Trace from player location to all directions (SweepStep pixels at screen edge).
	// Mark all tiles hit (excl. the tiles that stopped us).
	// TODO Remember trace polygon.
	screen0 := w.ScrollPos.Sub(m.Delta{DX: GameWidth / 2, DY: GameHeight / 2}).Delta(player.Pos)
	screen1 := screen0.Add(m.Delta{DX: GameWidth - 1, DY: GameHeight - 1})
	for x := screen0.DX; x < screen1.DX+SweepStep; x += SweepStep {
		w.traceLineAndMark(player.Pos, m.Delta{DX: x, DY: screen0.DY})
		w.traceLineAndMark(player.Pos, m.Delta{DX: x, DY: screen1.DY})
	}
	for y := screen0.DY; y < screen1.DY+SweepStep; y += SweepStep {
		w.traceLineAndMark(player.Pos, m.Delta{DX: screen0.DX, DY: y})
		w.traceLineAndMark(player.Pos, m.Delta{DX: screen1.DX, DY: y})
	}

	// Also mark all neighbors of hit tiles hit (up to ExpandTiles).
	markedTiles := []m.Pos{}
	for tilePos, tile := range w.Tiles {
		if tile.SpawnMark == w.SpawnMark {
			markedTiles = append(markedTiles, tilePos)
		}
	}
	expand := m.Delta{DX: ExpandTiles, DY: ExpandTiles}
	for _, pos := range markedTiles {
		w.LoadTilesForTileBox(pos.Sub(expand), pos.Add(expand), pos)
	}

	// TODO Mark all entities on marked tiles hit.
	// TODO Delete all unmarked entities.
	// TODO Spawn all entities on marked tiles if not already spawned.
	// TODO Mark all tiles on entities (this is NOT recursive, but entities may require the tiles they are on to be loaded so they can move).
	// (Somewhat tricky as entities may stand on warps; we have to walk from a known tile).

	// Delete all unmarked tiles.
	for pos, tile := range w.Tiles {
		if tile.SpawnMark != w.SpawnMark {
			delete(w.Tiles, pos)
		}
	}

	return nil
}

func (w *World) Draw(screen *ebiten.Image) {
	screen.Clear()

	// TODO Draw trace polygon to buffer.
	// TODO Expand and blur buffer (ExpandSize, BlurSize).

	// Draw all tiles.
	for pos, tile := range w.Tiles {
		screenPos := m.Pos{X: GameWidth / 2, Y: GameHeight / 2}.Add(pos.Scale(TileSize, 1).Delta(w.ScrollPos))
		opts := ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeCopy,
			Filter:        ebiten.FilterNearest,
		}
		opts.GeoM.SetElement(0, 0, float64(tile.Orientation.Right.DX))
		opts.GeoM.SetElement(0, 1, float64(tile.Orientation.Right.DY))
		opts.GeoM.SetElement(1, 0, float64(tile.Orientation.Down.DX))
		opts.GeoM.SetElement(1, 1, float64(tile.Orientation.Down.DY))
		opts.GeoM.Translate(float64(screenPos.X), float64(screenPos.Y))
		screen.DrawImage(tile.Image, &opts)
	}

	// TODO Draw all entities.
	// NOTE: if an entity is on a tile seen twice, render only once.
	// INTENTIONAL GLITCH (avoids rendering player twice and player-player collision). Entities live in tile coordinates, not world coordinates. "Looking away" can despawn these entities and respawn at their new location.
	// Makes wrap-around rooms somewhat less obvious.
	// Only way to fix seems to be making everything live in "universal covering" coordinates with orientation? Seems not worth it.
	// TODO: Decide if to keep this.
	// Multiply screen with buffer.
	// Invert buffer.
	// Multiply with previous screen, scroll pos delta applied.
	// Blur and darken buffer.
	// Add buffer to screen.
}

// LoadTile loads the next tile into the current world based on a currently
// known tile and its neighbor. Respects and applies warps.
func (w *World) LoadTile(p m.Pos, d m.Delta) m.Pos {
	newPos := p.Add(d)
	if tile, found := w.Tiles[newPos]; found {
		// Already loaded.
		tile.SpawnMark = w.SpawnMark
		return newPos
	}
	neighborTile := w.Tiles[p]
	t := neighborTile.Transform
	newLevelPos := neighborTile.LevelPos.Add(t.Apply(d))
	newLevelTile, found := w.Level.Tiles[newLevelPos]
	if !found {
		log.Panicf("Trying to load nonexisting tile at %v when moving from %v (%v) by %v (%v)",
			newLevelPos, p, neighborTile.LevelPos, d, t.Apply(d))
	}
	newTile := newLevelTile.Tile
	newTile.Transform = t.Concat(newTile.Transform)
	newTile.Orientation = t.Concat(newTile.Orientation)
	newTile.SpawnMark = w.SpawnMark
	w.Tiles[newPos] = &newTile
	return newPos
}

// LoadTilesForBox loads all tiles in the given box (p, d), assuming tile tp is already loaded.
func (w *World) LoadTilesForBox(p m.Pos, d m.Delta, tp m.Pos) {
	// Convert box to tile positions.
	tp0 := p.Scale(1, TileSize)
	tp1 := p.Add(d).Add(m.Delta{DX: -1, DY: -1}).Scale(1, TileSize)
	w.LoadTilesForTileBox(tp0, tp1, tp)
}

// LoadTilesForTileBox loads all tiles in the given tile based box, assuming tile tp is already loaded.
func (w *World) LoadTilesForTileBox(tp0, tp1, tp m.Pos) {
	// In range, load all.
	for y := tp.Y; y > tp0.Y; y-- {
		w.LoadTile(m.Pos{X: tp.X, Y: y}, m.North())
	}
	for y := tp.Y; y < tp1.Y; y++ {
		w.LoadTile(m.Pos{X: tp.X, Y: y}, m.South())
	}
	for y := tp0.Y; y <= tp1.Y; y++ {
		for x := tp.X; x > tp0.X; x-- {
			w.LoadTile(m.Pos{X: x, Y: y}, m.West())
		}
		for x := tp.X; x < tp1.X; x++ {
			w.LoadTile(m.Pos{X: x, Y: y}, m.East())
		}
	}
}

type TraceOptions struct {
	// If NoTiles is set, we ignore hits against tiles.
	NoTiles bool
	// If NoEntities is set, we ignore hits against entities.
	NoEntities bool
	// If LoadTiles is set, not yet known tiles will be loaded in by the trace operation.
	// Otherwise hitting a not-yet-loaded tile will end the trace.
	LoadTiles bool
}

// TraceResult returns the status of a trace operation.
type TraceResult struct {
	// Delta is the distance actually travelled until the trace stopped.
	Vector m.Delta
	// Path is the set of tiles touched, not including what stopped the trace.
	// For a line trace, any two neighboring tiles here are adjacent.
	Path []m.Pos
	// Entities is the set of entities touched, not including what stopped the trace.
	Entities []Entity
	// hitSolidTilePos is the position of the tile that stopped the trace, if any.
	HitSolidTilePos *m.Pos
	// HitSolidTile is the tile that stopped the trace, if any.
	HitSolidTile *Tile
	// HitSolidEntity is the entity that stopped the trace, if any.
	HitSolidEntity Entity
	// HitFogOfWar is set if the trace ended by hitting an unloaded tile.
	HitFogOfWar bool
}

// TraceLine moves from x,y by dx,dy in pixel coordinates.
func (w *World) TraceLine(p m.Pos, d m.Delta, o TraceOptions) TraceResult {
	return w.TraceBox(p, m.Delta{}, d, o)
}

// TraceBox moves from x,y size sx,sy by dx,dy in pixel coordinates.
func (w *World) TraceBox(p m.Pos, s, d m.Delta, o TraceOptions) TraceResult {
	// TODO: Implement
	return TraceResult{}
}

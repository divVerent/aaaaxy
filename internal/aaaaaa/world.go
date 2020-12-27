package aaaaaa

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// World represents the current game state including its entities.
type World struct {
	// tiles are all tiles currently loaded.
	Tiles map[Pos]Tile
	// entities are all entities currently loaded.
	Entities map[EntityID]*Entity
	// scrollPos is the current screen scrolling position.
	ScrollPos Pos
	// scrollTarget is where we want to scroll to.
	ScrollTarget Pos
	// scrollSpeed is the speed of scrolling to ScrollTarget, or 0 if not aiming for a target.
	ScrollSpeed int
	// level is the current tilemap (universal covering with warpzones).
	Level *Level
}

func NewWorld() *World {
	// Load map.
	// Create player entity.
	// Load in the tile the player is standing on.
	return &World{}
}

func (w *World) Update() error {
	// Let all entities move/act. Fetch player position.
	// Update ScrollPos based on player position and scroll target.
	// Unmark all tiles and entities (just bump mark index).
	// Trace from player location to all directions (SweepStep pixels at screen edge).
	// Remember trace polygon.
	// Mark all tiles hit (excl. the tiles that stopped us).
	// Also mark all neighbors of hit tiles hit (up to ExpandTiles).
	// Mark all entities on marked tiles hit.
	// Delete all unmarked entities.
	// Spawn all entities on marked tiles if not already spawned.
	// Mark all tiles on entities (this is NOT recursive, but entities may require the tiles they are on to be loaded so they can move).
	// Delete all unmarked tiles.
	return nil
}

func (w *World) Draw(screen *ebiten.Image) {
	// Draw trace polygon to buffer.
	// Expand and blur buffer (ExpandSize, BlurSize).
	// Draw all tiles.
	// Draw all entities.
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
func (w *World) LoadTile(p Pos, d Delta) Pos {
	// TODO implement
	return Pos{}
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
	Vector Delta
	// Path is the set of tiles touched, not including what stopped the trace.
	// For a line trace, any two neighboring tiles here are adjacent.
	Path []Pos
	// Entities is the set of entities touched, not including what stopped the trace.
	Entities []Entity
	// hitSolidTilePos is the position of the tile that stopped the trace, if any.
	HitSolidTilePos *Pos
	// HitSolidTile is the tile that stopped the trace, if any.
	HitSolidTile *Tile
	// HitSolidEntity is the entity that stopped the trace, if any.
	HitSolidEntity Entity
	// HitFogOfWar is set if the trace ended by hitting an unloaded tile.
	HitFogOfWar bool
}

// TraceLine moves from x,y by dx,dy in pixel coordinates.
func (w *World) TraceLine(p Pos, d Delta, o TraceOptions) TraceResult {
	// TODO: Optimize?
	return w.TraceBox(p, Delta{}, d, o)
}

// TraceBox moves from x,y size sx,sy by dx,dy in pixel coordinates.
func (w *World) TraceBox(p Pos, s, d Delta, o TraceOptions) TraceResult {
	// TODO: Implement
	return TraceResult{}
}

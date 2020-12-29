package aaaaaa

import (
	"fmt"
	"image/color"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"

	m "github.com/divVerent/aaaaaa/internal/math"
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
	// VisibilityMark is the current mark value to detect visible tiles/objects.
	VisibilityMark uint
	// DebugFont is the font to use for debug messages.
	DebugFont font.Face
}

func NewWorld() *World {
	// Load font.
	debugFont, err := truetype.Parse(gomono.TTF)
	if err != nil {
		log.Panicf("Could not load font: %v", err)
	}

	// Load map.
	level, err := LoadLevel("level")
	if err != nil {
		log.Panicf("Could not load level: %v", err)
	}
	w := World{
		Tiles:    map[m.Pos]*Tile{},
		Entities: map[EntityID]*Entity{},
		Level:    level,
		DebugFont: truetype.NewFace(debugFont, &truetype.Options{
			Size:    5,
			Hinting: font.HintingFull,
		}),
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
	tile.Transform = m.Identity()
	w.Tiles[w.Level.Player.LevelPos] = &tile
	w.LoadTilesForBox(w.Entities[w.PlayerID].Pos, w.Entities[w.PlayerID].Size, w.Level.Player.LevelPos)

	return &w
}

func (w *World) traceLineAndMark(from, to m.Pos) TraceResult {
	result := w.TraceLine(from, to, TraceOptions{
		LoadTiles: true,
	})
	for _, tilePos := range result.Path {
		w.Tiles[tilePos].VisibilityMark = w.VisibilityMark
	}
	return result
}

func (w *World) Update() error {
	// TODO Let all entities move/act. Fetch player position.
	player := w.Entities[w.PlayerID]
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		player.Pos.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		player.Pos.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		player.Pos.Y += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		player.Pos.X += 1
	}

	// Update ScrollPos based on player position and scroll target.
	w.ScrollPos = player.Pos
	log.Printf("player at %v", player.Pos)

	// Unmark all tiles and entities (just bump mark index).
	w.VisibilityMark++
	log.Printf("Updating visibility")

	// Trace from player location to all directions (SweepStep pixels at screen edge).
	// Mark all tiles hit (excl. the tiles that stopped us).
	// TODO Remember trace polygon.
	screen0 := w.ScrollPos.Sub(m.Delta{DX: GameWidth / 2, DY: GameHeight / 2})
	screen1 := screen0.Add(m.Delta{DX: GameWidth - 1, DY: GameHeight - 1})
	for x := screen0.X; x < screen1.X+SweepStep; x += SweepStep {
		w.traceLineAndMark(player.Pos, m.Pos{X: x, Y: screen0.Y})
		w.traceLineAndMark(player.Pos, m.Pos{X: x, Y: screen1.Y})
	}
	for y := screen0.Y; y < screen1.Y+SweepStep; y += SweepStep {
		w.traceLineAndMark(player.Pos, m.Pos{X: screen0.X, Y: y})
		w.traceLineAndMark(player.Pos, m.Pos{X: screen1.X, Y: y})
	}

	// Also mark all neighbors of hit tiles hit (up to ExpandTiles).
	markedTiles := []m.Pos{}
	for tilePos, tile := range w.Tiles {
		if tile.VisibilityMark == w.VisibilityMark {
			markedTiles = append(markedTiles, tilePos)
		}
	}
	expand := m.Delta{DX: ExpandTiles, DY: ExpandTiles}
	for _, pos := range markedTiles {
		w.LoadTilesForTileBox(pos.Sub(expand), pos.Add(expand), pos)
		for y := pos.Y - ExpandTiles; y <= pos.Y+ExpandTiles; y++ {
			for x := pos.X - ExpandTiles; x <= pos.X+ExpandTiles; x++ {
				w.Tiles[m.Pos{X: x, Y: y}].VisibilityMark = w.VisibilityMark
			}
		}
	}

	// TODO Mark all entities on marked tiles hit.
	// TODO Delete all unmarked entities.
	// TODO Spawn all entities on marked tiles if not already spawned.
	// TODO Mark all tiles on entities (this is NOT recursive, but entities may require the tiles they are on to be loaded so they can move).
	// (Somewhat tricky as entities may stand on warps; we have to walk from a known tile).

	// Delete all unmarked tiles.
	for pos, tile := range w.Tiles {
		if tile.VisibilityMark != w.VisibilityMark {
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
		/*
			opts.GeoM.SetElement(0, 0, float64(tile.Orientation.Right.DX))
			opts.GeoM.SetElement(0, 1, float64(tile.Orientation.Right.DY))
			opts.GeoM.SetElement(1, 0, float64(tile.Orientation.Down.DX))
			opts.GeoM.SetElement(1, 1, float64(tile.Orientation.Down.DY))
		*/
		opts.GeoM.Translate(float64(screenPos.X), float64(screenPos.Y))
		screen.DrawImage(tile.Image, &opts)

		text.Draw(screen, fmt.Sprintf("%d,%d", tile.LevelPos.X, tile.LevelPos.Y), w.DebugFont, screenPos.X, screenPos.Y+TileSize-1, color.Gray{128})
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
	if _, found := w.Tiles[newPos]; found {
		// Already loaded.
		return newPos
	}
	neighborTile := w.Tiles[p]
	if neighborTile == nil {
		log.Panicf("Trying to load with nonexisting neighbor tile at %v", p)
	}
	t := neighborTile.Transform
	newLevelPos := neighborTile.LevelPos.Add(t.Apply(d))
	newLevelTile, found := w.Level.Tiles[newLevelPos]
	if !found {
		log.Panicf("Trying to load nonexisting tile at %v when moving from %v (%v) by %v (%v)",
			newLevelPos, p, neighborTile.LevelPos, d, t.Apply(d))
	}
	if newLevelTile.Warpzone != nil {
		log.Printf("warping by %v", newLevelTile.Warpzone)
		t = newLevelTile.Warpzone.Transform.Concat(t)
		tile := w.Level.Tiles[newLevelTile.Warpzone.ToTile]
		if tile == nil {
			log.Panicf("nil new tile after warping to %v", newLevelTile.Warpzone)
		}
		newLevelTile = tile
	}
	newTile := newLevelTile.Tile
	newTile.Transform = t
	newTile.Orientation = t.Concat(newTile.Orientation)
	if newTile.Image == nil {
		log.Panicf("nil tile image at %v: %v", newLevelPos, newLevelTile)
	}
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
	// Entities is the set of entities touched, not including what stopped the trace.
	Entities []Entity
	// hitSolidTilePos is the position of the tile that stopped the trace, if any.
	HitSolidTilePos *m.Pos
	// HitSolidTile is the tile that stopped the trace, if any.
	HitSolidTile *Tile
	// HitSolidEntity is the entity that stopped the trace, if any.
	HitSolidEntity *Entity
	// HitFogOfWar is set if the trace ended by hitting an unloaded tile.
	HitFogOfWar bool
}

// TraceLine moves from x,y by dx,dy in pixel coordinates.
func (w *World) TraceLine(from, to m.Pos, o TraceOptions) TraceResult {
	return w.TraceBox(from, m.Delta{DX: 1, DY: 1}, to, o)
}

// tileStop locates the coordinate of the next tile _entry_ position in the given direction.
func tileStop(from, size, to int) int {
	if to >= from {
		// Tile pos p so that:
		//   p > from.
		//   (p + size) mod TileSize = 0.
		return from + TileSize - m.Mod(from+size, TileSize)
	} else {
		// Tile pos p so that:
		//   p < from
		//   (p + 1) mod TileSize = 0.
		return from - 1 - m.Mod(from, TileSize)
	}
}

// TraceBox moves from x,y size sx,sy by dx,dy in pixel coordinates.
func (w *World) TraceBox(from m.Pos, size m.Delta, to m.Pos, o TraceOptions) TraceResult {
	isLine := size == m.Delta{1, 1}
	result := TraceResult{
		EndPos:          to,
		Path:            nil,
		Entities:        nil,
		HitSolidTilePos: nil,
		HitSolidTile:    nil,
		HitSolidEntity:  nil,
		HitFogOfWar:     false,
	}
	if !o.NoTiles {
		prevTile := from.Scale(1, TileSize)
		// Sweep from from towards to, hitting tile boundaries as needed.
		pos := from
	TRACELOOP:
		for {
			// Where is the next stop?
			// TODO can cache the stop we didn't advance.
			xstopx := tileStop(pos.X, size.DX, to.X)
			xstop := m.Pos{xstopx, pos.Y}
			if to.X != pos.X {
				xstop.Y = pos.Y + (to.Y-pos.Y)*(xstop.X-pos.X)/(to.X-pos.X)
			}
			ystopy := tileStop(pos.Y, size.DY, to.Y)
			ystop := m.Pos{pos.X, ystopy}
			if to.Y != pos.Y {
				ystop.X = pos.X + (to.X-pos.X)*(ystop.Y-pos.Y)/(to.Y-pos.Y)
			}
			var stop m.Pos
			// Which stop comes first?
			if xstop.Delta(pos).Norm1() < ystop.Delta(pos).Norm1() {
				stop = xstop
			} else {
				stop = ystop
			}
			// Have we exceeded the goal?
			if stop.Delta(pos).Norm1() > to.Delta(pos).Norm1() {
				break
			}
			// Identify the "front" tile of the trace. This is the tile most likely to stop us.
			front, back := stop, stop
			move := m.Delta{}
			if to.X > pos.X {
				front.X += size.DX - 1
				move.DX = 1
			} else if to.X < pos.X {
				back.X += size.DX - 1
				move.DX = -1
			}
			front = front.Scale(1, TileSize)
			if to.Y > pos.Y {
				front.Y += size.DY - 1
				move.DY = 1
			} else if to.Y < pos.Y {
				back.Y += size.DY - 1
				move.DY = -1
			}
			back = back.Scale(1, TileSize)
			// Loading: walk from previous front to new front.
			if o.LoadTiles && isLine {
				w.LoadTile(prevTile, front.Delta(prevTile))
			}
			prevTile = front
			// Collision: hit the entire front.
			stopend := stop.Add(m.Delta{DX: size.DX - 1, DY: size.DY - 1})
			stopTile := stop.Scale(1, TileSize)
			stopendTile := stopend.Scale(1, TileSize)
			for y := stopTile.Y; y <= stopendTile.Y; y++ {
				for x := stopTile.X; x <= stopendTile.X; x++ {
					if x != front.X && y != front.Y {
						continue
					}
					tilePos := m.Pos{X: x, Y: y}
					tile := w.Tiles[tilePos]
					if tile == nil {
						result.HitFogOfWar = true
						result.EndPos = stop.Sub(move)
						break TRACELOOP
					}
					if tile.Solid {
						result.HitSolidTilePos = &tilePos
						result.HitSolidTile = tile
						result.EndPos = stop.Sub(move)
						break TRACELOOP
					}
					if isLine {
						result.Path = append(result.Path, tilePos)
					}
				}
			}
			pos = stop
		}
	}
	if !o.NoEntities {
		for _, ent := range w.Entities {
			ent = ent
			// Clip trace to ent.
			// If we hit an entity, we must also cut down the Path.
		}
	}
	return result
}

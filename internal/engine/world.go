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
	"encoding/json"
	"errors"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/centerprint"
	"github.com/divVerent/aaaaaa/internal/flag"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
	"github.com/divVerent/aaaaaa/internal/music"
	"github.com/divVerent/aaaaaa/internal/timing"
	"github.com/divVerent/aaaaaa/internal/vfs"
)

var (
	debugInitialOrientation = flag.String("debug_initial_orientation", "ES", "initial orientation of the game (BREAKS THINGS)")
	debugInitialCheckpoint  = flag.String("debug_initial_checkpoint", "", "initial checkpoint")
	debugTileWindowSize     = flag.Bool("debug_check_window_size", false, "if set, we verify that the tile window size is set high enough")
)

// World represents the current game state including its entities.
type World struct {
	renderer renderer

	// tiles are all tiles currently loaded.
	tiles []*level.Tile
	// Entities are all entities currently loaded.
	Entities map[EntityIncarnation]*Entity
	// PlayerIncarnation is the incarnation ID of the player entity.
	Player *Entity
	// Level is the current tilemap (universal covering with warpZones).
	Level *level.Level
	// Frame since last spawn. Used to let the world slowly "fade in".
	FramesSinceSpawn int
	// WarpZoneStates is the set of current overrides of warpzone state.
	// WarpZones can be turned on/off at will, as long as they are offscreen.
	WarpZoneStates map[string]bool

	// Properties that can in theory be regenerated from the above and thus do not
	// need serialization support.

	// scrollPos is the current screen scrolling position.
	scrollPos m.Pos

	// bottomRightTile is the tile at scrollPos.
	bottomRightTile m.Pos
	// visibilityMark is the current mark value to detect visible tiles/objects.
	visibilityMark uint

	// respawned is set if the player got respawned this frame.
	respawned bool

	// traceLineAndMarkPath receives the path from tracing visibility.
	// Exists to reduce memory allocation.
	traceLineAndMarkPath []m.Pos
}

// Initialized returns whether Init() has been called on this World before.
func (w *World) Initialized() bool {
	return w.tiles != nil
}

func (w *World) tileIndex(pos m.Pos) int {
	i := m.Mod(pos.X, tileWindowWidth) + m.Mod(pos.Y, tileWindowHeight)*tileWindowWidth
	if *debugTileWindowSize {
		p := w.tilePos(i)
		if p != pos {
			log.Panicf("accessed out of range tile: got %v, want near scroll tile %v", pos, w.bottomRightTile)
		}
	}
	return i
}

func (w *World) tilePos(i int) m.Pos {
	x := i % tileWindowWidth
	y := i / tileWindowWidth
	return m.Pos{
		X: x + tileWindowWidth*m.Div(w.bottomRightTile.X-x, tileWindowWidth),
		Y: y + tileWindowHeight*m.Div(w.bottomRightTile.Y-y, tileWindowHeight),
	}
}

func (w *World) Tile(pos m.Pos) *level.Tile {
	return w.tiles[w.tileIndex(pos)]
}

func (w *World) setTile(pos m.Pos, t *level.Tile) {
	w.tiles[w.tileIndex(pos)] = t
}

func (w *World) clearTile(pos m.Pos) {
	w.tiles[w.tileIndex(pos)] = nil
}

func (w *World) forEachTile(f func(pos m.Pos, t *level.Tile)) {
	for i, t := range w.tiles {
		if t == nil {
			continue
		}
		f(w.tilePos(i), t)
	}
}

// Init brings a world into a working state.
// Can be called more than once to reset _everything_.
func (w *World) Init() error {
	// Load map.
	lvl, err := level.Load("level")
	if err != nil {
		log.Panicf("Could not load level: %v", err)
	}

	// Allow reiniting if already done.
	for _, ent := range w.Entities {
		if ent == w.Player {
			continue
		}
		ent.Impl.Despawn()
	}

	*w = World{
		tiles:    make([]*level.Tile, tileWindowWidth*tileWindowHeight),
		Entities: map[EntityIncarnation]*Entity{},
		Level:    lvl,
	}
	w.renderer.Init(w)

	// Load tile the player starts on.
	w.setScrollPos(w.Level.Player.LevelPos.Mul(level.TileSize)) // Needed so we can set the tile.
	tile := w.Level.Tile(w.Level.Player.LevelPos).Tile
	tile.Transform = m.Identity()
	w.setTile(w.Level.Player.LevelPos, &tile)

	// Create player entity.
	w.Player, err = w.Spawn(w.Level.Player, w.Level.Player.LevelPos, &tile)
	if err != nil {
		log.Panicf("Could not spawn player: %v", err)
	}

	// Respawn the player at the desired start location (includes other startup).
	w.RespawnPlayer("")

	return nil
}

// Load loads the current savegame.
// If this fails, the world may be in an undefined state; call w.Init() or w.Load() to resume.
func (w *World) Load() error {
	state, err := vfs.ReadState(vfs.SavedGames, "save.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil // Not loading anything due to there being no state to load is OK.
		}
		return err
	}
	save := level.SaveGame{}
	err = json.Unmarshal(state, &save)
	if err != nil {
		return err
	}
	err = w.Level.LoadGame(save)
	if err != nil {
		return err
	}
	cpName := w.Level.Player.PersistentState["last_checkpoint"]
	w.RespawnPlayer(cpName)
	return nil
}

// Save saves the current savegame.
func (w *World) Save() error {
	save, err := w.Level.SaveGame()
	if err != nil {
		return err
	}
	state, err := json.MarshalIndent(save, "", "\t")
	if err != nil {
		return err
	}
	return vfs.WriteState(vfs.SavedGames, "save.json", state)
}

// SpawnPlayer spawns the player in a newly initialized world.
// As a side effect, it unloads all tiles.
// Spawning at checkpoint "" means the initial player location.
func (w *World) RespawnPlayer(checkpointName string) {
	// Load whether we've seen this checkpoint in flipped state.
	flipped := w.Level.Player.PersistentState["checkpoint_seen."+checkpointName] == "FlipX"

	if *debugInitialCheckpoint != "" {
		checkpointName = *debugInitialCheckpoint
	}
	cpSp := w.Level.Checkpoints[checkpointName]
	if cpSp == nil {
		log.Panicf("Could not spawn player: checkpoint %q not found", checkpointName)
	}

	cpTransform := m.Identity()
	cpTransformStr := cpSp.Properties["required_orientation"]
	if cpTransformStr != "" {
		var err error
		cpTransform, err = m.ParseOrientation(cpTransformStr)
		if err != nil {
			log.Panicf("Could not parse checkpoint orientation: %v", err)
		}
	}
	if flipped {
		cpTransform = cpTransform.Concat(m.FlipX())
	}

	w.visibilityMark++

	// First spawn the tile on the checkpoint.
	tile := w.Level.Tile(cpSp.LevelPos).Tile
	var err error
	tile.Transform, err = m.ParseOrientation(*debugInitialOrientation)
	if err != nil {
		log.Panicf("Could not parse initial orientation: %v", err)
	}
	tile.Transform = cpTransform.Concat(tile.Transform)
	tile.Orientation = tile.Transform.Inverse().Concat(tile.Orientation)

	// Build a new world around the CP tile and the player.
	w.visibilityMark = 0
	tile.VisibilityMark = w.visibilityMark
	for _, ent := range w.Entities {
		if ent == w.Player {
			continue
		}
		ent.Impl.Despawn()
	}
	w.Entities = map[EntityIncarnation]*Entity{
		w.Player.Incarnation: w.Player,
	}
	w.tiles = make([]*level.Tile, tileWindowWidth*tileWindowHeight)
	w.setScrollPos(cpSp.LevelPos.Mul(level.TileSize)) // Scroll the tile into view.
	w.setTile(cpSp.LevelPos, &tile)

	// Spawn the CP.
	cp, err := w.Spawn(cpSp, cpSp.LevelPos, &tile)
	if err != nil {
		log.Panicf("Could not spawn checkpoint: %v", err)
	}

	// Move the player to the center of the checkpoint.
	w.Player.Rect.Origin = cp.Rect.Origin.Add(cp.Rect.Size.Div(2)).Sub(w.Player.Rect.Size.Div(2))

	// Load all the stuff that the CP and player need.
	w.LoadTilesForRect(w.Player.Rect, cpSp.LevelPos)
	w.LoadTilesForRect(cp.Rect, cpSp.LevelPos)

	// Move the player down as far as possible.
	trace := w.TraceBox(w.Player.Rect, w.Player.Rect.Origin.Add(m.Delta{DX: 0, DY: 1024}), TraceOptions{
		NoEntities: true,
		LoadTiles:  true,
		IgnoreEnt:  w.Player,
		ForEnt:     w.Player,
	})
	w.Player.Rect.Origin = trace.EndPos

	w.LoadTilesForRect(w.Player.Rect, cpSp.LevelPos)
	w.visibilityMark++

	// Reset all warpzones.
	w.WarpZoneStates = map[string]bool{}

	// Make sure respawning always gets back to this CP.
	w.Level.Player.PersistentState["last_checkpoint"] = checkpointName

	// Notify the player, reset animation state.
	w.Player.Impl.(PlayerEntityImpl).Respawned()

	// Show the fade in.
	w.FramesSinceSpawn = 0

	// Scroll the player in view right away.
	w.setScrollPos(w.Player.Impl.(PlayerEntityImpl).LookPos())

	// Load the configured music.
	music.Switch(cpSp.Properties["music"])

	// Skip updating.
	w.respawned = true
}

func (w *World) traceLineAndMark(from, to m.Pos) TraceResult {
	result := w.TraceLine(from, to, TraceOptions{
		Mode:      HitOpaque,
		LoadTiles: true,
		IgnoreEnt: w.Player,
		ForEnt:    w.Player,
		PathOut:   &w.traceLineAndMarkPath,
	})
	for _, tilePos := range w.traceLineAndMarkPath {
		w.Tile(tilePos).VisibilityMark = w.visibilityMark
	}
	return result
}

func (w *World) updateEntities() {
	w.respawned = false
	for _, ent := range w.Entities {
		ent.Impl.Update()
		if w.respawned {
			// Once respawned, stop further processing to avoid
			// entities to interact with the respawned player.
			return
		}
	}
}

// updateScrollPos updates the current scroll position.
func (w *World) updateScrollPos(target m.Pos) {
	// Slowly move towards focus point.
	targetDelta := target.Delta(w.scrollPos)
	scrollDelta := targetDelta.MulFloat(scrollPerFrame)
	if scrollDelta.DX == 0 {
		if targetDelta.DX > 0 {
			scrollDelta.DX = +1
		}
		if targetDelta.DX < 0 {
			scrollDelta.DX = -1
		}
	}
	if scrollDelta.DY == 0 {
		if targetDelta.DY > 0 {
			scrollDelta.DY = +1
		}
		if targetDelta.DY < 0 {
			scrollDelta.DY = -1
		}
	}
	target = w.scrollPos.Add(scrollDelta)
	// Ensure player is onscreen.
	if target.X < w.Player.Rect.OppositeCorner().X-GameWidth/2+scrollMinDistance {
		target.X = w.Player.Rect.OppositeCorner().X - GameWidth/2 + scrollMinDistance
	}
	if target.X > w.Player.Rect.Origin.X+GameWidth/2-scrollMinDistance {
		target.X = w.Player.Rect.Origin.X + GameWidth/2 - scrollMinDistance
	}
	if target.Y < w.Player.Rect.OppositeCorner().Y-GameHeight/2+scrollMinDistance {
		target.Y = w.Player.Rect.OppositeCorner().Y - GameHeight/2 + scrollMinDistance
	}
	if target.Y > w.Player.Rect.Origin.Y+GameHeight/2-scrollMinDistance {
		target.Y = w.Player.Rect.Origin.Y + GameHeight/2 - scrollMinDistance
	}
	w.setScrollPos(target)
}

func (w *World) setScrollPos(pos m.Pos) {
	w.scrollPos = pos
	w.bottomRightTile = pos.Div(level.TileSize).Add(m.Delta{DX: tileWindowWidth / 2, DY: tileWindowHeight / 2})
}

// updateVisibility loads all visible tiles and discards all tiles not visible right now.
func (w *World) updateVisibility(eye m.Pos, maxDist int) {
	if maxDist < 1 {
		// Require at least 1 pixel trace distance, or else our polygon can't be correctly expanded.
		maxDist = 1
	}
	defer timing.Group()()

	// Delete all tiles merely marked for expanding.
	// TODO can we preserve but recheck them instead?
	timing.Section("cleanup_expanded")
	prevVisibilityMark := w.visibilityMark - 1
	eyePos := eye.Div(level.TileSize)
	w.forEachTile(func(pos m.Pos, tile *level.Tile) {
		if tile.VisibilityMark != prevVisibilityMark && pos != eyePos {
			w.clearTile(pos)
		}
	})

	// Unmark all tiles and entities (just bump mark index).
	w.visibilityMark++
	visibilityMark := w.visibilityMark

	// Trace from player location to all directions (sweepStep pixels at screen edge).
	// Mark all tiles hit (excl. the tiles that stopped us).
	timing.Section("trace")
	screen0 := w.scrollPos.Sub(m.Delta{DX: GameWidth / 2, DY: GameHeight / 2})
	screen1 := screen0.Add(m.Delta{DX: GameWidth - 1, DY: GameHeight - 1})
	w.renderer.visiblePolygonCenter = eye
	w.renderer.visiblePolygon = w.renderer.visiblePolygon[0:0]
	addTrace := func(rawTarget m.Pos) {
		delta := rawTarget.Delta(w.scrollPos)
		// Diamond shape (cooler?). Could otherwise easily use Length2 here.
		targetDist := delta.Length2()
		if targetDist > maxDist*maxDist {
			f := float64(maxDist) / math.Sqrt(float64(targetDist))
			delta = delta.MulFloat(f)
		}
		target := w.scrollPos.Add(delta)
		trace := w.traceLineAndMark(eye, target)
		w.renderer.visiblePolygon = append(w.renderer.visiblePolygon, trace.EndPos)
	}
	for x := screen0.X; x < screen1.X; x += sweepStep {
		addTrace(m.Pos{X: x, Y: screen0.Y})
	}
	for y := screen0.Y; y < screen1.Y; y += sweepStep {
		addTrace(m.Pos{X: screen1.X, Y: y})
	}
	for x := screen1.X; x > screen0.X; x -= sweepStep {
		addTrace(m.Pos{X: x, Y: screen1.Y})
	}
	for y := screen1.Y; y > screen0.Y; y -= sweepStep {
		addTrace(m.Pos{X: screen0.X, Y: y})
	}
	if *expandUsingVertices {
		if *expandUsingVerticesAccurately {
			w.renderer.expandedVisiblePolygon = expandMinkowski(w.renderer.visiblePolygon, expandSize)
		} else {
			w.renderer.expandedVisiblePolygon = expandSimple(w.renderer.visiblePolygonCenter, w.renderer.visiblePolygon, expandSize)
		}
	} else {
		w.renderer.expandedVisiblePolygon = w.renderer.visiblePolygon
	}

	// Also mark all neighbors of hit tiles hit (up to expandTiles).
	// For multiple expansion, need to do this in steps so initially we only base expansion on visible tiles.
	timing.Section("expand")
	markedTiles := []m.Pos{}
	w.forEachTile(func(tilePos m.Pos, tile *level.Tile) {
		if tile.VisibilityMark == visibilityMark {
			markedTiles = append(markedTiles, tilePos)
		}
	})
	w.visibilityMark++
	expansionMark := w.visibilityMark
	numExpandSteps := (2*expandTiles+1)*(2*expandTiles+1) - 1
	for i := 0; i < numExpandSteps; i++ {
		step := &expandSteps[i]
		for _, pos := range markedTiles {
			from := pos.Add(step.from)
			to := pos.Add(step.to)
			w.LoadTile(from, to.Delta(from))
			if w.Tile(to).VisibilityMark != visibilityMark {
				w.Tile(to).VisibilityMark = expansionMark
			}
		}
	}

	timing.Section("spawn_search")
	w.forEachTile(func(pos m.Pos, tile *level.Tile) {
		if tile.VisibilityMark != expansionMark && tile.VisibilityMark != visibilityMark {
			return
		}
		for _, spawnable := range tile.Spawnables {
			timing.Section("spawn")
			_, err := w.Spawn(spawnable, pos, tile)
			if err != nil {
				log.Panicf("Could not spawn entity %v: %v", spawnable, err)
			}
			timing.Section("spawn_search")
		}
	})

	timing.Section("despawn_search")
	for id, ent := range w.Entities {
		tp0, tp1 := tilesBox(ent.Rect)
		var pos *m.Pos
	DESPAWN_SEARCH:
		for y := tp0.Y; y <= tp1.Y; y++ {
			for x := tp0.X; x <= tp1.X; x++ {
				tp := m.Pos{X: x, Y: y}
				tile := w.Tile(tp)
				if tile == nil {
					continue
				}
				if tile.VisibilityMark == expansionMark || tile.VisibilityMark == visibilityMark {
					pos = &tp
					break DESPAWN_SEARCH
				}
			}
		}
		if pos != nil {
			timing.Section("load_entity_tiles")
			w.LoadTilesForTileBox(tp0, tp1, *pos)
			for y := tp0.Y; y <= tp1.Y; y++ {
				for x := tp0.X; x <= tp1.X; x++ {
					tp := m.Pos{X: x, Y: y}
					if w.Tile(tp).VisibilityMark != visibilityMark {
						w.Tile(tp).VisibilityMark = expansionMark
					}
				}
			}
			timing.Section("despawn_search")
		} else {
			timing.Section("despawn")
			ent.Impl.Despawn()
			delete(w.Entities, id)
			timing.Section("despawn_search")
		}
	}

	// Delete all unmarked tiles.
	timing.Section("cleanup_unmarked")
	w.forEachTile(func(pos m.Pos, tile *level.Tile) {
		if tile.VisibilityMark != expansionMark && tile.VisibilityMark != visibilityMark {
			w.clearTile(pos)
		}
	})
}

func (w *World) Update() error {
	defer timing.Group()()
	w.FramesSinceSpawn++

	// Let everything move.
	timing.Section("entities")
	w.updateEntities()

	// Fetch the player entity.
	playerImpl := w.Player.Impl.(PlayerEntityImpl)

	// Scroll towards the focus point.
	w.updateScrollPos(playerImpl.LookPos())

	// Update visibility and spawn/despawn entities.
	timing.Section("visibility")
	w.updateVisibility(playerImpl.EyePos(), w.FramesSinceSpawn*pixelsPerSpawnFrame)

	// Update centerprints.
	centerprint.Update()

	w.renderer.needPrevImage = true
	return nil
}

// SetWarpZoneState overrides the enabled/disabled state of a warpzone.
// This state resets on respawn.
func (w *World) SetWarpZoneState(name string, state bool) {
	w.WarpZoneStates[name] = state
}

// LoadTile loads the next tile into the current world based on a currently
// known tile and its neighbor. Respects and applies warps.
func (w *World) LoadTile(p m.Pos, d m.Delta) m.Pos {
	newPos := p.Add(d)
	if w.Tile(newPos) != nil {
		// Already loaded.
		return newPos
	}
	neighborTile := w.Tile(p)
	if neighborTile == nil {
		log.Panicf("Trying to load with nonexisting neighbor tile at %v", p)
	}
	t := neighborTile.Transform
	neighborLevelPos := neighborTile.LevelPos
	newLevelPos := neighborLevelPos.Add(t.Apply(d))
	newLevelTile := w.Level.Tile(newLevelPos)
	if newLevelTile == nil {
		// log.Printf("Trying to load nonexisting tile at %v when moving from %v (%v) by %v (%v)", newLevelPos, p, neighborLevelPos, d, t.Apply(d))
		newTile := level.Tile{
			LevelPos:           newLevelPos,
			Transform:          t,
			LoadedFromNeighbor: p,
		}
		w.setTile(newPos, &newTile)
		return newPos
	}
	warped := false
	for _, warp := range newLevelTile.WarpZones {
		// Don't enter warps from behind.
		if warp.PrevTile != neighborLevelPos {
			continue
		}
		// Honor the warpzone toggle state.
		state, overridden := w.WarpZoneStates[warp.Name]
		if !overridden {
			state = warp.InitialState
		}
		if !state {
			continue
		}
		if warped {
			log.Panicf("More than one active warpzone on %v", newLevelTile)
		}
		warped = true
		t = warp.Transform.Concat(t)
		tile := w.Level.Tile(warp.ToTile)
		if tile == nil {
			log.Panicf("Nil new tile after warping to %v", warp)
		}
		newLevelTile = tile
	}
	newTile := newLevelTile.Tile
	newTile.Transform = t
	// Orientation is inverse of the transform, as the transform is for loading
	// new tiles ("which tilemap direction is looking right on the screen") and
	// the orientation is for rendering ("how to rotate the sprite").
	newTile.Orientation = t.Inverse().Concat(newTile.Orientation)
	newTile.LoadedFromNeighbor = p
	w.setTile(newPos, &newTile)
	return newPos
}

// tilesBox returns corner coordinates for all tiles in a given box.
func tilesBox(r m.Rect) (m.Pos, m.Pos) {
	tp0 := r.Origin.Div(level.TileSize)
	tp1 := r.OppositeCorner().Div(level.TileSize)
	return tp0, tp1
}

// LoadTilesForRect loads all tiles in the given box (p, d), assuming tile tp is already loaded.
func (w *World) LoadTilesForRect(r m.Rect, tp m.Pos) {
	tp0, tp1 := tilesBox(r)
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

// TraceLine moves from x,y by dx,dy in pixel coordinates.
func (w *World) TraceLine(from, to m.Pos, o TraceOptions) TraceResult {
	return traceLine(w, from, to, o)
}

// TraceBox moves from x,y size sx,sy by dx,dy in pixel coordinates.
func (w *World) TraceBox(from m.Rect, to m.Pos, o TraceOptions) TraceResult {
	return traceBox(w, from, to, o)
}

func (w *World) Draw(screen *ebiten.Image) {
	w.renderer.Draw(screen)
}

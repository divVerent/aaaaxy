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
	"fmt"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/centerprint"
	"github.com/divVerent/aaaaxy/internal/demo"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/playerstate"
	"github.com/divVerent/aaaaxy/internal/splash"
	"github.com/divVerent/aaaaxy/internal/timing"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

var (
	debugCountTiles                  = flag.Bool("debug_count_tiles", false, "count tiles set/cleared")
	debugNeighborLoadingOptimization = flag.Bool("debug_neighbor_loading_optimization", true, "load tiles faster from the same neighbor tile (maybe incorrect, but faster)")
	debugCheckTileWindowSize         = flag.Bool("debug_check_tile_window_size", false, "if set, we verify that the tile window size is set high enough")
)

// World represents the current game state including its entities.
type World struct {
	renderer renderer

	// tiles are all tiles currently loaded.
	tiles []*level.Tile
	// markedTilesBuffer has the same size as tiles and is used when updating visibility.
	markedTilesBuffer []m.Pos
	// incarnations are all currently existing entity incarnations.
	incarnations map[EntityIncarnation]struct{}
	// entities are all entities currently loaded.
	entities entityList
	// entitiesByZ are all entities, grouped by Z index.
	entitiesByZ []entityList
	// opaqueEntities are all opaque entities currently loaded.
	opaqueEntities entityList
	// Player is the player entity.
	Player *Entity
	// PlayerState is the managed persistent state of the player.
	PlayerState playerstate.PlayerState
	// Level is the current tilemap (universal covering with warpZones).
	Level *level.Level
	// Frame since last spawn. Used to let the world slowly "fade in".
	FramesSinceSpawn int
	// WarpZoneStates is the set of current overrides of warpzone state.
	// WarpZones can be turned on/off at will, as long as they are offscreen.
	WarpZoneStates map[string]bool
	// warpzoneStatesChanged is set if warpzone state changed during this frame.
	warpzoneStatesChanged bool
	// TimerStarted is set on first input after game launch or reset.
	TimerStarted bool
	// TimerStopped is set when game time is paused.
	TimerStopped bool
	// MaxVisiblePixels is the max amount of pixels displayed from player origin.
	MaxVisiblePixels int
	// ForceCredits is set when we want to jump to credits.
	ForceCredits bool
	// GlobalColorM is a color matrix to apply to everything. Reset on every frame.
	GlobalColorM ebiten.ColorM

	// Properties that can in theory be regenerated from the above and thus do not
	// need serialization support.

	// scrollPos is the current screen scrolling position.
	scrollPos m.Pos

	// bottomRightTile is the tile at scrollPos.
	bottomRightTile m.Pos
	// frameVis is the current mark value to detect visible tiles/objects.
	frameVis level.VisibilityFlags
	// visTracing is set while tracing visibility and enables loading conflict detection
	visTracing bool

	// respawned is set if the player got respawned this frame.
	respawned bool

	// traceLineAndMarkPath receives the path from tracing visibility.
	// Exists to reduce memory allocation.
	traceLineAndMarkPath []m.Pos

	// Tile counter.
	tilesSet, tilesCleared int

	// Checkpoint spawn offset.
	prevCpID     level.EntityID
	prevCpOrigin m.Pos

	// Name of the save state.
	saveState int
}

// Initialized returns whether Init() has been called on this World before.
func (w *World) Initialized() bool {
	return w.tiles != nil
}

func (w *World) tileIndex(pos m.Pos) int {
	i := m.Mod(pos.X, tileWindowWidth) + m.Mod(pos.Y, tileWindowHeight)*tileWindowWidth
	if *debugCheckTileWindowSize {
		topLeftTile := w.bottomRightTile.Sub(m.Delta{DX: tileWindowWidth - 1, DY: tileWindowHeight - 1})
		if pos.X < topLeftTile.X || pos.X > w.bottomRightTile.X || pos.Y < topLeftTile.Y || pos.Y > w.bottomRightTile.Y {
			log.Fatalf("accessed out of range tile: requested pos %v, want a pos in window [%v, %v]", pos, topLeftTile, w.bottomRightTile)
		}
		p := w.tilePos(i)
		if p != pos {
			log.Fatalf("accessed out of range tile: got actual pos %v, want requested pos %v, near bottom-right scroll tile %v", p, pos, w.bottomRightTile)
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
	w.tilesSet++
}

func (w *World) clearTile(pos m.Pos) {
	w.tiles[w.tileIndex(pos)] = nil
	w.tilesCleared++
}

func (w *World) forEachTile(f func(pos m.Pos, t *level.Tile)) {
	for i, t := range w.tiles {
		if t == nil {
			continue
		}
		f(w.tilePos(i), t)
	}
}

func (w *World) ForEachEntity(f func(e *Entity)) {
	w.entities.forEach(func(e *Entity) error {
		f(e)
		return nil
	})
}

func (w *World) PreDespawn() {
	w.ForEachEntity(func(e *Entity) {
		if ed, ok := e.Impl.(PreDespawner); ok {
			ed.PreDespawn()
		}
	})
}

func (w *World) clearEntities() {
	w.entities.forEach(func(e *Entity) error {
		if e != w.Player {
			e.Impl.Despawn()
		}
		w.unlink(e)
		return nil
	})
}

var loadLevelCache *level.Level

func loadLevel() (*level.Level, error) {
	if loadLevelCache == nil {
		return nil, fmt.Errorf("trying to load level but nothing has been precached")
	}
	// Verify that the level hasn't changed.
	// If this hits when resetting the game, most likely Clone doesn't properly clone some state.
	err := loadLevelCache.VerifyHash()
	if err != nil {
		return nil, err
	}
	return loadLevelCache.Clone(), nil
}

var levelLoader *level.Loader = level.NewLoader("level")

func Precache(s *splash.State) (splash.Status, error) {
	status, err := s.Enter("loading level", "failed to load level", levelLoader.LoadStepwise)
	if status != splash.Continue {
		return status, err
	}

	status, err = s.Enter("precaching entities", "failed to precache entities", splash.Single(func() error {
		return precacheEntities(levelLoader.Level())
	}))
	if status != splash.Continue {
		return status, err
	}

	loadLevelCache = levelLoader.Level()
	levelLoader = nil // After returning Continue, this will never be called again.
	return splash.Continue, nil
}

func PaletteChanged() error {
	loaded, err := level.NewLoader("level").Load()
	if err != nil {
		return err
	}
	err = precacheEntities(loaded)
	if err != nil {
		return fmt.Errorf("failed to precache entities: %w", err)
	}
	loadLevelCache = loaded
	// Note: need to reinit world to make this actually take effect.
	return nil
}

// Init brings a world into a working state.
// Can be called more than once to reset _everything_.
func (w *World) Init(saveState int) error {
	lvl, err := loadLevel()
	if err != nil {
		return err
	}

	// Allow reiniting if already done.
	w.clearEntities()

	*w = World{
		tiles:          make([]*level.Tile, tileWindowWidth*tileWindowHeight),
		incarnations:   map[EntityIncarnation]struct{}{},
		entities:       makeList(allList),
		opaqueEntities: makeList(opaqueList),
		Level:          lvl,
		PlayerState: playerstate.PlayerState{
			Level: lvl,
		},
		prevCpID:  level.InvalidEntityID,
		saveState: saveState,
	}
	w.PlayerState.Init()
	w.renderer.Init(w)

	// Load tile the player starts on.
	w.setScrollPos(w.Level.Player.LevelPos.Mul(level.TileSize)) // Needed so we can set the tile.
	tile := w.Level.Tile(w.Level.Player.LevelPos).Tile
	tile.Transform = m.Identity()
	w.setTile(w.Level.Player.LevelPos, &tile)

	// Create player entity.
	w.Player, err = w.Spawn(w.Level.Player, w.Level.Player.LevelPos, &tile)
	if err != nil {
		return fmt.Errorf("could not spawn player: %w", err)
	}

	// Respawn the player at the desired start location (includes other startup).
	return w.RespawnPlayer("", true)
}

// Load loads the current savegame.
// If this fails, the world may be in an undefined state; call w.Init() or w.Load() to resume.
func (w *World) Load() error {
	saveName := fmt.Sprintf("save-%d.json", w.saveState)
	err := w.loadUnchecked(saveName)
	if errors.Is(err, os.ErrNotExist) {
		// No save game? Just reinit the world.
		return w.Init(w.saveState)
	}
	if err != nil {
		// Other error? Blow away the save game and reinit the world.
		if demo.Playing() {
			return err // No blowing away while playing demos as playing demos should not write.
		}
		log.Errorf("moving away save game due to error: %v", err)
		err = vfs.MoveAwayState(vfs.SavedGames, saveName)
		if err != nil {
			return err
		}
		return w.Init(w.saveState)
	}
	return nil
}

func (w *World) loadUnchecked(saveName string) error {
	save, intercepted := demo.InterceptPreLoadGame()
	if !intercepted {
		state, err := vfs.ReadState(vfs.SavedGames, saveName)
		if err != nil {
			return err
		}
		// Normal loading.
		save = &level.SaveGame{}
		err = json.Unmarshal(state, save)
		if err != nil {
			return err
		}
	}

	// Make sure that demo playback will also go back to this save.
	demo.InterceptPostLoadGame(save)

	if save == nil {
		// Nothing to load? Send an error upwards; this will reinit the world.
		return os.ErrNotExist
	}

	err := w.Level.LoadGame(save)
	if err != nil {
		return err
	}
	w.PlayerState.Init()
	return w.RespawnPlayer(w.PlayerState.LastCheckpoint(), true)
}

// Save saves the current savegame.
func (w *World) Save() error {
	save, err := w.Level.SaveGame()
	if err != nil {
		return err
	}
	if demo.InterceptSaveGame(save) {
		return nil
	}
	state, err := json.MarshalIndent(save, "", "\t")
	if err != nil {
		return err
	}
	if is, cheats := flag.Cheating(); is {
		return fmt.Errorf("not saving, as cheats are enabled: %s", cheats)
	}
	return vfs.WriteState(vfs.SavedGames, fmt.Sprintf("save-%d.json", w.saveState), state)
}

// SpawnPlayer spawns the player in a newly initialized world.
// As a side effect, it unloads all tiles.
// Spawning at checkpoint "" means the initial player location.
func (w *World) RespawnPlayer(checkpointName string, newGameSection bool) error {
	// Load whether we've seen this checkpoint in flipped state.
	flipped := w.PlayerState.CheckpointSeen(checkpointName) == playerstate.SeenFlipped

	cpSp := w.Level.Checkpoints[checkpointName]
	if cpSp == nil {
		return fmt.Errorf("could not spawn player: checkpoint %q not found", checkpointName)
	}

	cpTransform := m.Identity()
	cpTransformStr := cpSp.Properties["required_orientation"]
	if cpTransformStr != "" {
		cpTransforms, err := m.ParseOrientations(cpTransformStr)
		if err != nil {
			return fmt.Errorf("could not parse checkpoint orientation: %w", err)
		}
		cpTransform = cpTransforms[0]
	}
	if flipped {
		cpTransform = cpTransform.Concat(m.FlipX())
	}

	// First spawn the tile on the checkpoint.
	tile := w.Level.Tile(cpSp.LevelPos).Tile
	tile.Transform = cpTransform
	tile.Orientation = tile.Transform.Inverse().Concat(tile.Orientation)
	tile.ResolveImage()

	// Build a new world around the CP tile and the player.
	w.frameVis = 0
	tile.VisibilityFlags = w.frameVis
	w.clearEntities()
	w.link(w.Player)
	w.tiles = make([]*level.Tile, tileWindowWidth*tileWindowHeight)
	w.markedTilesBuffer = make([]m.Pos, tileWindowWidth*tileWindowHeight)
	w.setScrollPos(cpSp.LevelPos.Mul(level.TileSize)) // Scroll the tile into view.
	w.setTile(cpSp.LevelPos, &tile)

	// Spawn the CP.
	var cp *Entity
	if cpSp == w.Level.Player {
		cp = w.Player
	} else {
		var err error
		cp, err = w.Spawn(cpSp, cpSp.LevelPos, &tile)
		if err != nil {
			return fmt.Errorf("could not spawn checkpoint: %w", err)
		}
	}

	if newGameSection {
		// Clear all centerprints.
		// But only when coming from menu, not when respawning/teleporting in game.
		centerprint.Reset()
	}

	// Reset the ending stuff.
	w.TimerStopped = false
	w.MaxVisiblePixels = math.MaxInt32
	w.ForceCredits = false

	// Reset all warpzones.
	w.WarpZoneStates = map[string]bool{}

	// Move the player to the center of the checkpoint.
	w.Player.Rect.Origin = cp.Rect.Origin.Add(cp.Rect.Size.Div(2)).Sub(w.Player.Rect.Size.Div(2))

	// Load all the stuff that the player needs.
	w.LoadTilesForRect(w.Player.Rect, cpSp.LevelPos)

	// Show the fade in.
	w.FramesSinceSpawn = 0

	// Move the player down as far as possible.
	if cpSp.Properties["downtrace_on_spawn"] != "false" {
		var dir m.Delta
		if onGroundVecStr := cpSp.Properties["vvvvvv_gravity_direction"]; onGroundVecStr != "" {
			_, err := fmt.Sscanf(onGroundVecStr, "%d %d", &dir.DX, &dir.DY)
			if err != nil {
				return fmt.Errorf("invalid vvvvvv_gravity_direction: %w", err)
			}
		}
		if dir.IsZero() {
			dir = m.Delta{DX: 0, DY: 1}
		}
		trace := w.TraceBox(w.Player.Rect, w.Player.Rect.Origin.Add(dir.Mul(1024)), TraceOptions{
			Contents:   level.PlayerSolidContents,
			NoEntities: true,
			LoadTiles:  true,
			ForEnt:     w.Player,
		})
		w.Player.Rect.Origin = trace.EndPos
	}

	// Note that TraceBox must have loaded all tiles the player needs.
	// w.LoadTilesForRect(w.Player.Rect, cpSp.LevelPos)
	w.frameVis ^= level.FrameVis

	// Make sure respawning always gets back to this CP.
	w.PlayerState.RecordCheckpoint(checkpointName, flipped)

	// Notify the player, reset animation state.
	w.Player.Impl.(PlayerEntityImpl).Respawned()

	// Scroll the player in view right away.
	w.setScrollPos(w.Player.Impl.(PlayerEntityImpl).LookPos())

	// Adjust previous scroll position by how much the CP "moved".
	// That way, respawning right after touching a CP will retain CP-near screen content.
	// Mainly useful for VVVVVV mode.
	if w.prevCpID == cp.Incarnation.ID {
		cpDelta := cp.Rect.Origin.Delta(w.prevCpOrigin)
		w.renderer.prevScrollPos = w.renderer.prevScrollPos.Add(cpDelta)
		w.prevCpOrigin = cp.Rect.Origin
	}

	// Initialize whatever the checkpoint wants to do.
	w.TouchEvent(cp, []*Entity{w.Player})

	// Skip updating.
	w.respawned = true
	return nil
}

// TouchEvent notifies both entities that they touched the other.
// nil entities will be skipped as LHS of touches, and if bs is empty, it's
// treated as touching a nil (which inidcates touching world).
func (w *World) TouchEvent(a *Entity, bs []*Entity) {
	for _, b := range bs {
		if a != nil {
			a.Impl.Touch(b)
		}
		if b != nil {
			b.Impl.Touch(a)
		}
	}
	if len(bs) == 0 {
		if a != nil {
			a.Impl.Touch(nil)
		}
	}
}

func (w *World) PlayerTouchedCheckpoint(cp *Entity) {
	w.prevCpID = cp.Incarnation.ID
	w.prevCpOrigin = cp.Rect.Origin
}

func (w *World) traceLineAndMark(from, to m.Pos, pathStore *[]m.Pos) TraceResult {
	result := w.TraceLine(from, to, TraceOptions{
		Contents:  level.OpaqueContents,
		LoadTiles: true,
		ForEnt:    w.Player,
		PathOut:   pathStore,
	})
	for _, tilePos := range *pathStore {
		w.Tile(tilePos).VisibilityFlags = w.frameVis | level.TracedVis
	}
	return result
}

func (w *World) updateEntities() {
	// Entities may update these.
	w.warpzoneStatesChanged = false
	w.respawned = false
	w.GlobalColorM.Reset()

	w.entities.forEach(func(ent *Entity) error {
		ent.Impl.Update()
		if w.respawned {
			// Once respawned, stop further processing to avoid
			// entities to interact with the respawned player.
			return breakError
		}
		return nil
	})

	// Clean up newly spawned or despawned stuff.
	w.entities.compact()
	for i := range w.entitiesByZ {
		w.entitiesByZ[i].compact()
	}
	w.opaqueEntities.compact()
}

// updateScrollPos updates the current scroll position.
func (w *World) updateScrollPos(target m.Pos) {
	// Slowly move towards focus point.
	targetDelta := target.Delta(w.scrollPos)
	scrollDelta := targetDelta.MulFixed(m.NewFixedFloat64(scrollPerFrame))
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

	// Reset visibility!
	w.frameVis ^= level.FrameVis

	// Always mark the eye tile.
	timing.Section("eye")
	eyePos := eye.Div(level.TileSize)
	w.Tile(eyePos).VisibilityFlags = w.frameVis | level.TracedVis

	// Trace from player location to all directions (sweepStep pixels at screen edge).
	// Mark all tiles hit (excl. the tiles that stopped us).
	timing.Section("trace")
	screen0 := w.scrollPos.Sub(m.Delta{DX: GameWidth / 2, DY: GameHeight / 2})
	screen1 := screen0.Add(m.Delta{DX: GameWidth - 1, DY: GameHeight - 1})
	w.renderer.visiblePolygonCenter = eye
	// Pick xLen so that it is the SMALLEST xLen so that screen0+sweepStep*xLen>=screen1.
	xLen := (screen1.X - screen0.X + sweepStep - 1) / sweepStep
	yLen := (screen1.Y - screen0.Y + sweepStep - 1) / sweepStep
	wantLen := 2*xLen + 2*yLen
	if len(w.renderer.visiblePolygon) != wantLen {
		if w.renderer.visiblePolygon != nil {
			log.Infof("visible polygon size changed - unexpected but harmless; optimize for resizing and remove this check then, I guess: got %v, want %v", len(w.renderer.visiblePolygon), wantLen)
		}
		w.renderer.visiblePolygon = make([]m.Pos, wantLen)
	}
	addTrace := func(rawTarget m.Pos, index int) {
		delta := rawTarget.Delta(w.scrollPos).WithMaxLengthFixed(m.NewFixed(maxDist))
		target := w.scrollPos.Add(delta)
		trace := w.traceLineAndMark(eye, target, &w.traceLineAndMarkPath)
		w.renderer.visiblePolygon[index] = trace.EndPos
	}
	w.visTracing = true
	for i := 0; i < xLen; i++ {
		addTrace(m.Pos{X: screen0.X + sweepStep*i, Y: screen0.Y}, i)
		addTrace(m.Pos{X: screen1.X - sweepStep*i, Y: screen1.Y}, xLen+yLen+i)
	}
	for i := 0; i < yLen; i++ {
		addTrace(m.Pos{X: screen1.X, Y: screen0.Y + sweepStep*i}, xLen+i)
		addTrace(m.Pos{X: screen0.X, Y: screen1.Y - sweepStep*i}, 2*xLen+yLen+i)
	}
	w.visTracing = false
	if *expandUsingVertices {
		if *expandUsingVerticesAccurately {
			w.renderer.expandedVisiblePolygon = expandMinkowski(w.renderer.visiblePolygon, expandSize)
		} else {
			w.renderer.expandedVisiblePolygon = expandSimple(w.renderer.visiblePolygonCenter, w.renderer.visiblePolygon, expandSize)
		}
	} else {
		w.renderer.expandedVisiblePolygon = w.renderer.visiblePolygon
	}
	// BUG: the above also loads tiles (but doesn't mark) if their path was blocked by an entity.
	// Workaround: mark them as if they were previous frame's tiles, so they're not a basis for loading and get cleared at the end if needed.
	timing.Section("untrace_workaround")
	w.forEachTile(func(pos m.Pos, tile *level.Tile) {
		if tile.VisibilityFlags == w.frameVis {
			tile.VisibilityFlags ^= level.FrameVis
		}
	})

	// Also mark all neighbors of hit tiles hit (up to expandTiles).
	// For multiple expansion, need to do this in steps so initially we only base expansion on visible tiles.
	timing.Section("expand")
	markedTiles := w.markedTilesBuffer[:0]
	justTraced := w.frameVis | level.TracedVis
	w.forEachTile(func(tilePos m.Pos, tile *level.Tile) {
		if tile.VisibilityFlags == justTraced {
			markedTiles = append(markedTiles, tilePos)
		}
	})
	numExpandSteps := (2*expandTiles+1)*(2*expandTiles+1) - 1
	for i := 0; i < numExpandSteps; i++ {
		step := &expandSteps[i]
		for _, pos := range markedTiles {
			from := pos.Add(step.from)
			to := pos.Add(step.to)
			// It's not OK to load from an opaque tile, as that may sidestep warpzones.
			// So, pick an alternative in that case, or skip expanding entirely.
			if w.Tile(from).Contents.Opaque() {
				from = pos.Add(step.from2)
				if w.Tile(from).Contents.Opaque() {
					from = pos.Add(step.from3)
					if w.Tile(from).Contents.Opaque() {
						continue
					}
				}
			}
			w.LoadTile(from, to, to.Delta(from))
		}
	}

	timing.Section("spawn_search")
	w.forEachTile(func(pos m.Pos, tile *level.Tile) {
		if tile.VisibilityFlags&level.FrameVis != w.frameVis {
			return
		}
		for _, spawnable := range tile.Spawnables {
			_, err := w.Spawn(spawnable, pos, tile)
			if err != nil {
				log.Errorf("could not spawn entity %v: %v", spawnable, err)
			}
		}
	})

	timing.Section("despawn_search")
	w.entities.forEach(func(ent *Entity) error {
		tp0, tp1 := tilesBox(ent.Rect.Grow(ent.SpawnTilesGrowth))
		var pos m.Pos
		havePos := false
	DESPAWN_SEARCH:
		for y := tp0.Y; y <= tp1.Y; y++ {
			for x := tp0.X; x <= tp1.X; x++ {
				tp := m.Pos{X: x, Y: y}
				tile := w.Tile(tp)
				if tile == nil {
					continue
				}
				if tile.VisibilityFlags&level.FrameVis == w.frameVis {
					pos = tp
					havePos = true
					break DESPAWN_SEARCH
				}
			}
		}
		if havePos {
			if ent.RequireTiles {
				w.LoadTilesForTileBox(tp0, tp1, pos)
				for y := tp0.Y; y <= tp1.Y; y++ {
					for x := tp0.X; x <= tp1.X; x++ {
						tp := m.Pos{X: x, Y: y}
						tile := w.Tile(tp)
						if tile == nil {
							continue
						}
						if tile.VisibilityFlags&level.FrameVis != w.frameVis {
							tile.VisibilityFlags = w.frameVis
						}
					}
				}
			}
		} else {
			ent.Impl.Despawn()
			w.unlink(ent)
		}
		return nil
	})

	// Delete all unmarked tiles.
	timing.Section("cleanup_unmarked")
	w.forEachTile(func(pos m.Pos, tile *level.Tile) {
		if tile.VisibilityFlags&level.FrameVis != w.frameVis {
			w.clearTile(pos)
		}
	})
}

func (w *World) AssumeChanged() {
	w.renderer.worldChanged = true
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
	pixels := w.FramesSinceSpawn * pixelsPerSpawnFrame
	if pixels > w.MaxVisiblePixels {
		pixels = w.MaxVisiblePixels
	}
	w.updateVisibility(playerImpl.EyePos(), pixels)

	// Update centerprints.
	centerprint.Update()

	if *debugCountTiles {
		log.Infof("%d tiles set, %d tiles cleared", w.tilesSet, w.tilesCleared)
	}
	w.tilesSet, w.tilesCleared = 0, 0

	w.AssumeChanged()
	return nil
}

// SetWarpZoneState overrides the enabled/disabled state of a warpzone.
// This state resets on respawn.
func (w *World) SetWarpZoneState(name string, state bool) {
	w.WarpZoneStates[name] = state
	w.warpzoneStatesChanged = true
}

// LoadTile loads the next tile into the current world based on a currently
// known tile and its neighbor. Respects and applies warps.
func (w *World) LoadTile(p, newPos m.Pos, d m.Delta) *level.Tile {
	tile := w.Tile(newPos)
	if tile != nil {
		if tile.VisibilityFlags&level.FrameVis == w.frameVis {
			// Already loaded this frame.
			return tile
		}
		// From now on we know it doesn't have the same FrameVis.
		if *debugNeighborLoadingOptimization && !w.warpzoneStatesChanged && tile.LoadedFromNeighbor == p {
			// Loading from same neighbor as before is OK.
			// Note: this is INCORRECT if during the last frame, a warpzone changed status.
			tile.VisibilityFlags = w.frameVis
			return tile
		}
	}
	neighborTile := w.Tile(p)
	if neighborTile == nil {
		log.Errorf("trying to load with nonexisting neighbor tile at %v", p)
		return nil // Can't load.
	}
	if neighborTile.Contents.Opaque() {
		log.Errorf("trying to load from an opaque tile at %v", p)
		return nil // Can't load.
	}
	t := neighborTile.Transform
	neighborLevelPos := neighborTile.LevelPos
	newLevelPos := neighborLevelPos.Add(t.Apply(d))
	newLevelTile := w.Level.Tile(newLevelPos)
	if newLevelTile == nil {
		log.Debugf("trying to load nonexisting tile at %v when moving from %v (%v) by %v (%v)", newLevelPos, p, neighborLevelPos, d, t.Apply(d))
		newTile := &level.Tile{
			LevelPos:           newLevelPos,
			Transform:          t,
			LoadedFromNeighbor: p,
			VisibilityFlags:    w.frameVis,
		}
		w.setTile(newPos, newTile)
		return newTile
	}
	warped := false
	for _, warp := range newLevelTile.WarpZones {
		// Don't enter warps from behind.
		if warp.PrevTile != neighborLevelPos {
			continue
		}
		// Honor the warpzone toggle state.
		if warp.Switchable {
			if w.WarpZoneStates[warp.Name] == warp.Invert {
				continue
			}
		}
		if warped {
			log.Errorf("more than one active warpzone on %v", newLevelTile)
			return nil // Can't load.
		}
		warped = true
		t = warp.Transform.Concat(t)
		tile := w.Level.Tile(warp.ToTile)
		if tile == nil {
			log.Errorf("nil new tile after warping to %v", warp)
			return nil // Can't load.
		}
		newLevelTile = tile
	}
	if tile != nil {
		if tile.LevelPos == newLevelTile.Tile.LevelPos && tile.Transform == t {
			// Same tile as we had before. Can skip the copying.
			tile.LoadedFromNeighbor = p
			tile.VisibilityFlags = w.frameVis
			return tile
		}
	}
	newTile := newLevelTile.Tile
	newTile.Transform = t
	// Orientation is inverse of the transform, as the transform is for loading
	// new tiles ("which tilemap direction is looking right on the screen") and
	// the orientation is for rendering ("how to rotate the sprite").
	newTile.Orientation = t.Inverse().Concat(newTile.Orientation)
	newTile.ResolveImage()
	newTile.LoadedFromNeighbor = p
	newTile.VisibilityFlags = w.frameVis
	w.setTile(newPos, &newTile)
	return &newTile
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
		pos := m.Pos{X: tp.X, Y: y}
		w.LoadTile(pos, pos.Add(m.North()), m.North())
	}
	for y := tp.Y; y < tp1.Y; y++ {
		pos := m.Pos{X: tp.X, Y: y}
		w.LoadTile(pos, pos.Add(m.South()), m.South())
	}
	for y := tp0.Y; y <= tp1.Y; y++ {
		for x := tp.X; x > tp0.X; x-- {
			pos := m.Pos{X: x, Y: y}
			w.LoadTile(pos, pos.Add(m.West()), m.West())
		}
		for x := tp.X; x < tp1.X; x++ {
			pos := m.Pos{X: x, Y: y}
			w.LoadTile(pos, pos.Add(m.East()), m.East())
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

func (w *World) Draw(screen *ebiten.Image, blurFactor float64) {
	w.renderer.Draw(screen, blurFactor)
}

func encodeZ(z int) int {
	if z < 0 {
		return -1 - 2*z
	} else {
		return 2 * z
	}
}

func zBounds(l int) (min, max int) {
	last := l - 1
	max = last / 2
	min = -max - last%2
	return
}

func (w *World) unlink(e *Entity) {
	z := encodeZ(e.zIndex)
	w.entitiesByZ[z].remove(e)
	if e.contents.Opaque() {
		w.opaqueEntities.remove(e)
	}
	w.entities.remove(e)
	if e.Incarnation.IsValid() {
		delete(w.incarnations, e.Incarnation)
	}
}

func (w *World) link(e *Entity) {
	if e.Incarnation.IsValid() {
		w.incarnations[e.Incarnation] = struct{}{}
	}
	w.entities.insert(e)
	if e.contents.Opaque() {
		w.opaqueEntities.insert(e)
	}
	z := encodeZ(e.zIndex)
	for len(w.entitiesByZ) <= z {
		w.entitiesByZ = append(w.entitiesByZ, makeList(zList))
	}
	w.entitiesByZ[z].insert(e)
}

func (w *World) EntityIsAlive(incarnation EntityIncarnation) bool {
	_, found := w.incarnations[incarnation]
	return found
}

func (w *World) FindName(name string) []*Entity {
	var out []*Entity
	w.entities.forEach(func(ent *Entity) error {
		if ent.name == name {
			out = append(out, ent)
		}
		return nil
	})
	return out
}

func (w *World) FindContents(c level.Contents) []*Entity {
	if c == level.OpaqueContents {
		return w.opaqueEntities.items
	}
	var out []*Entity
	w.entities.forEach(func(ent *Entity) error {
		if ent.contents&c != 0 {
			out = append(out, ent)
		}
		return nil
	})
	return out
}

func (w *World) ScrollPos() m.Pos {
	return w.scrollPos
}

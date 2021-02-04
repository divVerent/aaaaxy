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
	"flag"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/divVerent/aaaaaa/internal/centerprint"
	"github.com/divVerent/aaaaaa/internal/font"
	m "github.com/divVerent/aaaaaa/internal/math"
	"github.com/divVerent/aaaaaa/internal/music"
	"github.com/divVerent/aaaaaa/internal/timing"
	"github.com/divVerent/aaaaaa/internal/vfs"
)

var (
	debugShowNeighbors      = flag.Bool("debug_show_neighbors", false, "show the neighbors tiles got loaded from")
	debugShowCoords         = flag.Bool("debug_show_coords", false, "show the level coordinates of each tile")
	debugShowOrientations   = flag.Bool("debug_show_orientations", false, "show the orientation of each tile")
	debugShowTransforms     = flag.Bool("debug_show_transforms", false, "show the transform of each tile")
	debugShowBboxes         = flag.Bool("debug_show_bboxes", false, "show the bounding boxes of all entities")
	debugInitialOrientation = flag.String("debug_initial_orientation", "ES", "initial orientation of the game (BREAKS THINGS)")
	debugInitialCheckpoint  = flag.String("debug_initial_checkpoint", "", "initial checkpoint")
	drawOutside             = flag.Bool("draw_outside", true, "draw outside of the visible area; requires draw_visibility_mask")
	drawVisibilityMask      = flag.Bool("draw_visibility_mask", true, "draw visibility mask (if disabled, all loaded tiles are shown")
	expandUsingVertices     = flag.Bool("expand_using_vertices", false, "expand using polygon math (just approximate, simplifies rendering)")
	debugShowTrace          = flag.String("debug_show_trace", "", "if set, the screen coordinates to trace towards and show trace info")
)

// World represents the current game state including its entities.
type World struct {
	// Tiles are all tiles currently loaded.
	Tiles map[m.Pos]*Tile
	// Entities are all entities currently loaded.
	Entities map[EntityIncarnation]*Entity
	// PlayerIncarnation is the incarnation ID of the player entity.
	Player *Entity
	// Level is the current tilemap (universal covering with warpZones).
	Level *Level
	// Frame since last spawn. Used to let the world slowly "fade in".
	FramesSinceSpawn int
	// WarpZoneStates is the set of current overrides of warpzone state.
	// WarpZones can be turned on/off at will, as long as they are offscreen.
	WarpZoneStates map[string]bool

	// Properties that can in theory be regenerated from the above and thus do not
	// need serialization support.

	// scrollPos is the current screen scrolling position.
	scrollPos m.Pos
	// visibilityMark is the current mark value to detect visible tiles/objects.
	visibilityMark uint
	// visiblePolygonCenter is the current eye position.
	visiblePolygonCenter m.Pos
	// visiblePolygon is the currently visible polygon.
	visiblePolygon []m.Pos
	// needPrevImage is set whenever the last call was Update.
	needPrevImage bool

	// Images retained across frames.

	// whiteImage is a single white pixel.
	whiteImage *ebiten.Image
	// prevImage is the previous screen content.
	prevImage *ebiten.Image
	// offScreenBuffer is the previous screen content after masking.
	offScreenBuffer *ebiten.Image
	// prevScrollPos is previous frame's scroll pos.
	prevScrollPos m.Pos
	// The shader for drawing visibility masks.
	visibilityMaskShader *ebiten.Shader

	// Temp storage within frames.

	// blurImage is an offscreen image used for blurring.
	blurImage *ebiten.Image
	// visibilityMaskImage is an offscreen image used for masking the visible area.
	visibilityMaskImage *ebiten.Image
	// respawned is set if the player got respawned this frame.
	respawned bool
}

// Initialized returns whether Init() has been called on this World before.
func (w *World) Initialized() bool {
	return w.Tiles != nil
}

// Init brings a world into a working state.
// Can be called more than once to reset _everything_.
func (w *World) Init() error {
	// Load map.
	level, err := LoadLevel("level")
	if err != nil {
		log.Panicf("Could not load level: %v", err)
	}

	*w = World{
		Tiles:               map[m.Pos]*Tile{},
		Entities:            map[EntityIncarnation]*Entity{},
		Level:               level,
		whiteImage:          ebiten.NewImage(1, 1),
		blurImage:           ebiten.NewImage(GameWidth, GameHeight),
		prevImage:           ebiten.NewImage(GameWidth, GameHeight),
		offScreenBuffer:     ebiten.NewImage(GameWidth, GameHeight),
		visibilityMaskImage: ebiten.NewImage(GameWidth, GameHeight),
	}
	w.whiteImage.Fill(color.Gray{255})
	w.prevImage.Fill(color.Gray{0})
	w.offScreenBuffer.Fill(color.Gray{0})

	if *debugUseShaders {
		w.visibilityMaskShader, err = loadShader("visibility_mask.go")
		if err != nil {
			log.Printf("could not load visibility mask shader: %v", err)
		}
	}

	// Load tile the player starts on.
	tile := w.Level.Tiles[w.Level.Player.LevelPos].Tile
	tile.Transform = m.Identity()
	w.Tiles[w.Level.Player.LevelPos] = &tile

	// Create player entity.
	w.Player, err = w.Level.Player.Spawn(w, w.Level.Player.LevelPos, &tile)
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
	save := SaveGame{}
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
	tile := w.Level.Tiles[cpSp.LevelPos].Tile
	var err error
	tile.Transform, err = m.ParseOrientation(*debugInitialOrientation)
	if err != nil {
		log.Panicf("Could not parse initial orientation: %v", err)
	}
	tile.Transform = cpTransform.Concat(tile.Transform)
	tile.Orientation = tile.Transform.Inverse().Concat(tile.Orientation)

	// Build a new world around the CP tile and the player.
	w.visibilityMark = 0
	tile.visibilityMark = w.visibilityMark
	for _, ent := range w.Entities {
		if ent == w.Player {
			continue
		}
		ent.Impl.Despawn()
	}
	w.Entities = map[EntityIncarnation]*Entity{
		w.Player.Incarnation: w.Player,
	}
	w.Tiles = map[m.Pos]*Tile{
		cpSp.LevelPos: &tile,
	}

	// Spawn the CP.
	cp, err := cpSp.Spawn(w, cpSp.LevelPos, &tile)
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
	w.scrollPos = w.Player.Impl.(PlayerEntityImpl).LookPos()

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
	})
	for _, tilePos := range result.Path {
		w.Tiles[tilePos].visibilityMark = w.visibilityMark
	}
	return result
}

func expandPolygon(center m.Pos, polygon []m.Pos, shift int) {
	orig := append([]m.Pos{}, polygon...)
	for i, v1 := range orig {
		// Rather approximate polygon expanding: just push each vertex shift away from the center.
		// Unlike correct polygon expansion perpendicular to sides,
		// this way ensures that we never include more than distance shift from the polugon.
		// However this is just approximate and causes artifacts when close to a wall.
		d := v1.Delta(center)
		l := d.Length()
		if l <= 0 {
			continue
		}
		f := float64(shift) / l
		polygon[i] = v1.Add(d.MulFloat(f))
	}
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
	w.scrollPos = target
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
	eyePos := eye.Div(TileSize)
	for pos, tile := range w.Tiles {
		if tile.visibilityMark != prevVisibilityMark && pos != eyePos {
			delete(w.Tiles, pos)
		}
	}

	// Unmark all tiles and entities (just bump mark index).
	w.visibilityMark++
	visibilityMark := w.visibilityMark

	// Trace from player location to all directions (sweepStep pixels at screen edge).
	// Mark all tiles hit (excl. the tiles that stopped us).
	timing.Section("trace")
	screen0 := w.scrollPos.Sub(m.Delta{DX: GameWidth / 2, DY: GameHeight / 2})
	screen1 := screen0.Add(m.Delta{DX: GameWidth - 1, DY: GameHeight - 1})
	w.visiblePolygonCenter = eye
	w.visiblePolygon = w.visiblePolygon[0:0]
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
		w.visiblePolygon = append(w.visiblePolygon, trace.EndPos)
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
		expandPolygon(w.visiblePolygonCenter, w.visiblePolygon, expandSize)
	}

	// Also mark all neighbors of hit tiles hit (up to expandTiles).
	// For multiple expansion, need to do this in steps so initially we only base expansion on visible tiles.
	timing.Section("expand")
	markedTiles := []m.Pos{}
	for tilePos, tile := range w.Tiles {
		if tile.visibilityMark == visibilityMark {
			markedTiles = append(markedTiles, tilePos)
		}
	}
	w.visibilityMark++
	expansionMark := w.visibilityMark
	numExpandSteps := (2*expandTiles+1)*(2*expandTiles+1) - 1
	for i := 0; i < numExpandSteps; i++ {
		step := &expandSteps[i]
		for _, pos := range markedTiles {
			from := pos.Add(step.from)
			to := pos.Add(step.to)
			w.LoadTile(from, to.Delta(from))
			if w.Tiles[to].visibilityMark != visibilityMark {
				w.Tiles[to].visibilityMark = expansionMark
			}
		}
	}

	timing.Section("spawn_search")
	for pos, tile := range w.Tiles {
		if tile.visibilityMark != expansionMark && tile.visibilityMark != visibilityMark {
			continue
		}
		for _, spawnable := range tile.Spawnables {
			timing.Section("spawn")
			_, err := spawnable.Spawn(w, pos, tile)
			if err != nil {
				log.Panicf("Could not spawn entity %v: %v", spawnable, err)
			}
			timing.Section("spawn_search")
		}
	}

	timing.Section("despawn_search")
	for id, ent := range w.Entities {
		tp0, tp1 := tilesBox(ent.Rect)
		var pos *m.Pos
	DESPAWN_SEARCH:
		for y := tp0.Y; y <= tp1.Y; y++ {
			for x := tp0.X; x <= tp1.X; x++ {
				tp := m.Pos{X: x, Y: y}
				tile := w.Tiles[tp]
				if tile == nil {
					continue
				}
				if tile.visibilityMark == expansionMark || tile.visibilityMark == visibilityMark {
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
					if w.Tiles[tp].visibilityMark != visibilityMark {
						w.Tiles[tp].visibilityMark = expansionMark
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
	for pos, tile := range w.Tiles {
		if tile.visibilityMark != expansionMark && tile.visibilityMark != visibilityMark {
			delete(w.Tiles, pos)
		}
	}
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

	w.needPrevImage = true
	return nil
}

func setGeoM(geoM *ebiten.GeoM, pos m.Pos, resize bool, entSize, imgSize m.Delta, orientation m.Orientation) {
	// Note: the logic here is rather inefficient but easy to verify.
	// If this turns out to be performance relevant, let's optimize.

	// Step 1: compute the raw corners at source and destination.
	rectI := m.Rect{Origin: m.Pos{}, Size: imgSize}
	var rectR m.Rect
	var scaledImgSize m.Delta
	if resize {
		scaledImgSize = entSize
		if orientation.Right.DX == 0 {
			scaledImgSize.DX, scaledImgSize.DY = scaledImgSize.DY, scaledImgSize.DX
		}
		rectR = m.Rect{Origin: pos, Size: entSize}
	} else {
		scaledImgSize = imgSize
		flippedSize := imgSize
		if orientation.Right.DX == 0 {
			flippedSize.DX, flippedSize.DY = flippedSize.DY, flippedSize.DX
		}
		rectR = m.Rect{Origin: pos, Size: flippedSize}
	}

	// Step 2: actually match up image corners with destination.
	rectIR := orientation.ApplyToRect2(m.Pos{}, rectI)
	rectIRS := orientation.ApplyToRect2(m.Pos{}, m.Rect{Origin: m.Pos{}, Size: scaledImgSize})

	// Step 3: rotate the image first.
	geoM.SetElement(0, 0, float64(orientation.Right.DX))
	geoM.SetElement(1, 0, float64(orientation.Right.DY))
	geoM.SetElement(0, 1, float64(orientation.Down.DX))
	geoM.SetElement(1, 1, float64(orientation.Down.DY))

	// Step 4: scale the image to the intended size.
	geoM.Scale(float64(rectR.Size.DX)/float64(rectIR.Size.DX),
		float64(rectR.Size.DY)/float64(rectIR.Size.DY))

	// Step 5: translate the image to the intended position.
	geoM.Translate(float64(rectR.Origin.X-rectIRS.Origin.X),
		float64(rectR.Origin.Y-rectIRS.Origin.Y))
}

func (w *World) drawTiles(screen *ebiten.Image, scrollDelta m.Delta) {
	for pos, tile := range w.Tiles {
		if tile.Image == nil {
			continue
		}
		screenPos := pos.Mul(TileSize).Add(scrollDelta)
		opts := ebiten.DrawImageOptions{
			// Note: could be CompositeModeCopy, but that can't be merged with entities pass.
			CompositeMode: ebiten.CompositeModeSourceOver,
			Filter:        ebiten.FilterNearest,
		}
		renderOrientation, renderImage := tile.Orientation, tile.Image
		if len(tile.ImageByOrientation) > 0 {
			// Locate pre-rotated tiles for better effect.
			o := tile.Transform.Concat(tile.Orientation)
			i := o.Inverse().Concat(tile.Orientation)
			img := tile.ImageByOrientation[i]
			if img != nil {
				renderOrientation, renderImage = o, img
			}
		}
		setGeoM(&opts.GeoM, screenPos, false, m.Delta{DX: TileSize, DY: TileSize}, m.Delta{DX: TileSize, DY: TileSize}, renderOrientation)
		screen.DrawImage(renderImage, &opts)
	}
}

func (w *World) drawEntities(screen *ebiten.Image, scrollDelta m.Delta) {
	zEnts := map[int][]*Entity{}
	for _, ent := range w.Entities {
		zEnts[ent.ZIndex] = append(zEnts[ent.ZIndex], ent)
	}
	for z := MinZIndex; z <= MaxZIndex; z++ {
		for _, ent := range zEnts[z] {
			if ent.Image == nil || ent.Alpha == 0 {
				continue
			}
			screenPos := ent.Rect.Origin.Add(scrollDelta).Add(ent.RenderOffset)
			opts := ebiten.DrawImageOptions{
				CompositeMode: ebiten.CompositeModeSourceOver,
				Filter:        ebiten.FilterNearest,
			}
			w, h := ent.Image.Size()
			imageSize := m.Delta{DX: w, DY: h}
			setGeoM(&opts.GeoM, screenPos, ent.ResizeImage, ent.Rect.Size, imageSize, ent.Orientation)
			opts.ColorM.Scale(1.0, 1.0, 1.0, ent.Alpha)
			screen.DrawImage(ent.Image, &opts)
		}
	}
}

func (w *World) drawDebug(screen *ebiten.Image, scrollDelta m.Delta) {
	for pos, tile := range w.Tiles {
		screenPos := pos.Mul(TileSize).Add(scrollDelta)
		if *debugShowNeighbors {
			neighborScreenPos := tile.LoadedFromNeighbor.Mul(TileSize).Add(scrollDelta)
			startx := float64(neighborScreenPos.X) + TileSize/2
			starty := float64(neighborScreenPos.Y) + TileSize/2
			endx := float64(screenPos.X) + TileSize/2
			endy := float64(screenPos.Y) + TileSize/2
			arrowpx := (startx + endx*2) / 3
			arrowpy := (starty + endy*2) / 3
			arrowdx := (endx - startx) / 6
			arrowdy := (endy - starty) / 6
			// Right only (1 0): left side goes by (-1, -1), right side by (-1, 1)
			// Down right (1 1): left side goes by (0, -2), right side by (-2, 0)
			// Down only (0 1): left side goes by (1, -1), right side by (-1, -1)
			// ax + by
			arrowlx := arrowpx - arrowdx + arrowdy
			arrowly := arrowpy - arrowdx - arrowdy
			arrowrx := arrowpx - arrowdx - arrowdy
			arrowry := arrowpy + arrowdx - arrowdy
			c := color.Gray{64}
			if tile.visibilityMark == w.visibilityMark {
				c = color.Gray{192}
			}
			ebitenutil.DrawLine(screen, startx, starty, endx, endy, c)
			ebitenutil.DrawLine(screen, arrowlx, arrowly, arrowpx, arrowpy, c)
			ebitenutil.DrawLine(screen, arrowrx, arrowry, arrowpx, arrowpy, c)
		}
		if *debugShowCoords {
			c := color.Gray{128}
			text.Draw(screen, fmt.Sprintf("%d,%d", tile.LevelPos.X, tile.LevelPos.Y), font.DebugSmall, screenPos.X, screenPos.Y+TileSize-1, c)
		}
		if *debugShowOrientations {
			midx := float64(screenPos.X) + TileSize/2
			midy := float64(screenPos.Y) + TileSize/2
			dx := tile.Orientation.Apply(m.Delta{DX: 4, DY: 0})
			ebitenutil.DrawLine(screen, midx, midy, midx+float64(dx.DX), midy+float64(dx.DY), color.NRGBA{R: 255, G: 0, B: 0, A: 255})
			dy := tile.Orientation.Apply(m.Delta{DX: 0, DY: 4})
			ebitenutil.DrawLine(screen, midx, midy, midx+float64(dy.DX), midy+float64(dy.DY), color.NRGBA{R: 0, G: 255, B: 0, A: 255})
		}
		if *debugShowTransforms {
			midx := float64(screenPos.X) + TileSize/2
			midy := float64(screenPos.Y) + TileSize/2
			dx := tile.Transform.Apply(m.Delta{DX: 4, DY: 0})
			ebitenutil.DrawLine(screen, midx, midy, midx+float64(dx.DX), midy+float64(dx.DY), color.NRGBA{R: 255, G: 0, B: 0, A: 255})
			dy := tile.Transform.Apply(m.Delta{DX: 0, DY: 4})
			ebitenutil.DrawLine(screen, midx, midy, midx+float64(dy.DX), midy+float64(dy.DY), color.NRGBA{R: 0, G: 255, B: 0, A: 255})
		}
	}
	for _, ent := range w.Entities {
		if *debugShowBboxes {
			boxColor := color.NRGBA{R: 128, G: 128, B: 128, A: 128}
			if ent.Solid {
				boxColor.R = 255
			}
			if ent.Opaque {
				boxColor.B = 255
			}
			ebitenutil.DrawRect(screen, float64(ent.Rect.Origin.X+scrollDelta.DX), float64(ent.Rect.Origin.Y+scrollDelta.DY), float64(ent.Rect.Size.DX), float64(ent.Rect.Size.DY), boxColor)
		}
	}
	if *debugShowTrace != "" {
		traces := strings.Split(*debugShowTrace, " ")
		for i := 0; i+1 < len(traces); i += 2 {
			var tracePos m.Delta
			var err error
			tracePos.DX, err = strconv.Atoi(traces[i])
			if err != nil {
				log.Panicf("invalid debug_show_trace %q: %v", traces[i], err)
			}
			tracePos.DY, err = strconv.Atoi(traces[i+1])
			if err != nil {
				log.Panicf("invalid debug_show_trace %q: %v", traces[i+1], err)
			}
			traceFrom := m.Pos{X: GameWidth / 2, Y: GameHeight / 2}.Sub(scrollDelta)
			traceTo := traceFrom.Add(tracePos)
			trace := w.TraceLine(traceFrom, traceTo, TraceOptions{})
			log.Print(trace)
			traceFromR := traceFrom.Add(scrollDelta)
			traceToR := traceTo.Add(scrollDelta)
			traceEndR := trace.EndPos.Add(scrollDelta)
			if i == 0 {
				for _, pos := range trace.Path {
					posR := pos.Mul(TileSize).Add(scrollDelta)
					a := float64(TileSize / 8)
					b := float64(TileSize) - 2*a
					ebitenutil.DrawRect(screen, float64(posR.X)+a, float64(posR.Y)+a, b, b, color.NRGBA{R: 0, G: 255, B: 0, A: 255})
				}
			}
			ebitenutil.DrawLine(screen, float64(traceFromR.X), float64(traceFromR.Y), float64(traceEndR.X), float64(traceEndR.Y), color.NRGBA{R: 255, G: 0, B: 0, A: 255})
			ebitenutil.DrawLine(screen, float64(traceEndR.X), float64(traceEndR.Y), float64(traceToR.X), float64(traceToR.Y), color.NRGBA{R: 255, G: 255, B: 0, A: 255})
		}
	}
}

func (w *World) rawDrawDest(screen *ebiten.Image) *ebiten.Image {
	if *drawVisibilityMask && *drawOutside {
		return w.offScreenBuffer
	}
	return screen
}

func (w *World) drawVisibilityMask(screen, drawDest *ebiten.Image, scrollDelta m.Delta) {
	// Draw trace polygon to buffer.
	geoM := ebiten.GeoM{}
	geoM.Translate(float64(scrollDelta.DX), float64(scrollDelta.DY))

	if w.needPrevImage {
		// Optimization note:
		// - This isn't optimal. Visibility mask maybe shouldn't even exist?
		// - If screen were a separate image, we could instead copy image to screen masked by polygon.
		// - Would remove one render call.
		// - Wouldn't allow blur though...?
		// Note: we put the mask on ALL four channels.
		w.visibilityMaskImage.Fill(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
		drawPolygonAround(w.visibilityMaskImage, w.visiblePolygonCenter, w.visiblePolygon, w.whiteImage, geoM, &ebiten.DrawTrianglesOptions{
			Address: ebiten.AddressRepeat,
		})

		e := expandSize
		if *expandUsingVertices {
			e = 0
		}
		BlurExpandImage(w.visibilityMaskImage, w.blurImage, w.visibilityMaskImage, blurSize, e, 1.0)
	}

	if *drawOutside {
		if *debugUseShaders && false {
			delta := w.scrollPos.Delta(w.prevScrollPos)
			screen.DrawRectShader(GameWidth, GameHeight, w.visibilityMaskShader, &ebiten.DrawRectShaderOptions{
				CompositeMode: ebiten.CompositeModeCopy,
				Uniforms: map[string]interface{}{
					"Scroll": []float32{float32(delta.DX) / GameWidth, float32(delta.DY) / GameHeight},
				},
				Images: [4]*ebiten.Image{
					w.visibilityMaskImage,
					drawDest,
					w.prevImage,
					nil,
				},
			})
		} else {
			// First set the alpha channel to the visibility mask.
			drawDest.DrawImage(w.visibilityMaskImage, &ebiten.DrawImageOptions{
				CompositeMode: ebiten.CompositeModeMultiply,
				Filter:        ebiten.FilterNearest,
			})

			// Then draw the background.
			delta := w.scrollPos.Delta(w.prevScrollPos)
			screen.DrawTriangles([]ebiten.Vertex{
				{
					DstX: 0, DstY: 0,
					SrcX: float32(delta.DX), SrcY: float32(delta.DY),
					ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1,
				},
				{
					DstX: GameWidth, DstY: 0,
					SrcX: GameWidth + float32(delta.DX), SrcY: float32(delta.DY),
					ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1,
				},
				{
					DstX: 0, DstY: GameHeight,
					SrcX: float32(delta.DX), SrcY: GameHeight + float32(delta.DY),
					ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1,
				},
				{
					DstX: GameWidth, DstY: GameHeight,
					SrcX: GameWidth + float32(delta.DX), SrcY: GameHeight + float32(delta.DY),
					ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1,
				},
			}, []uint16{0, 1, 2, 1, 2, 3}, w.prevImage, &ebiten.DrawTrianglesOptions{
				CompositeMode: ebiten.CompositeModeCopy,
				Filter:        ebiten.FilterNearest,
				Address:       ebiten.AddressClampToZero,
			})

			// Finally put the masked foreground on top.
			screen.DrawImage(drawDest, &ebiten.DrawImageOptions{
				CompositeMode: ebiten.CompositeModeSourceOver,
				Filter:        ebiten.FilterNearest,
			})
		}
		if w.needPrevImage {
			// Remember last image. Only do this once per update.
			BlurImage(screen, w.blurImage, w.prevImage, frameBlurSize, frameDarkenAlpha)
			w.prevScrollPos = w.scrollPos
		}
	} else {
		screen.DrawImage(w.visibilityMaskImage, &ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeDestinationIn,
			Filter:        ebiten.FilterNearest,
		})
	}

	w.needPrevImage = false
}

func (w *World) drawOverlays(screen *ebiten.Image, scrollDelta m.Delta) {
	zEnts := map[int][]*Entity{}
	for _, ent := range w.Entities {
		zEnts[ent.ZIndex] = append(zEnts[ent.ZIndex], ent)
	}
	for z := MinZIndex; z <= MaxZIndex; z++ {
		for _, ent := range zEnts[z] {
			ent.Impl.DrawOverlay(screen, scrollDelta)
		}
	}
}

func (w *World) Draw(screen *ebiten.Image) {
	scrollDelta := m.Pos{X: GameWidth / 2, Y: GameHeight / 2}.Delta(w.scrollPos)

	dest := w.rawDrawDest(screen)
	dest.Fill(color.Gray{0})
	w.drawTiles(dest, scrollDelta)
	w.drawEntities(dest, scrollDelta)
	if *drawVisibilityMask {
		w.drawVisibilityMask(screen, dest, scrollDelta)
	}
	w.drawOverlays(screen, scrollDelta)
	centerprint.Draw(screen)

	// Debug stuff comes last.
	w.drawDebug(screen, scrollDelta)
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
	if _, found := w.Tiles[newPos]; found {
		// Already loaded.
		return newPos
	}
	neighborTile := w.Tiles[p]
	if neighborTile == nil {
		log.Panicf("Trying to load with nonexisting neighbor tile at %v", p)
	}
	t := neighborTile.Transform
	neighborLevelPos := neighborTile.LevelPos
	newLevelPos := neighborLevelPos.Add(t.Apply(d))
	newLevelTile, found := w.Level.Tiles[newLevelPos]
	if !found {
		// log.Printf("Trying to load nonexisting tile at %v when moving from %v (%v) by %v (%v)", newLevelPos, p, neighborLevelPos, d, t.Apply(d))
		newTile := Tile{
			LevelPos:           newLevelPos,
			Transform:          t,
			LoadedFromNeighbor: p,
		}
		w.Tiles[newPos] = &newTile
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
		tile := w.Level.Tiles[warp.ToTile]
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
	w.Tiles[newPos] = &newTile
	return newPos
}

// tilesBox returns corner coordinates for all tiles in a given box.
func tilesBox(r m.Rect) (m.Pos, m.Pos) {
	tp0 := r.Origin.Div(TileSize)
	tp1 := r.OppositeCorner().Div(TileSize)
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

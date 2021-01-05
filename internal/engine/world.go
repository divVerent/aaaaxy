package engine

import (
	"flag"
	"fmt"
	"image/color"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"

	m "github.com/divVerent/aaaaaa/internal/math"
	"github.com/divVerent/aaaaaa/internal/timing"
)

var (
	debugShowNeighbors    = flag.Bool("debug_show_neighbors", false, "show the neighbors tiles got loaded from")
	debugShowCoords       = flag.Bool("debug_show_coords", false, "show the level coordinates of each tile")
	debugShowOrientations = flag.Bool("debug_show_orientations", false, "show the orientation of each tile")
	debugShowTransforms   = flag.Bool("debug_show_transforms", false, "show the transform of each tile")
	drawBlurs             = flag.Bool("draw_blurs", true, "perform blur effects; requires draw_visibility_mask")
	drawOutside           = flag.Bool("draw_outside", true, "draw outside of the visible area; requires draw_visibility_mask")
	drawVisibilityMask    = flag.Bool("draw_visibility_mask", true, "draw visibility mask (if disabled, all loaded tiles are shown")
	expandUsingVertices   = flag.Bool("expand_using_vertices", false, "expand using polygon math (just approximate, simplifies rendering)")
)

// World represents the current game state including its entities.
type World struct {
	// tiles are all tiles currently loaded.
	Tiles map[m.Pos]*Tile
	// entities are all entities currently loaded.
	Entities map[EntityID]*Entity
	// Player is the ID of the player entity.
	PlayerID EntityID
	// level is the current tilemap (universal covering with warpZones).
	Level *Level

	// Properties that can in theory be regenerated from the above and thus do not
	// need serialization support.

	// scrollPos is the current screen scrolling position.
	scrollPos m.Pos
	// visibilityMark is the current mark value to detect visible tiles/objects.
	visibilityMark uint
	// debugFont is the font to use for debug messages.
	debugFont font.Face
	// visiblePolygonCenter is the current eye position.
	visiblePolygonCenter m.Pos
	// visiblePolygon is the currently visible polygon.
	visiblePolygon []m.Pos
	// needPrevImageMasked is set whenever the last call was Update.
	needPrevImageMasked bool

	// Images retained across frames.

	// whiteImage is a single white pixel.
	whiteImage *ebiten.Image
	// prevImage is the previous screen content.
	prevImage *ebiten.Image
	// prevImageMasked is the previous screen content after masking.
	prevImageMasked *ebiten.Image
	// prevScrollPos is previous frame's scroll pos.
	prevScrollPos m.Pos

	// Temp storage within frames.

	// blurImage is an offscreen image used for blurring.
	blurImage *ebiten.Image
	// visibilityMaskImage is an offscreen image used for masking the visible area.
	visibilityMaskImage *ebiten.Image
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
		debugFont: truetype.NewFace(debugFont, &truetype.Options{
			Size:    5,
			Hinting: font.HintingFull,
		}),
		whiteImage:          ebiten.NewImage(1, 1),
		blurImage:           ebiten.NewImage(GameWidth, GameHeight),
		prevImage:           ebiten.NewImage(GameWidth, GameHeight),
		prevImageMasked:     ebiten.NewImage(GameWidth, GameHeight),
		visibilityMaskImage: ebiten.NewImage(GameWidth, GameHeight),
	}
	w.whiteImage.Fill(color.Gray{255})
	w.prevImage.Fill(color.Gray{0})
	w.prevImageMasked.Fill(color.Gray{0})

	// Load tile the player starts on.
	tile := w.Level.Tiles[w.Level.Player.LevelPos].Tile
	tile.Transform = m.Identity()
	w.Tiles[w.Level.Player.LevelPos] = &tile

	// Create player entity.
	w.PlayerID = w.Level.Player.ID
	playerEnt, err := w.Level.Player.Spawn(&w, w.Level.Player.LevelPos, &tile)
	if err != nil {
		log.Panicf("could not spawn player: %v", err)
	}
	w.Entities[w.PlayerID] = playerEnt

	// Load the other tiles that the player touches.
	w.LoadTilesForRect(w.Entities[w.PlayerID].Rect, w.Level.Player.LevelPos)
	w.visibilityMark++

	return &w
}

func (w *World) traceLineAndMark(from, to m.Pos) TraceResult {
	result := w.TraceLine(from, to, TraceOptions{
		Mode:      HitOpaque,
		LoadTiles: true,
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

// updateEntities lets all entities move/act.
func (w *World) updateEntities() {
	for _, ent := range w.Entities {
		ent.Impl.Update()
	}
}

// updateScrollPos updates the current scroll position.
func (w *World) updateScrollPos(player *Entity, target m.Pos) {
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
	if target.X < player.Rect.OppositeCorner().X-GameWidth/2+scrollMinDistance {
		target.X = player.Rect.OppositeCorner().X - GameWidth/2 + scrollMinDistance
	}
	if target.X > player.Rect.Origin.X+GameWidth/2-scrollMinDistance {
		target.X = player.Rect.Origin.X + GameWidth/2 - scrollMinDistance
	}
	if target.Y < player.Rect.OppositeCorner().Y-GameHeight/2+scrollMinDistance {
		target.Y = player.Rect.OppositeCorner().Y - GameHeight/2 + scrollMinDistance
	}
	if target.Y > player.Rect.Origin.Y+GameHeight/2-scrollMinDistance {
		target.Y = player.Rect.Origin.Y + GameHeight/2 - scrollMinDistance
	}
	w.scrollPos = target
}

// updateVisibility loads all visible tiles and discards all tiles not visible right now.
func (w *World) updateVisibility(eye m.Pos) {
	defer timing.Group()()

	// Delete all tiles merely marked for expanding.
	// TODO can we preserve but recheck them instead?
	timing.Section("cleanup_expanded")
	prevVisibilityMark := w.visibilityMark - 1
	for pos, tile := range w.Tiles {
		if tile.visibilityMark != prevVisibilityMark {
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
	for x := screen0.X; x < screen1.X; x += sweepStep {
		trace := w.traceLineAndMark(eye, m.Pos{X: x, Y: screen0.Y})
		w.visiblePolygon = append(w.visiblePolygon, trace.EndPos)
	}
	for y := screen0.Y; y < screen1.Y; y += sweepStep {
		trace := w.traceLineAndMark(eye, m.Pos{X: screen1.X, Y: y})
		w.visiblePolygon = append(w.visiblePolygon, trace.EndPos)
	}
	for x := screen1.X; x > screen0.X; x -= sweepStep {
		trace := w.traceLineAndMark(eye, m.Pos{X: x, Y: screen1.Y})
		w.visiblePolygon = append(w.visiblePolygon, trace.EndPos)
	}
	for y := screen1.Y; y > screen0.Y; y -= sweepStep {
		trace := w.traceLineAndMark(eye, m.Pos{X: screen0.X, Y: y})
		w.visiblePolygon = append(w.visiblePolygon, trace.EndPos)
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
			if w.Entities[spawnable.ID] == nil {
				timing.Section("spawn")
				ent, err := spawnable.Spawn(w, pos, tile)
				if err != nil {
					log.Panicf("could not spawn entity %v: %v", spawnable, err)
				}
				w.Entities[spawnable.ID] = ent
				timing.Section("spawn_search")
			}
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
				if w.Tiles[tp].visibilityMark == expansionMark || w.Tiles[tp].visibilityMark == visibilityMark {
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

	// Let everything move.
	timing.Section("entities")
	w.updateEntities()

	// Fetch the player entity.
	player := w.Entities[w.PlayerID]
	playerImpl := player.Impl.(PlayerEntityImpl)

	// Scroll towards the focus point.
	w.updateScrollPos(player, playerImpl.LookPos())

	// Update visibility and spawn/despawn entities.
	timing.Section("visibility")
	w.updateVisibility(playerImpl.EyePos())

	w.needPrevImageMasked = true
	return nil
}

func setGeoM(geoM *ebiten.GeoM, pos m.Pos, size m.Delta, orientation m.Orientation) {
	// Set the rotation.
	geoM.SetElement(0, 0, float64(orientation.Right.DX))
	geoM.SetElement(1, 0, float64(orientation.Right.DY))
	geoM.SetElement(0, 1, float64(orientation.Down.DX))
	geoM.SetElement(1, 1, float64(orientation.Down.DY))
	// Set the translation.
	// Note that in ebiten, the coordinate is the original origin, while we think in screenspace origin.
	a := m.Delta{} // Actually orientation.Apply(m.Delta{})
	d := orientation.Apply(size)
	if a.DX > d.DX {
		a.DX = d.DX
	}
	if a.DY > d.DY {
		a.DY = d.DY
	}
	geoM.Translate(float64(pos.X-a.DX), float64(pos.Y-a.DY))
}

func (w *World) drawTiles(screen *ebiten.Image, scrollDelta m.Delta) {
	for pos, tile := range w.Tiles {
		if tile.Image == nil {
			continue
		}
		screenPos := pos.Mul(TileSize).Add(scrollDelta)
		opts := ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeCopy,
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
		setGeoM(&opts.GeoM, screenPos, m.Delta{DX: TileSize, DY: TileSize}, renderOrientation)
		screen.DrawImage(renderImage, &opts)
	}
}

func (w *World) drawEntities(screen *ebiten.Image, scrollDelta m.Delta) {
	for _, ent := range w.Entities {
		screenPos := ent.Rect.Origin.Add(scrollDelta)
		opts := ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeSourceAtop,
			Filter:        ebiten.FilterNearest,
		}
		setGeoM(&opts.GeoM, screenPos, ent.Rect.Size, ent.Orientation)
		screen.DrawImage(ent.Image, &opts)
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
			text.Draw(screen, fmt.Sprintf("%d,%d", tile.LevelPos.X, tile.LevelPos.Y), w.debugFont, screenPos.X, screenPos.Y+TileSize-1, c)
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
}

func (w *World) drawVisibilityMask(screen *ebiten.Image, scrollDelta m.Delta) {
	// Draw trace polygon to buffer.
	geoM := ebiten.GeoM{}
	geoM.Translate(float64(scrollDelta.DX), float64(scrollDelta.DY))
	w.visibilityMaskImage.Fill(color.Gray{0})
	drawPolygonAround(w.visibilityMaskImage, w.visiblePolygonCenter, w.visiblePolygon, w.whiteImage, geoM, &ebiten.DrawTrianglesOptions{
		Address: ebiten.AddressRepeat,
	})

	if !*expandUsingVertices {
		expandImage(w.visibilityMaskImage, w.blurImage, expandSize, 1.0)
	}
	if *drawBlurs {
		expandImage(w.visibilityMaskImage, w.blurImage, blurSize, 0.5)
	}

	screen.DrawImage(w.visibilityMaskImage, &ebiten.DrawImageOptions{
		CompositeMode: ebiten.CompositeModeMultiply,
		Filter:        ebiten.FilterNearest,
	})

	if *drawOutside {
		delta := w.scrollPos.Delta(w.prevScrollPos)
		if w.needPrevImageMasked {
			// Make a scrolled copy of the last frame.
			w.prevImageMasked.Fill(color.Gray{0})
			opts := ebiten.DrawImageOptions{
				CompositeMode: ebiten.CompositeModeCopy,
				Filter:        ebiten.FilterNearest,
			}
			opts.GeoM.Translate(float64(-delta.DX), float64(-delta.DY))
			w.prevImageMasked.DrawImage(w.prevImage, &opts)

			// Blur and darken last image.
			if *drawBlurs {
				expandImage(w.prevImageMasked, w.blurImage, frameBlurSize, 0.5)
			}

			// Mask out the parts we've already drawn.
			opts = ebiten.DrawImageOptions{
				CompositeMode: ebiten.CompositeModeMultiply,
				Filter:        ebiten.FilterNearest,
			}
			opts.ColorM.Scale(-frameDarkenAlpha, -frameDarkenAlpha, -frameDarkenAlpha, 0)
			opts.ColorM.Translate(frameDarkenAlpha, frameDarkenAlpha, frameDarkenAlpha, 1)
			w.prevImageMasked.DrawImage(w.visibilityMaskImage, &opts)
		}

		// Add it to what we see.
		screen.DrawImage(w.prevImageMasked, &ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeLighter,
			Filter:        ebiten.FilterNearest,
		})

		if w.needPrevImageMasked {
			// Remember last image. Only do this once per update.
			w.prevImage.DrawImage(screen, &ebiten.DrawImageOptions{
				CompositeMode: ebiten.CompositeModeCopy,
				Filter:        ebiten.FilterNearest,
			})
			w.prevScrollPos = w.scrollPos
			w.needPrevImageMasked = false
		}
	}
}

func (w *World) Draw(screen *ebiten.Image) {
	scrollDelta := m.Pos{X: GameWidth / 2, Y: GameHeight / 2}.Delta(w.scrollPos)

	screen.Fill(color.Gray{0})
	w.drawTiles(screen, scrollDelta)
	w.drawEntities(screen, scrollDelta)
	if *drawVisibilityMask {
		w.drawVisibilityMask(screen, scrollDelta)
	}

	// Debug stuff comes last.
	w.drawDebug(screen, scrollDelta)
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
		log.Printf("Trying to load nonexisting tile at %v when moving from %v (%v) by %v (%v)",
			newLevelPos, p, neighborTile.LevelPos, d, t.Apply(d))
		newTile := Tile{
			LevelPos:           newLevelPos,
			Transform:          t,
			LoadedFromNeighbor: p,
		}
		w.Tiles[newPos] = &newTile
		return newPos
	}
	if newLevelTile.WarpZone != nil {
		t = newLevelTile.WarpZone.Transform.Concat(t)
		tile := w.Level.Tiles[newLevelTile.WarpZone.ToTile]
		if tile == nil {
			log.Panicf("nil new tile after warping to %v", newLevelTile.WarpZone)
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

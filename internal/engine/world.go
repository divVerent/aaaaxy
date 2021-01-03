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
	// VisiblePolygonCenter is the current eye position.
	VisiblePolygonCenter m.Pos
	// VisiblePolygon is the currently visible polygon.
	VisiblePolygon []m.Pos
	// NeedPrevImageMasked is set whenever the last call was Update.
	NeedPrevImageMasked bool

	// Images retained across frames.

	// WhiteImage is a single white pixel.
	WhiteImage *ebiten.Image
	// PrevImage is the previous screen content.
	PrevImage *ebiten.Image
	// PrevImageMasked is the previous screen content after masking.
	PrevImageMasked *ebiten.Image
	// PrevScrollPos is previous frame's scroll pos.
	PrevScrollPos m.Pos

	// Temp storage within frames.

	// BlurImage is an offscreen image used for blurring.
	BlurImage *ebiten.Image
	// VisibilityMaskImage is an offscreen image used for masking the visible area.
	VisibilityMaskImage *ebiten.Image
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
		WhiteImage:          ebiten.NewImage(1, 1),
		BlurImage:           ebiten.NewImage(GameWidth, GameHeight),
		PrevImage:           ebiten.NewImage(GameWidth, GameHeight),
		PrevImageMasked:     ebiten.NewImage(GameWidth, GameHeight),
		VisibilityMaskImage: ebiten.NewImage(GameWidth, GameHeight),
	}
	w.WhiteImage.Fill(color.Gray{255})
	w.PrevImage.Fill(color.Gray{0})
	w.PrevImageMasked.Fill(color.Gray{0})

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
	w.VisibilityMark++

	return &w
}

func (w *World) traceLineAndMark(from, to m.Pos) TraceResult {
	result := w.TraceLine(from, to, TraceOptions{
		Mode:      HitOpaque,
		LoadTiles: true,
	})
	for _, tilePos := range result.Path {
		w.Tiles[tilePos].VisibilityMark = w.VisibilityMark
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

func (w *World) Update() error {
	// Let all entities move/act. Fetch player position.
	for _, ent := range w.Entities {
		ent.Impl.Update()
	}

	// Player entity has special treatment.
	player := w.Entities[w.PlayerID]
	playerImpl := player.Impl.(PlayerEntityImpl)

	// Update scroll position.
	targetScrollPos := playerImpl.LookPos()
	// Slowly move towards focus point.
	// TODO Even converge to center when ScrollPerFrame too low. Somehow?
	targetScrollPos = w.ScrollPos.Add(targetScrollPos.Delta(w.ScrollPos).MulFloat(ScrollPerFrame))
	// Ensure player is onscreen.
	if targetScrollPos.X < player.Rect.OppositeCorner().X-GameWidth/2+ScrollMinDistance {
		targetScrollPos.X = player.Rect.OppositeCorner().X - GameWidth/2 + ScrollMinDistance
	}
	if targetScrollPos.X > player.Rect.Origin.X+GameWidth/2-ScrollMinDistance {
		targetScrollPos.X = player.Rect.Origin.X + GameWidth/2 - ScrollMinDistance
	}
	if targetScrollPos.Y < player.Rect.OppositeCorner().Y-GameHeight/2+ScrollMinDistance {
		targetScrollPos.Y = player.Rect.OppositeCorner().Y - GameHeight/2 + ScrollMinDistance
	}
	if targetScrollPos.Y > player.Rect.Origin.Y+GameHeight/2-ScrollMinDistance {
		targetScrollPos.Y = player.Rect.Origin.Y + GameHeight/2 - ScrollMinDistance
	}
	w.ScrollPos = targetScrollPos

	// Delete all tiles merely marked for expanding.
	// TODO can we preserve but recheck them instead?
	expansionMark := w.VisibilityMark
	for pos, tile := range w.Tiles {
		if tile.VisibilityMark == expansionMark {
			delete(w.Tiles, pos)
		}
	}

	// Unmark all tiles and entities (just bump mark index).
	w.VisibilityMark++
	visibilityMark := w.VisibilityMark

	// Trace from player location to all directions (SweepStep pixels at screen edge).
	// Mark all tiles hit (excl. the tiles that stopped us).
	// TODO Remember trace polygon.
	screen0 := w.ScrollPos.Sub(m.Delta{DX: GameWidth / 2, DY: GameHeight / 2})
	screen1 := screen0.Add(m.Delta{DX: GameWidth - 1, DY: GameHeight - 1})
	eye := playerImpl.EyePos()
	w.VisiblePolygonCenter = eye
	w.VisiblePolygon = w.VisiblePolygon[0:0]
	for x := screen0.X; x < screen1.X; x += SweepStep {
		trace := w.traceLineAndMark(eye, m.Pos{X: x, Y: screen0.Y})
		w.VisiblePolygon = append(w.VisiblePolygon, trace.EndPos)
	}
	for y := screen0.Y; y < screen1.Y; y += SweepStep {
		trace := w.traceLineAndMark(eye, m.Pos{X: screen1.X, Y: y})
		w.VisiblePolygon = append(w.VisiblePolygon, trace.EndPos)
	}
	for x := screen1.X; x > screen0.X; x -= SweepStep {
		trace := w.traceLineAndMark(eye, m.Pos{X: x, Y: screen1.Y})
		w.VisiblePolygon = append(w.VisiblePolygon, trace.EndPos)
	}
	for y := screen1.Y; y > screen0.Y; y -= SweepStep {
		trace := w.traceLineAndMark(eye, m.Pos{X: screen0.X, Y: y})
		w.VisiblePolygon = append(w.VisiblePolygon, trace.EndPos)
	}
	if *expandUsingVertices {
		expandPolygon(w.VisiblePolygonCenter, w.VisiblePolygon, ExpandSize)
	}

	// Also mark all neighbors of hit tiles hit (up to ExpandTiles).
	// For multiple expansion, need to do this in steps so initially we only base expansion on visible tiles.
	markedTiles := []m.Pos{}
	for tilePos, tile := range w.Tiles {
		if tile.VisibilityMark == visibilityMark {
			markedTiles = append(markedTiles, tilePos)
		}
	}
	w.VisibilityMark++
	expansionMark = w.VisibilityMark
	numExpandSteps := (2*ExpandTiles+1)*(2*ExpandTiles+1) - 1
	for i := 0; i < numExpandSteps; i++ {
		step := &ExpandSteps[i]
		for _, pos := range markedTiles {
			from := pos.Add(step.from)
			to := pos.Add(step.to)
			w.LoadTile(from, to.Delta(from))
			if w.Tiles[to].VisibilityMark != visibilityMark {
				w.Tiles[to].VisibilityMark = expansionMark
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
		if tile.VisibilityMark != expansionMark && tile.VisibilityMark != visibilityMark {
			delete(w.Tiles, pos)
		}
	}

	w.NeedPrevImageMasked = true

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

func (w *World) Draw(screen *ebiten.Image) {
	screen.Fill(color.Gray{0})

	scrollDelta := m.Pos{X: GameWidth / 2, Y: GameHeight / 2}.Delta(w.ScrollPos)

	if *drawVisibilityMask {
		// Draw trace polygon to buffer.
		geoM := ebiten.GeoM{}
		geoM.Translate(float64(scrollDelta.DX), float64(scrollDelta.DY))
		w.VisibilityMaskImage.Fill(color.Gray{0})
		DrawPolygonAround(w.VisibilityMaskImage, w.VisiblePolygonCenter, w.VisiblePolygon, w.WhiteImage, geoM, &ebiten.DrawTrianglesOptions{
			Address: ebiten.AddressRepeat,
		})

		// TODO Expand and blur buffer (ExpandSize, BlurSize).
		if !*expandUsingVertices {
			ExpandImage(w.VisibilityMaskImage, w.BlurImage, ExpandSize, 1.0)
		}
		if *drawBlurs {
			ExpandImage(w.VisibilityMaskImage, w.BlurImage, BlurSize, 0.5)
		}
	}

	// Draw all tiles.
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
			if tile.VisibilityMark == w.VisibilityMark {
				c = color.Gray{192}
			}
			ebitenutil.DrawLine(screen, startx, starty, endx, endy, c)
			ebitenutil.DrawLine(screen, arrowlx, arrowly, arrowpx, arrowpy, c)
			ebitenutil.DrawLine(screen, arrowrx, arrowry, arrowpx, arrowpy, c)
		}
		if *debugShowCoords {
			c := color.Gray{128}
			text.Draw(screen, fmt.Sprintf("%d,%d", tile.LevelPos.X, tile.LevelPos.Y), w.DebugFont, screenPos.X, screenPos.Y+TileSize-1, c)
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

	// TODO Draw all entities.
	for _, ent := range w.Entities {
		screenPos := ent.Rect.Origin.Add(scrollDelta)
		opts := ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeSourceAtop,
			Filter:        ebiten.FilterNearest,
		}
		setGeoM(&opts.GeoM, screenPos, ent.Rect.Size, ent.Orientation)
		screen.DrawImage(ent.Image, &opts)
	}

	// NOTE: if an entity is on a tile seen twice, render only once.
	// INTENTIONAL GLITCH (avoids rendering player twice and player-player collision). Entities live in tile coordinates, not world coordinates. "Looking away" can despawn these entities and respawn at their new location.
	// Makes wrap-around rooms somewhat less obvious.
	// Only way to fix seems to be making everything live in "universal covering" coordinates with orientation? Seems not worth it.
	// TODO: Decide if to keep this.

	// Mum.Delta{} // Actually ltiply screen with buffer.
	if *drawVisibilityMask {
		screen.DrawImage(w.VisibilityMaskImage, &ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeMultiply,
			Filter:        ebiten.FilterNearest,
		})

		if *drawOutside {
			delta := w.ScrollPos.Delta(w.PrevScrollPos)
			if w.NeedPrevImageMasked {
				// Make a scrolled copy of the last frame.
				w.PrevImageMasked.Fill(color.Gray{0})
				opts := ebiten.DrawImageOptions{
					CompositeMode: ebiten.CompositeModeCopy,
					Filter:        ebiten.FilterNearest,
				}
				opts.GeoM.Translate(float64(-delta.DX), float64(-delta.DY))
				w.PrevImageMasked.DrawImage(w.PrevImage, &opts)

				// Blur and darken last image.
				if *drawBlurs {
					ExpandImage(w.PrevImageMasked, w.BlurImage, FrameBlurSize, 0.5)
				}

				// Mask out the parts we've already drawn.
				opts = ebiten.DrawImageOptions{
					CompositeMode: ebiten.CompositeModeMultiply,
					Filter:        ebiten.FilterNearest,
				}
				opts.ColorM.Scale(-FrameDarkenAlpha, -FrameDarkenAlpha, -FrameDarkenAlpha, 0)
				opts.ColorM.Translate(FrameDarkenAlpha, FrameDarkenAlpha, FrameDarkenAlpha, 1)
				w.PrevImageMasked.DrawImage(w.VisibilityMaskImage, &opts)
			}

			// Add it to what we see.
			screen.DrawImage(w.PrevImageMasked, &ebiten.DrawImageOptions{
				CompositeMode: ebiten.CompositeModeLighter,
				Filter:        ebiten.FilterNearest,
			})

			if w.NeedPrevImageMasked {
				// Remember last image. Only do this once per update.
				w.PrevImage.DrawImage(screen, &ebiten.DrawImageOptions{
					CompositeMode: ebiten.CompositeModeCopy,
					Filter:        ebiten.FilterNearest,
				})
				w.PrevScrollPos = w.ScrollPos
			}
		}
	}

	w.NeedPrevImageMasked = false
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
	if newLevelTile.Warpzone != nil {
		t = newLevelTile.Warpzone.Transform.Concat(t)
		tile := w.Level.Tiles[newLevelTile.Warpzone.ToTile]
		if tile == nil {
			log.Panicf("nil new tile after warping to %v", newLevelTile.Warpzone)
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

// LoadTilesForRect loads all tiles in the given box (p, d), assuming tile tp is already loaded.
func (w *World) LoadTilesForRect(r m.Rect, tp m.Pos) {
	// Convert box to tile positions.
	tp0 := r.Origin.Div(TileSize)
	tp1 := r.OppositeCorner().Div(TileSize)
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
	return TraceLine(w, from, to, o)
}

// TraceBox moves from x,y size sx,sy by dx,dy in pixel coordinates.
func (w *World) TraceBox(from m.Rect, to m.Pos, o TraceOptions) TraceResult {
	return TraceBox(w, from, to, o)
}

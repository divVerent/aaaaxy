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
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/divVerent/aaaaxy/internal/centerprint"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/image"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/shader"
)

var (
	debugShowNeighbors            = flag.Bool("debug_show_neighbors", false, "show the neighbors tiles got loaded from")
	debugShowCoords               = flag.Bool("debug_show_coords", false, "show the level coordinates of each tile")
	debugShowOrientations         = flag.Bool("debug_show_orientations", false, "show the orientation of each tile")
	debugShowTransforms           = flag.Bool("debug_show_transforms", false, "show the transform of each tile")
	cheatShowBboxes               = flag.Bool("cheat_show_bboxes", false, "show the bounding boxes of all entities")
	debugShowVisiblePolygon       = flag.Bool("debug_show_visible_polygon", false, "show the visibility polygon")
	drawOutside                   = flag.Bool("draw_outside", flag.SystemDefault(map[string]interface{}{"js/*": false, "*/*": true}).(bool), "draw outside of the visible area; requires draw_visibility_mask")
	drawVisibilityMask            = flag.Bool("draw_visibility_mask", true, "draw visibility mask (if disabled, all loaded tiles are shown")
	expandUsingVertices           = flag.Bool("expand_using_vertices", true, "expand using polygon math (simplifies rendering)")
	expandUsingVerticesAccurately = flag.Bool("expand_using_vertices_accurately", true, "expand using simpler polygon math (just approximate, removes a render pass)")
)

type renderer struct {
	world *World

	// visiblePolygonCenter is the current eye position.
	visiblePolygonCenter m.Pos
	// visiblePolygon is the currently visible polygon.
	visiblePolygon []m.Pos
	// expandedVisiblePolygon is the visible polygon, expanded to show some walls.
	expandedVisiblePolygon []m.Pos
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
}

func (r *renderer) Init(w *World) {
	r.world = w
	r.whiteImage = ebiten.NewImage(1, 1)
	r.whiteImage = ebiten.NewImage(1, 1)
	r.whiteImage.Fill(color.Gray{255})
	r.blurImage = ebiten.NewImage(GameWidth, GameHeight)
	r.prevImage = ebiten.NewImage(GameWidth, GameHeight)
	r.prevImage.Fill(color.Gray{0})
	r.offScreenBuffer = ebiten.NewImage(GameWidth, GameHeight)
	r.offScreenBuffer.Fill(color.Gray{0})
	r.visibilityMaskImage = ebiten.NewImage(GameWidth, GameHeight)

	if *debugUseShaders {
		var err error
		r.visibilityMaskShader, err = shader.Load("visibility_mask.kage", nil)
		if err != nil {
			log.Errorf("BROKEN RENDERER, WILL FALLBACK: could not load visibility mask shader: %v", err)
			r.visibilityMaskShader = nil
		}
	}
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

func (r *renderer) drawTiles(screen *ebiten.Image, scrollDelta m.Delta) {
	r.world.forEachTile(func(pos m.Pos, tile *level.Tile) {
		if tile.ImageSrc == "" {
			return
		}
		screenPos := pos.Mul(level.TileSize).Add(scrollDelta)
		img, err := image.Load("tiles", tile.ImageSrc)
		if err != nil {
			log.Errorf("could not load already cached image %q for tile: %v", tile.ImageSrc, err)
			return
		}
		opts := ebiten.DrawImageOptions{
			// Note: could be CompositeModeCopy, but that can't be merged with entities pass.
			CompositeMode: ebiten.CompositeModeSourceOver,
			Filter:        ebiten.FilterNearest,
		}
		setGeoM(&opts.GeoM, screenPos, false, m.Delta{DX: level.TileSize, DY: level.TileSize}, m.Delta{DX: level.TileSize, DY: level.TileSize}, tile.Orientation)
		opts.ColorM = r.world.GlobalColorM
		screen.DrawImage(img, &opts)
	})
}

func (r *renderer) drawEntities(screen *ebiten.Image, scrollDelta m.Delta) {
	minZ, maxZ := zBounds(len(r.world.entitiesByZ))
	for z := minZ; z <= maxZ; z++ {
		for _, colormods := range []bool{true, false} {
			r.world.entitiesByZ[encodeZ(z)].forEach(func(ent *Entity) error {
				if ent.Image == nil || ent.Alpha == 0 || (ent.ColorAdd != [4]float64{0, 0, 0, 0}) != colormods {
					return nil
				}
				screenPos := ent.Rect.Origin.Add(scrollDelta).Add(ent.RenderOffset)
				opts := ebiten.DrawImageOptions{
					CompositeMode: ebiten.CompositeModeSourceOver,
					Filter:        ebiten.FilterNearest,
				}
				w, h := ent.Image.Size()
				imageSize := m.Delta{DX: w, DY: h}
				setGeoM(&opts.GeoM, screenPos, ent.ResizeImage, ent.Rect.Size, imageSize, ent.Orientation)
				opts.ColorM.Scale(ent.ColorMod[0], ent.ColorMod[1], ent.ColorMod[2], ent.ColorMod[3])
				opts.ColorM.Translate(ent.ColorAdd[0], ent.ColorAdd[1], ent.ColorAdd[2], ent.ColorAdd[3])
				opts.ColorM.Scale(1.0, 1.0, 1.0, ent.Alpha)
				opts.ColorM.Concat(r.world.GlobalColorM)
				screen.DrawImage(ent.Image, &opts)
				return nil
			})
		}
	}
}

func (r *renderer) drawDebug(screen *ebiten.Image, scrollDelta m.Delta) {
	r.world.forEachTile(func(pos m.Pos, tile *level.Tile) {
		screenPos := pos.Mul(level.TileSize).Add(scrollDelta)
		if *debugShowNeighbors {
			neighborScreenPos := tile.LoadedFromNeighbor.Mul(level.TileSize).Add(scrollDelta)
			startx := float64(neighborScreenPos.X) + level.TileSize/2
			starty := float64(neighborScreenPos.Y) + level.TileSize/2
			endx := float64(screenPos.X) + level.TileSize/2
			endy := float64(screenPos.Y) + level.TileSize/2
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
			if tile.VisibilityFlags&level.FrameVis == r.world.frameVis {
				c = color.Gray{192}
			}
			ebitenutil.DrawLine(screen, startx, starty, endx, endy, c)
			ebitenutil.DrawLine(screen, arrowlx, arrowly, arrowpx, arrowpy, c)
			ebitenutil.DrawLine(screen, arrowrx, arrowry, arrowpx, arrowpy, c)
		}
		if *debugShowCoords {
			c := color.Gray{128}
			font.DebugSmall.Draw(screen, fmt.Sprintf("%d,%d", tile.LevelPos.X, tile.LevelPos.Y), screenPos.Add(m.Delta{
				DX: 0,
				DY: level.TileSize - 1,
			}), false, c, color.Transparent)
		}
		if *debugShowOrientations {
			midx := float64(screenPos.X) + level.TileSize/2
			midy := float64(screenPos.Y) + level.TileSize/2
			dx := tile.Orientation.Apply(m.Delta{DX: 4, DY: 0})
			ebitenutil.DrawLine(screen, midx, midy, midx+float64(dx.DX), midy+float64(dx.DY), color.NRGBA{R: 255, G: 0, B: 0, A: 255})
			dy := tile.Orientation.Apply(m.Delta{DX: 0, DY: 4})
			ebitenutil.DrawLine(screen, midx, midy, midx+float64(dy.DX), midy+float64(dy.DY), color.NRGBA{R: 0, G: 255, B: 0, A: 255})
		}
		if *debugShowTransforms {
			midx := float64(screenPos.X) + level.TileSize/2
			midy := float64(screenPos.Y) + level.TileSize/2
			dx := tile.Transform.Apply(m.Delta{DX: 4, DY: 0})
			ebitenutil.DrawLine(screen, midx, midy, midx+float64(dx.DX), midy+float64(dx.DY), color.NRGBA{R: 255, G: 0, B: 0, A: 255})
			dy := tile.Transform.Apply(m.Delta{DX: 0, DY: 4})
			ebitenutil.DrawLine(screen, midx, midy, midx+float64(dy.DX), midy+float64(dy.DY), color.NRGBA{R: 0, G: 255, B: 0, A: 255})
		}
	})
	if *cheatShowBboxes {
		r.world.entities.forEach(func(ent *Entity) error {
			boxColor := color.NRGBA{R: 128, G: 128, B: 128, A: 128}
			if ent.contents.PlayerSolid() {
				boxColor.R = 255
			}
			if ent.contents.ObjectSolid() {
				boxColor.G = 255
			}
			if ent.contents.Opaque() {
				boxColor.B = 255
			}
			ebitenutil.DrawRect(screen, float64(ent.Rect.Origin.X+scrollDelta.DX), float64(ent.Rect.Origin.Y+scrollDelta.DY), float64(ent.Rect.Size.DX), float64(ent.Rect.Size.DY), boxColor)
			if ent.BorderPixels > 0 {
				boxColor.A = 255
				ebitenutil.DrawRect(screen, float64(ent.Rect.Origin.X+scrollDelta.DX-ent.BorderPixels), float64(ent.Rect.Origin.Y+scrollDelta.DY-ent.BorderPixels), float64(ent.Rect.Size.DX+ent.BorderPixels), float64(ent.BorderPixels), boxColor)
				ebitenutil.DrawRect(screen, float64(ent.Rect.Origin.X+scrollDelta.DX+ent.Rect.Size.DX), float64(ent.Rect.Origin.Y+scrollDelta.DY-ent.BorderPixels), float64(ent.BorderPixels), float64(ent.Rect.Size.DY+ent.BorderPixels), boxColor)
				ebitenutil.DrawRect(screen, float64(ent.Rect.Origin.X+scrollDelta.DX), float64(ent.Rect.Origin.Y+scrollDelta.DY+ent.Rect.Size.DY), float64(ent.Rect.Size.DX+ent.BorderPixels), float64(ent.BorderPixels), boxColor)
				ebitenutil.DrawRect(screen, float64(ent.Rect.Origin.X+scrollDelta.DX-ent.BorderPixels), float64(ent.Rect.Origin.Y+scrollDelta.DY), float64(ent.BorderPixels), float64(ent.Rect.Size.DY+ent.BorderPixels), boxColor)
			}
			return nil
		})
	}

	if *debugShowVisiblePolygon {
		texM := ebiten.GeoM{}
		texM.Scale(0, 0)

		adjustedPolygon := make([]m.Pos, len(r.expandedVisiblePolygon))
		for i, pos := range r.expandedVisiblePolygon {
			adjustedPolygon[i] = pos.Add(scrollDelta)
		}
		DrawPolyLine(screen, 3, adjustedPolygon, r.whiteImage, color.NRGBA{R: 255, G: 0, B: 0, A: 255}, &texM, &ebiten.DrawTrianglesOptions{})

		adjustedPolygon = make([]m.Pos, len(r.visiblePolygon))
		for i, pos := range r.visiblePolygon {
			adjustedPolygon[i] = pos.Add(scrollDelta)
		}
		DrawPolyLine(screen, 3, adjustedPolygon, r.whiteImage, color.NRGBA{R: 0, G: 0, B: 255, A: 255}, &texM, &ebiten.DrawTrianglesOptions{})
	}
}

func (r *renderer) rawDrawDest(screen *ebiten.Image) *ebiten.Image {
	if *drawVisibilityMask && *drawOutside {
		return r.offScreenBuffer
	}
	return screen
}

func (r *renderer) drawVisibilityMask(screen, drawDest *ebiten.Image, scrollDelta m.Delta) {
	// Draw trace polygon to buffer.
	geoM := ebiten.GeoM{}
	geoM.Translate(float64(scrollDelta.DX), float64(scrollDelta.DY))
	texM := ebiten.GeoM{}
	texM.Scale(0, 0)

	if *expandUsingVertices && !*expandUsingVerticesAccurately && !*drawBlurs && !*drawOutside {
		drawAntiPolygonAround(screen, r.visiblePolygonCenter, r.expandedVisiblePolygon, r.whiteImage, color.Gray{0}, geoM, texM, &ebiten.DrawTrianglesOptions{})
		return
	}

	if r.needPrevImage {
		// Optimization note:
		// - This isn't optimal. Visibility mask maybe shouldn't even exist?
		// - If screen were a separate image, we could instead copy image to screen masked by polygon.
		// - Would remove one render call.
		// - Wouldn't allow blur though...?
		// Note: we put the mask on ALL four channels.
		r.visibilityMaskImage.Fill(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
		drawPolygonAround(r.visibilityMaskImage, r.visiblePolygonCenter, r.expandedVisiblePolygon, r.whiteImage, color.Gray{255}, geoM, texM, &ebiten.DrawTrianglesOptions{})

		e := expandSize
		if *expandUsingVertices {
			e = 0
		}
		BlurExpandImage(r.visibilityMaskImage, r.blurImage, r.visibilityMaskImage, blurSize, e, 1.0, 0.0)
	}

	if *drawOutside {
		if r.visibilityMaskShader != nil {
			delta := r.world.scrollPos.Delta(r.prevScrollPos)
			screen.DrawRectShader(GameWidth, GameHeight, r.visibilityMaskShader, &ebiten.DrawRectShaderOptions{
				CompositeMode: ebiten.CompositeModeCopy,
				Uniforms: map[string]interface{}{
					"Scroll": []float32{float32(delta.DX) / GameWidth, float32(delta.DY) / GameHeight},
				},
				Images: [4]*ebiten.Image{
					r.visibilityMaskImage,
					drawDest,
					r.prevImage,
					nil,
				},
			})
		} else {
			// First set the alpha channel to the visibility mask.
			drawDest.DrawImage(r.visibilityMaskImage, &ebiten.DrawImageOptions{
				CompositeMode: ebiten.CompositeModeMultiply,
				Filter:        ebiten.FilterNearest,
			})

			// Then draw the background.
			delta := r.world.scrollPos.Delta(r.prevScrollPos)
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
			}, []uint16{0, 1, 2, 1, 2, 3}, r.prevImage, &ebiten.DrawTrianglesOptions{
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
		if r.needPrevImage {
			// Remember last image. Only do this once per update.
			BlurImage(screen, r.blurImage, r.prevImage, frameBlurSize, frameDarkenAlpha, frameDarkenAmount, 1.0)
			r.prevScrollPos = r.world.scrollPos
		}
	} else {
		screen.DrawImage(r.visibilityMaskImage, &ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeDestinationIn,
			Filter:        ebiten.FilterNearest,
		})
	}

	r.needPrevImage = false
}

func (r *renderer) Draw(screen *ebiten.Image) {
	scrollDelta := m.Pos{X: GameWidth / 2, Y: GameHeight / 2}.Delta(r.world.scrollPos)

	dest := r.rawDrawDest(screen)
	dest.Fill(color.Gray{0})
	r.drawTiles(dest, scrollDelta)
	r.drawEntities(dest, scrollDelta)
	if *drawVisibilityMask {
		r.drawVisibilityMask(screen, dest, scrollDelta)
	}
	centerprint.Draw(screen)

	// Debug stuff comes last.
	r.drawDebug(screen, scrollDelta)
}

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
	"github.com/hajimehoshi/ebiten/v2/colorm"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/offscreen"
	"github.com/divVerent/aaaaxy/internal/shader"
)

var (
	drawBlurs = flag.Bool("draw_blurs", true, "perform blur effects; requires draw_visibility_mask")
)

const (
	roundColorToNearest = 1.0 / 510.0
)

func blurPassFixedFunction(img, out *ebiten.Image, mode ebiten.Blend, dx, dy int, scale, addR, addG, addB float64) {
	opts := colorm.DrawImageOptions{
		Blend:  mode,
		Filter: ebiten.FilterNearest,
	}
	var colorM colorm.ColorM
	colorM.Scale(scale, scale, scale, 1)
	colorM.Translate(addR+roundColorToNearest, addG+roundColorToNearest, addB+roundColorToNearest, 0)
	opts.GeoM.Translate(float64(dx), float64(dy))
	colorm.DrawImage(out, img, colorM, &opts)
}

func blurImageFixedFunction(name string, img, out *ebiten.Image, size int, scale, darkenR, darkenG, darkenB, darkenToR, darkenToG, darkenToB float64) {
	// Only power-of-two blurs look good with this approach, so let's scale down the blur as much as needed.
	size++
	for size&(size-1) != 0 {
		size--
	}

	sz := img.Bounds().Size()

	if offscreen.AvoidReuse() {
		src := img // Only in first pass. Otherwise it is a temp image.
		for size > 1 {
			size /= 2
			tmp := offscreen.New(fmt.Sprintf("%s.Horiz.%d", name, size), sz.X, sz.Y)
			tmp.Fill(color.Gray{0})
			blurPassFixedFunction(src, tmp, ebiten.BlendCopy, -size, 0, 0.5, 0, 0, 0)
			blurPassFixedFunction(src, tmp, ebiten.BlendLighter, size, 0, 0.5, 0, 0, 0)
			if src != img {
				// Not first pass.
				offscreen.Dispose(src)
			}
			dst := out
			dstScale := 0.5
			dstAddR, dstAddG, dstAddB := 0.0, 0.0, 0.0
			if size > 1 {
				// Not last pass.
				dst = offscreen.New(fmt.Sprintf("%s.Vert.%d", name, size), sz.X, sz.Y)
			} else {
				// Last pass.
				dstScale *= scale
				dstAddR, dstAddG, dstAddB = (-darkenR+darkenToR*(1-scale))*0.5, (-darkenG+darkenToG*(1-scale))*0.5, (-darkenB+darkenToB*(1-scale))*0.5
				dst = out
			}
			dst.Fill(color.Gray{0})
			blurPassFixedFunction(tmp, dst, ebiten.BlendCopy, 0, -size, dstScale, dstAddR, dstAddG, dstAddB)
			blurPassFixedFunction(tmp, dst, ebiten.BlendLighter, 0, size, dstScale, dstAddR, dstAddG, dstAddB)
			offscreen.Dispose(tmp)
			src = dst
		}
	} else {
		tmp := offscreen.New(name, sz.X, sz.Y)
		src := img
		for size > 1 {
			size /= 2
			tmp.Fill(color.Gray{0})
			blurPassFixedFunction(src, tmp, ebiten.BlendCopy, -size, 0, 0.5, 0, 0, 0)
			blurPassFixedFunction(src, tmp, ebiten.BlendLighter, size, 0, 0.5, 0, 0, 0)
			dstScale := 0.5
			dstAddR, dstAddG, dstAddB := 0.0, 0.0, 0.0
			if size <= 1 {
				dstScale *= scale
				dstAddR, dstAddG, dstAddB = (-darkenR+darkenToR*(1-scale))*0.5, (-darkenG+darkenToG*(1-scale))*0.5, (-darkenB+darkenToB*(1-scale))*0.5
			}
			out.Fill(color.Gray{0})
			blurPassFixedFunction(tmp, out, ebiten.BlendCopy, 0, -size, dstScale, dstAddR, dstAddG, dstAddB)
			blurPassFixedFunction(tmp, out, ebiten.BlendLighter, 0, size, dstScale, dstAddR, dstAddG, dstAddB)
			src = out
		}
		offscreen.Dispose(tmp)
	}
}

func BlurExpandImage(name string, img, out *ebiten.Image, blurSize, expandSize int, scale, darken float64) {
	// Blurring and expanding can be done in a single step by doing a regular blur then scaling up at the last step.
	if !*drawBlurs {
		blurSize = 0
	}
	size := blurSize + expandSize
	scale *= (2*float64(size) + 1) / (2*float64(blurSize) + 1)
	BlurImage(name, img, out, size, scale, darken, color.Gray{0}, 1.0)
}

var (
	blurBroken = false
)

func BlurImage(name string, img, out *ebiten.Image, size int, scale, darken float64, darkenTo color.Color, blurFade float64) {
	sz := img.Bounds().Size()
	scale *= blurFade
	scale += 1 - blurFade
	darken *= blurFade
	darkenToRi, darkenToGi, darkenToBi, _ := darkenTo.RGBA()
	darkenToR, darkenToG, darkenToB := float64(darkenToRi)/65535.0, float64(darkenToGi)/65535.0, float64(darkenToBi)/65535.0
	darkenR, darkenG, darkenB := darken, darken, darken
	if darkenToR >= 0.5 {
		darkenR = -darkenR
	}
	if darkenToG >= 0.5 {
		darkenG = -darkenG
	}
	if darkenToB >= 0.5 {
		darkenB = -darkenB
	}
	if !*drawBlurs && scale <= 1 {
		// Blurs can be globally turned off.
		if img == out {
			if scale == 1.0 && darken == 0.0 {
				return
			}
			copyOptions := &ebiten.DrawImageOptions{
				Blend:  ebiten.BlendCopy,
				Filter: ebiten.FilterNearest,
			}
			tmp := offscreen.New(name, sz.X, sz.Y)
			defer offscreen.Dispose(tmp)
			tmp.DrawImage(img, copyOptions)
			options := &colorm.DrawImageOptions{
				Blend:  ebiten.BlendCopy,
				Filter: ebiten.FilterNearest,
			}
			var colorM colorm.ColorM
			colorM.Scale(scale, scale, scale, 1.0)
			colorM.Translate(-darkenR+darkenToR*(1-scale)+roundColorToNearest, -darkenG+darkenToG*(1-scale)+roundColorToNearest, -darkenB+darkenToB*(1-scale)+roundColorToNearest, 0.0)
			colorm.DrawImage(out, tmp, colorM, options)
		} else {
			options := &colorm.DrawImageOptions{
				Blend:  ebiten.BlendCopy,
				Filter: ebiten.FilterNearest,
			}
			var colorM colorm.ColorM
			colorM.Scale(scale, scale, scale, 1.0)
			colorM.Translate(-darkenR+darkenToR*(1-scale)+roundColorToNearest, -darkenG+darkenToG*(1-scale)+roundColorToNearest, -darkenB+darkenToB*(1-scale)+roundColorToNearest, 0.0)
			colorm.DrawImage(out, img, colorM, options)
		}
		return
	}
	if blurBroken {
		blurImageFixedFunction(name, img, out, size, scale, darkenR, darkenG, darkenB, darkenToR, darkenToG, darkenToB)
		return
	}
	// Too bad we can't have integer uniforms, so we need to templatize this
	// shader instead. Should be faster than having conditionals inside the
	// shader code.
	blurShader, err := shader.Load("blur.kage.tmpl", map[string]string{
		"Size": fmt.Sprint(size),
	})
	if err != nil {
		log.Errorf("BROKEN RENDERER, WILL FALLBACK: could not load blur shader: %v", err)
		blurBroken = true
		blurImageFixedFunction(name, img, out, size, scale, darkenR, darkenG, darkenB, darkenToR, darkenToG, darkenToB)
		return
	}
	centerScale := 1.0 / (2*float64(size)*blurFade + 1)
	otherScale := blurFade * centerScale
	tmp := offscreen.New(fmt.Sprintf("%s.Horiz", name), sz.X, sz.Y)
	defer offscreen.Dispose(tmp)
	tmp.DrawRectShader(sz.X, sz.Y, blurShader, &ebiten.DrawRectShaderOptions{
		Blend: ebiten.BlendCopy,
		Uniforms: map[string]interface{}{
			"Step":        []float32{1, 0},
			"CenterScale": float32(centerScale),
			"OtherScale":  float32(otherScale),
			"Add":         []float32{roundColorToNearest, roundColorToNearest, roundColorToNearest, 0.0},
		},
		Images: [4]*ebiten.Image{
			img,
			nil,
			nil,
			nil,
		},
	})
	out.DrawRectShader(sz.X, sz.Y, blurShader, &ebiten.DrawRectShaderOptions{
		Blend: ebiten.BlendCopy,
		Uniforms: map[string]interface{}{
			"Step":        []float32{0, 1},
			"CenterScale": float32(centerScale * scale),
			"OtherScale":  float32(otherScale * scale),
			"Add":         []float32{float32(-darkenR + darkenToR*(1-scale) + roundColorToNearest), float32(-darkenG + darkenToG*(1-scale) + roundColorToNearest), float32(-darkenB + darkenToB*(1-scale) + roundColorToNearest), 0.0},
		},
		Images: [4]*ebiten.Image{
			tmp,
			nil,
			nil,
			nil,
		},
	})
}

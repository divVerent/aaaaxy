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

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/offscreen"
	"github.com/divVerent/aaaaxy/internal/shader"
)

var (
	drawBlurs = flag.Bool("draw_blurs", true, "perform blur effects; requires draw_visibility_mask")
)

func blurPassFixedFunction(img, out *ebiten.Image, mode ebiten.Blend, dx, dy int, scale, darken float64) {
	opts := ebiten.DrawImageOptions{
		Blend:  mode,
		Filter: ebiten.FilterNearest,
	}
	opts.ColorM.Scale(1, 1, 1, scale)
	opts.ColorM.Translate(-darken, -darken, -darken, 0)
	opts.GeoM.Translate(float64(dx), float64(dy))
	out.DrawImage(img, &opts)
}

func blurImageFixedFunction(name string, img, out *ebiten.Image, size int, scale, darken float64) {
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
			blurPassFixedFunction(src, tmp, ebiten.BlendCopy, -size, 0, 0.5, 0)
			blurPassFixedFunction(src, tmp, ebiten.BlendLighter, size, 0, 0.5, 0)
			if src != img {
				// Not first pass.
				offscreen.Dispose(src)
			}
			dst := out
			dstScale := 0.5
			dstDarken := 0.0
			if size > 1 {
				// Not last pass.
				dst = offscreen.New(fmt.Sprintf("%s.Vert.%d", name, size), sz.X, sz.Y)
			} else {
				// Last pass.
				dstScale *= scale
				dstDarken = darken
				dst = out
			}
			blurPassFixedFunction(tmp, dst, ebiten.BlendCopy, -size, 0, dstScale, dstDarken)
			blurPassFixedFunction(tmp, dst, ebiten.BlendLighter, size, 0, dstScale, dstDarken)
			offscreen.Dispose(tmp)
			src = dst
		}
	} else {
		tmp := offscreen.New(name, sz.X, sz.Y)
		src := img
		for size > 1 {
			size /= 2
			blurPassFixedFunction(src, tmp, ebiten.BlendCopy, -size, 0, 0.5, 0)
			blurPassFixedFunction(src, tmp, ebiten.BlendLighter, size, 0, 0.5, 0)
			dstScale := 0.5
			dstDarken := 0.0
			if size <= 1 {
				dstScale *= scale
				dstDarken = darken
			}
			blurPassFixedFunction(tmp, out, ebiten.BlendCopy, -size, 0, dstScale, dstDarken)
			blurPassFixedFunction(tmp, out, ebiten.BlendLighter, size, 0, dstScale, dstDarken)
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
	BlurImage(name, img, out, size, scale, darken, 1.0)
}

var (
	blurBroken = false
)

func BlurImage(name string, img, out *ebiten.Image, size int, scale, darken, blurFade float64) {
	sz := img.Bounds().Size()
	scale *= scale * blurFade
	scale += 1 - blurFade
	darken *= blurFade
	if !*drawBlurs && scale <= 1 {
		// Blurs can be globally turned off.
		if img == out {
			if scale == 1.0 && darken == 0.0 {
				return
			}
			options := &ebiten.DrawImageOptions{
				Blend:  ebiten.BlendCopy,
				Filter: ebiten.FilterNearest,
			}
			tmp := offscreen.New(name, sz.X, sz.Y)
			defer offscreen.Dispose(tmp)
			tmp.DrawImage(img, options)
			options.ColorM.Scale(scale, scale, scale, 1.0)
			options.ColorM.Translate(-darken, -darken, -darken, 0.0)
			out.DrawImage(tmp, options)
		} else {
			options := &ebiten.DrawImageOptions{
				Blend:  ebiten.BlendCopy,
				Filter: ebiten.FilterNearest,
			}
			options.ColorM.Scale(scale, scale, scale, 1.0)
			options.ColorM.Translate(-darken, -darken, -darken, 0.0)
			out.DrawImage(img, options)
		}
		return
	}
	if blurBroken {
		blurImageFixedFunction(name, img, out, size, scale, darken)
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
		blurImageFixedFunction(name, img, out, size, scale, darken)
		return
	}
	centerScale := 1.0 / (2*float64(size)*blurFade + 1)
	otherScale := blurFade * centerScale
	tmp := offscreen.New(fmt.Sprintf("%s.Horiz", name), sz.X, sz.Y)
	defer offscreen.Dispose(tmp)
	tmp.DrawRectShader(sz.X, sz.Y, blurShader, &ebiten.DrawRectShaderOptions{
		Blend: ebiten.BlendCopy,
		Uniforms: map[string]interface{}{
			"Step":        []float32{1 / float32(sz.X), 0},
			"CenterScale": float32(centerScale),
			"OtherScale":  float32(otherScale),
			"Add":         []float32{float32(-darken), float32(-darken), float32(-darken), 0.0},
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
			"Step":        []float32{0, 1 / float32(sz.Y)},
			"CenterScale": float32(centerScale * scale),
			"OtherScale":  float32(otherScale * scale),
			"Add":         []float32{float32(-darken), float32(-darken), float32(-darken), 0.0},
		},
		Images: [4]*ebiten.Image{
			tmp,
			nil,
			nil,
			nil,
		},
	})
}

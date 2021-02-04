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
	"flag"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/vfs"
)

var (
	debugUseShaders = flag.Bool("debug_use_shaders", true, "enable use of custom shaders")
	drawBlurs       = flag.Bool("draw_blurs", true, "perform blur effects; requires draw_visibility_mask")
)

func blurImageFixedFunction(img, tmp, out *ebiten.Image, size int, weight, scale float64) {
	opts := ebiten.DrawImageOptions{
		CompositeMode: ebiten.CompositeModeLighter,
		Filter:        ebiten.FilterNearest,
	}
	opts.ColorM.Scale(weight, weight, weight, 1)
	size++
	src := img
	for size > 1 {
		size /= 2
		tmp.Fill(color.Gray{0})
		opts.GeoM.Reset()
		opts.GeoM.Translate(-float64(size), 0)
		tmp.DrawImage(src, &opts)
		opts.GeoM.Reset()
		opts.GeoM.Translate(float64(size), 0)
		tmp.DrawImage(src, &opts)
		src = out
		out.Fill(color.Gray{0})
		if size <= 1 {
			opts.ColorM.Scale(scale, scale, scale, 1)
		}
		opts.GeoM.Reset()
		opts.GeoM.Translate(0, -float64(size))
		out.DrawImage(tmp, &opts)
		opts.GeoM.Reset()
		opts.GeoM.Translate(0, float64(size))
		out.DrawImage(tmp, &opts)
	}
}

func loadShader(name string) (*ebiten.Shader, error) {
	shaderReader, err := vfs.Load("shaders", name)
	if err != nil {
		return nil, fmt.Errorf("could not open shader %q: %v", name, err)
	}
	defer shaderReader.Close()
	shaderCode, err := ioutil.ReadAll(shaderReader)
	if err != nil {
		return nil, fmt.Errorf("could not read shader %q: %v", name, err)
	}
	shader, err := ebiten.NewShader(shaderCode)
	if err != nil {
		return nil, fmt.Errorf("could not compile shader %q: %v", name, err)
	}
	return shader, nil
}

var (
	blurShader *ebiten.Shader
)

func BlurExpandImage(img, tmp, out *ebiten.Image, blurSize, expandSize int, scale float64) {
	// Blurring and expanding can be done in a single step by doing a regular blur then scaling up at the last step.
	size := blurSize + expandSize
	scale *= (2*float64(size) + 1) / (2*float64(blurSize) + 1)
	BlurImage(img, tmp, out, size, scale)
}

func BlurImage(img, tmp, out *ebiten.Image, size int, scale float64) {
	if !*drawBlurs && scale <= 1 {
		// Blurs can be globally turned off.
		if img != out {
			options := &ebiten.DrawImageOptions{
				CompositeMode: ebiten.CompositeModeCopy,
				Filter:        ebiten.FilterNearest,
			}
			options.ColorM.Scale(scale, scale, scale, 1.0)
			out.DrawImage(img, options)
		}
		return
	}
	if !*debugUseShaders {
		blurImageFixedFunction(img, tmp, out, size, 0.5, scale)
		return
	}
	if blurShader == nil {
		var err error
		blurShader, err = loadShader("blur.go")
		if err != nil {
			log.Panicf("could not load blur shader: %v", err)
		}
	}
	w, h := img.Size()
	scale = math.Sqrt(scale) / (2*float64(size) + 1)
	tmp.DrawRectShader(w, h, blurShader, &ebiten.DrawRectShaderOptions{
		CompositeMode: ebiten.CompositeModeCopy,
		Uniforms: map[string]interface{}{
			"Size":  float32(size),
			"Step":  []float32{1 / float32(w), 0},
			"Scale": float32(scale),
		},
		Images: [4]*ebiten.Image{
			img,
			nil,
			nil,
			nil,
		},
	})
	out.DrawRectShader(w, h, blurShader, &ebiten.DrawRectShaderOptions{
		CompositeMode: ebiten.CompositeModeCopy,
		Uniforms: map[string]interface{}{
			"Size":  float32(size),
			"Step":  []float32{0, 1 / float32(h)},
			"Scale": float32(scale),
		},
		Images: [4]*ebiten.Image{
			tmp,
			nil,
			nil,
			nil,
		},
	})
}

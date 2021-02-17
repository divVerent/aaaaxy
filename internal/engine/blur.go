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
	"bytes"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/flag"
	"github.com/divVerent/aaaaaa/internal/vfs"
)

var (
	debugUseShaders = flag.Bool("debug_use_shaders", true, "enable use of custom shaders")
	drawBlurs       = flag.Bool("draw_blurs", true, "perform blur effects; requires draw_visibility_mask")
)

func blurImageFixedFunction(img, tmp, out *ebiten.Image, size int, scale float64) {
	opts := ebiten.DrawImageOptions{
		CompositeMode: ebiten.CompositeModeLighter,
		Filter:        ebiten.FilterNearest,
	}
	size++
	// Only power-of-two blurs look good with this approach, so let's scale down the blur as much as needed.
	for size&(size-1) != 0 {
		size--
	}
	src := img
	opts.ColorM.Scale(1, 1, 1, 0.5)
	for size > 1 {
		size /= 2
		tmp.Fill(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
		opts.CompositeMode = ebiten.CompositeModeCopy
		opts.GeoM.Reset()
		opts.GeoM.Translate(-float64(size), 0)
		tmp.DrawImage(src, &opts)
		opts.CompositeMode = ebiten.CompositeModeLighter
		opts.GeoM.Reset()
		opts.GeoM.Translate(float64(size), 0)
		tmp.DrawImage(src, &opts)
		src = out
		if size <= 1 {
			opts.ColorM.Scale(1, 1, 1, scale)
		}
		out.Fill(color.NRGBA{R: 0, G: 0, B: 0, A: 0})
		opts.CompositeMode = ebiten.CompositeModeCopy
		opts.GeoM.Reset()
		opts.GeoM.Translate(0, -float64(size))
		out.DrawImage(tmp, &opts)
		opts.CompositeMode = ebiten.CompositeModeLighter
		opts.GeoM.Reset()
		opts.GeoM.Translate(0, float64(size))
		out.DrawImage(tmp, &opts)
	}
}

func loadShader(name string, params map[string]string) (*ebiten.Shader, error) {
	shaderReader, err := vfs.Load("shaders", name)
	if err != nil {
		return nil, fmt.Errorf("could not open shader %q: %v", name, err)
	}
	defer shaderReader.Close()
	shaderCode, err := ioutil.ReadAll(shaderReader)
	if err != nil {
		return nil, fmt.Errorf("could not read shader %q: %v", name, err)
	}
	// Add some basic templating so we can remove branches from the shaders.
	// Not using text/template so that shader files can still be processed by gofmt.
	for name, value := range params {
		shaderCode = bytes.ReplaceAll(shaderCode, []byte("PARAMS[\""+name+"\"]"), []byte("(("+value+"))"))
	}
	shader, err := ebiten.NewShader(shaderCode)
	if err != nil {
		return nil, fmt.Errorf("could not compile shader %q: %v", name, err)
	}
	return shader, nil
}

var (
	blurShaders = map[int]*ebiten.Shader{}
)

func BlurExpandImage(img, tmp, out *ebiten.Image, blurSize, expandSize int, scale float64) {
	// Blurring and expanding can be done in a single step by doing a regular blur then scaling up at the last step.
	if !*drawBlurs {
		blurSize = 0
	}
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
		blurImageFixedFunction(img, tmp, out, size, scale)
		return
	}
	blurShader := blurShaders[size]
	if blurShader == nil {
		var err error
		// Too bad we can't have integer uniforms, so we need to templatize this
		// shader instead. Should be faster than having conditionals inside the
		// shader code.
		blurShader, err = loadShader("blur.kage", map[string]string{
			"Size": fmt.Sprint(size),
		})
		if err != nil {
			log.Panicf("could not load blur shader: %v", err)
		}
		blurShaders[size] = blurShader
	}
	w, h := img.Size()
	scaleX := 1 / (2*float64(size) + 1)
	tmp.DrawRectShader(w, h, blurShader, &ebiten.DrawRectShaderOptions{
		CompositeMode: ebiten.CompositeModeCopy,
		Uniforms: map[string]interface{}{
			"Step":  []float32{1 / float32(w), 0},
			"Scale": float32(scaleX),
		},
		Images: [4]*ebiten.Image{
			img,
			nil,
			nil,
			nil,
		},
	})
	scaleY := scale / (2*float64(size) + 1)
	out.DrawRectShader(w, h, blurShader, &ebiten.DrawRectShaderOptions{
		CompositeMode: ebiten.CompositeModeCopy,
		Uniforms: map[string]interface{}{
			"Size":  float32(size),
			"Step":  []float32{0, 1 / float32(h)},
			"Scale": float32(scaleY),
		},
		Images: [4]*ebiten.Image{
			tmp,
			nil,
			nil,
			nil,
		},
	})
}

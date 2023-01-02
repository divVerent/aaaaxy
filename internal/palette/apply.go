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

package palette

import (
	"errors"
	"fmt"
	"image"
	"image/color"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	debugCheckImagePalette = flag.Bool("debug_check_image_palette", false, "log errors if images or object colors mismatch palette colors")
)

var (
	remapping bool = false
)

// SetCurrent changes the current palette. Returns whether the remapping table changed.
func SetCurrent(pal *Palette, doRemap bool) bool {
	if pal == current {
		return false
	}
	var prevRemap map[uint32]uint32
	if current != nil && remapping {
		prevRemap = current.remap
	}
	if pal != nil && len(pal.remap) != 0 {
		log.Infof("note: remapping %d colors (slow)", len(pal.remap))
	}
	current = pal
	remapping = doRemap
	if current == nil {
		if len(prevRemap) != 0 {
			return true
		}
	} else {
		var newRemap map[uint32]uint32
		if doRemap {
			newRemap = current.remap
		}
		if len(newRemap) != len(prevRemap) {
			return true
		}
		for from, to := range newRemap {
			if prevTo, found := prevRemap[from]; !found || prevTo != to {
				return true
			}
		}
	}
	return false
}

func Current() *Palette {
	return current
}

// ApplyToImage applies this palette's associated color remapping to an image.
func (p *Palette) ApplyToImage(img image.Image, name string) image.Image {
	if p == nil || len(p.remap) == 0 || !remapping {
		return img
	}
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			rgba := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
			newImg.SetRGBA(x, y, p.ApplyToRGBA(rgba, name))
		}
	}
	return newImg
}

// ApplyToRGBA applies this palette's associated col remapping to a single col.
func (p *Palette) ApplyToRGBA(col color.RGBA, name string) color.RGBA {
	if *debugCheckImagePalette {
		rgb := (uint32(col.R) << 16) | (uint32(col.G) << 8) | uint32(col.B)
		if !egaColorsSet[rgb] {
			log.Warningf("invalid color %v in %v: color is not in palette", col, name)
		}
		if col.A != 255 && col.A != 0 {
			log.Warningf("invalid color %v in %v: premultiplied color is neither fully opaque nor fully transparent", col, name)
		}
	}
	if p == nil || len(p.remap) == 0 || !remapping {
		return col
	}
	// Color is premultiplied - can't handle that well.
	// So for now, only remap if alpha is 255.
	if col.A != 255 {
		return col
	}
	// Remap rgb.
	rgb := (uint32(col.R) << 16) | (uint32(col.G) << 8) | uint32(col.B)
	if *debugCheckImagePalette && !egaColorsSet[rgb] {
		log.Warningf("invalid color in %v: color is not in palette", name)
	}
	if rgbM, found := p.remap[rgb]; found {
		return toRGB(rgbM).toRGBA()
	}
	return col
}

// ApplyToNRGBA applies this palette's associated col remapping to a single col.
func (p *Palette) ApplyToNRGBA(col color.NRGBA, name string) color.NRGBA {
	if *debugCheckImagePalette {
		rgb := (uint32(col.R) << 16) | (uint32(col.G) << 8) | uint32(col.B)
		if !egaColorsSet[rgb] {
			log.Warningf("invalid color in %v: color is not in palette", name)
		}
	}
	if p == nil || len(p.remap) == 0 || !remapping {
		return col
	}
	// Remap rgb.
	rgb := (uint32(col.R) << 16) | (uint32(col.G) << 8) | uint32(col.B)
	if rgbM, found := p.remap[rgb]; found {
		colM := toRGB(rgbM).toNRGBA()
		// Preserve alpha.
		colM.A = col.A
		return colM
	}
	return col
}

// rawEGA gets the named EGA color.
func (p *Palette) rawEGA(i EGAIndex) uint32 {
	if p == nil {
		return egaColors[i]
	}
	return p.ega[i]
}

func EGA(i EGAIndex, a uint8) color.NRGBA {
	u := current.rawEGA(i)
	return color.NRGBA{
		R: uint8(u >> 16),
		G: uint8((u >> 8) & 0xFF),
		B: uint8(u & 0xFF),
		A: a,
	}
}

func Parse(s string, name string) (color.NRGBA, error) {
	if s == "" {
		return color.NRGBA{}, errors.New("no color specified")
	}
	// Trailing ! means "do not map according to palette".
	p := len(s) - 1
	doApply := s[p] != '!'
	if !doApply {
		s = s[:p]
	}
	var r, g, b, a uint8
	if _, err := fmt.Sscanf(s, "#%02x%02x%02x%02x", &a, &r, &g, &b); err != nil {
		return color.NRGBA{}, err
	}
	c := color.NRGBA{R: r, G: g, B: b, A: a}
	if doApply {
		return current.ApplyToNRGBA(c, name), nil
	}
	return c, nil
}

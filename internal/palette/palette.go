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
	"fmt"
	"image"
	"image/color"
	"sort"
	"strings"

	"github.com/divVerent/aaaaxy/internal/log"
)

// Palette encapsulates a color palette.
type Palette struct {
	// size is the number of colors this palette has. Is > 0 for any valid palette.
	size int

	// egaIndices is a list with one entry per protected color whose value is the EGA color index.
	egaIndices []int

	// protected is the number of protected colors.
	// This is also used to compute the Bayer pattern size.
	protected int

	// colors are the palette colors.
	colors []uint32

	// remap is the color remapping.
	remap map[uint32]uint32

	// ega is the set of EGA colors after remapping.
	ega [EGACount]uint32
}

var current *Palette

func newPalette(egaIndices []int, c0 []uint32) *Palette {
	// Keep only unique colors beyond egaIndices.
	h := make(map[uint32]struct{}, len(c0))
	c := make([]uint32, 0, len(c0))
	for i, u := range c0 {
		if _, found := h[u]; found && i >= len(egaIndices) {
			continue
		}
		h[u] = struct{}{}
		c = append(c, u)
	}

	protected := len(egaIndices)
	if protected == 0 {
		protected = len(c)
	}
	ega := egaColors
	remap := map[uint32]uint32{}
	for thisIdx, egaIdx := range egaIndices {
		from := toRGB(egaColors[egaIdx]).toUint32()
		to := toRGB(c[thisIdx]).toUint32()
		if from != to {
			remap[from] = to
			ega[egaIdx] = to
		}
	}
	if len(remap) == 0 {
		remap = nil
	}
	pal := &Palette{
		size:       len(c),
		egaIndices: egaIndices,
		protected:  protected,
		colors:     c,
		remap:      remap,
		ega:        ega,
	}
	return pal
}

// Names returns the names of all palettes, in quoted comma separated for, for inclusion in a flag description.
func Names() string {
	l := make([]string, 0, len(data))
	for p := range data {
		l = append(l, p)
	}
	sort.Strings(l)
	return "'" + strings.Join(l, "', '") + "'"
}

// ByName returns the PalData for the given palette. Do not modify the returned object.
func ByName(name string) *Palette {
	return data[name]
}

// ApplyToImage applies this palette's associated color remapping to an image.
func (p *Palette) ApplyToImage(img image.Image) image.Image {
	if p == nil || len(p.remap) == 0 {
		return img
	}
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			rgba := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
			newImg.SetRGBA(x, y, p.ApplyToRGBA(rgba))
		}
	}
	return newImg
}

// ApplyToRGBA applies this palette's associated col remapping to a single col.
func (p *Palette) ApplyToRGBA(col color.RGBA) color.RGBA {
	if p == nil || len(p.remap) == 0 {
		return col
	}
	// Color is premultiplied - can't handle that well.
	// So for now, only remap if alpha is 255.
	if col.A != 255 {
		return col
	}
	// Remap rgb.
	rgb := (uint32(col.R) << 16) | (uint32(col.G) << 8) | uint32(col.B)
	if rgbM, found := p.remap[rgb]; found {
		return toRGB(rgbM).toRGBA()
	}
	return col
}

// ApplyToNRGBA applies this palette's associated col remapping to a single col.
func (p *Palette) ApplyToNRGBA(col color.NRGBA) color.NRGBA {
	if p == nil || len(p.remap) == 0 {
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

// SetCurrent changes the current palette. Returns whether the remapping table changed.
func SetCurrent(pal *Palette) bool {
	if pal == current {
		return false
	}
	var prevRemap map[uint32]uint32
	if current != nil {
		prevRemap = current.remap
	}
	if pal != nil && len(pal.remap) != 0 {
		log.Infof("note: remapping %d colors (slow)", len(pal.remap))
	}
	current = pal
	if current == nil {
		if len(prevRemap) != 0 {
			return true
		}
	} else {
		if len(current.remap) != len(prevRemap) {
			return true
		}
		for from, to := range current.remap {
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

func EGA(i EGAIndex, a uint8) color.NRGBA {
	u := current.rawEGA(i)
	return color.NRGBA{
		R: uint8(u >> 16),
		G: uint8((u >> 8) & 0xFF),
		B: uint8(u & 0xFF),
		A: a,
	}
}

func Parse(s string) (color.NRGBA, error) {
	var r, g, b, a uint8
	if _, err := fmt.Sscanf(s, "#%02x%02x%02x%02x", &a, &r, &g, &b); err != nil {
		return color.NRGBA{}, err
	}
	return current.ApplyToNRGBA(color.NRGBA{R: r, G: g, B: b, A: a}), nil
}

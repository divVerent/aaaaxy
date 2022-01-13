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

package aaaaxy

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
	m "github.com/divVerent/aaaaxy/internal/math"
)

var (
	paletteColordist = flag.String("palette_colordist", "redmean", "color distance function to use; one of 'weighted', 'redmean'")
)

type rgb [3]float32

func (c rgb) diff(other rgb) float32 {
	switch *paletteColordist {
	case "weighted":
		dr := c[0] - other[0]
		dg := c[1] - other[1]
		db := c[2] - other[2]
		return 0.3*dr*dr + 0.59*dg*dg + 0.11*db*db
	case "redmean":
		dr := c[0] - other[0]
		dg := c[1] - other[1]
		db := c[2] - other[2]
		rr := (c[0] + other[0]) / 2 * 255.0 / 256.0 // Some inaccuracy to match Wikipedia.
		return (2+rr)*dr*dr + 4*dg*dg + (3-rr)*db*db
	default:
		*paletteColordist = "redmean"
		return c.diff(other)
	}
}

func (c rgb) toColor() color.Color {
	return color.NRGBA{
		R: uint8(m.Rint(float64(c[0] * 255))),
		G: uint8(m.Rint(float64(c[1] * 255))),
		B: uint8(m.Rint(float64(c[2] * 255))),
		A: 255,
	}
}

func (p *palData) lookup(i int) rgb {
	return rgb{
		p.colors[3*i],
		p.colors[3*i+1],
		p.colors[3*i+2],
	}
}

func (p *palData) lookupNearest(c rgb) int {
	bestI := 0
	bestS := c.diff(p.lookup(0))
	for i := 1; i < p.size; i++ {
		s := c.diff(p.lookup(i))
		if s < bestS {
			bestI, bestS = i, s
		}
	}
	return bestI
}

func (p *palData) toLUT(img *ebiten.Image) (int, int) {
	bounds := img.Bounds()
	ox := bounds.Min.X
	oy := bounds.Min.Y
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y
	lutSize := int(math.Cbrt(float64(w) * float64(h)))
	var perRow int
	for {
		perRow = w / lutSize
		rows := (lutSize + perRow - 1) / perRow
		heightNeeded := rows * lutSize
		if heightNeeded <= h {
			break
		}
		lutSize--
	}
	for r := 0; r < lutSize; r++ {
		for g := 0; g < lutSize; g++ {
			for b := 0; b < lutSize; b++ {
				x := ox + (b%perRow)*lutSize + r
				y := oy + (b/perRow)*lutSize + g
				c := rgb{
					(float32(r) + 0.5) / float32(lutSize),
					(float32(g) + 0.5) / float32(lutSize),
					(float32(b) + 0.5) / float32(lutSize),
				}
				i := p.lookupNearest(c)
				cNew := p.lookup(i)
				img.Set(x, y, cNew.toColor())
			}
		}
	}
	return lutSize, perRow
}

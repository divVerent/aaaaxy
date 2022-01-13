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
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	paletteColordist = flag.String("palette_colordist", "redmean", "color distance function to use; one of 'weighted', 'redmean'")
)

type rgb [3]uint8

func (c rgb) diff(other rgb) float64 {
	switch *paletteColordist {
	case "weighted":
		dr := float64(c[0]) - float64(other[0])
		dg := float64(c[1]) - float64(other[1])
		db := float64(c[2]) - float64(other[2])
		return 0.3*dr*dr + 0.59*dg*dg + 0.11*db*db
	case "redmean":
		dr := float64(c[0]) - float64(other[0])
		dg := float64(c[1]) - float64(other[1])
		db := float64(c[2]) - float64(other[2])
		rr := (float64(c[0]) + float64(other[0])) / 2
		return (2+rr/256)*dr*dr + 4*dg*dg + (2+(255-rr)/256)*db*db
	default:
		*paletteColordist = "redmean"
		return c.diff(other)
	}
}

func (c rgb) toColor() color.Color {
	return color.NRGBA{
		R: uint8(c[0]),
		G: uint8(c[1]),
		B: uint8(c[2]),
		A: 255,
	}
}

func (p *palData) lookup(i int) rgb {
	u := p.colors[i]
	return rgb{
		uint8(u >> 16),
		uint8((u >> 8) & 0xFF),
		uint8(u & 0xFF),
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
	defer func(t0 time.Time) {
		dt := time.Since(t0)
		log.Infof("building palette LUT took %v", dt)
	}(time.Now())
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
					uint8(255*(float64(r)+0.5)/float64(lutSize) + 0.5),
					uint8(255*(float64(g)+0.5)/float64(lutSize) + 0.5),
					uint8(255*(float64(b)+0.5)/float64(lutSize) + 0.5),
				}
				i := p.lookupNearest(c)
				cNew := p.lookup(i)
				img.Set(x, y, cNew.toColor())
			}
		}
	}
	return lutSize, perRow
}

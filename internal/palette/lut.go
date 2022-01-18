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
	"image"
	"image/color"
	"math"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	paletteColordist = flag.String("palette_colordist", "redmean", "color distance function to use; one of 'weighted', 'redmean'")
)

type rgb [3]float64 // Actually integers from 0 to 255, but storing as float64 allows fastest math.

func (c rgb) diff(other rgb) float64 {
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
		rr := (c[0] + other[0]) / 2
		return (2+rr)*dr*dr + 4*dg*dg + (2+255/256.0-rr)*db*db
	default:
		*paletteColordist = "redmean"
		return c.diff(other)
	}
}

func (c rgb) toColor() color.NRGBA {
	return color.NRGBA{
		R: uint8(c[0]*255 + 0.5),
		G: uint8(c[1]*255 + 0.5),
		B: uint8(c[2]*255 + 0.5),
		A: 255,
	}
}

func (p *Palette) lookup(i int) rgb {
	u := p.colors[i]
	return rgb{
		float64(u>>16) / 255,
		float64((u>>8)&0xFF) / 255,
		float64(u&0xFF) / 255,
	}
}

// lookupNearest returns the palette color nearest to c.
// hint should be the palette index of a "nearby" color, if possible.
// Providing this helps branch prediction a LOT (on my computer, 770ms -> 440ms for the VGA palette).
func (p *Palette) lookupNearest(c rgb, hint int) int {
	bestI := hint
	bestS := c.diff(p.lookup(hint))
	for i := 0; i < p.size; i++ {
		s := c.diff(p.lookup(i))
		if s < bestS {
			bestI, bestS = i, s
		}
	}
	return bestI
}

func (p *Palette) ToLUT(img *ebiten.Image) (int, int) {
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
	var perRow, heightNeeded, widthNeeded int
	for {
		perRow = w / lutSize
		widthNeeded = perRow * lutSize
		rows := (lutSize + perRow - 1) / perRow
		heightNeeded = rows * lutSize
		if heightNeeded <= h {
			break
		}
		lutSize--
	}
	// Note: creating a temp image, and copying to that, so this does not invoke
	// thread synchronization as writing to an ebiten.Image would.
	rect := image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: widthNeeded,
			Y: heightNeeded,
		},
	}
	tmp := image.NewNRGBA(rect)
	var wg sync.WaitGroup
	for y := 0; y < heightNeeded; y++ {
		wg.Add(1)
		go func(y int) {
			g := y % lutSize
			gFloat := (float64(g) + 0.5) / float64(lutSize)
			bY := (y / lutSize) * perRow
			i := 0
			for x := 0; x < widthNeeded; x++ {
				r := x % lutSize
				rFloat := (float64(r) + 0.5) / float64(lutSize)
				b := bY + x/lutSize
				bFloat := (float64(b) + 0.5) / float64(lutSize)
				c := rgb{rFloat, gFloat, bFloat}
				i = p.lookupNearest(c, i)
				cNew := p.lookup(i)
				tmp.SetNRGBA(x+ox, y+oy, cNew.toColor())
			}
			wg.Done()
		}(y)
	}
	wg.Wait()
	img.SubImage(rect).(*ebiten.Image).ReplacePixels(tmp.Pix)
	return lutSize, perRow
}

// BayerPattern computes the Bayer pattern for this palette.
func (p *Palette) BayerPattern(size int) []float32 {
	// New palette also needs new Bayer pattern.
	sizeSquare := size * size
	bits := 0
	if size > 1 {
		bits = math.Ilogb(float64(size-1)) + 1
	}
	sizeCeil := 1 << bits
	sizeCeilSquare := sizeCeil * sizeCeil
	scale := p.bayerScale / float64(sizeCeilSquare)
	offset := float64(sizeCeilSquare-1) / 2.0
	bayern := make([]float32, sizeSquare)
	for i := range bayern {
		x := i % size
		y := i / size
		z := x ^ y
		b := 0
		for bit := 1; bit < size; bit *= 2 {
			b *= 4
			if y&bit != 0 {
				b += 1
			}
			if z&bit != 0 {
				b += 2
			}
		}
		bayern[i] = float32((float64(b) - offset) * scale)
	}
	return bayern
}

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
	"math"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	paletteColordist = flag.String("palette_colordist", "redmean", "color distance function to use; one of 'weighted', 'redmean', 'cielab', 'cieluv'")
)

type rgb [3]float64 // Range is from 0 to 1 in sRGB color space.

func (c rgb) String() string {
	n := c.toNRGBA()
	return fmt.Sprintf("#%02X%02X%02X", n.R, n.G, n.B)
}

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
	case "cielab":
		return c.toColorful().DistanceLab(other.toColorful())
	case "cieluv":
		return c.toColorful().DistanceLuv(other.toColorful())
	default:
		*paletteColordist = "redmean"
		return c.diff(other)
	}
}

func (c rgb) toNRGBA() color.NRGBA {
	return color.NRGBA{
		R: uint8(c[0]*255 + 0.5),
		G: uint8(c[1]*255 + 0.5),
		B: uint8(c[2]*255 + 0.5),
		A: 255,
	}
}

func (c rgb) toColorful() colorful.Color {
	return colorful.Color{
		R: c[0],
		G: c[1],
		B: c[2],
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
func (p *Palette) lookupNearest(c rgb) int {
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
			for x := 0; x < widthNeeded; x++ {
				r := x % lutSize
				rFloat := (float64(r) + 0.5) / float64(lutSize)
				b := bY + x/lutSize
				bFloat := (float64(b) + 0.5) / float64(lutSize)
				c := rgb{rFloat, gFloat, bFloat}
				i := p.lookupNearest(c)
				cNew := p.lookup(i)
				tmp.SetNRGBA(x+ox, y+oy, cNew.toNRGBA())
			}
			wg.Done()
		}(y)
	}
	wg.Wait()
	img.SubImage(rect).(*ebiten.Image).ReplacePixels(tmp.Pix)

	return lutSize, perRow
}

func sizeBayer(size int, bayerScale float64) (sizeSquare, sizeCeilSquare int, scale, offset float64) {
	sizeSquare = size * size
	bits := 0
	if size > 1 {
		bits = math.Ilogb(float64(size-1)) + 1
	}
	sizeCeil := 1 << bits
	sizeCeilSquare = sizeCeil * sizeCeil
	scale = bayerScale / float64(sizeCeilSquare)
	offset = float64(sizeCeilSquare-1) / 2.0
	return
}

func clamp(a, mi, ma float64) float64 {
	if a < mi {
		return mi
	}
	if a > ma {
		return ma
	}
	return a
}

func (p *Palette) CheckProtectedColors(base *Palette, lutSize, bayerSize int) {
	defer func(t0 time.Time) {
		dt := time.Since(t0)
		log.Infof("checking palette LUT took %v", dt)
	}(time.Now())
	_, sizeCeilSquare, scale, offset := sizeBayer(bayerSize, p.bayerScale)
	for ref := 0; ref < base.size; ref++ {
		cRef := base.lookup(ref)
		for b := 0; b < sizeCeilSquare; b++ {
			shift := (float64(b) - offset) * scale
			cLuttered := rgb{
				clamp(math.Floor((cRef[0]+shift)*float64(lutSize)), 0, float64(lutSize-1)),
				clamp(math.Floor((cRef[1]+shift)*float64(lutSize)), 0, float64(lutSize-1)),
				clamp(math.Floor((cRef[2]+shift)*float64(lutSize)), 0, float64(lutSize-1)),
			}
			cShifted := rgb{
				(cLuttered[0] + 0.5) / float64(lutSize),
				(cLuttered[1] + 0.5) / float64(lutSize),
				(cLuttered[2] + 0.5) / float64(lutSize),
			}
			j := p.lookupNearest(cShifted)
			cNew := p.lookup(j)
			if cNew != cRef {
				log.Warningf("protected color got broken: %v (%v) + %v = %v -> %v (%v)",
					cRef, ref, shift, cShifted, cNew, j)
			}
		}
	}
}

// BayerPattern computes the Bayer pattern for this palette.
func (p *Palette) BayerPattern(size int) []float32 {
	sizeSquare, _, scale, offset := sizeBayer(size, p.bayerScale)
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

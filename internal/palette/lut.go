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
	"sort"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
)

var (
	paletteColordist             = flag.String("palette_colordist", "weighted", "color distance function to use; one of 'weighted', 'weightedL', 'rgbL', 'redmean', 'cielab', 'cieluv'")
	paletteMaxCycles             = flag.Float64("palette_max_cycles", 640*360*128, "maximum number of cycles to spend on palette generation")
	palettePsychovisualFactor    = flag.Float64("palette_psychovisual_factor", 0.02, "factor by which to include the psychovisual model when generating a two-color palette LUT")
	palettePsychovisualDampening = flag.Float64("palette_psychovisual_dampening", 0.5, "factor by which to dampen the psychovisual model when mixing evenly")
)

type rgb [3]float64 // Range is from 0 to 1 in sRGB color space.

func (c rgb) String() string {
	n := c.toNRGBA()
	return fmt.Sprintf("#%02X%02X%02X", n.R, n.G, n.B)
}

func (c rgb) mix(other rgb, f float64) rgb {
	return rgb{
		c[0] + (other[0]-c[0])*f,
		c[1] + (other[1]-c[1])*f,
		c[2] + (other[2]-c[2])*f,
	}
}

func (c rgb) computeF(c0, c1 rgb) float64 {
	// See computeF in the shader.
	ur := c[0] - c0[0]
	ug := c[1] - c0[1]
	ub := c[2] - c0[2]
	vr := c1[0] - c0[0]
	vg := c1[1] - c0[1]
	vb := c1[2] - c0[2]
	duv := ur*vr*3 + ug*vg*4 + ub*vb*2
	dvv := vr*vr*3 + vg*vg*4 + vb*vb*2
	return duv / dvv
}

func (c rgb) diff2(other rgb) float64 {
	switch *paletteColordist {
	case "weighted":
		dr := c[0] - other[0]
		dg := c[1] - other[1]
		db := c[2] - other[2]
		return 3*dr*dr + 4*dg*dg + 2*db*db
	case "weightedL":
		// Adapted from https: //bisqwit.iki.fi/story/howto/dither/jy/#PsychovisualModel
		dr := c[0] - other[0]
		dg := c[1] - other[1]
		db := c[2] - other[2]
		dl := 3*dr + 4*dg + 2*db
		return 3*dr*dr + 4*dg*dg + 2*db*db + 13*dl*dl
	case "rgbL":
		// Directly from https: //bisqwit.iki.fi/story/howto/dither/jy/#PsychovisualModel
		dr := c[0] - other[0]
		dg := c[1] - other[1]
		db := c[2] - other[2]
		dl := 0.299*dr + 0.587*dg + 0.114*db
		return (0.299*dr*dr+0.587*dg*dg+0.114*db*db)*0.75 + dl*dl
	case "redmean":
		dr := c[0] - other[0]
		dg := c[1] - other[1]
		db := c[2] - other[2]
		rr := (c[0] + other[0]) / 2
		return (2+rr)*dr*dr + 4*dg*dg + (2+255/256.0-rr)*db*db
	case "cielab":
		return math.Pow(c.toColorful().DistanceLab(other.toColorful()), 2)
	case "cieluv":
		return math.Pow(c.toColorful().DistanceLuv(other.toColorful()), 2)
	default:
		*paletteColordist = "weighted"
		return c.diff2(other)
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
	bestS := c.diff2(p.lookup(0))
	for i := 1; i < p.size; i++ {
		s := c.diff2(p.lookup(i))
		if s < bestS {
			bestI, bestS = i, s
		}
	}
	return bestI
}

func (p *Palette) nearestTwoChoices() (int, int, int) {
	nI := p.protected
	if nI <= 0 || nI >= p.size {
		nI = p.size - 1
	}
	nJMin := p.size - nI
	nJMax := p.size - 1
	return nI, nJMin, nJMax
}

// lookupNearestTwo returns the pair of distinct palette colors nearest to c.
func (p *Palette) lookupNearestTwo(c rgb) (int, int) {
	bestI := 0
	bestJ := 0
	bestS := math.Inf(+1)
	// TODO different protect logic: if c is "near" a protected color (i.e. if a protected color maps to c - check at caller and pass in here), force i to be that protected color and j can be any other.
	// Otherwise, pick at will.
	// Swap if backwards ordered, though!
	// That will mean we need to consider all options in the cycle count again, though.
	nI, _, _ := p.nearestTwoChoices()
	for i := 0; i < nI; i++ {
		for j := i + 1; j < p.size; j++ {
			c0 := p.lookup(i)
			c1 := p.lookup(j)
			f := c.computeF(c0, c1)
			if f < 0 {
				f = 0
			}
			if f > 1 {
				f = 1
			}
			c_ := c0.mix(c1, f)
			// Including c0.diff2(c1) as per https://bisqwit.iki.fi/story/howto/dither/jy/#PsychovisualModel
			// We seem to need a lower factor for this game's content though.
			s := c_.diff2(c) + *palettePsychovisualFactor*c0.diff2(c1)*(1.0-*palettePsychovisualDampening*(1.0-2.0*math.Abs(f-0.5)))
			if s < bestS {
				bestI, bestJ, bestS = i, j, s
			}
		}
	}
	return bestI, bestJ
}

func computeLUTSize(w, h int, maxEntries float64) (int, int, int) {
	pixels := float64(w * h)
	if pixels < maxEntries {
		maxEntries = pixels
	}
	size := int(math.Cbrt(maxEntries))
	var perRow, rowCount int
	for size > 0 {
		perRow = w / size
		rowCount = (size + perRow - 1) / perRow
		heightNeeded := rowCount * size
		if heightNeeded <= h {
			break
		}
		// Can just brute force the best size, we're dealing with low numbers here in the first place.
		size--
	}
	return size, perRow, rowCount
}

func (p *Palette) computeNearestLUT(lutSize, perRow, lutWidth, lutHeight, lutStride int, pix []byte) {
	var wg sync.WaitGroup
	for y := 0; y < lutHeight; y++ {
		wg.Add(1)
		go func(y int) {
			g := y % lutSize
			gFloat := (float64(g) + 0.5) / float64(lutSize)
			bY := (y / lutSize) * perRow
			for x := 0; x < lutWidth; x++ {
				r := x % lutSize
				rFloat := (float64(r) + 0.5) / float64(lutSize)
				b := bY + x/lutSize
				if b >= lutSize {
					break
				}
				bFloat := (float64(b) + 0.5) / float64(lutSize)
				c := rgb{rFloat, gFloat, bFloat}
				i := p.lookupNearest(c)
				cNew := p.lookup(i)
				rgba := cNew.toNRGBA()
				o := y*lutStride + x*4
				pix[o] = rgba.R
				pix[o+1] = rgba.G
				pix[o+2] = rgba.B
				pix[o+3] = 255
			}
			wg.Done()
		}(y)
	}
	wg.Wait()
}

func (p *Palette) computeBayerScaleLUT(lutSize, perRow, lutWidth, lutHeight, lutStride int, pix []byte) {
	// Also compute for each pixel the distance to the next color when adding or subtracting to all of r,g,b.
	// Use this to compute a dynamic Bayer scale.
	// At points, Bayer scale should be the MIN of the distances to next colors.
	// Elsewhere, Bayer scale ideally should be those values interpolated.
	// What can we practically get?
	// Store that data in the alpha channel.

	// For each protected palette index, find its ideal bayer scale.
	scales := make([]int, p.protected)
	var wg sync.WaitGroup
	for i := 0; i < p.protected; i++ {
		wg.Add(1)
		go func(i int) {
			c := p.lookup(i).toNRGBA()
			scale := 1
		FoundScale:
			for scale < 256 {
				for d := -1; d <= 1; d += 2 {
					rr := int(c.R) + scale*d
					gg := int(c.G) + scale*d
					bb := int(c.B) + scale*d
					r := rr * lutSize / 255
					g := gg * lutSize / 255
					b := bb * lutSize / 255
					if r < 0 {
						r = 0
					}
					if r >= lutSize {
						r = lutSize - 1
					}
					if g < 0 {
						g = 0
					}
					if g >= lutSize {
						g = lutSize - 1
					}
					if b < 0 {
						b = 0
					}
					if b >= lutSize {
						b = lutSize - 1
					}
					x := r + lutSize*(b%perRow)
					y := g + lutSize*(b/perRow)
					o := y*lutStride + x*4
					if pix[o] != c.R || pix[o+1] != c.G || pix[o+2] != c.B {
						break FoundScale
					}
				}
				scale++
			}
			scale--
			// Make all scales one LUT entry lower.
			// This fixes pathological gradients due to a roundoff error
			// in the color right next to a palette color.
			scale -= (255 + lutSize - 1) / lutSize
			if scale < 0 {
				scale = 0
			}
			scales[i] = scale
			wg.Done()
		}(i)
	}
	wg.Wait()

	// Set alpha channel to best Bayer scale for each pixel.
	for i := 0; i < p.protected; i++ {
		c := p.lookup(i).toNRGBA()
		rr := int(c.R)
		gg := int(c.G)
		bb := int(c.B)
		r := rr * lutSize / 255
		g := gg * lutSize / 255
		b := bb * lutSize / 255
		if r >= lutSize {
			r = lutSize - 1
		}
		if g >= lutSize {
			g = lutSize - 1
		}
		if b >= lutSize {
			b = lutSize - 1
		}
		x := r + lutSize*(b%perRow)
		y := g + lutSize*(b/perRow)
		o := y*lutStride + x*4
		pix[o+3] = uint8(scales[i])
	}
	for y := 0; y < lutHeight; y++ {
		wg.Add(1)
		go func(y int) {
			g := y % lutSize
			gFloat := (float64(g) + 0.5) / float64(lutSize)
			bY := (y / lutSize) * perRow
			for x := 0; x < lutWidth; x++ {
				o := y*lutStride + x*4
				if pix[o+3] != 255 {
					continue
				}
				r := x % lutSize
				rFloat := (float64(r) + 0.5) / float64(lutSize)
				b := bY + x/lutSize
				if b >= lutSize {
					break
				}
				bFloat := (float64(b) + 0.5) / float64(lutSize)
				c := rgb{rFloat, gFloat, bFloat}
				sum, weight := 0.0, 0.0
				for i, scale := range scales {
					c2 := p.lookup(i)
					f := 1 / c.diff2(c2)
					sum += f * float64(scale)
					weight += f
				}
				scale := m.Rint(sum / weight)
				pix[o+3] = uint8(scale)
			}
			wg.Done()
		}(y)
	}
	wg.Wait()
}

func (p *Palette) computeNearestTwoLUT(lutSize, perRow, lutWidth, lutHeight, lutStride int, pix []byte) {
	lut2 := lutWidth * 4

	var wg sync.WaitGroup
	for y := 0; y < lutHeight; y++ {
		wg.Add(1)
		go func(y int) {
			g := y % lutSize
			gFloat := (float64(g) + 0.5) / float64(lutSize)
			bY := (y / lutSize) * perRow
			for x := 0; x < lutWidth; x++ {
				r := x % lutSize
				rFloat := (float64(r) + 0.5) / float64(lutSize)
				b := bY + x/lutSize
				if b >= lutSize {
					break
				}
				bFloat := (float64(b) + 0.5) / float64(lutSize)
				c := rgb{rFloat, gFloat, bFloat}
				i, j := p.lookupNearestTwo(c)
				cI, cJ := p.lookup(i), p.lookup(j)
				rgbaI, rgbaJ := cI.toNRGBA(), cJ.toNRGBA()
				o := y*lutStride + x*4
				pix[o] = rgbaI.R
				pix[o+1] = rgbaI.G
				pix[o+2] = rgbaI.B
				pix[o+3] = 255
				o += lut2
				pix[o] = rgbaJ.R
				pix[o+1] = rgbaJ.G
				pix[o+2] = rgbaJ.B
				pix[o+3] = 255
			}
			wg.Done()
		}(y)
	}
	wg.Wait()
}

func (p *Palette) ToLUT(numLUTs int, img *ebiten.Image) (int, int, int) {
	var lutSize int
	defer func(t0 time.Time) {
		dt := time.Since(t0)
		log.Infof("building palette LUT of size %d took %v", lutSize, dt)
	}(time.Now())
	bounds := img.Bounds()
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y

	var timePerEntry float64
	switch numLUTs {
	case 1:
		// Finding nearest color is brute force, trying every palette index.
		timePerEntry = float64(p.size)
	case 2:
		// Algorithmic steps * measured time fraction.
		nI, nJMin, nJMax := p.nearestTwoChoices()
		timePerEntry = float64(nI) * float64(nJMin+nJMax) / 2 * 156.4 / 87.1
	default:
		log.Fatalf("unsupported LUT count: got %v, want 1 or 2", numLUTs)
	}
	maxEntries := *paletteMaxCycles / timePerEntry
	if maxEntries < 8 {
		maxEntries = 8
	}

	lutSize, perRow, rowCount := computeLUTSize(w/numLUTs, h, maxEntries)

	// Note: creating a temp image, and copying to that, so this does not invoke
	// thread synchronization as writing to an ebiten.Image would.
	lutWidth := lutSize * perRow
	lutHeight := lutSize * rowCount
	lutStride := lutWidth * numLUTs * 4
	pix := make([]uint8, lutStride*lutHeight)

	switch numLUTs {
	case 1:
		p.computeNearestLUT(lutSize, perRow, lutWidth, lutHeight, lutStride, pix)
		p.computeBayerScaleLUT(lutSize, perRow, lutWidth, lutHeight, lutStride, pix)
	case 2:
		p.computeNearestTwoLUT(lutSize, perRow, lutWidth, lutHeight, lutStride, pix)
	default:
		log.Fatalf("unsupported LUT count: got %v, want 1 or 2", numLUTs)
	}

	rect := image.Rectangle{
		Min: bounds.Min,
		Max: image.Point{
			X: bounds.Min.X + lutWidth*numLUTs,
			Y: bounds.Min.Y + lutHeight,
		},
	}
	img.SubImage(rect).(*ebiten.Image).ReplacePixels(pix)

	return lutSize, perRow, lutWidth
}

func sizeBayer(size int) (sizeSquare int, scale, offset float64) {
	sizeSquare = size * size
	bits := 0
	if size > 1 {
		bits = math.Ilogb(float64(size-1)) + 1
	}
	sizeCeil := 1 << bits
	sizeCeilSquare := sizeCeil * sizeCeil
	// Map to [0..1] _exclusive_ ranges.
	// Not _perfect_, but way nicer to work with.
	offset = 0.5 / float64(sizeCeilSquare)
	scale = 1.0 / float64(sizeCeilSquare)
	return
}

func sizeHalftone(size int) (sizeSquare int, scale, offset float64) {
	sizeSquare = size * size
	// Map to [0..1] _exclusive_ ranges.
	// Not _perfect_, but way nicer to work with.
	offset = 0.5 / float64(sizeSquare)
	scale = 1.0 / float64(sizeSquare)
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

// BayerPattern computes the Bayer pattern for this palette.
func (p *Palette) BayerPattern(size int) []float32 {
	sizeSquare, scale, offset := sizeBayer(size)
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
		bayern[i] = float32((float64(b) + offset) * scale)
	}
	return bayern
}

// HalftonePattern computes the Halftone pattern for this palette.
func (p *Palette) HalftonePattern(size int) []float32 {
	sizeSquare, scale, offset := sizeHalftone(size)
	type index struct {
		i         int
		distance2 int
		angle     float64
	}
	weighted := make([]index, sizeSquare)
	for i := range weighted {
		x := i % size
		y := i / size
		// Take distance from top left pixel corner.
		dx := 2*x + 1
		dy := 2*y + 1
		if dx > size {
			dx -= 2 * size
		}
		if dy > size {
			dy -= 2 * size
		}
		d := dx*dx + dy*dy
		// Compute angle as tie breaker.
		// Negate Y to get mathematically positive angles and not clockwise.
		a := math.Atan2(-float64(dy), float64(dx))
		if a < 0 {
			a += 2 * math.Pi
		}
		weighted[i] = index{i, d, a}
	}
	// Note: sort in reverse; the innermost pixel thus gets filled first.
	sort.Slice(weighted, func(i, j int) bool {
		do := weighted[i].distance2 - weighted[j].distance2
		if do != 0 {
			return do < 0
		}
		da := weighted[i].angle - weighted[j].angle
		if da != 0 {
			return da < 0
		}
		log.Fatalf("unreachable code: same distance and angle should be impossible, but happened at %v and %v", weighted[i].i, weighted[j].i)
		return false
	})
	bayern := make([]float32, sizeSquare)
	for b, idx := range weighted {
		bayern[idx.i] = float32((float64(b) + offset) * scale)
	}
	return bayern
}

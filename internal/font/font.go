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

package font

import (
	"fmt"
	"image"
	"sort"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/splash"
)

var (
	pinFontsToCache           = flag.Bool("pin_fonts_to_cache", false, "pin all fonts to glyph cache")
	pinFontsToCacheAtStart    = flag.Bool("pin_fonts_to_cache_at_start", false, "pin all fonts to glyph cache right at startup (otherwise this work is spread across the first few frames)")
	pinFontsToCacheBaseWeight = flag.Int("pin_fonts_to_cache_base_weight", 1, "base weight for English characters when font pinning")
	pinFontsToCacheCount      = flag.Int("pin_fonts_to_cache_count", 512, "maximum number of characters to pin")
	pinFontsToCacheFraction   = flag.Int("pin_fonts_to_cache_fraction", 30, "fraction of all characters to cache per frame")
	fontThreshold             = flag.Int("font_threshold", 0x7800, "threshold for font rendering; lower values are bolder; 0 means antialias as usual; threshold range is 1 to 65535 inclusive; set to 0 to use smooth font rendering instead")
	fontExtraSpacing          = flag.Int("font_extra_spacing", 31, "additional spacing for fonts in 64th pixels; should help with outline effect")
	fontFractionalSpacing     = flag.Bool("font_fractional_spacing", false, "allow fractional font spacing; looks better but may be slower; makes --pin_fonts_to_cache less effective")
	debugFontOverride         = flag.String("debug_font_override", "", "name of font to use instead of the intended font")
	debugFontProfiling        = flag.Bool("debug_font_profiling", false, "measure how long font caching took")
)

// Face is an alias to font.Face so users do not need to import the font package.
type Face struct {
	Face    font.Face
	Outline font.Face
}

func makeFace(f font.Face, size int) *Face {
	effect := &fontEffects{
		Face:       f,
		LineHeight: size,
	}
	outline := &fontOutline{effect}
	face := &Face{
		Face:    effect,
		Outline: outline,
	}
	return face
}

var fontProfilingTotal time.Duration

func LoadIntoCacheStepwise() func(s *splash.State) (splash.Status, error) {
	if !*pinFontsToCache {
		return func(s *splash.State) (splash.Status, error) {
			return splash.Continue, nil
		}
	}
	charSetStr := string(charSet)
	done := map[*Face]struct{}{}
	names := make([]string, 0, len(ByName))
	for name := range ByName {
		names = append(names, name)
	}
	sort.Strings(names)
	return func(s *splash.State) (splash.Status, error) {
		for _, name := range names {
			status, err := s.Enter(fmt.Sprintf("precaching %s", name), locale.G.Get("precaching %s", name), fmt.Sprintf("could not precache %v", name), splash.Single(func() error {
				f := ByName[name]
				if _, found := done[f]; found {
					return nil
				}
				done[f] = struct{}{}
				if *pinFontsToCacheAtStart {
					var t0 time.Time
					if *debugFontProfiling {
						t0 = time.Now()
					}
					f.precache(charSetStr)
					if *debugFontProfiling {
						dt := time.Since(t0)
						fontProfilingTotal += dt
						log.Infof("caching font %v: %v (total: %v)", name, dt, fontProfilingTotal)
					}
				}
				return nil
			}))
			if status != splash.Continue {
				return status, err
			}
		}
		return splash.Continue, nil
	}
}

// We always keep the game character set in cache.
// This has to be repeated regularly as Ebitengine expires unused cache entries.
func KeepInCache() {
	if !*pinFontsToCache {
		return
	}
	f := *pinFontsToCacheFraction
	l := len(charSet)
	low := charSetPos * l / f
	charSetPos++
	high := charSetPos * l / f
	if charSetPos == f {
		charSetPos = 0
	}
	charSubSet := charSet[low:high]
	charSubSetStr := string(charSubSet)
	done := map[*Face]struct{}{}
	for _, f := range ByName {
		if _, found := done[f]; found {
			continue
		}
		done[f] = struct{}{}
		f.precache(charSubSetStr)
	}
}

var (
	ByFont      = map[string]map[string]*Face{}
	ByName      = map[string]*Face(nil)
	currentFont string
)

type fontEffects struct {
	font.Face
	LineHeight int
}

func roundFixed(f fixed.Int26_6) fixed.Int26_6 {
	if *fontFractionalSpacing {
		return f
	}
	return (f + 32) & ^63
}

func (e *fontEffects) Glyph(dot fixed.Point26_6, r rune) (
	image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
	dr, mask, maskp, advance, ok := e.Face.Glyph(dot, r)
	return dr, fontEffectsMask(mask), maskp, advance, ok
}

func (e *fontEffects) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	adv, ok := e.Face.GlyphAdvance(r)
	if adv != 0 {
		adv += fixed.Int26_6(*fontExtraSpacing)
	}
	return roundFixed(adv), ok
}

func (e *fontEffects) Metrics() font.Metrics {
	// Override the line height to match what go-freetype did.
	// Narrower lines really look better here.
	m := e.Face.Metrics()
	m.Height = fixed.Int26_6(e.LineHeight) << 6
	return m
}

func (e *fontEffects) Kern(r0, r1 rune) fixed.Int26_6 {
	kern := e.Face.Kern(r0, r1)
	return roundFixed(kern)
}

func fontEffectsMask(src image.Image) image.Image {
	if *fontThreshold <= 0 {
		return src
	}
	r := src.Bounds()
	dst := image.NewAlpha(r)
	pr := 0
	for y := r.Min.Y; y < r.Max.Y; y++ {
		p := pr
		for x := r.Min.X; x < r.Max.X; x++ {
			_, _, _, a := src.At(x, y).RGBA()
			if int(a) >= *fontThreshold {
				dst.Pix[p] = 0xFF
			}
			p++
		}
		pr += dst.Stride
	}
	return dst
}

type fontOutline struct {
	font.Face
}

func (o *fontOutline) Glyph(dot fixed.Point26_6, r rune) (
	image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
	dr, mask, maskp, advance, ok := o.Face.Glyph(dot, r)
	drExpanded := image.Rectangle{
		Min: image.Point{
			X: dr.Min.X - 1,
			Y: dr.Min.Y - 1,
		},
		Max: image.Point{
			X: dr.Max.X + 1,
			Y: dr.Max.Y + 1,
		},
	}
	maskpExpanded := image.Point{
		X: maskp.X - 1,
		Y: maskp.Y - 1,
	}
	return drExpanded, fontOutlineMask(mask), maskpExpanded, advance, ok
}

func (o *fontOutline) GlyphBounds(r rune) (fixed.Rectangle26_6, fixed.Int26_6, bool) {
	bounds, advance, ok := o.Face.GlyphBounds(r)
	bounds.Min.X -= 1 << 6
	bounds.Min.Y -= 1 << 6
	bounds.Max.X += 1 << 6
	bounds.Max.Y += 1 << 6
	return bounds, advance, ok
}

func (o *fontOutline) Metrics() font.Metrics {
	m := o.Face.Metrics()
	m.Height += 2 << 6
	m.Ascent += 1 << 6
	m.Descent += 1 << 6
	return m
}

func outlineLine(a []uint8, n, stride int) {
	// In-place replace every value by the max of a[i-1]-1, a[i], a[i+1]-1.

	x := uint8(0)
	y := a[0]
	z := a[stride]

	p := 0
	m := n - 2
	for i := 0; i < n; i++ {
		mx := x
		my := y
		mz := z
		if mx > my {
			my = mx - 1
		}
		if mz > my {
			my = mz - 1
		}
		a[p] = my
		p += stride
		x, y = y, z
		if i < m {
			z = a[p+stride]
		} else {
			// Would be out of range.
			z = 0
		}
	}
}

func fontOutlineMask(src image.Image) image.Image {
	// The outline is:
	// - Transparent where the font is fully opaque (only if antialiasing is off).
	//   This fixes alpha blending of "font atop outline".
	// - Otherwise it's the max of the bordering pixels.

	// First make a copy. This saves At() calls.
	// Image must have 0-based rectangle.
	srcR := src.Bounds()
	r := image.Rectangle{
		Min: image.Point{
			X: srcR.Min.X - 1,
			Y: srcR.Min.Y - 1,
		},
		Max: image.Point{
			X: srcR.Max.X + 1,
			Y: srcR.Max.Y + 1,
		},
	}
	dst := image.NewAlpha(r)
	pr := dst.Stride
	for y := srcR.Min.Y; y < srcR.Max.Y; y++ {
		p := pr
		p++
		for x := srcR.Min.X; x < srcR.Max.X; x++ {
			_, _, _, a := src.At(x, y).RGBA()
			dst.Pix[p] = uint8((a + 128) / 257)
			p++
		}
		pr += dst.Stride
	}

	// Then replace every value by the max of the eight values around them - 1, or the self value.
	// This is done as a separable operation.

	pr = 0
	for y := r.Min.Y; y < r.Max.Y; y++ {
		outlineLine(dst.Pix[pr:], r.Max.X-r.Min.X, 1)
		pr += dst.Stride
	}

	pr = 0
	for x := r.Min.X; x < r.Max.X; x++ {
		outlineLine(dst.Pix[pr:], r.Max.Y-r.Min.Y, dst.Stride)
		pr++
	}

	// Finally, if NOT antialiasing, remap pixel values.
	if *fontThreshold > 0 {
		pr = 0
		for y := r.Min.Y; y < r.Max.Y; y++ {
			p := pr
			for x := r.Min.X; x < r.Max.X; x++ {
				switch dst.Pix[p] {
				case 0, 0xFF:
					dst.Pix[p] = 0
				default:
					dst.Pix[p] = 0xFF
				}
				p++
			}
			pr += dst.Stride
		}
	}

	return dst
}

func SetFont(font string) error {
	charSet = locale.CharSet(charSetBase, *pinFontsToCacheBaseWeight, *pinFontsToCacheCount)
	log.Infof("charset pinned: %v", string(charSet))
	charSetPos = 0
	if *debugFontOverride != "" {
		font = *debugFontOverride
	}
	if font == currentFont {
		return nil
	}
	ByName = ByFont[font]
	currentFont = font
	if ByName == nil {
		ByName = map[string]*Face{}
		ByFont[font] = ByName
		switch font {
		case "bitmapfont":
			return initBitmapfont()
		case "unifont":
			return initUnifont()
		default:
			return initGoFont()
		}
	}
	return nil
}

func CurrentFont() string {
	return currentFont
}

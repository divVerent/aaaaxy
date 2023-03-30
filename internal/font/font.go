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
	"image/color"
	"sort"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/splash"
)

var (
	pinFontsToCache           = flag.Bool("pin_fonts_to_cache", true, "pin all fonts to glyph cache")
	pinFontsToCacheBaseWeight = flag.Int("pin_fonts_to_cache_base_weight", 1, "base weight for English characters when font pinning")
	pinFontsToCacheCount      = flag.Int("pin_fonts_to_cache_count", 512, "maximum number of characters to pin")
	pinFontsToCacheFraction   = flag.Int("pin_fonts_to_cache_fraction", 30, "fraction of all characters to cache per frame")
	fontThreshold             = flag.Int("font_threshold", 0x7800, "threshold for font rendering; lower values are bolder; 0 means antialias as usual; threshold range is 1 to 65535 inclusive; set to 0 to use smooth font rendering instead")
	fontExtraSpacing          = flag.Int("font_extra_spacing", 31, "additional spacing for fonts in 64th pixels; should help with outline effect")
	fontFractionalSpacing     = flag.Bool("font_fractional_spacing", false, "allow fractional font spacing; looks better but may be slower; makes --pin_fonts_to_cache less effective")
	debugFontOverride         = flag.String("debug_font_override", "", "name of font to use instead of the intended font")
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

func LoadIntoCacheStepwise(s *splash.State) (splash.Status, error) {
	if !*pinFontsToCache {
		return splash.Continue, nil
	}
	charSetStr := string(charSet)
	done := map[*Face]struct{}{}
	names := make([]string, 0, len(ByName))
	for name := range ByName {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		status, err := s.Enter(fmt.Sprintf("precaching %s", name), locale.G.Get("precaching %s", name), fmt.Sprintf("could not precache %v", name), splash.Single(func() error {
			f := ByName[name]
			if _, found := done[f]; found {
				return nil
			}
			done[f] = struct{}{}
			f.precache(charSetStr)
			return nil
		}))
		if status != splash.Continue {
			return status, err
		}
	}
	return splash.Continue, nil
}

// We always keep the game character set in cache.
// This has to be repeated regularly as Ebitengine expires unused cache entries.
func KeepInCache() {
	if !*pinFontsToCache {
		return
	}
	charSubSet := charSet
	if charSetCached {
		f := *pinFontsToCacheFraction
		l := len(charSet)
		low := charSetPos * l / f
		charSetPos++
		high := charSetPos * l / f
		if charSetPos == f {
			charSetPos = 0
		}
		charSubSet = charSet[low:high]
	} else {
		charSetCached = true
		charSetPos = 0
	}
	charSubSetStr := string(charSubSet)
	done := map[*Face]struct{}{}
	for _, f := range ByName {
		if _, found := done[f]; found {
			continue
		}
		done[f] = struct{}{}
		f.precache(charSubSetStr)
	}
	return
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
	return dr, &fontEffectsMask{mask}, maskp, advance, ok
}

func (e *fontEffects) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	adv, ok := e.Face.GlyphAdvance(r)
	adv += fixed.Int26_6(*fontExtraSpacing)
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

type fontEffectsMask struct {
	image.Image
}

func (e *fontEffectsMask) At(x, y int) color.Color {
	// FYI: This yields a mask, and the draw function will only ever look at its alpha channel.
	base := e.Image.At(x, y)
	if *fontThreshold <= 0 {
		return base
	}
	_, _, _, a := base.RGBA()
	if int(a) < *fontThreshold {
		return color.Transparent
	}
	return color.Opaque
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
		X: maskp.X,
		Y: maskp.Y,
	}
	return drExpanded, &fontOutlineMask{
		Image: mask,
		Rect: m.Rect{
			Origin: m.Pos{
				X: maskp.X,
				Y: maskp.Y,
			},
			Size: m.Delta{
				DX: dr.Max.X - dr.Min.X,
				DY: dr.Max.Y - dr.Min.Y,
			},
		},
	}, maskpExpanded, advance, ok
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

type fontOutlineMask struct {
	image.Image
	Rect m.Rect
}

func (o *fontOutlineMask) Bounds() image.Rectangle {
	r := o.Image.Bounds()
	r.Max.X += 2
	r.Max.Y += 2
	return r
}

func (o *fontOutlineMask) atRaw(x, y int) color.Color {
	if o.Rect.DeltaPos(m.Pos{X: x, Y: y}).IsZero() {
		return o.Image.At(x, y)
	}
	return color.Transparent
}

func (o *fontOutlineMask) At(x, y int) color.Color {
	// The outline is:
	// - Transparent where the font is fully opaque (only if antialiasing is off).
	//   This fixes alpha blending of "font atop outline".
	if *fontThreshold > 0 {
		_, _, _, a := o.atRaw(x-1, y-1).RGBA()
		if a == 0xFFFF {
			return color.Transparent
		}
	}
	// - Otherwise it's the max of the bordering pixels.
	var maxA uint32
	for dy := -2; dy <= 0; dy++ {
		for dx := -2; dx <= 0; dx++ {
			if dx == -1 && dy == -1 {
				continue
			}
			_, _, _, a := o.atRaw(x+dx, y+dy).RGBA()
			if a > maxA {
				maxA = a
			}
		}
	}
	return color.Alpha16{A: uint16(maxA)}
}

func SetFont(font string) error {
	charSet = locale.CharSet(charSetBase, *pinFontsToCacheBaseWeight, *pinFontsToCacheCount)
	log.Infof("charset pinned: %v", string(charSet))
	charSetCached = false
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
		case "unifont":
			return initUnifont()
		default:
			return initGoFont()
		}
	}
	return nil
}

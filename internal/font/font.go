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
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/divVerent/aaaaxy/internal/flag"
	m "github.com/divVerent/aaaaxy/internal/math"
)

var (
	pinFontsToCache       = flag.Bool("pin_fonts_to_cache", true, "pin all fonts to glyph cache")
	fontThreshold         = flag.Int("font_threshold", 0x8000, "threshold for font rendering; lower values are bolder; 0 means antialias as usual; threshold range is 1 to 65535 inclusive; set to 0 to use smooth font rendering instead")
	fontExtraSpacing      = flag.Int("font_extra_spacing", 31, "additional spacing for fonts in 64th pixels; should help with outline effect")
	fontFractionalSpacing = flag.Bool("font_fractional_spacing", false, "allow fractional font spacing; looks better but may be slower; makes --pin_fonts_to_cache less effective")
)

// Face is an alias to font.Face so users do not need to import the font package.
type Face struct {
	Face    font.Face
	Outline font.Face
}

func makeFace(f font.Face, size int) (Face, error) {
	effect := &fontEffects{
		Face:       f,
		LineHeight: size,
	}
	outline := &fontOutline{effect}
	face := Face{
		Face:    effect,
		Outline: outline,
	}
	return face, nil
}

// We always keep the game character set in cache.
// This has to be repeated regularly as Ebitengine expires unused cache entries.
func KeepInCache(dst *ebiten.Image) {
	if *pinFontsToCache {
		for _, f := range ByName {
			f.precache(charSet())
		}
	}
}

var (
	ByName      = map[string]Face{}
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
	if font == currentFont {
		return nil
	}
	currentFont = font
	switch font {
	default:
		return initGoFont()
	}
}

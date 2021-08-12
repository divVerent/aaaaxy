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

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/gofont/gosmallcaps"
	"golang.org/x/image/math/fixed"

	"github.com/divVerent/aaaaxy/internal/flag"
	m "github.com/divVerent/aaaaxy/internal/math"
)

var (
	pinFontsToCache       = flag.Bool("pin_fonts_to_cache", true, "Pin all fonts to glyph cache.")
	pinFontsToCacheHarder = flag.Bool("pin_fonts_to_cache_harder", false, "Do a dummy draw command to pin fonts to glyph cache harder.")
	fontThreshold         = flag.Int("font_threshold", 0x5E00, "Threshold for font rendering; lower values are bolder. 0 means antialias as usual; threshold range is 1 to 65535 inclusive.")
	fontExtraSpacing      = flag.Int("font_extra_spacing", 32, "Additional spacing for fonts in 64th pixels; should help with outline effect.")
)

// Face is an alias to font.Face so users do not need to import the font package.
type Face struct {
	Face    font.Face
	Outline font.Face
}

func makeFace(f font.Face) Face {
	effect := &fontEffects{f}
	outline := &fontOutline{effect}
	face := Face{
		Face:    effect,
		Outline: outline,
	}
	all = append(all, face)
	return face
}

// cacheChars are all characters the game uses. ASCII plus all Unicode our map file contains.
var cacheChars = " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~τπö¾©"

// We always keep the game character set in cache.
// This has to be repeated regularly as ebiten expires unused cache entries.
func KeepInCache(dst *ebiten.Image) {
	if *pinFontsToCacheHarder {
		for _, f := range all {
			f.precache(dst, cacheChars)
		}
	}
	if *pinFontsToCache {
		for _, f := range all {
			f.recache(cacheChars)
		}
	}
}

var (
	all            = []Face{}
	ByName         = map[string]Face{}
	Centerprint    Face
	CenterprintBig Face
	DebugSmall     Face
	MenuBig        Face
	Menu           Face
	MenuSmall      Face
)

func Init() error {
	// Load the fonts.
	regular, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return fmt.Errorf("could not load goitalic font: %v", err)
	}
	italic, err := truetype.Parse(goitalic.TTF)
	if err != nil {
		return fmt.Errorf("could not load goitalic font: %v", err)
	}
	bold, err := truetype.Parse(gobold.TTF)
	if err != nil {
		return fmt.Errorf("could not load gosmallcaps font: %v", err)
	}
	mono, err := truetype.Parse(gomono.TTF)
	if err != nil {
		return fmt.Errorf("could not load gomono font: %v", err)
	}
	smallcaps, err := truetype.Parse(gosmallcaps.TTF)
	if err != nil {
		return fmt.Errorf("could not load gosmallcaps font: %v", err)
	}

	ByName["Small"] = makeFace(truetype.NewFace(regular, &truetype.Options{
		Size:       10,
		Hinting:    font.HintingFull,
		SubPixelsX: 1,
		SubPixelsY: 1,
	}))
	ByName["Regular"] = makeFace(truetype.NewFace(regular, &truetype.Options{
		Size:       14,
		Hinting:    font.HintingFull,
		SubPixelsX: 1,
		SubPixelsY: 1,
	}))
	ByName["Italic"] = makeFace(truetype.NewFace(italic, &truetype.Options{
		Size:       14,
		Hinting:    font.HintingFull,
		SubPixelsX: 1,
		SubPixelsY: 1,
	}))
	ByName["Bold"] = makeFace(truetype.NewFace(bold, &truetype.Options{
		Size:       14,
		Hinting:    font.HintingFull,
		SubPixelsX: 1,
		SubPixelsY: 1,
	}))
	ByName["Mono"] = makeFace(truetype.NewFace(mono, &truetype.Options{
		Size:       14,
		Hinting:    font.HintingFull,
		SubPixelsX: 1,
		SubPixelsY: 1,
	}))
	ByName["SmallCaps"] = makeFace(truetype.NewFace(smallcaps, &truetype.Options{
		Size:       14,
		Hinting:    font.HintingFull,
		SubPixelsX: 1,
		SubPixelsY: 1,
	}))
	Centerprint = makeFace(truetype.NewFace(italic, &truetype.Options{
		Size:       14,
		Hinting:    font.HintingFull,
		SubPixelsX: 1,
		SubPixelsY: 1,
	}))
	CenterprintBig = makeFace(truetype.NewFace(smallcaps, &truetype.Options{
		Size:       24,
		Hinting:    font.HintingFull,
		SubPixelsX: 1,
		SubPixelsY: 1,
	}))
	DebugSmall = makeFace(truetype.NewFace(regular, &truetype.Options{
		Size:       9,
		Hinting:    font.HintingFull,
		SubPixelsX: 1,
		SubPixelsY: 1,
	}))
	Menu = makeFace(truetype.NewFace(smallcaps, &truetype.Options{
		Size:       18,
		Hinting:    font.HintingFull,
		SubPixelsX: 1,
		SubPixelsY: 1,
	}))
	MenuBig = makeFace(truetype.NewFace(smallcaps, &truetype.Options{
		Size:       24,
		Hinting:    font.HintingFull,
		SubPixelsX: 1,
		SubPixelsY: 1,
	}))
	MenuSmall = makeFace(truetype.NewFace(smallcaps, &truetype.Options{
		Size:       12,
		Hinting:    font.HintingFull,
		SubPixelsX: 1,
		SubPixelsY: 1,
	}))

	return nil
}

type fontEffects struct {
	font.Face
}

func (e *fontEffects) Glyph(dot fixed.Point26_6, r rune) (
	image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
	dr, mask, maskp, advance, ok := e.Face.Glyph(dot, r)
	return dr, &fontEffectsMask{mask}, maskp, advance, ok
}

func (e *fontEffects) Kern(r0, r1 rune) fixed.Int26_6 {
	return e.Face.Kern(r0, r1) + fixed.Int26_6(*fontExtraSpacing)
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

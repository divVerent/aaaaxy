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
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"

	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/m"
)

var (
	precacheImg *ebiten.Image
)

// boundString returns the bounding rectangle of the given text.
func (f Face) boundString(str string) m.Rect {
	var r m.Rect
	bounds, _ := font.BoundString(f.Outline.GoX, str)
	x0 := bounds.Min.X.Floor()
	y0 := bounds.Min.Y.Floor()
	x1 := bounds.Max.X.Ceil()
	y1 := bounds.Max.Y.Ceil()
	r = m.Rect{
		Origin: m.Pos{
			X: x0,
			Y: y0,
		},
		Size: m.Delta{
			DX: x1 - x0,
			DY: y1 - y0,
		},
	}
	if r.Size.DX <= 0 {
		r.Size.DX = 1
	}
	if r.Size.DY <= 0 {
		r.Size.DY = 1
	}
	return r
}

// BoundString returns the bounding rectangle of the given text.
func (f Face) BoundString(str string) m.Rect {
	lines := strings.Split(str, "\n")
	var totalBounds m.Rect
	lineHeight := f.Outline.GoX.Metrics().Height.Ceil()
	y := 0
	for i, line := range lines {
		lines[i] = locale.ActiveShape(line)
	}
	for _, line := range lines {
		bounds := f.boundString(line)
		bounds.Origin.Y += y
		totalBounds = totalBounds.Union(bounds)
		y += lineHeight
	}
	return totalBounds
}

// drawLine draws one line of text.
func drawLine(f *faceWrapper, dst *ebiten.Image, line string, x, y int, align text.Align, fg color.Color) {
	// Use Ebitengine's glyph cache.
	options := &text.DrawOptions{
		LayoutOptions: text.LayoutOptions{
			LineSpacing:    0,
			PrimaryAlign:   align,
			SecondaryAlign: text.AlignStart,
		},
	}
	options.GeoM.Translate(float64(x), float64(y)-float64(f.GoX.Metrics().Ascent)/float64(1<<6))
	options.ColorScale.ScaleWithColor(fg)
	text.Draw(dst, line, f.Ebi, options)
}

type Align int

const (
	Left Align = iota
	Center
	Right
)

// Draw draws the given text.
func (f Face) Draw(dst *ebiten.Image, str string, pos m.Pos, boxAlign Align, fg, bg color.Color) {
	// We need to do our own line splitting because
	// we always want to center and Ebitengine would left adjust.
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		lines[i] = locale.ActiveShape(line)
	}
	y := pos.Y
	lineHeight := f.Outline.GoX.Metrics().Height.Ceil()
	var align text.Align
	switch boxAlign {
	case Left:
		align = text.AlignStart
	case Center:
		align = text.AlignCenter
	case Right:
		align = text.AlignEnd
	}
	for _, line := range lines {
		if _, _, _, a := bg.RGBA(); a != 0 {
			drawLine(f.Outline, dst, line, pos.X, y, align, bg)
		}
		drawLine(f.Face, dst, line, pos.X, y, align, fg)
		y += lineHeight
	}
}

func (f Face) precache(chars string) {
	if *fontFractionalSpacing {
		text.CacheGlyphs(chars, f.Face.Ebi)
		text.CacheGlyphs(chars, f.Outline.Ebi)
	} else {
		// Always cache at position 0 only.
		if precacheImg == nil {
			precacheImg = ebiten.NewImage(1, 1)
		}
		options := &text.DrawOptions{}
		text.Draw(precacheImg, chars, f.Face.Ebi, options)
		text.Draw(precacheImg, chars, f.Outline.Ebi, options)
	}
}

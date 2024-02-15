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
	"image/draw"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	"github.com/divVerent/aaaaxy/internal/locale"
	m "github.com/divVerent/aaaaxy/internal/math"
)

// boundString returns the bounding rectangle of the given text.
func (f Face) boundString(str string) m.Rect {
	var r m.Rect
	bounds, _ := font.BoundString(f.Outline, str)
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
	lineHeight := f.Outline.Metrics().Height.Ceil()
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
func drawLine(f font.Face, dst draw.Image, line string, x, y int, fg color.Color) {
	if locale.ActiveUsesEbitenText() {
		dst, ok := dst.(*ebiten.Image)
		if ok {
			// Use Ebitengine's glyph cache.
			text.Draw(dst, line, f, x, y, fg)
			return
		}
	}
	// No glyph cache.
	d := font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(fg),
		Face: f,
		Dot:  fixed.P(x, y),
	}
	d.DrawString(line)
}

type Align int

const (
	AsBounds Align = iota
	Left
	Center
	Right
)

// Draw draws the given text.
func (f Face) Draw(dst draw.Image, str string, pos m.Pos, boxAlign Align, fg, bg color.Color) {
	// We need to do our own line splitting because
	// we always want to center and Ebitengine would left adjust.
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		lines[i] = locale.ActiveShape(line)
	}
	bounds := make([]m.Rect, len(lines))
	for i, line := range lines {
		bounds[i] = f.boundString(line)
	}
	var totalMin, totalMax int
	for _, lineBounds := range bounds {
		xMin := lineBounds.Origin.X
		xMax := xMin + lineBounds.Size.DX
		if xMin < totalMin {
			totalMin = xMin
		}
		if xMax > totalMax {
			totalMax = xMax
		}
	}
	// AsBounds: offset := pos.X + totalBounds.Size.DX/2 + totalBounds.Origin.X
	// Center: offset := pos.X
	// Left: offset := pos.X + totalBounds.Size.DX/2
	// Right: offset := pos.X - (totalBounds.Size.DX+1)/2
	offset := pos.X
	switch boxAlign {
	case AsBounds:
		offset += (totalMin + totalMax) / 2
	case Left:
		offset += (totalMax - totalMin) / 2
	case Right:
		offset -= (totalMax - totalMin + 1) / 2
	}
	y := pos.Y
	lineHeight := f.Outline.Metrics().Height.Ceil()
	for i, line := range lines {
		lineBounds := bounds[i]
		// totalBounds: tX size tDX
		// lineBouds: lX size lDX
		// Want lX+d .. lX+lDX+d centered in tX .. tX+tDX
		// Thus: lX+d - tX = tX+tDX - (lX+lDX+d)
		// d = tX - lX + (tDX - lDX)/2.
		x := offset - lineBounds.Origin.X - lineBounds.Size.DX/2
		if _, _, _, a := bg.RGBA(); a != 0 {
			drawLine(f.Outline, dst, line, x, y, bg)
		}
		// Draw the text itself.
		drawLine(f.Face, dst, line, x, y, fg)
		y += lineHeight
	}
}

func (f Face) precache(chars string) {
	text.CacheGlyphs(f.Face, chars)
	text.CacheGlyphs(f.Outline, chars)
}

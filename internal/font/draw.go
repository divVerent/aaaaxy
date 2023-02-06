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

	m "github.com/divVerent/aaaaxy/internal/math"
)

// BoundString returns the bounding rectangle of the given text.
func (f Face) BoundString(str string) m.Rect {
	rect := text.BoundString(f.Outline, str)
	return m.Rect{
		Origin: m.Pos{
			X: rect.Min.X,
			Y: rect.Min.Y,
		},
		Size: m.Delta{
			DX: rect.Max.X - rect.Min.X,
			DY: rect.Max.Y - rect.Min.Y,
		},
	}
}

// drawLine draws one line of text.
func drawLine(f font.Face, dst draw.Image, line string, x, y int, fg color.Color) {
	switch dst := dst.(type) {
	case *ebiten.Image:
		// Use Ebitengine's glyph cache.
		text.Draw(dst, line, f, x, y, fg)
	default:
		// No glyph cache.
		d := font.Drawer{
			Dst:  dst,
			Src:  image.NewUniform(fg),
			Face: f,
			Dot:  fixed.P(x, y),
		}
		d.DrawString(line)
	}
}

type Align int

const (
	Left Align = iota
	Center
	Right
)

// Draw draws the given text.
func (f Face) Draw(dst draw.Image, str string, pos m.Pos, boxAlign Align, fg, bg color.Color) {
	// We need to do our own line splitting because
	// we always want to center and Ebitengine would left adjust.
	totalBounds := f.BoundString(str)
	// Center: offset := pos.X
	// Left: offset := pos.X + totalBounds.Size.DX/2
	// Right: offset := pos.X - (totalBounds.Size.DX+1)/2
	offset := pos.X
	switch boxAlign {
	case Left:
		offset += totalBounds.Size.DX / 2
	case Right:
		offset -= (totalBounds.Size.DX + 1) / 2
	}
	fy := fixed.I(pos.Y)
	for _, line := range strings.Split(str, "\n") {
		lineBounds := f.BoundString(line)
		// totalBounds: tX size tDX
		// lineBouds: lX size lDX
		// Want lX+d .. lX+lDX+d centered in tX .. tX+tDX
		// Thus: lX+d - tX = tX+tDX - (lX+lDX+d)
		// d = tX - lX + (tDX - lDX)/2.
		x := offset - lineBounds.Origin.X - lineBounds.Size.DX/2
		y := fy.Floor()
		if _, _, _, a := bg.RGBA(); a != 0 {
			drawLine(f.Outline, dst, line, x, y, bg)
		}
		// Draw the text itself.
		drawLine(f.Face, dst, line, x, y, fg)
		fy += f.Outline.Metrics().Height + 1 // Line height is 1 pixel above font height.
	}
}

func (f Face) precache(chars string) {
	text.CacheGlyphs(f.Face, chars)
	text.CacheGlyphs(f.Outline, chars)
}

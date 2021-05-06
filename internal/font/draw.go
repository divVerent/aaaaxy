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
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/math/fixed"

	m "github.com/divVerent/aaaaaa/internal/math"
)

// BoundString returns the bounding rectangle of the given text.
func (f Face) BoundString(str string) m.Rect {
	rect := text.BoundString(f.Face, str)
	// Note: we expand the rect by 1 to include the outline.
	return m.Rect{
		Origin: m.Pos{
			X: rect.Min.X - 1,
			Y: rect.Min.Y - 1,
		},
		Size: m.Delta{
			DX: rect.Max.X - rect.Min.X + 2,
			DY: rect.Max.Y - rect.Min.Y + 2,
		},
	}
}

// Draw draws the given text.
func (f Face) Draw(dst *ebiten.Image, str string, pos m.Pos, centerX bool, fg, bg color.Color) {
	// We need to do our own line splitting because
	// we always want to center and ebiten would left adjust.
	var totalBounds m.Rect
	if !centerX {
		totalBounds = f.BoundString(str)
	}
	fy := fixed.I(pos.Y)
	for _, line := range strings.Split(str, "\n") {
		lineBounds := f.BoundString(line)
		// totalBounds: tX size tDX
		// lineBouds: lX size lDX
		// Want lX+d .. lX+lDX+d centered in tX .. tX+tDX
		// Thus: lX+d - tX = tX+tDX - (lX+lDX+d)
		// d = tX - lX + (tDX - lDX)/2.
		x := pos.X + totalBounds.Origin.X - lineBounds.Origin.X + (totalBounds.Size.DX-lineBounds.Size.DX)/2
		y := fy.Floor()
		if _, _, _, a := bg.RGBA(); a != 0 {
			// Draw the outline.
			for dx := -1; dx <= +1; dx++ {
				for dy := -1; dy <= +1; dy++ {
					if dx == 0 && dy == 0 {
						continue
					}
					text.Draw(dst, line, f.Face, x+dx, y+dy, bg)
				}
			}
		}
		// Draw the text itself.
		text.Draw(dst, line, f.Face, x, y, fg)
		fy += f.Face.Metrics().Height
	}
}

func (f Face) precache(dst *ebiten.Image, chars string) {
	text.Draw(dst, chars, f.Face, 0, 0, color.NRGBA{R: 0, G: 0, B: 0, A: 0})
}

func (f Face) recache(chars string) {
	text.CacheGlyphs(f.Face, chars)
}

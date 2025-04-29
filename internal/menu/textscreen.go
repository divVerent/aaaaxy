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

package menu

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/m"
)

// textScreenScrollInPos is the position where sure nothing can be seen.
func textScreenScrollInPos(text []string, lineHeight int) int {
	return engine.GameHeight + lineHeight
}

// textScreenStartPos is the position where the start of the text shows.
func textScreenStartPos(text []string, lineHeight int) int {
	return lineHeight
}

// textScreenAdjustScrollUp performs scrolling up.
func textScreenAdjustScrollUp(text []string, y, d int, lineHeight int) int {
	if y > lineHeight {
		return y
	}
	if y+d > lineHeight {
		return lineHeight
	}
	return y + d
}

// textScreenAdjustScrollDown performs scrolling down.
func textScreenAdjustScrollDown(text []string, y, d int, lineHeight int) int {
	t := textScreenEndPos(text, lineHeight)
	if y < t {
		return y
	}
	if y-d < t {
		return t
	}
	return y - d
}

// textScreenEndPos is the position where the end of the text shows.
func textScreenEndPos(text []string, lineHeight int) int {
	return -lineHeight*len(text) + engine.GameHeight
}

func renderTextScreen(dst *ebiten.Image, titleFont, normalFont *font.Face, text []string, pos m.Pos, align font.Align, lineHeight int, titleFG, titleBG, normalFG, normalBG color.Color) {
	x := pos.X
	nextIsTitle := true
	for i, line := range text {
		if line == "" {
			nextIsTitle = true
			continue
		}
		isTitle := nextIsTitle
		nextIsTitle = false
		y := lineHeight*i + pos.Y
		if y < 0 || y >= engine.GameHeight+lineHeight {
			continue
		}
		if isTitle {
			titleFont.Draw(dst, line, m.Pos{X: x, Y: y}, align, titleFG, titleBG)
		} else {
			normalFont.Draw(dst, line, m.Pos{X: x, Y: y}, align, normalFG, normalBG)
		}
	}
}

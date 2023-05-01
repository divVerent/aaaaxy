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

package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/divVerent/aaaaxy/internal/flag"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/palette"
)

type editMode int

const (
	editNone editMode = iota
	editMin
	editBoth
	editMax
)

type touchEditInfo struct {
	active    bool
	frames    int
	rect      *m.Rect
	xMode     editMode
	yMode     editMode
	startPos  m.Pos
	startRect m.Rect
}

const (
	gridSize                = 8
	touchToggleActionFrames = 120
)

var (
	touchEditPad bool = false

	touchReservedArea = m.Rect{
		Origin: m.Pos{X: 192, Y: 64},
		Size:   m.Delta{DX: 640 - 192 - 192, DY: 360 - 64},
	}

	snaps = []m.Delta{
		// Distance: 0
		{DX: 0, DY: 0},
		// Distance: 1, always enumerated clockwise starting at 12.
		{DX: 0, DY: -1},
		{DX: 1, DY: 0},
		{DX: 0, DY: 1},
		{DX: -1, DY: 0},
		// Distance: sqrt(2), always enumerated clockwise starting at 12.
		{DX: 1, DY: -1},
		{DX: 1, DY: 1},
		{DX: -1, DY: 1},
		{DX: -1, DY: -1},
		// Distance: 2, always enumerated clockwise starting at 12.
		{DX: 0, DY: -2},
		{DX: 2, DY: 0},
		{DX: 0, DY: 2},
		{DX: -2, DY: 0},
		// Distance: sqrt(5), always enumerated clockwise starting at 12.
		{DX: 1, DY: -2},
		{DX: 2, DY: -1},
		{DX: 2, DY: 1},
		{DX: 1, DY: 2},
		{DX: -1, DY: 2},
		{DX: -2, DY: 1},
		{DX: -2, DY: -1},
		{DX: -1, DY: -2},
		// Distance: sqrt(8), always enumerated clockwise starting at 12.
		{DX: 2, DY: -2},
		{DX: 2, DY: 2},
		{DX: -2, DY: 2},
		{DX: -2, DY: -2},
	}
)

func touchEditMode(g int) editMode {
	switch g {
	case 0:
		return editMin
	case 3:
		return editMax
	default:
		return editNone
	}
}

func touchEditOrigin(mode editMode, o int, dp int) int {
	switch mode {
	case editMin, editBoth:
		return o + dp
	default:
		return o
	}
}

func touchEditSize(mode editMode, s int, dp int) int {
	switch mode {
	case editMin:
		return s - dp
	case editMax:
		return s + dp
	default:
		return s
	}
}

func touchEditAllowed(rect *m.Rect, replacement m.Rect, gameWidth, gameHeight int) bool {
	if replacement.Origin.X < 0 || replacement.Origin.Y < 0 {
		return false
	}
	if replacement.Size.DX < 64 || replacement.Size.DY < 64 {
		return false
	}
	if replacement.OppositeCorner().X >= gameWidth || replacement.OppositeCorner().Y >= gameHeight {
		return false
	}
	if touchReservedArea.Delta(replacement).IsZero() {
		return false
	}
	for _, i := range impulses {
		if i.touchRect == nil || i.touchRect.Size.IsZero() || i.touchRect == rect {
			continue
		}
		if replacement.Delta(*i.touchRect).IsZero() {
			return false
		}
	}
	return true
}

func touchToggleActionButton(gameWidth, gameHeight int) {
	if !Action.touchRect.Size.IsZero() {
		Action.touchRect.Size = m.Delta{}
		return
	}
	// Start with the square inscribed into Jump.
	jumpSize := Jump.touchRect.Size
	for range []int{0, 1} {
		// Find an empty place close to Jump, and make it the same size.
		for _, neighbor := range []m.Delta{
			{DX: 1, DY: 0},
			{DX: 0, DY: -1},
			{DX: -1, DY: 0},
			{DX: 0, DY: 1},
		} {
			initialSize := jumpSize
			initialOrigin := Jump.touchRect.Origin
			// Minimize size in the direction moving in.
			if neighbor.DX == 0 {
				initialSize.DY = 64
			} else {
				initialSize.DX = 64
			}
			// Set initial origin.
			if neighbor.DX < 0 {
				initialOrigin.X -= initialSize.DX
			} else if neighbor.DX > 0 {
				initialOrigin.X += jumpSize.DX
			} else if neighbor.DY < 0 {
				initialOrigin.Y -= initialSize.DY
			} else if neighbor.DY > 0 {
				initialOrigin.Y += jumpSize.DY
			}
			newRect := m.Rect{
				Origin: initialOrigin,
				Size:   initialSize,
			}
			if touchEditAllowed(Action.touchRect, newRect, gameWidth, gameHeight) {
				sizeStep := neighbor.Mul(gridSize)
				originStep := m.Delta{}
				if sizeStep.DX < 0 {
					originStep.DX = sizeStep.DX
					sizeStep.DX = -sizeStep.DX
				}
				if sizeStep.DY < 0 {
					originStep.DY = sizeStep.DY
					sizeStep.DY = -sizeStep.DY
				}
				// Expand as far as we can.
				for {
					expandedRect := m.Rect{
						Origin: newRect.Origin.Add(originStep),
						Size:   newRect.Size.Add(sizeStep),
					}
					if !touchEditAllowed(Action.touchRect, expandedRect, gameWidth, gameHeight) {
						break
					}
					newRect = expandedRect
				}
				*Action.touchRect = newRect
				return
			}
		}
		// Didn't work? Swap size and try again.
		jumpSize.DX, jumpSize.DY = jumpSize.DY, jumpSize.DX
	}
}

func touchEditUpdate(gameWidth, gameHeight int) bool {
	if !touchEditPad {
		for _, t := range touches {
			t.edit.active = false
		}
		return false
	}
	eatTouches := false
	for _, t := range touches {
		if !t.hit {
			continue
		}
		if !touchReservedArea.DeltaPos(t.pos).IsZero() {
			eatTouches = true
		}
		if t.edit.active {
			// Move what is being hit.
			t.edit.frames++
			if t.edit.rect == nil {
				// Hold for some time to toggle action button existence.
				if t.edit.frames == touchToggleActionFrames {
					touchToggleActionButton(gameWidth, gameHeight)
				}
				continue
			}
			newRect := *t.edit.rect
			// The truncate rounding in m.Div slightly prefers the same coordinate. Good.
			dx := gridSize * m.Div(t.pos.X-t.edit.startPos.X, gridSize)
			dy := gridSize * m.Div(t.pos.Y-t.edit.startPos.Y, gridSize)
			for _, snap := range snaps {
				newRect.Origin.X = touchEditOrigin(t.edit.xMode, t.edit.startRect.Origin.X, dx+gridSize*snap.DX)
				newRect.Origin.Y = touchEditOrigin(t.edit.yMode, t.edit.startRect.Origin.Y, dy+gridSize*snap.DY)
				newRect.Size.DX = touchEditSize(t.edit.xMode, t.edit.startRect.Size.DX, dx+gridSize*snap.DX)
				newRect.Size.DY = touchEditSize(t.edit.yMode, t.edit.startRect.Size.DY, dy+gridSize*snap.DY)
				if touchEditAllowed(t.edit.rect, newRect, gameWidth, gameHeight) {
					*t.edit.rect = newRect
					break
				}
			}
		} else {
			t.edit.active = true
			t.edit.rect = nil
			t.edit.startPos = t.pos
			// Identify what is hit, set flag, xMode, yMode appropriately.
			// Just set active if nothing is hit.
			for _, i := range impulses {
				if i.touchRect == nil || i.touchRect.Size.IsZero() {
					continue
				}
				gx, gy := i.touchRect.GridPos(t.pos, 4, 4)
				if gx < 0 || gy < 0 || gx >= 4 || gy >= 4 {
					continue
				}
				// Hit, so start active this rectangle.
				t.edit.rect = i.touchRect
				t.edit.startRect = *t.edit.rect
				t.edit.xMode = touchEditMode(gx)
				t.edit.yMode = touchEditMode(gy)
				if t.edit.xMode == editNone && t.edit.yMode == editNone {
					t.edit.xMode = editBoth
					t.edit.yMode = editBoth
				}
				break
			}
		}
	}
	return eatTouches
}

func touchEditDraw(screen *ebiten.Image) {
	if !touchEditPad {
		return
	}
	for _, i := range impulses {
		if i.touchRect == nil || i.touchRect.Size.IsZero() {
			continue
		}
		boxColor := palette.EGA(palette.White, 255)
		vector.DrawFilledRect(screen, float32(i.touchRect.Origin.X), float32(i.touchRect.Origin.Y), float32(i.touchRect.Size.DX), float32(i.touchRect.Size.DY), boxColor, false)
		innerColor := palette.EGA(palette.DarkGrey, 255)
		vector.DrawFilledRect(screen, float32(i.touchRect.Origin.X+1), float32(i.touchRect.Origin.Y+1), float32(i.touchRect.Size.DX-2), float32(i.touchRect.Size.DY-2), innerColor, false)
	}
	gridColor := palette.EGA(palette.LightGrey, 32)
	sz := screen.Bounds().Size()
	for x := 0; x < sz.X/gridSize; x++ {
		for y := 0; y < sz.Y/gridSize; y++ {
			r := m.Rect{
				Origin: m.Pos{X: x * gridSize, Y: y * gridSize},
				Size:   m.Delta{DX: 8, DY: 8},
			}
			if touchReservedArea.Delta(r).IsZero() {
				continue
			}
			vector.DrawFilledRect(screen, float32(x*gridSize+1), float32(y*gridSize+1), 6, 6, gridColor, false)
		}
	}
	touchPadDraw(screen)
}

func touchSetEditor(want bool) {
	if touchEditPad == want {
		return
	}
	touchCancelClicks()
	touchEditPad = want
}

func TouchResetEditor() {
	flag.ResetFlagToDefault("touch_rect_left")
	flag.ResetFlagToDefault("touch_rect_right")
	flag.ResetFlagToDefault("touch_rect_down")
	flag.ResetFlagToDefault("touch_rect_up")
	flag.ResetFlagToDefault("touch_rect_jump")
	flag.ResetFlagToDefault("touch_rect_action")
	flag.ResetFlagToDefault("touch_rect_exit")
}

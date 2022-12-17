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
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/palette"
)

// TODO(divVerent):
// Make each rect a command line option.
// Only store the touch rect - the draw rect shall be the largest box of correct aspect that fits inside.
// Then make an editor for these.
// Idea: put the edit mode in here, but make it impossible to cover the center.
// Also no button overlap.
// Use an 8x8 grid (gcd).
// Controls:
// - Each active finger has a state and a start pos.
// - State is what object it moves and which corner/side.
// - 4x4 grid.
// - left, none, none, right
// - however, if both x and y are none, move both
// - do not allow overlaps or outside
// - do moves in steps to move as much as needed
// - in center of screen, menu items/buttons to exit input edit mode
// - min size of each button: 64x64

var touchEditPad bool

type editMode int

const (
	editNone editMode = iota
	editMin
	editBoth
	editMax
)

type touchEditInfo struct {
	active bool
	rect   *m.Rect
	xMode  editMode
	yMode  editMode
}

func touchEditUpdate(gameWidth, gameHeight int) {
	if !touchEditPad {
		return
	}
	for _, t := range touches {
		if t.active {
			// Move what is being hit.
		} else {
			// Identify what is hit, set flag, xMode, yMode appropriately.
			// Just set active if nothing is hit.
		}
	}
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
		ebitenutil.DrawRect(screen, float64(i.touchRect.Origin.X), float64(i.touchRect.Origin.Y), float64(i.touchRect.Size.DX), float64(i.touchRect.Size.DY), boxColor)
		innerColor := palette.EGA(palette.DarkGrey, 255)
		ebitenutil.DrawRect(screen, float64(i.touchRect.Origin.X+1), float64(i.touchRect.Origin.Y+1), float64(i.touchRect.Size.DX-2), float64(i.touchRect.Size.DY-2), innerColor)
	}
	gridColor := palette.EGA(palette.LightGrey, 32)
	w, h := screen.Size()
	for x := 0; x < w/8; x++ {
		for y := 0; y < h/8; y++ {
			ebitenutil.DrawRect(screen, float64(x*8+1), float64(y*8+1), 6, 6, gridColor)
		}
	}
}

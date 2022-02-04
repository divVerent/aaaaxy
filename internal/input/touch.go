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

	"github.com/divVerent/aaaaxy/internal/flag"
	m "github.com/divVerent/aaaaxy/internal/math"
)

var (
	touch      = flag.Bool("touch", true, "enable touch input")
	touchForce = flag.Bool("touch_force", false, "always show touch controls")
)

const (
	touchClickMaxFrames = 30
	touchPadFrames      = 300
)

type touchInfo struct {
	frames int
	pos    m.Pos
	hit    bool
}

var (
	touchWantPad  bool
	touches       map[ebiten.TouchID]*touchInfo
	touchIDs      []ebiten.TouchID
	touchHoverPos m.Pos
	touchPadFrame int
)

func touchUpdate(screenWidth, screenHeight, gameWidth, gameHeight int) {
	if !*touch {
		return
	}
	for _, t := range touches {
		t.hit = false
	}
	touchIDs = ebiten.AppendTouchIDs(touchIDs[:0])
	if len(touchIDs) > 0 {
		// Either support touch OR mouse. This prevents duplicate click events.
		mouseCancel()
		touchPadFrame = touchPadFrames
	} else if touchPadFrame > 0 {
		touchPadFrame--
	}
	for _, id := range touchIDs {
		t, found := touches[id]
		if !found {
			t = &touchInfo{}
			touches[id] = t
		}
		t.hit = true
		t.frames++
	}
	hoverAcc := m.Pos{}
	hoverCnt := 0
	for id, t := range touches {
		if !t.hit {
			if t.frames < touchClickMaxFrames {
				clickPos = &t.pos
			}
			delete(touches, id)
			continue
		}
		x, y := ebiten.TouchPosition(id)
		x = (x*gameWidth + screenWidth/2) / screenWidth
		y = (y*gameHeight + screenHeight/2) / screenHeight
		t.pos = m.Pos{X: x, Y: y}
		if t.frames < touchClickMaxFrames {
			hoverAcc = hoverAcc.Add(t.pos.Delta(m.Pos{}))
			hoverCnt++
		}
	}
	if hoverCnt > 0 {
		touchHoverPos = hoverAcc.Add(m.Delta{DX: hoverCnt / 2, DY: hoverCnt / 2}).Div(hoverCnt)
		hoverPos = &touchHoverPos
	}
}

func touchSetWantPad(want bool) {
	touchWantPad = want
}

func (i *impulse) touchPressed() InputMap {
	if !touchWantPad {
		return 0
	}
	if i.touchRect.Size.IsZero() {
		return 0
	}
	for _, t := range touches {
		if i.touchRect.DeltaPos(t.pos).IsZero() {
			return Touchscreen
		}
	}
	return 0
}

func touchDraw(screen *ebiten.Image) {
	if !touchWantPad {
		return
	}
	if !*touchForce && touchPadFrame > 0 {
		return
	}
	for _, i := range impulses {
		if i.touchRect.Size.IsZero() {
			continue
		}
		// TODO(divVerent): Draw it.
		// Should look similar to NES pad.
	}
}

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

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/image"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/palette"
)

var (
	touch      = flag.Bool("touch", true, "enable touch input")
	touchForce = flag.Bool("touch_force", flag.SystemDefault(map[string]interface{}{
		"android/*": true,
		"ios/*":     true,
		"js/*":      true,
		"*/*":       false,
	}).(bool), "always show touch controls")
	touchRectLeft   = flag.Text("touch_rect_left", m.Rect{Origin: m.Pos{X: 0, Y: 232}, Size: m.Delta{DX: 64, DY: 64}}, "touch rectangle for moving left")
	touchRectRight  = flag.Text("touch_rect_right", m.Rect{Origin: m.Pos{X: 64, Y: 232}, Size: m.Delta{DX: 64, DY: 64}}, "touch rectangle for moving right")
	touchRectDown   = flag.Text("touch_rect_down", m.Rect{Origin: m.Pos{X: 0, Y: 296}, Size: m.Delta{DX: 128, DY: 64}}, "touch rectangle for moving down")
	touchRectUp     = flag.Text("touch_rect_up", m.Rect{Origin: m.Pos{X: 0, Y: 168}, Size: m.Delta{DX: 128, DY: 64}}, "touch rectangle for moving up")
	touchRectJump   = flag.Text("touch_rect_jump", m.Rect{Origin: m.Pos{X: 576, Y: 296}, Size: m.Delta{DX: 64, DY: 64}}, "touch rectangle for jumping")
	touchRectAction = flag.Text("touch_rect_action", m.Rect{Origin: m.Pos{X: 576, Y: 0}, Size: m.Delta{DX: 64, DY: 296}}, "touch rectangle for performing an action")
	touchRectExit   = flag.Text("touch_rect_exit", m.Rect{Origin: m.Pos{X: 0, Y: 0}, Size: m.Delta{DX: 128, DY: 64}}, "touch rectangle for exiting")
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
	touchUsePad   bool
	touchShowPad  bool
	touchEditPad  bool
	touches       = map[ebiten.TouchID]*touchInfo{}
	touchIDs      []ebiten.TouchID
	touchHoverPos m.Pos
	touchPadFrame int
)

func touchUpdate(screenWidth, screenHeight, gameWidth, gameHeight int, crtK1, crtK2 float64) {
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
			if !touchEditPad && t.frames < touchClickMaxFrames {
				clickPos = &t.pos
			}
			delete(touches, id)
			continue
		}
		x, y := ebiten.TouchPosition(id)
		t.pos = pointerCoords(screenWidth, screenHeight, gameWidth, gameHeight, crtK1, crtK2, x, y)
		hoverAcc = hoverAcc.Add(t.pos.Delta(m.Pos{}))
		hoverCnt++
	}
	if !touchEditPad && hoverCnt > 0 {
		touchHoverPos = hoverAcc.Add(m.Delta{DX: hoverCnt / 2, DY: hoverCnt / 2}).Div(hoverCnt)
		hoverPos = &touchHoverPos
	}
}

func touchSetUsePad(want bool) {
	touchUsePad = want
}

func touchSetShowPad(want bool) {
	touchShowPad = want
}

func (i *impulse) touchPressed() InputMap {
	if touchEditPad || !touchUsePad {
		return 0
	}
	if i.touchRect == nil || i.touchRect.Size.IsZero() {
		return 0
	}
	for _, t := range touches {
		if i.touchRect.DeltaPos(t.pos).IsZero() {
			return Touchscreen
		}
	}
	return 0
}

func touchInit() error {
	var err error
	Left.touchImage, err = image.Load("sprites", "touch_left.png")
	if err != nil {
		return err
	}
	Right.touchImage, err = image.Load("sprites", "touch_right.png")
	if err != nil {
		return err
	}
	Up.touchImage, err = image.Load("sprites", "touch_up.png")
	if err != nil {
		return err
	}
	Down.touchImage, err = image.Load("sprites", "touch_down.png")
	if err != nil {
		return err
	}
	Jump.touchImage, err = image.Load("sprites", "touch_jump.png")
	if err != nil {
		return err
	}
	Action.touchImage, err = image.Load("sprites", "touch_action.png")
	if err != nil {
		return err
	}
	Exit.touchImage, err = image.Load("sprites", "touch_exit.png")
	if err != nil {
		return err
	}
	return nil
}

func touchDraw(screen *ebiten.Image) {
	if !touchShowPad {
		return
	}
	if !*touchForce && touchPadFrame <= 0 {
		return
	}
	for _, i := range impulses {
		if i.touchRect == nil || i.touchRect.Size.IsZero() {
			continue
		}
		if touchEditPad {
			boxColor := palette.EGA(palette.White, 255)
			ebitenutil.DrawRect(screen, float64(i.touchRect.Origin.X), float64(i.touchRect.Origin.Y), float64(i.touchRect.Size.DX), float64(i.touchRect.Size.DY), boxColor)
			innerColor := palette.EGA(palette.DarkGrey, 255)
			ebitenutil.DrawRect(screen, float64(i.touchRect.Origin.X+1), float64(i.touchRect.Origin.Y+1), float64(i.touchRect.Size.DX-2), float64(i.touchRect.Size.DY-2), innerColor)
		}
		img := i.touchImage
		if img == nil {
			continue
		}
		options := &ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeSourceOver,
			Filter:        ebiten.FilterNearest,
		}
		w, h := img.Size()
		ox := float64(i.touchRect.Origin.X)
		oy := float64(i.touchRect.Origin.Y)
		sw := float64(i.touchRect.Size.DX) / float64(w)
		sh := float64(i.touchRect.Size.DY) / float64(h)
		if sw < sh {
			oy += float64(h) * 0.5 * (sh - sw)
			sh = sw
		} else if sw > sh {
			ox += float64(w) * 0.5 * (sw - sh)
			sw = sh
		}
		options.GeoM.Scale(sw, sh)
		options.GeoM.Translate(ox, oy)
		if i.Held {
			options.ColorM.Scale(-1, -1, -1, 1)
			options.ColorM.Translate(1, 1, 1, 0)
		}
		screen.DrawImage(img, options)
	}
	if touchEditPad {
		gridColor := palette.EGA(palette.LightGrey, 32)
		w, h := screen.Size()
		for x := 0; x < w/8; x++ {
			for y := 0; y < h/8; y++ {
				ebitenutil.DrawRect(screen, float64(x*8+1), float64(y*8+1), 6, 6, gridColor)
			}
		}
	}
}

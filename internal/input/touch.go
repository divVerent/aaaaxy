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
	"github.com/hajimehoshi/ebiten/v2/colorm"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/image"
	m "github.com/divVerent/aaaaxy/internal/math"
)

var (
	touch      = flag.Bool("touch", true, "enable touch input")
	touchForce = flag.Bool("touch_force", flag.SystemDefault(map[string]bool{
		"android/*": true,
		"ios/*":     true,
		"js/*":      true,
		"*/*":       false,
	}), "always show touch controls")
	touchRectLeft   = flag.Text("touch_rect_left", m.Rect{Origin: m.Pos{X: 0, Y: 232}, Size: m.Delta{DX: 64, DY: 64}}, "touch rectangle for moving left")
	touchRectRight  = flag.Text("touch_rect_right", m.Rect{Origin: m.Pos{X: 64, Y: 232}, Size: m.Delta{DX: 64, DY: 64}}, "touch rectangle for moving right")
	touchRectDown   = flag.Text("touch_rect_down", m.Rect{Origin: m.Pos{X: 0, Y: 296}, Size: m.Delta{DX: 128, DY: 64}}, "touch rectangle for moving down")
	touchRectUp     = flag.Text("touch_rect_up", m.Rect{Origin: m.Pos{X: 0, Y: 168}, Size: m.Delta{DX: 128, DY: 64}}, "touch rectangle for moving up")
	touchRectJump   = flag.Text("touch_rect_jump", m.Rect{Origin: m.Pos{X: 576, Y: 296}, Size: m.Delta{DX: 64, DY: 64}}, "touch rectangle for jumping")
	touchRectAction = flag.Text("touch_rect_action", m.Rect{}, "touch rectangle for performing an action")
	touchRectExit   = flag.Text("touch_rect_exit", m.Rect{Origin: m.Pos{X: 0, Y: 0}, Size: m.Delta{DX: 128, DY: 64}}, "touch rectangle for exiting")
)

const (
	touchClickMaxFrames = 30
	touchPadFrames      = 300
)

type touchInfo struct {
	clickFrames    int
	clickCancelled bool
	pos            m.Pos
	prevPos        m.Pos
	hit            bool
	edit           touchEditInfo
}

var (
	touchUsePad           bool
	touchShowPad          bool
	touches               = map[ebiten.TouchID]*touchInfo{}
	touchIDs              []ebiten.TouchID
	touchHoverPos         m.Pos
	touchPadFrame         int
	touchPadUsed          bool = false
	actionButtonAvailable bool = false
)

func touchCancelClicks() {
	for _, t := range touches {
		t.clickCancelled = true
	}
}

func touchEmulateMouse() {
	hoverAcc := m.Pos{}
	hoverCnt := 0
	for _, t := range touches {
		if !t.hit {
			if t.clickFrames < touchClickMaxFrames && !t.clickCancelled {
				clickPos = &t.prevPos
			}
			continue
		}
		hoverAcc = hoverAcc.Add(t.pos.Delta(m.Pos{}))
		hoverCnt++
	}
	if hoverCnt > 0 {
		touchHoverPos = hoverAcc.Add(m.Delta{DX: hoverCnt / 2, DY: hoverCnt / 2}).Div(hoverCnt)
		hoverPos = &touchHoverPos
	}
}

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
		touchPadUsed = true
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
		t.clickFrames++
		t.prevPos = t.pos
		x, y := ebiten.TouchPosition(id)
		t.pos = pointerCoords(screenWidth, screenHeight, gameWidth, gameHeight, crtK1, crtK2, x, y)
	}
	if touchEditUpdate(gameWidth, gameHeight) {
		// log.Infof("touchEditUpdate returned true - not emulating mouse")
	} else {
		touchEmulateMouse()
	}
	for id, t := range touches {
		if !t.hit {
			delete(touches, id)
			continue
		}
	}
}

func touchSetUsePad(want bool) {
	if touchUsePad == want {
		return
	}
	touchCancelClicks()
	touchUsePad = want
}

func touchSetShowPad(want bool) {
	if touchShowPad == want {
		return
	}
	touchCancelClicks()
	touchShowPad = want
}

func (i *impulse) touchPressed() InputMap {
	if touchEditPad || !touchUsePad {
		return 0
	}
	if i.touchRect == nil {
		return 0
	}
	for _, t := range touches {
		if i.touchRect.Size.IsZero() {
			touched := false
			for _, other := range impulses {
				if other == i {
					continue
				}
				if other.touchRect == nil {
					continue
				}
				if other.touchRect.Size.IsZero() {
					continue
				}
				if other.touchRect.DeltaPos(t.pos).IsZero() {
					touched = true
					break
				}
			}
			if !touched {
				return Touchscreen
			}
		} else if i.touchRect.DeltaPos(t.pos).IsZero() {
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
	if touchEditPad {
		return
	}
	if !touchShowPad {
		return
	}
	if !*touchForce && touchPadFrame <= 0 {
		return
	}
	touchPadDraw(screen)
}

func touchPadDraw(screen *ebiten.Image) {
	for _, i := range impulses {
		r := i.touchRect
		if r == nil {
			continue
		}
		img := i.touchImage
		if img == nil {
			continue
		}
		options := &colorm.DrawImageOptions{
			Blend:  ebiten.BlendSourceOver,
			Filter: ebiten.FilterNearest,
		}
		var colorM colorm.ColorM
		if r.Size.IsZero() {
			if !i.Held || !actionButtonAvailable {
				continue
			}
			sz := screen.Bounds().Size()
			r = &m.Rect{
				Origin: m.Pos{X: (sz.X - 32) / 2, Y: sz.Y - 32},
				Size:   m.Delta{DX: 32, DY: 32},
			}
			colorM.Scale(1, 1, 1, 1.0/3)
		}
		sz := img.Bounds().Size()
		ox := float64(r.Origin.X)
		oy := float64(r.Origin.Y)
		sw := float64(r.Size.DX) / float64(sz.X)
		sh := float64(r.Size.DY) / float64(sz.Y)
		if sw < sh {
			oy += float64(sz.Y) * 0.5 * (sh - sw)
			sh = sw
		} else if sw > sh {
			ox += float64(sz.X) * 0.5 * (sw - sh)
			sw = sh
		}
		options.GeoM.Scale(sw, sh)
		options.GeoM.Translate(ox, oy)
		if i.Held {
			colorM.Scale(-1, -1, -1, 1)
			colorM.Translate(1, 1, 1, 0)
		}
		colorm.DrawImage(screen, img, colorM, options)
	}
}

func HaveTouch() bool {
	return *touchForce || touchPadUsed
}

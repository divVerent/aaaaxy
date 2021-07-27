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

package centerprint

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/font"
	m "github.com/divVerent/aaaaxy/internal/math"
)

type Centerprint struct {
	text       string
	bounds     m.Rect
	color      color.Color
	waitScroll bool
	waitFade   bool
	face       font.Face
	pos        InitialPosition

	alphaFrames int
	alphaFrame  int
	scrollPos   int
	fadeOut     bool
	active      bool
}

var (
	screenWidth, screenHeight int
	centerprints              []*Centerprint
)

type Importance int

const (
	Important Importance = iota
	NotImportant
)

type InitialPosition int

const (
	Top InitialPosition = iota
	Middle
)

func NormalFont() font.Face {
	return font.Centerprint
}

func BigFont() font.Face {
	return font.CenterprintBig
}

func Reset() {
	centerprints = centerprints[:0]
}

func New(txt string, imp Importance, pos InitialPosition, face font.Face, color color.Color, fadeTime time.Duration) *Centerprint {
	frames := int(fadeTime * 60 / time.Second)
	if frames < 1 {
		frames = 1
	}
	cp := &Centerprint{
		text:        txt,
		color:       color,
		waitScroll:  imp == Important,
		waitFade:    true,
		face:        face,
		pos:         pos,
		alphaFrame:  1,
		alphaFrames: frames,
		active:      true,
	}
	cp.bounds = cp.face.BoundString(txt)
	if pos == Middle {
		cp.scrollPos = cp.targetPos()
	}
	if len(centerprints) != 0 {
		height := cp.bounds.Size.DY + 1 // Leave one pixel between lines.
		if centerprints[0].scrollPos < height {
			cp.scrollPos = centerprints[0].scrollPos - height
		}
	}
	centerprints = append(centerprints, cp)
	return cp
}

func (cp *Centerprint) SetFadeOut(fadeOut bool) {
	cp.fadeOut = fadeOut
}

func (cp *Centerprint) targetPos() int {
	switch cp.pos {
	case Top:
		return screenHeight / 4
	case Middle:
		return screenHeight / 3
	default:
		log.Panicf("invalid initial position: %v", cp.pos)
		return 0
	}
}

func (cp *Centerprint) update() bool {
	if cp.scrollPos < cp.targetPos() {
		cp.scrollPos++
	} else {
		cp.waitScroll = false
	}
	if cp.waitFade || cp.waitScroll || !cp.fadeOut {
		if cp.scrollPos > 0 {
			if cp.alphaFrame < cp.alphaFrames {
				cp.alphaFrame++
			} else {
				cp.waitFade = false
			}
		}
	} else {
		if cp.alphaFrame > 0 {
			cp.alphaFrame--
		}
	}
	if cp.alphaFrame == 0 {
		cp.active = false
		return false
	}
	return true
}

func (cp *Centerprint) draw(screen *ebiten.Image) {
	a := float64(cp.alphaFrame) / float64(cp.alphaFrames)
	if a == 0 {
		return
	}
	var alphaM ebiten.ColorM
	alphaM.Scale(1.0, 1.0, 1.0, a)
	fg := alphaM.Apply(cp.color)
	bg := color.NRGBA{R: 0, G: 0, B: 0, A: uint8(a * 255)}
	x := screenWidth / 2
	y := cp.scrollPos - cp.bounds.Size.DY - cp.bounds.Origin.Y
	cp.face.Draw(screen, cp.text, m.Pos{X: x, Y: y}, true, fg, bg)
}

func (cp *Centerprint) Active() bool {
	return cp != nil && cp.active
}

func Update() {
	offscreens := 0
	for i, cp := range centerprints {
		if !cp.update() && i == offscreens {
			offscreens++
		}
	}
	centerprints = centerprints[offscreens:]
}

func Draw(screen *ebiten.Image) {
	screenWidth, screenHeight = screen.Size()
	for _, cp := range centerprints {
		cp.draw(screen)
	}
}

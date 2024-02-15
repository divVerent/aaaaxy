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
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"

	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/palette"
)

type Centerprint struct {
	text       string
	bounds     m.Rect
	bgColor    color.Color
	fgColor    color.Color
	waitScroll bool
	waitFade   bool
	face       *font.Face
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

func (i Importance) MarshalText() ([]byte, error) {
	switch i {
	case Important:
		return []byte("Important"), nil
	case NotImportant:
		return []byte("NotImportant"), nil
	}
	return nil, fmt.Errorf("could not marshal Importance %d", i)
}

func (i *Importance) UnmarshalText(text []byte) error {
	switch string(text) {
	case "Important":
		*i = Important
		return nil
	case "NotImportant":
		*i = NotImportant
		return nil
	default:
		return fmt.Errorf("unexpected Importance value: %q", string(text))
	}
}

type InitialPosition int

const (
	Top InitialPosition = iota
	Middle
)

func (i InitialPosition) MarshalText() ([]byte, error) {
	switch i {
	case Top:
		return []byte("Top"), nil
	case Middle:
		return []byte("Middle"), nil
	}
	return nil, fmt.Errorf("could not marshal InitialPosition %d", i)
}

func (i *InitialPosition) UnmarshalText(text []byte) error {
	switch string(text) {
	case "Top":
		*i = Top
		return nil
	case "Middle":
		*i = Middle
		return nil
	default:
		return fmt.Errorf("unexpected InitialPosition value: %q", string(text))
	}
}

func NormalFont() *font.Face {
	return font.ByName["Centerprint"]
}

func BigFont() *font.Face {
	return font.ByName["CenterprintBig"]
}

func Reset() {
	centerprints = centerprints[:0]
}

func New(txt string, imp Importance, pos InitialPosition, face *font.Face, fgColor color.Color, fadeTime time.Duration) *Centerprint {
	return NewWithBG(txt, imp, pos, face, palette.EGA(palette.Black, 255), fgColor, fadeTime)
}

func NewWithBG(txt string, imp Importance, pos InitialPosition, face *font.Face, bgColor, fgColor color.Color, fadeTime time.Duration) *Centerprint {
	frames := int(fadeTime * 60 / time.Second)
	if frames < 1 {
		frames = 1
	}
	cp := &Centerprint{
		text:        txt,
		bgColor:     bgColor,
		fgColor:     fgColor,
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
	centerprints = append(centerprints, cp)
	return cp
}

func (cp *Centerprint) SetFadeOut(fadeOut bool) {
	cp.fadeOut = fadeOut
}

func (cp *Centerprint) height() int {
	return cp.bounds.Size.DY + 1 // Leave one pixel between lines.
}

func (cp *Centerprint) targetPos() int {
	var t int
	switch cp.pos {
	case Top:
		t = screenHeight / 4
	case Middle:
		t = screenHeight / 3
	default:
		log.TraceErrorf("invalid initial position: %v", cp.pos)
		t = screenHeight / 3
	}
	if t < cp.height() {
		t = cp.height()
	}
	return t
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
	var alphaM colorm.ColorM
	alphaM.Scale(1.0, 1.0, 1.0, a)
	fg := alphaM.Apply(cp.fgColor)
	bg := alphaM.Apply(cp.bgColor)
	x := screenWidth / 2
	y := cp.scrollPos - cp.bounds.Size.DY - cp.bounds.Origin.Y
	cp.face.Draw(screen, cp.text, m.Pos{X: x, Y: y}, font.Center, fg, bg)
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
	sz := screen.Bounds().Size()
	screenWidth, screenHeight = sz.X, sz.Y
	for _, cp := range centerprints {
		cp.draw(screen)
	}
}

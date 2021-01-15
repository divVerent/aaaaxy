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
	"image"
	"image/color"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomonobold"
)

const (
	alphaFrames = 64
)

type Centerprint struct {
	text   string
	bounds image.Rectangle
	color  color.Color
	force  bool
	face   font.Face

	alphaFrame int
	scrollPos  int
	fadeOut    bool
	active     bool
}

var (
	screenWidth, screenHeight int
	normalFace, bigFace       font.Face
	centerprints              []*Centerprint
)

func init() {
	normalFont, err := truetype.Parse(goitalic.TTF)
	if err != nil {
		log.Panicf("Could not load goitalic font: %v", err)
	}
	normalFace = truetype.NewFace(normalFont, &truetype.Options{
		Size:    16,
		Hinting: font.HintingFull,
	})
	bigFont, err := truetype.Parse(gomonobold.TTF)
	if err != nil {
		log.Panicf("Could not load gomonobold font: %v", err)
	}
	bigFace = truetype.NewFace(bigFont, &truetype.Options{
		Size:    24,
		Hinting: font.HintingFull,
	})
}

type Importance int

const (
	Important Importance = iota
	NotImportant
)

type InitialPosition int

const (
	Top = iota
	Middle
)

type Font int

const (
	NormalFont = iota
	BigFont
)

func New(txt string, imp Importance, pos InitialPosition, font Font, color color.Color) *Centerprint {
	cp := &Centerprint{
		text:       txt,
		color:      color,
		force:      imp == Important,
		alphaFrame: 1,
		active:     true,
	}
	switch font {
	case NormalFont:
		cp.face = normalFace
	case BigFont:
		cp.face = bigFace
	default:
		log.Panicf("Unknown centerprint font: %v", font)
	}
	cp.bounds = text.BoundString(cp.face, txt)
	if pos == Middle {
		cp.scrollPos = cp.targetPos()
	}
	if len(centerprints) != 0 {
		height := cp.bounds.Max.Y - cp.bounds.Min.Y
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
	return (screenHeight - (cp.bounds.Min.Y - cp.bounds.Max.Y)) / 4
}

func (cp *Centerprint) update() bool {
	if cp.scrollPos < cp.targetPos() {
		cp.scrollPos++
	} else if cp.alphaFrame >= alphaFrames {
		cp.force = false
	}
	if cp.force || !cp.fadeOut {
		if cp.scrollPos > 0 && cp.alphaFrame < alphaFrames {
			cp.alphaFrame++
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
	a := float64(cp.alphaFrame) / alphaFrames
	if a == 0 {
		return
	}
	var alphaM ebiten.ColorM
	alphaM.Scale(1.0, 1.0, 1.0, a)
	fg := alphaM.Apply(cp.color)
	bg := color.NRGBA{R: 0, G: 0, B: 0, A: uint8(a * 255)}
	x := (screenWidth-(cp.bounds.Max.X-cp.bounds.Min.X))/2 - cp.bounds.Min.X
	y := cp.scrollPos - cp.bounds.Max.Y
	for dx := -1; dx <= +1; dx++ {
		for dy := -1; dy <= +1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			text.Draw(screen, cp.text, cp.face, x+dx, y+dy, bg)
		}
	}
	text.Draw(screen, cp.text, cp.face, x, y, fg)
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

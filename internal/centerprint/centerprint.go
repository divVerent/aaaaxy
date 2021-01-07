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
)

const (
	alphaFrames = 64
)

type Centerprint struct {
	text   string
	bounds image.Rectangle
	color  color.Color
	force  bool

	alphaFrame int
	scrollPos  int
	fadeOut    bool
	active     bool
}

var (
	screenWidth, screenHeight int
	face                      font.Face
	centerprints              []*Centerprint
)

func init() {
	centerprintFont, err := truetype.Parse(goitalic.TTF)
	if err != nil {
		log.Panicf("could not load goitalic font: %v", err)
	}
	face = truetype.NewFace(centerprintFont, &truetype.Options{
		Size:    16,
		Hinting: font.HintingFull,
	})
}

func New(txt string, force bool, color color.Color) *Centerprint {
	cp := &Centerprint{
		text:       txt,
		bounds:     text.BoundString(face, txt),
		color:      color,
		force:      force,
		alphaFrame: 1,
		active:     true,
	}
	centerprints = append(centerprints, cp)
	return cp
}

func (cp *Centerprint) SetFadeOut(fadeOut bool) {
	cp.fadeOut = fadeOut
}

func (cp *Centerprint) update() bool {
	if cp.force || !cp.fadeOut {
		if cp.alphaFrame < alphaFrames {
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
	if cp.scrollPos < (screenHeight-(cp.bounds.Min.Y-cp.bounds.Max.Y))/4 {
		cp.scrollPos++
	} else {
		cp.force = false
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
	x := screenWidth - (cp.bounds.Max.X-cp.bounds.Min.X)/2 - cp.bounds.Min.X
	y := cp.scrollPos - cp.bounds.Max.Y
	for dx := -1; dx <= +1; dx++ {
		for dy := -1; dy <= +1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			text.Draw(screen, cp.text, face, x+dx, y+dy, bg)
		}
	}
	text.Draw(screen, cp.text, face, x, y, fg)
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

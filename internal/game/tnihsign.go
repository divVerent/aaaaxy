package game

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goitalic"

	"github.com/divVerent/aaaaaa/internal/engine"
	m "github.com/divVerent/aaaaaa/internal/math"
)

var tnihFace font.Face

func init() {
	tnihFont, err := truetype.Parse(goitalic.TTF)
	if err != nil {
		log.Panicf("could not load goitalic font: %v", err)
	}
	tnihFace = truetype.NewFace(tnihFont, &truetype.Options{
		Size:    20,
		Hinting: font.HintingFull,
	})
}

const (
	tnihAlphaFrames = 64
)

// TnihSign just displays a text and remembers that it was hit.
type TnihSign struct {
	World     *engine.World
	Spawnable *engine.Spawnable
	Entity    *engine.Entity

	Text      string
	TextSize  image.Rectangle
	SeenImage *ebiten.Image

	ScrollFrame int
	AlphaFrame  int
	ForceOn     bool
}

func (t *TnihSign) Spawn(w *engine.World, s *engine.Spawnable, e *engine.Entity) error {
	t.World = w
	t.Spawnable = s
	t.Entity = e
	var err error
	t.SeenImage, err = engine.LoadImage("sprites", "tnihsign_seen.png")
	if err != nil {
		return fmt.Errorf("could not load sign seen sprite: %v", err)
	}
	if s.PersistentState["seen"] == "true" {
		t.Entity.Image = t.SeenImage
	} else {
		t.Entity.Image, err = engine.LoadImage("sprites", "tnihsign.png")
		if err != nil {
			return fmt.Errorf("could not load sign sprite: %v", err)
		}
	}
	t.Entity.Orientation = m.Identity()
	t.Text = s.Properties["text"]
	t.TextSize = text.BoundString(tnihFace, t.Text)
	return nil
}

func (t *TnihSign) Despawn() {}

func (t *TnihSign) Update() {
	if t.ForceOn || (t.World.Player.Rect.Delta(t.Entity.Rect) == m.Delta{}) {
		if t.AlphaFrame < tnihAlphaFrames {
			t.AlphaFrame++
		}
	} else {
		if t.AlphaFrame > 0 {
			t.AlphaFrame--
		}
	}
	if t.AlphaFrame > 0 {
		if t.ScrollFrame < (engine.GameHeight-(t.TextSize.Min.Y-t.TextSize.Max.Y))/4 {
			t.ScrollFrame++
			if t.Spawnable.PersistentState["seen"] != "true" {
				t.ForceOn = true
				t.Spawnable.PersistentState["seen"] = "true"
				t.Entity.Image = t.SeenImage
			}
		} else {
			t.ForceOn = false
		}
	} else {
		t.ScrollFrame = 0
	}
}

func (t *TnihSign) Touch(other *engine.Entity) {}

func (t *TnihSign) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {
	// TODO turn into a real centerprint system...
	if t.AlphaFrame == 0 {
		return
	}
	a := uint8(float64(t.AlphaFrame) / tnihAlphaFrames * 255)
	x := (engine.GameWidth-(t.TextSize.Max.X-t.TextSize.Min.X))/2 - t.TextSize.Min.X
	y := t.ScrollFrame - t.TextSize.Max.Y
	for dx := -1; dx <= +1; dx++ {
		for dy := -1; dy <= +1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			text.Draw(screen, t.Text, tnihFace, x+dx, y+dy, color.NRGBA{R: 0, G: 0, B: 0, A: a})
		}
	}
	text.Draw(screen, t.Text, tnihFace, x, y, color.NRGBA{R: 255, G: 255, B: 255, A: a})
}

func init() {
	engine.RegisterEntityType(&TnihSign{})
}

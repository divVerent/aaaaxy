package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/centerprint"
	"github.com/divVerent/aaaaaa/internal/engine"
	m "github.com/divVerent/aaaaaa/internal/math"
)

// TnihSign just displays a text and remembers that it was hit.
type TnihSign struct {
	World     *engine.World
	Spawnable *engine.Spawnable
	Entity    *engine.Entity

	Text      string
	SeenImage *ebiten.Image

	Centerprint *centerprint.Centerprint
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
	return nil
}

func (t *TnihSign) Despawn() {}

func (t *TnihSign) Update() {
	if (t.World.Player.Rect.Delta(t.Entity.Rect) == m.Delta{}) {
		if t.Centerprint.Active() {
			t.Centerprint.SetFadeOut(false)
		} else {
			t.Centerprint = centerprint.New(t.Text, t.Spawnable.PersistentState["seen"] != "true", color.NRGBA{R: 255, G: 255, B: 85, A: 255})
			t.Spawnable.PersistentState["seen"] = "true"
		}
	} else {
		if t.Centerprint.Active() {
			t.Centerprint.SetFadeOut(true)
		}
	}
}

func (t *TnihSign) Touch(other *engine.Entity) {}

func (t *TnihSign) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&TnihSign{})
}

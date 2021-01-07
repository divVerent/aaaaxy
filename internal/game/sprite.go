package game

import (
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	m "github.com/divVerent/aaaaaa/internal/math"
)

// Sprite is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
type Sprite struct{}

func (s *Sprite) Spawn(w *engine.World, sp *engine.Spawnable, e *engine.Entity) error {
	var err error
	directory := sp.Properties["image_dir"]
	if directory == "" {
		directory = "sprites"
	}
	e.Image, err = engine.LoadImage(directory, sp.Properties["image"])
	if err != nil {
		return err
	}
	e.ResizeImage = true
	e.Solid = sp.Properties["solid"] != "false"
	e.Opaque = sp.Properties["opaque"] != "false"
	if sp.Properties["alpha"] != "" {
		e.Alpha, err = strconv.ParseFloat(sp.Properties["alpha"], 64)
		if err != nil {
			return fmt.Errorf("could not decode alpha %q: %v", sp.Properties["alpha"], err)
		}
	}
	if sp.Properties["z_index"] != "" {
		e.ZIndex, err = strconv.Atoi(sp.Properties["z_index"])
		if err != nil {
			return fmt.Errorf("could not decode z index %q: %v", sp.Properties["z_index"], err)
		}
	}
	return nil
}

func (s *Sprite) Despawn() {}

func (s *Sprite) Update() {}

func (s *Sprite) Touch(other *engine.Entity) {}

func (s *Sprite) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&Sprite{})
}

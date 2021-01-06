package game

import (
	"fmt"
	"strconv"

	"github.com/divVerent/aaaaaa/internal/engine"
)

// Sprite is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
type Sprite struct{}

func (p *Sprite) Spawn(w *engine.World, s *engine.Spawnable, e *engine.Entity) error {
	var err error
	directory := s.Properties["image_dir"]
	if directory == "" {
		directory = "sprites"
	}
	e.Image, err = engine.LoadImage(directory, s.Properties["image"])
	if err != nil {
		return err
	}
	e.ResizeImage = true
	e.Solid = s.Properties["solid"] != "false"
	e.Opaque = s.Properties["opaque"] != "false"
	if s.Properties["alpha"] != "" {
		e.Alpha, err = strconv.ParseFloat(s.Properties["alpha"], 64)
		if err != nil {
			return fmt.Errorf("could not decode alpha %q: %v", s.Properties["alpha"], err)
		}
	}
	return nil
}

func (p *Sprite) Despawn() {}

func (p *Sprite) Update() {}

func (p *Sprite) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&Sprite{})
}

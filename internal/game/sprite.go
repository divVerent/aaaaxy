package game

import (
	"github.com/divVerent/aaaaaa/internal/engine"
)

// Sprite is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
type Sprite struct{}

func (p *Sprite) Spawn(w *engine.World, s *engine.Spawnable, e *engine.Entity) error {
	var err error
	e.Image, err = engine.LoadImage("sprites", s.Properties["image"])
	if err != nil {
		return err
	}
	e.Solid = s.Properties["solid"] != "false"
	e.Opaque = s.Properties["opaque"] != "false"
	return nil
}

func (p *Sprite) Despawn() {}

func (p *Sprite) Update() {}

func (p *Sprite) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&Sprite{})
}

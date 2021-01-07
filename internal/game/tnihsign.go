package game

import (
	"github.com/divVerent/aaaaaa/internal/engine"
)

// TnihSign just displays a text and remembers that it was hit.
type TnihSign struct {
	Spawnable *engine.Spawnable
}

func (t *TnihSign) Spawn(w *engine.World, s *engine.Spawnable, e *engine.Entity) error {
	t.Spawnable = s
	// Load image.
	// Property: "text".
	return nil
}

func (t *TnihSign) Despawn() {}

func (t *TnihSign) Update() {}

func (t *TnihSign) Touch(other *engine.Entity) {
	// TODO.
}

func init() {
	engine.RegisterEntityType(&TnihSign{})
}

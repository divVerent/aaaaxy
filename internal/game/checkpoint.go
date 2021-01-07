package game

import (
	"github.com/divVerent/aaaaaa/internal/engine"
)

// Checkpoint remembers that it was hit and allows spawning from there again. Also displays a text.
type Checkpoint struct {
	World     *engine.World
	Spawnable *engine.Spawnable
}

func (t *Checkpoint) Spawn(w *engine.World, s *engine.Spawnable, e *engine.Entity) error {
	t.Spawnable = s
	t.World = w
	// Property: "name".
	return nil
}

func (t *Checkpoint) Despawn() {}

func (t *Checkpoint) Update() {}

func (t *Checkpoint) Touch(other *engine.Entity) {
	// TODO.
}

func init() {
	engine.RegisterEntityType(&Checkpoint{})
}

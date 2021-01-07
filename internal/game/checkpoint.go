package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	m "github.com/divVerent/aaaaaa/internal/math"
)

// Checkpoint remembers that it was hit and allows spawning from there again. Also displays a text.
type Checkpoint struct {
	World     *engine.World
	Spawnable *engine.Spawnable
}

func (c *Checkpoint) Spawn(w *engine.World, s *engine.Spawnable, e *engine.Entity) error {
	c.Spawnable = s
	c.World = w
	// Property: "name".
	return nil
}

func (c *Checkpoint) Despawn() {}

func (c *Checkpoint) Update() {}

func (c *Checkpoint) Touch(other *engine.Entity) {
	// TODO.
}

func (c *Checkpoint) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&Checkpoint{})
}

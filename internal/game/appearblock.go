package game

import (
	"github.com/divVerent/aaaaaa/internal/engine"
)

// AppearBlock is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
type AppearBlock struct {
	World  *engine.World
	Entity *engine.Entity

	AnimFrame int
}

const (
	AppearFrames         = 16
	AppearXDistance      = 2 * engine.TileSize
	AppearYDistance      = engine.TileSize / 4
	AppearSolidThreshold = 12
)

func (a *AppearBlock) Spawn(w *engine.World, s *engine.Spawnable, e *engine.Entity) error {
	a.World = w
	a.Entity = e

	var err error
	e.Image, err = engine.LoadImage("sprites", "appearblock.png")
	if err != nil {
		return err
	}
	e.Opaque = false
	e.Solid = false
	e.Alpha = 0.0
	return nil
}

func (a *AppearBlock) Despawn() {}

func (a *AppearBlock) Update() {
	delta := a.Entity.Rect.Delta(a.World.Player.Rect)
	if delta.DY > 0 && delta.DX <= AppearXDistance && delta.DX >= -AppearXDistance && delta.DY <= AppearYDistance && delta.DY >= -AppearYDistance {
		if a.AnimFrame < AppearFrames {
			a.AnimFrame++
		}
	} else {
		if a.AnimFrame > 0 {
			a.AnimFrame--
		}
	}
	a.Entity.Alpha = float64(a.AnimFrame) / AppearFrames
	// Make nonsolid if inside (to unstick player).
	a.Entity.Solid = a.AnimFrame >= AppearSolidThreshold && delta.DY > 0
}

func (a *AppearBlock) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&AppearBlock{})
}

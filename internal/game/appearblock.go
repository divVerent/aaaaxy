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
	AppearDistance       = engine.TileSize / 4
	AppearSolidThreshold = 8
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

func (a *AppearBlock) isNear(other *engine.Entity) bool {
	return a.Entity.Rect.Distance2(other.Rect) <= AppearDistance*AppearDistance
}

func (a *AppearBlock) Update() {
	if a.isNear(a.World.Player) {
		if a.AnimFrame < AppearFrames {
			a.AnimFrame++
		}
	} else {
		if a.AnimFrame > 0 {
			a.AnimFrame--
		}
	}
	a.Entity.Alpha = float64(a.AnimFrame) / AppearFrames
	a.Entity.Solid = a.AnimFrame >= AppearSolidThreshold
}

func (a *AppearBlock) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&AppearBlock{})
}

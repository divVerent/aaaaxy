package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	m "github.com/divVerent/aaaaaa/internal/math"
)

type Player struct {
	World  *engine.World
	Entity *engine.Entity

	OnGround bool
	Velocity m.Delta
	SubPixel m.Delta
}

// Player height is 30 px.
// So 30 px ~ 180 cm.
// Gravity is 9.81 m/s^2 = 163.5 px/s^2.
const (
	SubPixelScale  = 65536
	MaxGroundSpeed = 40 * SubPixelScale / engine.GameTPS
	GroundAccel    = 80 * SubPixelScale / engine.GameTPS / engine.GameTPS
	MaxAirSpeed    = 20 * SubPixelScale / engine.GameTPS
	AirAccel       = 40 * SubPixelScale / engine.GameTPS / engine.GameTPS
	JumpVelocity   = 200 * SubPixelScale / engine.GameTPS
	Gravity        = 160 * SubPixelScale / engine.GameTPS / engine.GameTPS

	KeyLeft  = ebiten.KeyLeft
	KeyRight = ebiten.KeyRight
	KeyUp    = ebiten.KeyUp
	KeyDown  = ebiten.KeyDown
	KeyJump  = ebiten.KeySpace
)

func (p *Player) Spawn(w *engine.World, s *engine.Spawnable, e *engine.Entity) error {
	p.World = w
	p.Entity = e
	var err error
	p.Entity.Image, err = engine.LoadImage("sprites", "player.png")
	if err != nil {
		return err
	}
	p.Entity.Rect.Size = m.Delta{DX: engine.PlayerWidth, DY: engine.PlayerHeight}
	return nil
}

func (p *Player) Despawn() {
	log.Panicf("the player should never despawn")
}

func (p *Player) Update() {
	log.Printf("initial velocity %v", p.Velocity)
	if p.OnGround {
		if ebiten.IsKeyPressed(KeyLeft) {
			p.Velocity.DX -= GroundAccel
			if p.Velocity.DX < -MaxGroundSpeed {
				p.Velocity.DX = -MaxGroundSpeed
			}
		}
		if ebiten.IsKeyPressed(KeyRight) {
			p.Velocity.DX += GroundAccel
			if p.Velocity.DX > MaxGroundSpeed {
				p.Velocity.DX = MaxGroundSpeed
			}
		}
		if ebiten.IsKeyPressed(KeyJump) {
			p.Velocity.DY += -JumpVelocity
			p.OnGround = false
		}
	} else {
		if ebiten.IsKeyPressed(KeyLeft) {
			p.Velocity.DX -= AirAccel
			if p.Velocity.DX < -MaxAirSpeed {
				p.Velocity.DX = -MaxAirSpeed
			}
		}
		if ebiten.IsKeyPressed(KeyRight) {
			p.Velocity.DX += AirAccel
			if p.Velocity.DX > MaxAirSpeed {
				p.Velocity.DX = MaxAirSpeed
			}
		}
	}
	p.Velocity.DY += Gravity
	log.Printf("final velocity %v", p.Velocity)
	p.SubPixel = p.SubPixel.Add(p.Velocity)
	move := p.SubPixel.Div(SubPixelScale)
	if move.DX != 0 {
		dest := p.Entity.Rect.Origin.Add(m.Delta{DX: move.DX})
		trace := p.World.TraceBox(p.Entity.Rect, dest, engine.TraceOptions{
			ForEnt: p.Entity,
		})
		if trace.EndPos == dest {
			// Nothing hit.
			p.SubPixel.DX -= move.DX * SubPixelScale
		} else {
			// Hit something. Move as far as we can in direction of the hit, but not farther than intended.
			if p.SubPixel.DX > SubPixelScale-1 {
				p.SubPixel.DX = SubPixelScale - 1
			} else if p.SubPixel.DX < 0 {
				p.SubPixel.DX = 0
			}
			p.Velocity.DX = 0
		}
		p.Entity.Rect.Origin = trace.EndPos
	}
	if move.DY != 0 {
		dest := p.Entity.Rect.Origin.Add(m.Delta{DY: move.DY})
		trace := p.World.TraceBox(p.Entity.Rect, dest, engine.TraceOptions{
			ForEnt: p.Entity,
		})
		if trace.EndPos == dest {
			// Nothing hit.
			p.SubPixel.DY -= move.DY * SubPixelScale
		} else {
			// Hit something. Move as far as we can in direction of the hit, but not farther than intended.
			if p.SubPixel.DY > SubPixelScale-1 {
				p.SubPixel.DY = SubPixelScale - 1
			} else if p.SubPixel.DY < 0 {
				p.SubPixel.DY = 0
			}
			p.Velocity.DY = 0
		}
		p.Entity.Rect.Origin = trace.EndPos
	}
	if p.OnGround {
		trace := p.World.TraceBox(p.Entity.Rect, p.Entity.Rect.Origin.Add(m.Delta{DX: 0, DY: 1}), engine.TraceOptions{
			ForEnt: p.Entity,
		})
		if trace.EndPos != p.Entity.Rect.Origin {
			log.Printf("left ground unexpectedly")
			p.OnGround = false
		}
	} else if p.SubPixel.DY == SubPixelScale-1 && p.Velocity.DY >= 0 {
		trace := p.World.TraceBox(p.Entity.Rect, p.Entity.Rect.Origin.Add(m.Delta{DX: 0, DY: 1}), engine.TraceOptions{
			ForEnt: p.Entity,
		})
		if trace.EndPos == p.Entity.Rect.Origin {
			log.Printf("hit ground unexpectedly")
			p.OnGround = true
		}
	}
}

func (p *Player) Touch(other *engine.Entity) {
	// Nothing happens; we rather handle this on other's Touch event.
}

func init() {
	engine.RegisterEntityType(&Player{})
}

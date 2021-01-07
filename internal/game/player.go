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
	Jumping  bool
	LookUp   bool
	LookDown bool
	Velocity m.Delta
	SubPixel m.Delta
}

// Player height is 30 px.
// So 30 px ~ 180 cm.
// Gravity is 9.81 m/s^2 = 163.5 px/s^2.
const (
	// PlayerWidth is the width of the player.
	PlayerWidth = engine.TileSize - 2
	// PlayerHeight is the height of the player.
	PlayerHeight = 2*engine.TileSize - 2
	// PlayerEyeDX is the X coordinate of the player's eye.
	PlayerEyeDX = engine.TileSize/2 - 1
	// PlayerEyeDY is the Y coordinate of the player's eye.
	PlayerEyeDY = engine.TileSize/2 - 1

	// LookTiles is how many tiles the player can look up/down.
	LookDistance = engine.TileSize * 4

	SubPixelScale = 65536

	// Nice run/jump speed.
	MaxGroundSpeed = 160 * SubPixelScale / engine.GameTPS
	GroundAccel    = 960 * SubPixelScale / engine.GameTPS / engine.GameTPS
	GroundFriction = 640 * SubPixelScale / engine.GameTPS / engine.GameTPS

	// Air control is lower than ground control.
	MaxAirSpeed = 120 * SubPixelScale / engine.GameTPS
	AirAccel    = 160 * SubPixelScale / engine.GameTPS / engine.GameTPS

	// We want 4.5 tiles high jumps, i.e. 72px high jumps (plus something).
	// Jump shall take 1 second.
	// Yields:
	// v0^2 / (2 * g) = 72
	// 2 v0 / g = 1
	// ->
	// v0 = 288
	// g = 576
	// Note: assuming 1px=6cm, this is actually 17.3m/s and 3.5x earth gravity.
	JumpVelocity = 288 * SubPixelScale / engine.GameTPS
	Gravity      = 576 * SubPixelScale / engine.GameTPS / engine.GameTPS
	MaxSpeed     = 2 * engine.TileSize * SubPixelScale

	// We want at least 19px high jumps so we can be sure a jump moves at least 2 tiles up.
	JumpExtraGravity = 72*Gravity/19 - Gravity

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
	p.Entity.Rect.Size = m.Delta{DX: PlayerWidth, DY: PlayerHeight}
	p.Entity.ZIndex = engine.MaxZIndex
	return nil
}

func (p *Player) Despawn() {
	log.Panicf("the player should never despawn")
}

func accelerate(vel *int, accel, max, dir int) {
	newVel := *vel + dir*accel
	if newVel*dir > max {
		newVel = max * dir
	}
	if newVel*dir > *vel*dir {
		*vel = newVel
	}
}

func friction(vel *int, friction int) {
	accelerate(vel, friction, 0, +1)
	accelerate(vel, friction, 0, -1)
}

func (p *Player) Update() {
	p.LookUp = ebiten.IsKeyPressed(KeyUp)
	p.LookDown = ebiten.IsKeyPressed(KeyDown)
	if ebiten.IsKeyPressed(KeyJump) {
		if !p.Jumping && p.OnGround {
			p.Velocity.DY -= JumpVelocity
			p.OnGround = false
			p.Jumping = true
		}
	} else {
		p.Jumping = false
	}
	if p.OnGround {
		friction(&p.Velocity.DX, GroundFriction)
		if ebiten.IsKeyPressed(KeyLeft) {
			accelerate(&p.Velocity.DX, GroundAccel, MaxGroundSpeed, -1)
		}
		if ebiten.IsKeyPressed(KeyRight) {
			accelerate(&p.Velocity.DX, GroundAccel, MaxGroundSpeed, +1)
		}
	} else {
		if ebiten.IsKeyPressed(KeyLeft) {
			accelerate(&p.Velocity.DX, AirAccel, MaxAirSpeed, -1)
		}
		if ebiten.IsKeyPressed(KeyRight) {
			accelerate(&p.Velocity.DX, AirAccel, MaxAirSpeed, +1)
		}
		if p.Velocity.DY < 0 && !p.Jumping {
			p.Velocity.DY += JumpExtraGravity
		}
	}
	p.Velocity.DY += Gravity
	if p.Velocity.DY > MaxSpeed {
		p.Velocity.DY = MaxSpeed
	}
	p.SubPixel = p.SubPixel.Add(p.Velocity)
	move := p.SubPixel.Div(SubPixelScale)
	if move.DX != 0 {
		dest := p.Entity.Rect.Origin.Add(m.Delta{DX: move.DX})
		trace := p.World.TraceBox(p.Entity.Rect, dest, engine.TraceOptions{
			IgnoreEnt: p.Entity,
		})
		p.handleTouch(trace)
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
			IgnoreEnt: p.Entity,
		})
		p.handleTouch(trace)
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
			p.OnGround = true
		}
		p.Entity.Rect.Origin = trace.EndPos
	} else if p.OnGround {
		trace := p.World.TraceBox(p.Entity.Rect, p.Entity.Rect.Origin.Add(m.Delta{DX: 0, DY: 1}), engine.TraceOptions{
			IgnoreEnt: p.Entity,
		})
		p.handleTouch(trace)
		if trace.EndPos != p.Entity.Rect.Origin {
			p.OnGround = false
		}
	}
}

func (p *Player) handleTouch(trace engine.TraceResult) {
	if trace.HitEntity != nil {
		trace.HitEntity.Impl.Touch(p.Entity)
	}
}

func (p *Player) Touch(other *engine.Entity) {
	// Nothing happens; we rather handle this on other's Touch event.
}

func (p *Player) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

// EyePos returns the position the player eye is at.
func (p *Player) EyePos() m.Pos {
	return p.Entity.Rect.Origin.Add(m.Delta{DX: PlayerEyeDX, DY: PlayerEyeDY})
}

// LookPos returns the position the player is focusing at.
func (p *Player) LookPos() m.Pos {
	focus := p.EyePos()
	if p.LookUp {
		focus.Y -= LookDistance
	}
	if p.LookDown {
		focus.Y += LookDistance
	}
	return focus
}

func init() {
	engine.RegisterEntityType(&Player{})
}

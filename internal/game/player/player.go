// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package player

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/animation"
	"github.com/divVerent/aaaaaa/internal/engine"
	m "github.com/divVerent/aaaaaa/internal/math"
	"github.com/divVerent/aaaaaa/internal/noise"
	"github.com/divVerent/aaaaaa/internal/sound"
)

type Player struct {
	World           *engine.World
	Entity          *engine.Entity
	PersistentState map[string]string

	OnGround   bool
	Jumping    bool
	JumpingUp  bool
	LookUp     bool
	LookDown   bool
	Velocity   m.Delta
	SubPixel   m.Delta
	Respawning bool

	Anim         animation.State
	JumpSound    *sound.Sound
	LandSound    *sound.Sound
	HitHeadSound *sound.Sound
}

// Player height is 30 px.
// So 30 px ~ 180 cm.
// Gravity is 9.81 m/s^2 = 163.5 px/s^2.
const (
	// PlayerWidth is the hitbox width of the player.
	// Actual width is 14 (one extra pixel to left and right).
	PlayerWidth = 12
	// PlayerHeight is the hitbox height of the player.
	// Actual height is 30 (three extra pixels at the top).
	PlayerHeight = 27
	// PlayerEyeDX is the X coordinate of the player's eye.
	PlayerEyeDX = 6
	// PlayerEyeDY is the Y coordinate of the player's eye.
	PlayerEyeDY = 4
	// PlayerOffsetDX is the player's render offset.
	PlayerOffsetDX = -1
	// PlayerOffsetDY is the player's render offset.
	PlayerOffsetDY = -3

	// LookTiles is how many tiles the player can look up/down.
	LookDistance = engine.TileSize * 4

	SubPixelScale = 65536

	// Nice run/jump speed.
	MaxGroundSpeed = 160 * SubPixelScale / engine.GameTPS
	GroundAccel    = 960 * SubPixelScale / engine.GameTPS / engine.GameTPS
	GroundFriction = 640 * SubPixelScale / engine.GameTPS / engine.GameTPS

	// Air max speed is lower than ground control, but same forward accel.
	MaxAirSpeed = 120 * SubPixelScale / engine.GameTPS
	AirAccel    = 320 * SubPixelScale / engine.GameTPS / engine.GameTPS

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

	NoiseMinSpeed = 384 * SubPixelScale / engine.GameTPS
	NoiseMaxSpeed = MaxSpeed

	// We want at least 19px high jumps so we can be sure a jump moves at least 2 tiles up.
	JumpExtraGravity = 72*Gravity/19 - Gravity

	// Animation tuning.
	AnimGroundSpeed = 20 * SubPixelScale / engine.GameTPS

	KeyLeft    = ebiten.KeyLeft
	KeyRight   = ebiten.KeyRight
	KeyUp      = ebiten.KeyUp
	KeyDown    = ebiten.KeyDown
	KeyJump    = ebiten.KeySpace
	KeyRespawn = ebiten.KeyR
)

func (p *Player) Spawn(w *engine.World, s *engine.Spawnable, e *engine.Entity) error {
	p.World = w
	p.Entity = e
	p.PersistentState = s.PersistentState
	p.Entity.Rect.Size = m.Delta{DX: PlayerWidth, DY: PlayerHeight}
	p.Entity.RenderOffset = m.Delta{DX: PlayerOffsetDX, DY: PlayerOffsetDY}
	p.Entity.ZIndex = engine.MaxZIndex

	p.Anim.Init("player", map[string]*animation.Group{
		"idle": {
			Frames:        2,
			FrameInterval: 172,
			NextInterval:  180,
			NextAnim:      "idle",
		},
		"walk": {
			Frames:        6,
			FrameInterval: 4,
			NextInterval:  4 * 6,
			NextAnim:      "walk",
		},
		"jump": {
			Frames:       1,
			NextInterval: 8,
			NextAnim:     "jump",
		},
		"land": {
			Frames:       1,
			NextInterval: 8,
			NextAnim:     "idle",
			WaitFinish:   true,
		},
		"hithead": {
			Frames:       1,
			NextInterval: 8,
			NextAnim:     "idle",
			WaitFinish:   true,
		}}, "idle")

	var err error
	p.JumpSound, err = sound.Load("jump.ogg")
	if err != nil {
		return fmt.Errorf("could not load jump sound: %v", err)
	}
	p.LandSound, err = sound.Load("land.ogg")
	if err != nil {
		return fmt.Errorf("could not load land sound: %v", err)
	}
	p.HitHeadSound, err = sound.Load("hithead.ogg")
	if err != nil {
		return fmt.Errorf("could not load hithead sound: %v", err)
	}

	return nil
}

func (p *Player) Despawn() {
	log.Panicf("The player should never despawn")
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
	if ebiten.IsKeyPressed(KeyRespawn) {
		if !p.Respawning {
			// TODO remove this debug hack, menu will do this instead. Maybe also a "death" routine.
			cpName := p.PersistentState["last_checkpoint"]
			cpFlipped := p.PersistentState["checkpoint_seen."+cpName] == "FlipX"
			p.World.RespawnPlayer(cpName, cpFlipped)
			return
		}
	} else {
		p.Respawning = false
	}
	p.LookUp = ebiten.IsKeyPressed(KeyUp)
	p.LookDown = ebiten.IsKeyPressed(KeyDown)
	moveLeft := ebiten.IsKeyPressed(KeyLeft)
	moveRight := ebiten.IsKeyPressed(KeyRight)
	if ebiten.IsKeyPressed(KeyJump) {
		if !p.Jumping && p.OnGround {
			p.Velocity.DY -= JumpVelocity
			p.OnGround = false
			p.Jumping = true
			p.JumpingUp = true
			p.JumpSound.Play()
		}
	} else {
		p.Jumping = false
	}
	if p.OnGround {
		maxSpeed := MaxGroundSpeed + GroundFriction
		if moveLeft {
			accelerate(&p.Velocity.DX, GroundAccel, maxSpeed, -1)
		}
		if moveRight {
			accelerate(&p.Velocity.DX, GroundAccel, maxSpeed, +1)
		}
		friction(&p.Velocity.DX, GroundFriction)
	} else {
		if moveLeft {
			accelerate(&p.Velocity.DX, AirAccel, MaxAirSpeed, -1)
		}
		if moveRight {
			accelerate(&p.Velocity.DX, AirAccel, MaxAirSpeed, +1)
		}
		if p.Velocity.DY < 0 && p.JumpingUp && !p.Jumping {
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
			ForEnt:    p.Entity,
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
		p.handleTouch(trace)
	}
	if move.DY != 0 {
		dest := p.Entity.Rect.Origin.Add(m.Delta{DY: move.DY})
		trace := p.World.TraceBox(p.Entity.Rect, dest, engine.TraceOptions{
			IgnoreEnt: p.Entity,
			ForEnt:    p.Entity,
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
			// If moving down, set OnGround flag.
			if move.DY > 0 {
				if !p.OnGround {
					p.Anim.SetGroup("land")
					p.LandSound.Play()
				}
				p.OnGround = true
				p.JumpingUp = false
			} else {
				p.Anim.SetGroup("hithead")
				p.HitHeadSound.Play()
			}
		}
		p.Entity.Rect.Origin = trace.EndPos
		p.handleTouch(trace)
	} else if p.OnGround {
		trace := p.World.TraceBox(p.Entity.Rect, p.Entity.Rect.Origin.Add(m.Delta{DX: 0, DY: 1}), engine.TraceOptions{
			IgnoreEnt: p.Entity,
			ForEnt:    p.Entity,
		})
		if trace.EndPos != p.Entity.Rect.Origin {
			p.OnGround = false
		}
		p.handleTouch(trace)
	}

	if moveLeft && !moveRight {
		p.Entity.Orientation = m.Identity()
	}
	if moveRight && !moveLeft {
		p.Entity.Orientation = m.FlipX()
	}
	if p.OnGround {
		if p.Velocity.DX > -AnimGroundSpeed && p.Velocity.DX < AnimGroundSpeed {
			p.Anim.SetGroup("idle")
		} else {
			p.Anim.SetGroup("walk")
		}
	} else {
		p.Anim.SetGroup("jump")
	}
	p.Anim.Update(p.Entity)
	speed := math.Sqrt(float64(p.Velocity.Length2()))
	if speed >= NoiseMinSpeed {
		amount := (speed - NoiseMinSpeed) / (NoiseMaxSpeed - NoiseMinSpeed)
		noise.Set(amount)
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

// Respawned informs the player that the world moved/respawned it.
func (p *Player) Respawned() {
	p.OnGround = true                // Do not get landing anim right away.
	p.Jumping = true                 // Jump key must be hit again.
	p.JumpingUp = false              // Do not assume we're in the first half of a jump (fastfall).
	p.Velocity = m.Delta{}           // Stop moving.
	p.SubPixel = m.Delta{}           // Stop moving.
	p.Respawning = true              // Block the respawn key until released.
	p.Anim.ForceGroup("idle")        // Reset animation.
	p.Entity.Image = nil             // Hide player until next Update.
	p.Entity.Orientation = m.FlipX() // Default to looking right.
}

func init() {
	engine.RegisterEntityType(&Player{})
}

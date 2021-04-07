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
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/animation"
	"github.com/divVerent/aaaaaa/internal/centerprint"
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/mixins"
	"github.com/divVerent/aaaaaa/internal/input"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
	"github.com/divVerent/aaaaaa/internal/noise"
	"github.com/divVerent/aaaaaa/internal/sound"
)

type Player struct {
	mixins.Physics
	World           *engine.World
	Entity          *engine.Entity
	PersistentState map[string]string

	AirFrames     int // Number of frames since last leaving ground.
	LastGroundPos m.Pos
	Jumping       bool
	JumpingUp     bool
	LookUp        bool
	LookDown      bool
	Respawning    bool
	WasOnGround   bool

	// Controlling Riser objects.
	CanCarry bool
	CanPush  bool
	CanStand bool

	Anim            animation.State
	JumpSound       *sound.Sound
	LandSound       *sound.Sound
	HitHeadSound    *sound.Sound
	GotAbilitySound *sound.Sound
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
	LookDistance = level.TileSize * 4

	// Nice run/jump speed.
	MaxGroundSpeed = 160 * mixins.SubPixelScale / engine.GameTPS
	GroundAccel    = GroundFriction + AirAccel
	GroundFriction = 640 * mixins.SubPixelScale / engine.GameTPS / engine.GameTPS

	// Air max speed is lower than ground control, but same forward accel.
	MaxAirSpeed = 120 * mixins.SubPixelScale / engine.GameTPS
	AirAccel    = 480 * mixins.SubPixelScale / engine.GameTPS / engine.GameTPS

	// We want 4.5 tiles high jumps, i.e. 72px high jumps (plus something).
	// Jump shall take 1 second.
	// Yields:
	// v0^2 / (2 * g) = 72
	// 2 v0 / g = 1
	// ->
	// v0 = 288
	// g = 576
	// Note: assuming 1px=6cm, this is actually 17.3m/s and 3.5x earth gravity.
	JumpVelocity = 288 * mixins.SubPixelScale / engine.GameTPS
	Gravity      = 576 * mixins.SubPixelScale / engine.GameTPS / engine.GameTPS
	MaxSpeed     = 2 * level.TileSize * mixins.SubPixelScale

	NoiseMinSpeed = 384 * mixins.SubPixelScale / engine.GameTPS
	NoiseMaxSpeed = MaxSpeed
	NoisePower    = 2.0

	// We want at least 19px high jumps so we can be sure a jump moves at least 2 tiles up.
	JumpExtraGravity = 72*Gravity/19 - Gravity

	// Number of frames to allow jumping after leaving ground. This is an extra 1/30 sec.
	// 7 allows reliable walking over 2 tile gaps.
	// 1 allows reliable walking over 1 tile gaps.
	// 0 allows some walking over 1 tile gaps.
	ExtraGroundFrames = 4

	// Animation tuning.
	AnimGroundSpeed = 20 * mixins.SubPixelScale / engine.GameTPS

	KeyLeft    = ebiten.KeyLeft
	KeyRight   = ebiten.KeyRight
	KeyUp      = ebiten.KeyUp
	KeyDown    = ebiten.KeyDown
	KeyJump    = ebiten.KeySpace
	KeyRespawn = ebiten.KeyR
)

func (p *Player) HasAbility(name string) bool {
	key := "can_" + name
	return p.PersistentState[key] == "true"
}

func (p *Player) GiveAbility(name, text string) {
	key := "can_" + name
	if p.PersistentState[key] == "true" {
		return
	}

	p.PersistentState[key] = "true"
	err := p.World.Save()
	if err != nil {
		log.Printf("Could not save game: %v", err)
		return
	}
	p.reloadAbilities()

	centerprint.New(text, centerprint.Important, centerprint.Middle, centerprint.BigFont(), color.NRGBA{R: 190, G: 0, B: 0, A: 255}).SetFadeOut(true)
	p.GotAbilitySound.Play()
}

func (p *Player) reloadAbilities() {
	p.CanCarry = p.HasAbility("carry")
	p.CanPush = p.HasAbility("push")
	p.CanStand = p.HasAbility("stand")
}

func (p *Player) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	p.Physics.Init(w, e, p.handleTouch)
	p.World = w
	p.Entity = e
	p.PersistentState = s.PersistentState
	p.Entity.Rect.Size = m.Delta{DX: PlayerWidth, DY: PlayerHeight}
	p.Entity.RenderOffset = m.Delta{DX: PlayerOffsetDX, DY: PlayerOffsetDY}
	p.Entity.ZIndex = engine.MaxZIndex
	p.Entity.RequireTiles = true // We're tracing, so we need our tiles to be loaded.
	w.SetSolid(p.Entity, true)   // Needed so platforms don't let players fall through.
	p.reloadAbilities()

	err := p.Anim.Init("player", map[string]*animation.Group{
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
	if err != nil {
		return fmt.Errorf("could not initialize player animation: %v", err)
	}

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
	p.GotAbilitySound, err = sound.Load("got_ability.ogg")
	if err != nil {
		return fmt.Errorf("could not load got_ability sound: %v", err)
	}

	// Reset as if after respawn.
	p.Respawned()

	return nil
}

func (p *Player) Despawn() {
	log.Panic("tried to despawn the player - this should never happen")
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
	p.LookUp = input.Up.Held
	p.LookDown = input.Down.Held
	moveLeft := input.Left.Held
	moveRight := input.Right.Held
	if input.Jump.Held {
		if !p.Jumping && p.AirFrames <= ExtraGroundFrames {
			p.Velocity.DY -= JumpVelocity
			p.OnGround = false
			p.AirFrames = ExtraGroundFrames + 1
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
	if p.AirFrames > ExtraGroundFrames {
		// No gravity while we still can jump.
		p.Velocity.DY += Gravity
	}
	if p.Velocity.DY > MaxSpeed {
		p.Velocity.DY = MaxSpeed
	}

	// Run physics.
	p.WasOnGround = p.OnGround
	p.Physics.Update() // May call handleTouch.

	if moveLeft && !moveRight {
		p.Entity.Orientation = m.Identity()
	}
	if moveRight && !moveLeft {
		p.Entity.Orientation = m.FlipX()
	}
	if p.OnGround {
		p.LastGroundPos = p.Entity.Rect.Origin
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
		amount := math.Pow((speed-NoiseMinSpeed)/(NoiseMaxSpeed-NoiseMinSpeed), NoisePower)
		noise.Set(amount)
	}
	if p.OnGround {
		p.AirFrames = 0
	} else {
		p.AirFrames++
	}
}

func (p *Player) handleTouch(trace engine.TraceResult) {
	if trace.HitDelta.DY > 0 {
		p.JumpingUp = false
	}
	if p.OnGround && !p.WasOnGround {
		p.Anim.SetGroup("land")
		p.LandSound.Play()
	}
	p.WasOnGround = p.OnGround
	if trace.HitDelta.DY < 0 {
		p.Anim.SetGroup("hithead")
		p.HitHeadSound.Play()
	}
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
	return m.Pos{
		X: p.Entity.Rect.Origin.X + PlayerEyeDX,
		Y: p.Entity.Rect.Origin.Y + PlayerEyeDY,
	}
}

// LookPos returns the position the player is focusing at.
func (p *Player) LookPos() m.Pos {
	focus := m.Pos{
		X: p.Entity.Rect.Origin.X + PlayerEyeDX,
		Y: p.LastGroundPos.Y + PlayerEyeDY,
	}
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
	p.Physics.Reset()                // Stop moving.
	p.LastGroundPos = p.EyePos()     // Center the camera.
	p.AirFrames = 0                  // Assume on ground.
	p.WasOnGround = p.OnGround       // Back to ground.
	p.Jumping = true                 // Jump key must be hit again.
	p.JumpingUp = false              // Do not assume we're in the first half of a jump (fastfall).
	p.Respawning = true              // Block the respawn key until released.
	p.Anim.ForceGroup("idle")        // Reset animation.
	p.Entity.Image = nil             // Hide player until next Update.
	p.Entity.Orientation = m.FlipX() // Default to looking right.
	p.reloadAbilities()              // Abilities may have changed.
}

func (p *Player) ActionPressed() bool {
	return input.Action.Held
}

func init() {
	engine.RegisterEntityType(&Player{})
}

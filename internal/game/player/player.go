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
	"math"
	"time"

	"github.com/divVerent/aaaaxy/internal/animation"
	"github.com/divVerent/aaaaxy/internal/centerprint"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/game/constants"
	"github.com/divVerent/aaaaxy/internal/game/interfaces"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/noise"
	"github.com/divVerent/aaaaxy/internal/palette"
	"github.com/divVerent/aaaaxy/internal/sound"
)

var (
	cheatInAirJump = flag.Bool("cheat_in_air_jump", false, "allow jumping while in air (allows getting anywhere)")
)

type Player struct {
	mixins.Physics
	World  *engine.World
	Entity *engine.Entity

	CoyoteFrames   int // Number of frames w/o gravity and w/ jumping. Goes down to -1 (0 is just timed out, -1 is normal)
	LastGroundPos  m.Pos
	Jumping        bool
	JumpingUp      bool
	LookUp         bool
	LookDown       bool
	Respawning     bool
	WasOnGround    bool
	PrevVelocity   m.Delta
	VVVVVV         bool
	JustSpawned    bool
	Goal           *engine.Entity
	EasterEggCount int

	Anim animation.State

	JumpSound       *sound.Sound
	VVVVVVSound     *sound.Sound
	LandSound       *sound.Sound
	HitHeadSound    *sound.Sound
	HitWallSound    *sound.Sound
	GotAbilitySound *sound.Sound
}

var _ interfaces.Abilityer = &Player{}
var _ interfaces.ActionPresseder = &Player{}
var _ interfaces.VVVVVVer = &Player{}

// Player height is 30 px.
// So 30 px ~ 180 cm.
// Gravity is 9.81 m/s^2 = 163.5 px/s^2.
const (
	// PlayerWidth is the hitbox width of the player.
	// Actual width is 14 (one extra pixel to left and right).
	PlayerWidth = 14 + 2*PlayerOffsetDX
	// PlayerHeight is the hitbox height of the player.
	// Actual height is 30 (three extra pixels at the top).
	PlayerHeight = 30 + PlayerOffsetDY + PlayerFlippedOffsetDY
	// PlayerEyeDX is the X coordinate of the player's eye.
	PlayerEyeDX = 5
	// PlayerEyeDY is the Y coordinate of the player's eye.
	PlayerEyeDY = 3
	// PlayerOffsetDX is the player's render offset.
	PlayerOffsetDX = -2
	// PlayerOffsetDY is the player's render offset.
	PlayerOffsetDY = -4
	// PlayerFlippedOffsetDY is the player's render offset.
	PlayerFlippedOffsetDY = -1
	// PlayerBorderPixels is the size of the player's black border.
	PlayerBorderPixels = 1

	// LookTiles is how many tiles the player can look up/down.
	LookDistance = level.TileSize * 4

	// Nice run/jump speed.
	MaxGroundSpeed = 160 * constants.SubPixelScale / engine.GameTPS
	GroundAccel    = GroundFriction + AirAccel
	GroundFriction = 640 * constants.SubPixelScale / engine.GameTPS / engine.GameTPS

	// Air max speed is lower than ground control, but same forward accel.
	MaxAirSpeed = 120 * constants.SubPixelScale / engine.GameTPS
	AirAccel    = 480 * constants.SubPixelScale / engine.GameTPS / engine.GameTPS

	// We want 4.5 tiles high jumps, i.e. 72px high jumps (plus something).
	// Jump shall take 1 second.
	// Yields:
	// v0^2 / (2 * g) = 72
	// 2 v0 / g = 1
	// ->
	// v0 = 288
	// g = 576
	// Note: assuming 1px=6cm, this is actually 17.3m/s and 3.5x earth gravity.
	JumpVelocity = 288 * constants.SubPixelScale / engine.GameTPS
	MaxSpeed     = 2 * level.TileSize * constants.SubPixelScale

	// Scale noise by speed.
	NoiseMinSpeed = 384 * constants.SubPixelScale / engine.GameTPS
	NoiseMaxSpeed = MaxSpeed
	NoisePower    = 2.0

	// Scale hitwall sound by speed.
	HitWallMinSpeed = 40 * constants.SubPixelScale / engine.GameTPS
	HitWallMaxSpeed = 160 * constants.SubPixelScale / engine.GameTPS

	// We want at least 19px high jumps so we can be sure a jump moves at least 2 tiles up.
	JumpExtraGravity = 72*constants.Gravity/19 - constants.Gravity

	// Number of frames to allow jumping after leaving ground. This is an extra 1/30 sec.
	// 8 allows reliable walking over 2 tile gaps.
	// 2 allows reliable walking over 1 tile gaps.
	// 1 allows some walking over 1 tile gaps.
	ExtraGroundFrames = 5

	// Animation tuning.
	AnimGroundSpeed = 20 * constants.SubPixelScale / engine.GameTPS
)

func (p *Player) SetVVVVVV(vvvvvv bool, up m.Delta, factor float64) {
	if vvvvvv == p.VVVVVV {
		return
	}
	if !up.IsZero() {
		p.OnGroundVec = up
	}
	p.VVVVVV = vvvvvv
	if !p.JustSpawned {
		p.VVVVVVSound.Play()
	}
	if factor != 1.0 {
		n := p.OnGroundVec
		velUp := n.Mul(p.Velocity.Dot(n))
		velUpScaled := velUp.MulFixed(m.NewFixedFloat64(factor))
		p.Velocity = p.Velocity.Add(velUpScaled.Sub(velUp))
	}
	p.LastGroundPos = p.Entity.Rect.Origin
	// Note: NOT resetting JumpingUp here.
	// This allows for an exploit where hitting a gravity flip retains
	// the increased gravity amount if the jump key is released.
	// A place in Vae Victis is much easier this way.
}

func (p *Player) HasAbility(name string) bool {
	return p.World.PlayerState.HasAbility(name)
}

func (p *Player) GiveAbility(name, text string) {
	if !p.World.PlayerState.GiveAbility(name) {
		return
	}

	err := p.World.Save()
	if err != nil {
		log.Errorf("could not save game: %v", err)
		return
	}

	centerprint.New(text, centerprint.Important, centerprint.Middle, centerprint.BigFont(), palette.EGA(palette.Red, 255), time.Second).SetFadeOut(true)
	p.GotAbilitySound.Play()
}

func (p *Player) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	p.Physics.Init(w, e, level.PlayerSolidContents, p.handleTouch)
	p.World = w
	p.Entity = e
	p.Entity.Rect.Size = m.Delta{DX: PlayerWidth, DY: PlayerHeight}
	p.Entity.RenderOffset = m.Delta{DX: PlayerOffsetDX, DY: PlayerOffsetDY}
	p.Entity.BorderPixels = PlayerBorderPixels
	w.SetZIndex(p.Entity, constants.PlayerZ)
	w.SetSolid(p.Entity, true) // Needed so platforms don't let players fall through.

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
		return fmt.Errorf("could not initialize player animation: %w", err)
	}

	p.JumpSound, err = sound.Load("jump.ogg")
	if err != nil {
		return fmt.Errorf("could not load jump sound: %w", err)
	}
	p.VVVVVVSound, err = sound.Load("vvvvvv.ogg")
	if err != nil {
		return fmt.Errorf("could not load vvvvvv sound: %w", err)
	}
	p.LandSound, err = sound.Load("land.ogg")
	if err != nil {
		return fmt.Errorf("could not load land sound: %w", err)
	}
	p.HitHeadSound, err = sound.Load("hithead.ogg")
	if err != nil {
		return fmt.Errorf("could not load hithead sound: %w", err)
	}
	p.HitWallSound, err = sound.Load("hitwall.ogg")
	if err != nil {
		return fmt.Errorf("could not load hitwall sound: %w", err)
	}
	p.GotAbilitySound, err = sound.Load("got_ability.ogg")
	if err != nil {
		return fmt.Errorf("could not load got_ability sound: %w", err)
	}

	// Reset as if after respawn.
	p.Respawned()

	return nil
}

func (p *Player) Despawn() {
	log.Fatalf("tried to despawn the player - this should never happen")
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
	p.JustSpawned = false
	var moveLeft, moveRight, jump bool
	if p.Goal == nil {
		p.LookUp = input.Up.Held
		p.LookDown = input.Down.Held
		moveLeft = input.Left.Held
		moveRight = input.Right.Held
		jump = input.Jump.Held
		action := input.Action.Held
		if p.LookUp || p.LookDown || moveLeft || moveRight || jump || action {
			p.World.TimerStarted = true
		}
	} else {
		// Walk towards goal!
		p.LookUp = false
		p.LookDown = false
		delta := p.Goal.Rect.Center().Delta(p.Entity.Rect.Center())
		moveLeft = delta.DX < 0
		moveRight = delta.DX > 0
		jump = false
	}
	if jump {
		if !p.Jumping && (p.CoyoteFrames > 0 || *cheatInAirJump) {
			p.Velocity = p.Velocity.Add(p.OnGroundVec.Mul(-JumpVelocity))
			p.OnGround = false
			p.CoyoteFrames = -1
			p.Jumping = true
			p.JumpingUp = true
			if p.VVVVVV {
				p.OnGroundVec = p.OnGroundVec.Mul(-1)
			}
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
		if p.Velocity.Dot(p.OnGroundVec) < 0 && p.JumpingUp && !p.Jumping {
			p.Velocity = p.Velocity.Add(p.OnGroundVec.Mul(JumpExtraGravity))
		}
	}
	if p.CoyoteFrames <= 0 {
		// No gravity while we still can jump.
		p.Velocity = p.Velocity.Add(p.OnGroundVec.Mul(constants.Gravity))
	}
	p.Velocity = p.Velocity.WithMaxLengthFixed(m.NewFixed(MaxSpeed))

	// Run physics.
	p.WasOnGround = p.OnGround
	p.PrevVelocity = p.Velocity
	p.Physics.Update() // May call handleTouch.

	if moveLeft && !moveRight {
		p.Entity.Orientation = m.Identity()
	}
	if moveRight && !moveLeft {
		p.Entity.Orientation = m.FlipX()
	}
	if p.OnGroundVec.Dot(p.Entity.Orientation.Down) < 0 {
		p.Entity.Orientation = p.Entity.Orientation.Concat(m.FlipY())
	}
	if p.OnGroundVec.DY < 0 {
		p.Entity.RenderOffset.DY = PlayerFlippedOffsetDY
	} else {
		p.Entity.RenderOffset.DY = PlayerOffsetDY
	}
	if p.OnGround {
		p.LastGroundPos = p.Entity.Rect.Origin
		if p.Velocity.DX > -AnimGroundSpeed && p.Velocity.DX < AnimGroundSpeed {
			p.Anim.SetGroup("idle")
		} else {
			p.Anim.SetGroup("walk")
		}
	} else {
		if p.VVVVVV {
			// Always update the scroll pos while in flipping mode.
			p.LastGroundPos = p.Entity.Rect.Origin
		}
		p.Anim.SetGroup("jump")
	}
	p.Anim.Update(p.Entity)
	speed := p.Velocity.Length()
	if speed >= NoiseMinSpeed {
		amount := math.Pow((speed-NoiseMinSpeed)/(NoiseMaxSpeed-NoiseMinSpeed), NoisePower)
		noise.Set(amount)
	}
	if p.OnGround {
		p.CoyoteFrames = ExtraGroundFrames
	} else if p.CoyoteFrames >= 0 {
		p.CoyoteFrames--
	}

	// Easter egg.
	// Doing this in player code so it only runs while the game is active.
	if input.EasterEggJustHit() {
		p.EasterEggCount++
		if p.EasterEggCount%4 == 0 {
			centerprint.New(locale.G.Get("Fine, I give up, have it your way.\nAll cheats are documented in --help."), centerprint.Important, centerprint.Top, centerprint.BigFont(), palette.EGA(palette.LightBlue, 255), 5*time.Second).SetFadeOut(true)
			p.GotAbilitySound.Play()
		} else {
			centerprint.New(locale.G.Get("You really thought this would do something?"), centerprint.Important, centerprint.Middle, centerprint.BigFont(), palette.EGA(palette.LightCyan, 255), time.Second).SetFadeOut(true)
			// No sound. We really want to do nothing here.
		}
	}

	// Konami code. Grants 30 lives. Too bad this game does not use lives :)
	if input.KonamiCodeJustHit() {
		centerprint.New(locale.G.Get("You now have 30 lives. Enjoy!"), centerprint.Important, centerprint.Top, centerprint.BigFont(), palette.EGA(palette.LightMagenta, 255), 5*time.Second).SetFadeOut(true)
		p.GotAbilitySound.Play()
	}
}

func (p *Player) handleTouch(trace engine.TraceResult) {
	if trace.HitDelta.Dot(p.OnGroundVec) > 0 {
		p.JumpingUp = false
	}
	if p.OnGround && !p.WasOnGround && p.CoyoteFrames < 0 {
		p.Anim.SetGroup("land")
		p.LandSound.Play()
	}
	if trace.HitDelta.Dot(p.OnGroundVec) < 0 {
		p.Anim.SetGroup("hithead")
		p.HitHeadSound.Play()
	}
	if trace.HitDelta.Dot(p.OnGroundVec) == 0 {
		dv := p.Velocity.Sub(p.PrevVelocity)
		speed := dv.Norm1()
		if speed > HitWallMinSpeed {
			vol := float64(speed-HitWallMinSpeed) / float64(HitWallMaxSpeed-HitWallMinSpeed)
			if vol > 1 {
				vol = 1
			}
			p.HitWallSound.PlayAtVolume(vol)
		}
	}
	p.World.TouchEvent(p.Entity, trace.HitEntities)

	// Update so we can get more deltas.
	p.WasOnGround = true
	p.PrevVelocity = p.Velocity
}

func (p *Player) Touch(other *engine.Entity) {
	// Nothing happens; we rather handle this on other's Touch event.
}

func (p *Player) eyeDY() int {
	if p.OnGroundVec.DY < 0 {
		return p.Entity.Rect.Size.DY - 1 - PlayerEyeDY
	} else {
		return PlayerEyeDY
	}
}

// EyePos returns the position the player eye is at.
func (p *Player) EyePos() m.Pos {
	return m.Pos{
		X: p.Entity.Rect.Origin.X + PlayerEyeDX,
		Y: p.Entity.Rect.Origin.Y + p.eyeDY(),
	}
}

// LookPos returns the position the player is focusing at.
func (p *Player) LookPos() m.Pos {
	focus := m.Pos{
		X: p.Entity.Rect.Origin.X + PlayerEyeDX,
		Y: p.LastGroundPos.Y + p.eyeDY(),
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
	p.Physics.Reset()                      // Stop moving.
	p.LastGroundPos = p.Entity.Rect.Origin // Center the camera.
	p.CoyoteFrames = ExtraGroundFrames     // Assume on ground.
	p.WasOnGround = p.OnGround             // Back to ground.
	p.Jumping = true                       // Jump key must be hit again.
	p.VVVVVV = false                       // Normal physics.
	p.OnGroundVec = m.Delta{DX: 0, DY: 1}  // Gravity points down.
	p.JumpingUp = false                    // Do not assume we're in the first half of a jump (fastfall).
	p.Respawning = true                    // Block the respawn key until released.
	p.Anim.ForceGroup("idle")              // Reset animation.
	p.Entity.Image = nil                   // Hide player until next Update.
	p.Entity.Orientation = m.FlipX()       // Default to looking right.
	p.Goal = nil                           // Normal input.
	p.JustSpawned = true                   // Just respawned.
}

func (p *Player) ActionPressed() bool {
	if p.Goal != nil {
		return false
	}
	return input.Action.Held
}

func (p *Player) SetVelocityForJump(velocity m.Delta) {
	p.Physics.SetVelocityForJump(velocity)
	p.JumpingUp = false
	p.CoyoteFrames = -1
}

func (p *Player) SetGoal(goal *engine.Entity) {
	p.Goal = goal
}

func init() {
	engine.RegisterEntityType(&Player{})
}

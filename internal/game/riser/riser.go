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

package riser

import (
	"fmt"
	"time"

	"github.com/divVerent/aaaaxy/internal/animation"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/constants"
	"github.com/divVerent/aaaaxy/internal/game/interfaces"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/level"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/sound"
)

type riserState int

const (
	Inactive riserState = iota
	IdlingUp
	MovingUp
	MovingLeft
	MovingRight
	GettingCarried
)

type Riser struct {
	mixins.Physics
	World  *engine.World
	Entity *engine.Entity

	NormalSize  m.Delta
	CarriedSize m.Delta

	State riserState

	Anim      animation.State
	FadeFrame int

	PlayerOnGroundVec m.Delta

	carrySound, pushSound, riseSound riserSound
}

const (
	// SmallRiserWidth is the hitbox width of the riser.
	// Actual width is 16 (one extra pixel to left and right).
	SmallRiserWidth = 14
	// SmallRiserHeight is the hitbox height of the riser.
	// Actual height is 16 (one extra pixel to left and right).
	SmallRiserHeight = 14
	// SmallRiserOffsetDX is the riser's render offset.
	SmallRiserOffsetDX = -1
	// SmallRiserOffsetDY is the riser's render offset.
	SmallRiserOffsetDY = -1
	// LargeRiserWidth is the hitbox width of the riser.
	// Actual width is 16 (one extra pixel to left and right).
	LargeRiserWidth = 30
	// LargeRiserHeight is the hitbox height of the riser.
	// Actual height is 32 (one extra pixel to left and right).
	LargeRiserHeight = 14
	// LargeRiserOffsetDX is the riser's render offset.
	LargeRiserOffsetDX = -1
	// LargeRiserOffsetDY is the riser's render offset.
	LargeRiserOffsetDY = -1
	// RiserBorderPixels is the riser's border size.
	RiserBorderPixels = 1
	// CarriedSmallRiserWidth is the hitbox width of a riser being carried.
	CarriedSmallRiserWidth = 8
	// CarriedSmallRiserHeight is the hitbox height of a riser being carried.
	CarriedSmallRiserHeight = 8
	// CarriedLargeRiserWidth is the hitbox width of a riser being carried.
	CarriedLargeRiserWidth = 24
	// CarriedLargeRiserHeight is the hitbox height of a riser being carried.
	CarriedLargeRiserHeight = 8

	// IdleSpeed is the speed the riser moves upwards when not used.
	IdleSpeed = 15 * constants.SubPixelScale / engine.GameTPS

	// UpSpeed is the speed the riser moves upwards when the player is standing on it.
	UpSpeed = 60 * constants.SubPixelScale / engine.GameTPS

	// SideSpeed is the speed of the riser when pushed away.
	SideSpeed = 60 * constants.SubPixelScale / engine.GameTPS

	// FadeFrames is how many frames risers take to fade in or out.
	FadeFrames = 16

	// FollowFactor is how fast a riser should follow the player per second.
	FollowFactor = 24.0

	// FollowMaxDistance is the max distance allowed while following the player.
	// Hardest part: Multi-Party Authorization.
	FollowMaxDistance = 24

	// RepelSpeed is the speed at which risers repel each other until they no longer overlap.
	RepelSpeed = 15 * constants.SubPixelScale / engine.GameTPS
)

func (r *Riser) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	r.Physics.Init(w, e, level.ObjectSolidContents, r.handleTouch)
	r.World = w
	r.Entity = e

	if r.Entity.Rect.Size.DY > 16 {
		// We actually allow this, but create a _small_ riser centered at the same place.
		// We can't use a large one as it might be in solid.
		// Do NOT adjust Y; we keep the upper edge as is,
		// as the upper edge is what matters for standing on it.
		r.Entity.Rect.Size.DY = 16
	}
	if r.Entity.Rect.Size.DY != 16 {
		return fmt.Errorf("unexpected riser height: got %v, want 16", r.Entity.Rect.Size.DY)
	}

	var sprite string
	switch r.Entity.Rect.Size.DX {
	case 16:
		r.NormalSize = m.Delta{DX: SmallRiserWidth, DY: SmallRiserHeight}
		r.CarriedSize = m.Delta{DX: CarriedSmallRiserWidth, DY: CarriedSmallRiserHeight}
		r.Entity.RenderOffset = m.Delta{DX: SmallRiserOffsetDX, DY: SmallRiserOffsetDY}
		sprite = "riser_small"
	case 32:
		r.NormalSize = m.Delta{DX: LargeRiserWidth, DY: LargeRiserHeight}
		r.CarriedSize = m.Delta{DX: CarriedLargeRiserWidth, DY: CarriedLargeRiserHeight}
		r.Entity.RenderOffset = m.Delta{DX: LargeRiserOffsetDX, DY: LargeRiserOffsetDY}
		sprite = "riser_large"
	default:
		return fmt.Errorf("unexpected riser width: got %v, want 16 or 32", r.Entity.Rect.Size.DX)
	}
	r.Entity.BorderPixels = RiserBorderPixels
	r.Entity.Rect.Origin = r.Entity.Rect.Origin.Sub(r.Entity.RenderOffset)
	r.Entity.Rect.Size = r.NormalSize
	w.SetZIndex(r.Entity, constants.RiserMovingZ)
	r.Entity.Alpha = 0 // We fade in.
	r.State = Inactive
	r.Entity.Orientation = m.Identity()

	if sp.Properties["flipped"] == "true" {
		r.OnGroundVec = r.OnGroundVec.Mul(-1)
	}

	err := r.Anim.Init(sprite, map[string]*animation.Group{
		"inactive": {
			Frames:        1,
			FrameInterval: 16,
			NextInterval:  16,
			NextAnim:      "inactive",
		},
		"idle": {
			Frames:        1,
			FrameInterval: 16,
			NextInterval:  16,
			NextAnim:      "idle",
		},
		"left": {
			Frames:        2,
			FrameInterval: 16,
			NextInterval:  32,
			NextAnim:      "left",
		},
		"right": {
			Frames:        2,
			FrameInterval: 16,
			NextInterval:  32,
			NextAnim:      "right",
		},
		"up": {
			Frames:        2,
			FrameInterval: 16,
			NextInterval:  32,
			NextAnim:      "up",
		},
	}, "inactive")
	if err != nil {
		return fmt.Errorf("could not initialize riser animation: %v", err)
	}

	r.carrySound, err = newRiserSound(1, "riser_carry")
	if err != nil {
		return err
	}
	r.pushSound, err = newRiserSound(2, "riser_push")
	if err != nil {
		return err
	}
	r.riseSound, err = newRiserSound(2, "riser_rise")
	if err != nil {
		return err
	}

	return nil
}

func (r *Riser) Despawn() {
	r.pushSound.despawn()
	r.carrySound.despawn()
	r.riseSound.despawn()
}

type riserSound struct {
	// Constants.
	minFrames int
	on, off   *sound.GroupedSound

	// State.
	counter int
}

func newRiserSound(minFrames int, name string) (riserSound, error) {
	onSnd, err := sound.Load(name + ".ogg")
	if err != nil {
		return riserSound{}, fmt.Errorf("could not load %v sound: %v", name, err)
	}
	offSnd, err := sound.Load(name + "_stop.ogg")
	if err != nil {
		return riserSound{}, fmt.Errorf("could not load %v_stop sound: %v", name, err)
	}
	return riserSound{
		minFrames: minFrames,
		on:        onSnd.Grouped(),
		off:       offSnd.Grouped(),
	}, nil
}

func (s *riserSound) update(cur bool) {
	if cur {
		s.counter++
		if s.counter == s.minFrames {
			s.on.Play()
			if !s.on.IsPlaying() {
				s.on.Reset()
			}
		}
	} else {
		if s.counter >= s.minFrames {
			s.on.Close()
			s.off.Play()
			if s.off.Current() >= 100*time.Millisecond {
				s.off.Reset()
			}
		}
		s.counter = 0
	}
}

func (s *riserSound) despawn() {
	s.on.Close()
	s.off.Close()
	s.counter = 0
}

func (r *Riser) Update() {
	playerAbilities := r.World.Player.Impl.(interfaces.Abilityer)
	playerButtons := r.World.Player.Impl.(interfaces.ActionPresseder)
	playerPhysics := r.World.Player.Impl.(interfaces.Physics)
	canCarry := playerAbilities.HasAbility("carry")
	canPush := playerAbilities.HasAbility("push")
	canStand := playerAbilities.HasAbility("stand")
	actionPressed := playerButtons.ActionPressed()
	playerOnMe := playerPhysics.ReadGroundEntity() == r.Entity
	playerDelta := r.World.Player.Rect.Delta(r.Entity.Rect)
	playerAboveMe := playerDelta.DX == 0 && playerDelta.Dot(r.OnGroundVec) < 0

	if canCarry && !playerOnMe && actionPressed && (playerDelta.IsZero() || (r.State == GettingCarried && playerDelta.Norm1() <= FollowMaxDistance)) {
		r.State = GettingCarried
	} else if canPush && actionPressed {
		if r.World.Player.Rect.Center().X < r.Entity.Rect.Center().X {
			r.State = MovingRight
		} else {
			r.State = MovingLeft
		}
	} else if canStand && playerAboveMe {
		r.State = MovingUp
	} else if canCarry || canPush || canStand {
		r.State = IdlingUp
	} else {
		r.State = Inactive
	}

	switch r.State {
	case Inactive:
		r.Anim.SetGroup("inactive")
		r.Velocity = m.Delta{}
	case IdlingUp:
		r.Anim.SetGroup("idle")
		r.Velocity = r.OnGroundVec.Mul(-IdleSpeed)
	case MovingUp:
		r.Anim.SetGroup("up")
		r.Velocity = r.OnGroundVec.Mul(-UpSpeed)
	case MovingLeft:
		r.Anim.SetGroup("left")
		r.Velocity = r.OnGroundVec.Mul(-IdleSpeed).Add(m.Delta{DX: -SideSpeed, DY: 0})
	case MovingRight:
		r.Anim.SetGroup("right")
		r.Velocity = r.OnGroundVec.Mul(-IdleSpeed).Add(m.Delta{DX: SideSpeed, DY: 0})
	case GettingCarried:
		r.Anim.SetGroup("idle")
		// r.Velocity = playerPhysics.ReadVelocity() // Hacky carry physics; good enough?
		pxDelta := r.World.Player.Rect.Center().Delta(r.Entity.Rect.Center())
		subDelta := playerPhysics.ReadSubPixel().Sub(r.SubPixel)
		fullDelta := pxDelta.Mul(constants.SubPixelScale).Add(subDelta)
		r.Velocity = fullDelta.MulFloat(FollowFactor / engine.GameTPS)

		if r.PlayerOnGroundVec.IsZero() {
			// All OK, just need to initialize grabbing.
		} else if r.PlayerOnGroundVec != playerPhysics.ReadOnGroundVec() {
			// Player's onground vec changed. Apply the change to ours.
			// TODO(divVerent): Actually make this smarter? Like, actually apply the transform by which the player changed?
			r.OnGroundVec = r.OnGroundVec.Mul(-1)
		}
		r.PlayerOnGroundVec = playerPhysics.ReadOnGroundVec()
	}
	if r.State == GettingCarried {
		// Never solid during carrying.
		r.World.MutateContents(r.Entity, level.SolidContents, 0)
	} else if canStand && playerAboveMe {
		// Solid to player when player is above.
		r.World.MutateContents(r.Entity, level.SolidContents, level.SolidContents)
	} else {
		// Otherwise, only solid to objects.
		r.World.MutateContents(r.Entity, level.SolidContents, level.ObjectSolidContents)
	}
	if playerOnMe || !canStand || r.State == GettingCarried {
		r.Physics.IgnoreEnt = r.World.Player // Move upwards despite player standing on it.
	} else {
		r.Physics.IgnoreEnt = nil
	}
	if r.State == GettingCarried {
		r.World.SetZIndex(r.Entity, constants.RiserCarriedZ)
	} else {
		r.World.SetZIndex(r.Entity, constants.RiserMovingZ)
	}

	// Adjust hitbox size.
	targetSize := r.NormalSize
	if r.State == GettingCarried {
		targetSize = r.CarriedSize
	}
	targetSizeChange := targetSize.Sub(r.Entity.Rect.Size)
	if r.ModifyHitBoxCentered(targetSizeChange) != targetSizeChange {
		// TODO(divVerent): Could also track the missing size change instead. That will require hitbox modifying to not always be centered.
		r.World.Detach(r.Entity)
	}

	// Also, risers that touch each other repel each other.
	r.World.ForEachEntity(func(other *engine.Entity) {
		if other == r.Entity {
			return
		}
		otherR, ok := other.Impl.(*Riser)
		if !ok {
			return
		}
		dr := r.Entity.Rect.Delta(other.Rect)
		if dr.Norm0() <= RiserBorderPixels {
			pxDelta := r.Entity.Rect.Center().Delta(other.Rect.Center())
			subDelta := r.SubPixel.Sub(otherR.SubPixel)
			fullDelta := pxDelta.Mul(constants.SubPixelScale).Add(subDelta)
			left := m.Right().Apply(r.OnGroundVec)
			var scaledDelta m.Delta
			if fullDelta.IsZero() {
				// On full overlap, move one of them _left/right_ which is the most gameplay friendly direction.
				// This will cause the other one to get repelled to the right.
				// Why not down? Yes, these are easier to grab then,
				// but it helps more that one can stand between them and push them away from each other.
				// How to pick left/right deterministically? We want no randomness...
				// So we pick them alternatingly, so eventually repelling will work.
				if r.World.PlayerState.Frames()%2 == 0 {
					scaledDelta = left.Mul(RepelSpeed)
				} else {
					scaledDelta = left.Mul(-RepelSpeed)
				}
			} else {
				// Similarly, ensure that the left/right component is at least half the work.
				// We do this by assuming an increased L/R delta of exactly the value of the U/D delta.
				upDown := fullDelta.Sub(left.Mul(fullDelta.Dot(left))).Norm1()
				if fullDelta.Dot(left) < 0 {
					upDown = -upDown
				}
				fullDelta = fullDelta.Add(left.Mul(upDown))
				scaledDelta = fullDelta.WithLength(RepelSpeed)
			}
			r.Velocity = r.Velocity.Add(scaledDelta)
		}
	})

	// Run physics.
	if !r.Velocity.IsZero() {
		r.Physics.Update() // May call handleTouch.
	}

	// Carry is a clear state.
	r.carrySound.update(r.State == GettingCarried)
	// For push and moving up, decide sound by whether we're actually moving.
	r.pushSound.update((r.State == MovingLeft || r.State == MovingRight) && r.Velocity.DX != 0)
	r.riseSound.update(r.State == MovingUp && (r.Velocity.DY == UpSpeed || r.Velocity.DY == -UpSpeed))

	r.Anim.Update(r.Entity)

	if r.OnGroundVec.DY < 0 {
		r.Entity.Orientation = m.FlipY()
	} else {
		r.Entity.Orientation = m.Identity()
	}

	if r.Entity.Detached() {
		if r.FadeFrame > 0 {
			r.FadeFrame--
		}
		if r.FadeFrame == 0 {
			r.World.Despawn(r.Entity)
		}
	} else {
		if r.FadeFrame < FadeFrames {
			r.FadeFrame++
		}
	}
	r.Entity.Alpha = float64(r.FadeFrame) / float64(FadeFrames)
}

func (r *Riser) handleTouch(trace engine.TraceResult) {
	// Risers can touch stuff. Gonna use this for switches.
	if trace.HitEntity != nil {
		r.World.TouchEvent(r.Entity, trace.HitEntity)
	}
}

func (r *Riser) Touch(other *engine.Entity) {
	// Nothing happens; we rather handle this on other's Touch event.
}

func init() {
	engine.RegisterEntityType(&Riser{})
}

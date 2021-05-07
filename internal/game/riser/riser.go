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
	"math"
	"math/rand"

	"github.com/divVerent/aaaaaa/internal/animation"
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/constants"
	"github.com/divVerent/aaaaaa/internal/game/interfaces"
	"github.com/divVerent/aaaaaa/internal/game/mixins"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
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
	World           *engine.World
	Entity          *engine.Entity
	PersistentState map[string]string

	State riserState

	Anim       animation.State
	Despawning bool
	FadeFrame  int
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

	// IdleSpeed is the speed the riser moves upwards when not used.
	IdleSpeed = 15 * mixins.SubPixelScale / engine.GameTPS

	// UpSpeed is the speed the riser moves upwards when the player is standing on it.
	UpSpeed = 60 * mixins.SubPixelScale / engine.GameTPS

	// SideSpeed is the speed of the riser when pushed away.
	SideSpeed = 60 * mixins.SubPixelScale / engine.GameTPS

	// FadeFrames is how many frames risers take to fade in or out.
	FadeFrames = 16

	// FollowFactor is how fast a riser should follow the player per second.
	FollowFactor = 24.0

	// FollowMaxDistance is the max distance allowed while following the player.
	FollowMaxDistance = 32

	// RepelSpeed is the speed at which risers repel each other until they no longer overlap.
	RepelSpeed = 15 * mixins.SubPixelScale / engine.GameTPS
)

func (r *Riser) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	r.Physics.Init(w, e, level.ObjectSolidContents, r.handleTouch)
	r.World = w
	r.Entity = e

	if r.Entity.Rect.Size.DY != 16 {
		return fmt.Errorf("unexpected riser height: got %v, want 16", r.Entity.Rect.Size.DY)
	}

	var sprite string
	switch r.Entity.Rect.Size.DX {
	case 16:
		r.Entity.Rect.Size = m.Delta{DX: SmallRiserWidth, DY: SmallRiserHeight}
		r.Entity.RenderOffset = m.Delta{DX: SmallRiserOffsetDX, DY: SmallRiserOffsetDY}
		sprite = "riser_small"
	case 32:
		r.Entity.Rect.Size = m.Delta{DX: LargeRiserWidth, DY: LargeRiserHeight}
		r.Entity.RenderOffset = m.Delta{DX: LargeRiserOffsetDX, DY: LargeRiserOffsetDY}
		sprite = "riser_large"
	default:
		return fmt.Errorf("unexpected riser width: got %v, want 16 or 32", r.Entity.Rect.Size.DX)
	}

	r.Entity.Rect.Origin = r.Entity.Rect.Origin.Sub(r.Entity.RenderOffset)
	w.SetZIndex(r.Entity, constants.RiserMovingZ)
	r.Entity.RequireTiles = true // We're tracing, so we need our tiles to be loaded.
	r.Entity.Alpha = 0           // We fade in.
	r.State = Inactive
	r.Entity.Orientation = m.Identity()

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

	return nil
}

func (r *Riser) Despawn() {}

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
	playerAboveMe := playerDelta.DX == 0 && playerDelta.DY < 0

	if canCarry && !playerOnMe && actionPressed && (playerDelta == m.Delta{} || (r.State == GettingCarried && playerDelta.Norm1() <= FollowMaxDistance)) {
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
		r.Velocity = m.Delta{DX: 0, DY: -IdleSpeed}
	case MovingUp:
		r.Anim.SetGroup("up")
		r.Velocity = m.Delta{DX: 0, DY: -UpSpeed}
	case MovingLeft:
		r.Anim.SetGroup("left")
		r.Velocity = m.Delta{DX: -SideSpeed, DY: -IdleSpeed}
	case MovingRight:
		r.Anim.SetGroup("right")
		r.Velocity = m.Delta{DX: SideSpeed, DY: -IdleSpeed}
	case GettingCarried:
		r.Anim.SetGroup("idle")
		// r.Velocity = playerPhysics.ReadVelocity() // Hacky carry physics; good enough?
		pxDelta := r.World.Player.Rect.Center().Delta(r.Entity.Rect.Center())
		subDelta := playerPhysics.ReadSubPixel().Sub(r.SubPixel)
		fullDelta := pxDelta.Mul(mixins.SubPixelScale).Add(subDelta)
		r.Velocity = fullDelta.MulFloat(FollowFactor / engine.GameTPS)
	}
	r.World.MutateContentsBool(r.Entity, level.PlayerSolidContents, canStand && playerAboveMe && r.State != GettingCarried)
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
		if dr == (m.Delta{}) {
			pxDelta := r.Entity.Rect.Center().Delta(other.Rect.Center())
			subDelta := r.SubPixel.Sub(otherR.SubPixel)
			fullDelta := pxDelta.Mul(mixins.SubPixelScale).Add(subDelta)
			var scaledDelta m.Delta
			if fullDelta == (m.Delta{}) {
				angle := rand.Float64() * (2 * math.Pi)
				scaledDelta = m.Delta{
					DX: int(RepelSpeed * math.Cos(angle)),
					DY: int(RepelSpeed * math.Sin(angle)),
				}
			} else {
				scaledDelta = fullDelta.MulFloat(RepelSpeed / math.Sqrt(float64(fullDelta.Length2())))
			}
			r.Velocity = r.Velocity.Add(scaledDelta)
		}
	})

	// Run physics.
	if (r.Velocity != m.Delta{}) {
		r.Physics.Update() // May call handleTouch.
	}

	r.Anim.Update(r.Entity)

	if r.Despawning {
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
		trace.HitEntity.Impl.Touch(r.Entity)
	}
}

func (r *Riser) Touch(other *engine.Entity) {
	// Nothing happens; we rather handle this on other's Touch event.
}

func init() {
	engine.RegisterEntityType(&Riser{})
}

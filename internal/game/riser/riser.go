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

	Anim animation.State
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
)

func (r *Riser) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	r.Physics.Init(w, e, level.ObjectSolidContents, r.handleTouch)
	r.World = w
	r.Entity = e

	var sprite string
	if s.Properties["large"] == "true" {
		r.Entity.Rect.Size = m.Delta{DX: LargeRiserWidth, DY: LargeRiserHeight}
		r.Entity.RenderOffset = m.Delta{DX: LargeRiserOffsetDX, DY: LargeRiserOffsetDY}
		sprite = "riser_large"
	} else {
		r.Entity.Rect.Size = m.Delta{DX: SmallRiserWidth, DY: SmallRiserHeight}
		r.Entity.RenderOffset = m.Delta{DX: SmallRiserOffsetDX, DY: SmallRiserOffsetDY}
		sprite = "riser_small"
	}
	r.Entity.Rect.Origin = r.Entity.Rect.Origin.Sub(r.Entity.RenderOffset)
	w.SetZIndex(r.Entity, constants.RiserZ)
	r.Entity.RequireTiles = true // We're tracing, so we need our tiles to be loaded.
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

func (r *Riser) isBelow(other *engine.Entity) bool {
	return other.Rect.OppositeCorner().Y < r.Entity.Rect.Origin.Y
}

func (r *Riser) Update() {
	playerAbilities := r.World.Player.Impl.(interfaces.Abilityer)
	playerButtons := r.World.Player.Impl.(interfaces.ActionPresseder)
	playerPhysics := r.World.Player.Impl.(interfaces.Physics)
	canCarry := playerAbilities.HasAbility("carry")
	canPush := playerAbilities.HasAbility("push")
	canStand := playerAbilities.HasAbility("stand")
	actionPressed := playerButtons.ActionPressed()
	playerGround := playerPhysics.ReadGroundEntity()

	if canCarry && playerGround != r.Entity && (r.World.Player.Rect.Delta(r.Entity.Rect) == m.Delta{}) && actionPressed {
		r.State = GettingCarried
	} else if canPush && actionPressed {
		if r.World.Player.Rect.Center().X < r.Entity.Rect.Center().X {
			r.State = MovingRight
		} else {
			r.State = MovingLeft
		}
	} else if canStand && playerGround == r.Entity {
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
		r.World.MutateContentsBool(r.Entity, level.PlayerSolidContents, false)
		r.Physics.IgnoreEnt = nil // Should never hit player.
	case IdlingUp:
		r.Anim.SetGroup("idle")
		r.Velocity = m.Delta{DX: 0, DY: -IdleSpeed}
		r.World.MutateContentsBool(r.Entity, level.PlayerSolidContents, canStand && r.isBelow(r.World.Player))
		r.Physics.IgnoreEnt = nil // Stop at player!
	case MovingUp:
		r.Anim.SetGroup("up")
		r.Velocity = m.Delta{DX: 0, DY: -UpSpeed}
		r.World.MutateContentsBool(r.Entity, level.PlayerSolidContents, canStand && r.isBelow(r.World.Player))
		r.Physics.IgnoreEnt = r.World.Player // Move upwards despite player standing on it.
	case MovingLeft:
		r.Anim.SetGroup("left")
		r.Velocity = m.Delta{DX: -SideSpeed, DY: -IdleSpeed}
		r.World.MutateContentsBool(r.Entity, level.PlayerSolidContents, canStand && r.isBelow(r.World.Player))
		r.Physics.IgnoreEnt = r.World.Player // Player may be standing on it.
	case MovingRight:
		r.Anim.SetGroup("right")
		r.Velocity = m.Delta{DX: SideSpeed, DY: -IdleSpeed}
		r.World.MutateContentsBool(r.Entity, level.PlayerSolidContents, canStand && r.isBelow(r.World.Player))
		r.Physics.IgnoreEnt = r.World.Player // Player may be standing on it.
	case GettingCarried:
		r.Anim.SetGroup("idle")
		r.Velocity = playerPhysics.ReadVelocity() // Hacky carry physics; good enough?
		r.World.MutateContentsBool(r.Entity, level.PlayerSolidContents, false)
		r.Physics.IgnoreEnt = r.World.Player // May collide with player.
	}

	// Run physics.
	if (r.Velocity != m.Delta{}) {
		r.Physics.Update() // May call handleTouch.
	}

	r.Anim.Update(r.Entity)
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

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

package trigger

import (
	"fmt"
	"math"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/constants"
	"github.com/divVerent/aaaaxy/internal/game/interfaces"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/level"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/sound"
)

// JumpPad, when hit by the player, sends the player on path to set destination.
// Note that sadly, JumpPads are rarely ever useful in rooms that can be used in multiple orientations.
// May want to introduce required orientation like with checkpoints.
// Or could require player to hit jumppad from above.
type JumpPad struct {
	mixins.NonSolidTouchable
	World  *engine.World
	Entity *engine.Entity

	Destination m.Pos
	Height      int

	TouchedFrame int
	JumpSound    *sound.Sound
}

func (j *JumpPad) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	j.NonSolidTouchable.Init(w, e)
	j.World = w
	j.Entity = e
	w.SetOpaque(e, false)
	w.SetSolid(e, sp.Properties["solid"] != "false") // Default true.

	var delta m.Delta
	_, err := fmt.Sscanf(sp.Properties["delta"], "%d %d", &delta.DX, &delta.DY)
	if err != nil {
		return fmt.Errorf("failed to parse delta: %v", err)
	}
	var relDelta m.Delta
	_, err = fmt.Sscanf(sp.Properties["rel_delta"], "%d %d", &relDelta.DX, &relDelta.DY)
	if err != nil && sp.Properties["rel_delta"] != "" {
		return fmt.Errorf("failed to parse absolute delta: %v", err)
	}
	// Destination is actually measured from center of trigger; need to transform to worldspace.
	j.Destination = e.Rect.Center().Add(e.Transform.Inverse().Apply(delta)).Add(relDelta)
	_, err = fmt.Sscanf(sp.Properties["height"], "%d", &j.Height)
	if err != nil {
		return fmt.Errorf("failed to parse height: %v", err)
	}

	j.JumpSound, err = sound.Load("jump.ogg")
	if err != nil {
		return fmt.Errorf("could not load jump sound: %v", err)
	}

	return nil
}

func (j *JumpPad) Despawn() {}

func (j *JumpPad) Update() {
	j.NonSolidTouchable.Update()
	if j.TouchedFrame > 0 {
		j.TouchedFrame--
	}
}

func calculateJump(delta m.Delta, heightParam int) m.Delta {
	apexOutside := heightParam < 0
	// Convert to relative height above jump: always negative (up).
	var height int
	if apexOutside {
		height = heightParam
	} else {
		height = -heightParam
	}
	// Height is meant above entire jump.
	targetHigher := delta.DY < 0
	if targetHigher {
		height += delta.DY
	}
	// Requirements:
	// - vDY * tA + 1/2 * playerGravity * tA^2 = height * SubpixelScale
	//   -> vDY = -sqrt(-2 * height * playerGravity * SubpixelScale)
	// - vDY + playerGravity * tA = 0
	//   -> tA = -vDY / playerGravity
	vDY := -int(math.Sqrt(2 * float64(-height) * float64(constants.Gravity) * float64(constants.SubPixelScale)))
	// Actually move downwards if requested!
	if apexOutside && !targetHigher {
		vDY = -vDY
	}
	// Finally:
	// - vDY * t + 1/2 * playerGravity * t^2 = deltaDY * SubpixelScale
	// - vDX * t = deltaDX
	a := 0.5 * constants.Gravity
	b := float64(vDY)
	c := -float64(delta.DY) * constants.SubPixelScale
	u := -b / (2 * a)
	d := b*b - 4*a*c
	v := 0.0
	if d >= 0 {
		// Mathematically, D < 0 means the jump is impossible.
		// However usually it just implies a roundoff error,
		// especially when height==0. So let's just allow it.
		v = math.Sqrt(d) / (2 * a)
	}
	if apexOutside && !targetHigher {
		v = -v
	}
	t := u + v
	vDX := int(float64(delta.DX) * constants.SubPixelScale / t)
	return m.Delta{DX: vDX, DY: vDY}
}

func (j *JumpPad) Touch(other *engine.Entity) {
	// Do we actually touch the player?
	if other != j.World.Player {
		return
	}
	p := other.Impl.(interfaces.Physics)
	// Can not touch from below (not gonna work anyway).
	if other.Rect.Delta(j.Entity.Rect).Dot(p.ReadOnGroundVec()) > 0 {
		return
	}
	// Compute parameters for jump.
	source := other.Rect.Foot()
	dest := j.Destination
	delta := dest.Delta(source)
	// Can't jump from the "opposite side" of the jumppad (not gonna work either).
	if delta.DX >= 0 {
		if other.Rect.OppositeCorner().X < j.Entity.Rect.Origin.X {
			return
		}
	}
	if delta.DX <= 0 {
		if other.Rect.Origin.X > j.Entity.Rect.OppositeCorner().X {
			return
		}
	}
	// Require player to leave before jumping the player again.
	prevTouchedFrame := j.TouchedFrame
	j.TouchedFrame = 2
	if prevTouchedFrame > 0 {
		return
	}
	// Perform the jump.
	if p.ReadOnGroundVec().DY < 0 {
		// HACK: Can we rather support arbitrary OnGroundVec?
		p.SetVelocityForJump(m.FlipY().Apply(calculateJump(m.FlipY().Apply(delta), j.Height)))
	} else {
		p.SetVelocityForJump(calculateJump(delta, j.Height))
	}
	j.JumpSound.Play()
}

func init() {
	engine.RegisterEntityType(&JumpPad{})
}

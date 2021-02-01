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

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/mixins"
	"github.com/divVerent/aaaaaa/internal/game/player"
	m "github.com/divVerent/aaaaaa/internal/math"
	"github.com/divVerent/aaaaaa/internal/sound"
)

// JumpPad, when hit by the player, sends the player on path to set destination.
// Note that sadly, JumpPads are rarely ever useful in rooms that can be used in multiple orientations.
// May want to introduce required orientation like with checkpoints.
// Or could require player to hit jumppad from above.
type JumpPad struct {
	mixins.Settable
	World  *engine.World
	Entity *engine.Entity

	Destination m.Pos
	Height      int

	TouchedFrame int
	JumpSound    *sound.Sound
}

func (j *JumpPad) Spawn(w *engine.World, s *engine.Spawnable, e *engine.Entity) error {
	j.Settable.Init(s)
	j.World = w
	j.Entity = e
	e.Opaque = false
	e.Solid = true

	var delta m.Delta
	_, err := fmt.Sscanf(s.Properties["delta"], "%d %d", &delta.DX, &delta.DY)
	if err != nil {
		return fmt.Errorf("failed to parse delta: %v", err)
	}
	var relDelta m.Delta
	_, err = fmt.Sscanf(s.Properties["rel_delta"], "%d %d", &relDelta.DX, &relDelta.DY)
	if err != nil && s.Properties["rel_delta"] != "" {
		return fmt.Errorf("failed to parse absolute delta: %v", err)
	}
	// Destination is actually measured from center of trigger; need to transform to worldspace.
	j.Destination = e.Rect.Center().Add(e.Transform.Inverse().Apply(delta)).Add(relDelta)
	_, err = fmt.Sscanf(s.Properties["height"], "%d", &j.Height)
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
	vDY := -int(math.Sqrt(2 * float64(-height) * float64(player.Gravity) * float64(player.SubPixelScale)))
	// Actually move downwards if requested!
	if apexOutside && !targetHigher {
		vDY = -vDY
	}
	// Finally:
	// - vDY * t + 1/2 * playerGravity * t^2 = deltaDY * SubpixelScale
	// - vDX * t = deltaDX
	a := 0.5 * player.Gravity
	b := float64(vDY)
	c := -float64(delta.DY) * player.SubPixelScale
	u := -b / (2 * a)
	v := math.Sqrt(b*b-4*a*c) / (2 * a)
	if apexOutside && !targetHigher {
		v = -v
	}
	t := u + v
	vDX := int(float64(delta.DX) * player.SubPixelScale / t)
	return m.Delta{DX: vDX, DY: vDY}
}

func (j *JumpPad) Touch(other *engine.Entity) {
	if !j.Settable.State {
		return
	}
	// Do we actually touch the player?
	p, ok := other.Impl.(*player.Player)
	if !ok {
		return
	}
	// Can not touch from below (not gonna work anyway).
	if other.Rect.Origin.Y > j.Entity.Rect.OppositeCorner().Y {
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
	p.Velocity = calculateJump(delta, j.Height)
	p.OnGround = false
	j.JumpSound.Play()
}

func (j *JumpPad) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&JumpPad{})
}

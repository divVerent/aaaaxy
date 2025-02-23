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

package mixins

import (
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/constants"
	"github.com/divVerent/aaaaxy/internal/game/interfaces"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
)

type Physics struct {
	World  *engine.World
	Entity *engine.Entity

	// StepHeight is the number of pixels to allow stepping up/down.
	StepHeight int

	Contents        level.Contents
	OnGround        bool
	OnGroundVec     m.Delta // Vector that points "down" in gravity direction.
	GroundEntity    *engine.Entity
	Velocity        m.Delta // An input to be set changed by caller.
	SubPixel        m.Delta
	IgnoreEnt       *engine.Entity
	handleTouchFunc func(trace engine.TraceResult)
}

type trivialPhysics struct {
	engine.EntityImpl
	Physics
}

func (t *trivialPhysics) Update() {
	t.Physics.Update()
}

var _ interfaces.Physics = &trivialPhysics{}

func (p *Physics) Init(w *engine.World, e *engine.Entity, contents level.Contents, handleTouch func(trace engine.TraceResult)) {
	p.World = w
	p.Entity = e
	p.Contents = contents
	p.handleTouchFunc = handleTouch
	p.OnGroundVec = m.Delta{DX: 0, DY: 1}

	// We're tracing, so we need our tiles to be loaded.
	p.Entity.RequireTiles = true

	// Set initial subpixel to be in the center of the start pixel.
	p.SubPixel = m.Delta{DX: constants.SubPixelScale / 2, DY: constants.SubPixelScale / 2}
}

func (p *Physics) Reset() {
	p.OnGround = true
	p.GroundEntity = nil
	p.Velocity = m.Delta{}
	p.SubPixel = m.Delta{DX: constants.SubPixelScale / 2, DY: constants.SubPixelScale / 2}
}

func (p *Physics) traceMove(contents level.Contents, move m.Delta) engine.TraceResult {
	dest := p.Entity.Rect.Origin.Add(move)
	trace := p.World.TraceBox(p.Entity.Rect, dest, engine.TraceOptions{
		Contents:  contents,
		IgnoreEnt: p.IgnoreEnt,
		ForEnt:    p.Entity,
		LoadTiles: true,
	})
	return trace
}

func (p *Physics) tryMove(move m.Delta, stepping bool) (m.Delta, bool, *engine.TraceResult) {
	groundChecked := false
	trace := p.traceMove(p.Contents, move)
	if trace.HitDelta.IsZero() {
		// Nothing hit. We're done.
		if !stepping {
			p.SubPixel.DX -= move.DX * constants.SubPixelScale
			p.SubPixel.DY -= move.DY * constants.SubPixelScale
		}
		p.Entity.Rect.Origin = trace.EndPos
		if move.Dot(p.OnGroundVec) != 0 {
			// If move had a Y component, we're flying.
			p.OnGround, p.GroundEntity, groundChecked = false, nil, true
		}
		return m.Delta{DX: 0, DY: 0}, groundChecked, nil
	}
	var hitEntity *engine.Entity
	if len(trace.HitEntities) != 0 {
		hitEntity = trace.HitEntities[0]
	}
	if trace.HitDelta.DX != 0 {
		// An X hit. Just adjust X subpixel to be as close to the hit as possible.
		if p.SubPixel.DX > constants.SubPixelScale-1 {
			p.SubPixel.DX = constants.SubPixelScale - 1
		} else if p.SubPixel.DX < 0 {
			p.SubPixel.DX = 0
		}
		if !stepping {
			p.SubPixel.DY -= (trace.EndPos.Y - p.Entity.Rect.Origin.Y) * constants.SubPixelScale
			p.Velocity.DX = 0
		}
		move.DX = 0
		move.DY -= trace.EndPos.Y - p.Entity.Rect.Origin.Y
		p.Entity.Rect.Origin = trace.EndPos

		// Just in case we have left/right gravity... (not yet).
		if trace.HitDelta.Dot(p.OnGroundVec) > 0 {
			p.OnGround, p.GroundEntity, groundChecked = true, hitEntity, true
		} else if trace.HitDelta.Dot(p.OnGroundVec) < 0 {
			p.OnGround, p.GroundEntity, groundChecked = false, nil, true
		}
	} else if trace.HitDelta.DY != 0 {
		// A Y hit. Also update ground status.
		if p.SubPixel.DY > constants.SubPixelScale-1 {
			p.SubPixel.DY = constants.SubPixelScale - 1
		} else if p.SubPixel.DY < 0 {
			p.SubPixel.DY = 0
		}
		if !stepping {
			p.SubPixel.DX -= (trace.EndPos.X - p.Entity.Rect.Origin.X) * constants.SubPixelScale
			p.Velocity.DY = 0
		}
		move.DX -= trace.EndPos.X - p.Entity.Rect.Origin.X
		move.DY = 0
		p.Entity.Rect.Origin = trace.EndPos

		if trace.HitDelta.Dot(p.OnGroundVec) > 0 {
			p.OnGround, p.GroundEntity, groundChecked = true, hitEntity, true
		} else if trace.HitDelta.Dot(p.OnGroundVec) < 0 {
			p.OnGround, p.GroundEntity, groundChecked = false, nil, true
		}
	}
	return move, groundChecked, &trace
}

func (p *Physics) slideMove(move m.Delta) bool {
	groundChecked := false
	for !move.IsZero() {
		var ground bool
		var trace *engine.TraceResult
		move, ground, trace = p.tryMove(move, false)
		preTouchVel := p.Velocity
		if trace != nil {
			p.handleTouchFunc(*trace)
		}
		teleported := p.Velocity != preTouchVel
		groundChecked = groundChecked || ground

		if teleported {
			// This happens only if handleTouchFunc changed velocity.
			break
		}
	}
	return groundChecked
}

func (p *Physics) walkMove(move m.Delta) bool {
	if p.StepHeight == 0 {
		return p.slideMove(move)
	}

	prevOrigin := p.Entity.Rect.Origin

	wasOnGround := p.OnGround

	groundChecked := false
	for !move.IsZero() {
		var ground bool
		var trace *engine.TraceResult
		prevGoal := p.Entity.Rect.Origin.Add(move)
		prevVel := p.Velocity
		move, ground, trace = p.tryMove(move, false)
		preTouchVel := p.Velocity
		if trace != nil {
			p.handleTouchFunc(*trace)
		}
		teleported := p.Velocity != preTouchVel
		groundChecked = groundChecked || ground

		if teleported {
			// This happens only if handleTouchFunc changed velocity.
			break
		}

		if trace != nil && trace.HitDelta.Dot(p.OnGroundVec) == 0 {
			log.Debugf("walkMove: trying to upstep")
			groundChecked = false

			// Hit a vertical wall.
			stepDown := p.OnGroundVec.Mul(p.StepHeight)
			stepUp := stepDown.Mul(-1)
			stepDownTrace := stepDown.Add(p.OnGroundVec) // Need one extra as we're not actually stepping this far.

			// 1. Step up.
			traceResult := p.traceMove(p.Contents & ^level.PlayerWalkableSolidContents, stepUp)
			if !traceResult.HitDelta.IsZero() {
				log.Debugf("walkMove: blocked upwards")
				// Hit something.
				// Don't stairstep.
				continue
			}
			p.Entity.Rect.Origin = traceResult.EndPos
			p.OnGround, p.GroundEntity, groundChecked = false, nil, true // Sure in air now.

			// 2. Continue move with _previous_ velocity.
			move = prevGoal.Add(stepUp).Delta(p.Entity.Rect.Origin)
			p.Velocity = prevVel
			moveClipped, _, trace := p.tryMove(move, false)
			preTouchVel = p.Velocity
			if trace != nil {
				p.handleTouchFunc(*trace)
			}
			teleported := p.Velocity != preTouchVel
			moveRemaining := prevGoal.Add(stepUp).Delta(p.Entity.Rect.Origin)
			if moveRemaining == move {
				// If no progress was made, actually do clip this move.
				move = moveClipped
			} else {
				// If any progress was made, recompute the remaining move for next iteration.
				move = moveRemaining
				if !teleported {
					// Only if handleTouchFunc left velocity alone.
					p.Velocity = prevVel
				}
			}

			// If we hit a stair again, we must NOT lose velocity!

			// 3. Step down (always).
			traceResult = p.traceMove(p.Contents & ^level.PlayerWalkableSolidContents, stepDownTrace)
			if traceResult.HitDelta.IsZero() {
				// Nothing found. Go back to original height, which is still in air.
				log.Debugf("walkMove: didn't reach ground after upstepping, so stepped back to original height")
				p.Entity.Rect.Origin = p.Entity.Rect.Origin.Add(stepDown)
				p.OnGround, p.GroundEntity, groundChecked = false, nil, true
			} else {
				log.Debugf("walkMove: stepped up by %v, reached ground", traceResult.EndPos.Delta(p.Entity.Rect.Origin.Add(stepDown)))
				p.Entity.Rect.Origin = traceResult.EndPos
				var hitEntity *engine.Entity
				if len(traceResult.HitEntities) != 0 {
					hitEntity = traceResult.HitEntities[0]
				}
				log.Debugf("ground hit")
				p.OnGround, p.GroundEntity, groundChecked = true, hitEntity, true
			}

			if teleported {
				// This happens only if handleTouchFunc changed velocity.
				break
			}
		}
	}

	// Step down for walking down stairs (only when started on ground, ground nearby and not moving upwards).
	if wasOnGround && p.Velocity.Dot(p.OnGroundVec) >= 0 {
		moved := p.Entity.Rect.Origin.Delta(prevOrigin)
		side := moved.Sub(p.OnGroundVec.Mul(moved.Dot(p.OnGroundVec))).Norm1()
		// NOTE: This must perform all downstepping in one move!
		// So consider StepHeight an angle.
		// Even if side == 0, this makes sense:
		// it'll never actually move the player,
		// but will update the onground flag and touch the entity the player is standing on.
		stepDown := p.OnGroundVec.Mul(p.StepHeight*side + 1) // Need one extra as we're not actually stepping this far.
		traceResult := p.traceMove(p.Contents & ^level.PlayerWalkableSolidContents, stepDown)
		if traceResult.HitDelta.IsZero() {
			// Nothing found. Stay in air.
			p.OnGround, p.GroundEntity, groundChecked = false, nil, true
		} else {
			if traceResult.EndPos != p.Entity.Rect.Origin {
				log.Debugf("walkMove: stepped down by %v", traceResult.EndPos.Delta(p.Entity.Rect.Origin))
				p.Entity.Rect.Origin = traceResult.EndPos
			}
			var hitEntity *engine.Entity
			if len(traceResult.HitEntities) != 0 {
				hitEntity = traceResult.HitEntities[0]
			}
			p.OnGround, p.GroundEntity, groundChecked = true, hitEntity, true
			p.handleTouchFunc(traceResult)
		}
	}

	return groundChecked
}

func (p *Physics) Update() {
	oldOrigin := p.Entity.Rect.Origin

	p.SubPixel = p.SubPixel.Add(p.Velocity)
	move := p.SubPixel.Div(constants.SubPixelScale)

	groundChecked := p.walkMove(move)

	if p.OnGround {
		if p.Velocity.Dot(p.OnGroundVec) < 0 {
			// Can't be on ground while moving up.
			p.OnGround, p.GroundEntity = false, nil
		} else if !groundChecked && !p.OnGroundVec.IsZero() {
			trace := p.World.TraceBox(p.Entity.Rect, p.Entity.Rect.Origin.Add(p.OnGroundVec), engine.TraceOptions{
				Contents:  p.Contents,
				IgnoreEnt: p.IgnoreEnt,
				ForEnt:    p.Entity,
				LoadTiles: true,
			})
			if trace.EndPos != p.Entity.Rect.Origin {
				p.OnGround, p.GroundEntity = false, nil
			} else {
				// p.OnGround = true // Always has been.
				var hitEntity *engine.Entity
				if len(trace.HitEntities) != 0 {
					hitEntity = trace.HitEntities[0]
				}
				p.GroundEntity = hitEntity
				p.handleTouchFunc(trace)
			}
		}
	}

	// Now if I am the ground, push everyone on me.
	delta := p.Entity.Rect.Origin.Delta(oldOrigin)
	if !delta.IsZero() {
		p.World.ForEachEntity(func(other *engine.Entity) {
			otherP, ok := other.Impl.(interfaces.Physics)
			if !ok {
				return
			}
			if otherP.ReadGroundEntity() == p.Entity {
				trace := p.World.TraceBox(other.Rect, other.Rect.Origin.Add(delta), engine.TraceOptions{
					Contents:  otherP.ReadContents(),
					IgnoreEnt: p.IgnoreEnt,
					ForEnt:    other,
					LoadTiles: true,
				})
				other.Rect.Origin = trace.EndPos
				if !trace.HitDelta.IsZero() {
					otherP.HandleTouch(trace)
				}
			}
		})
	}
}

func (p *Physics) ReadGroundEntity() *engine.Entity {
	return p.GroundEntity
}

func (p *Physics) HandleTouch(trace engine.TraceResult) {
	p.handleTouchFunc(trace)
}

func (p *Physics) ReadVelocity() m.Delta {
	return p.Velocity
}

func (p *Physics) SetVelocity(velocity m.Delta) {
	p.Velocity = velocity
}

func (p *Physics) SetVelocityForJump(velocity m.Delta) {
	p.SetVelocity(velocity)
	p.OnGround = false
}

func (p *Physics) ReadContents() level.Contents {
	return p.Contents
}

func (p *Physics) ReadSubPixel() m.Delta {
	return p.SubPixel
}

func (p *Physics) ReadOnGround() bool {
	return p.OnGround
}

func (p *Physics) ReadOnGroundVec() m.Delta {
	return p.OnGroundVec
}

func (p *Physics) ModifyHitBoxCentered(bySize m.Delta) m.Delta {
	if bySize.IsZero() {
		// Skip processing if we would have nothing to do.
		// NOTE: Function should verifiably do nothing in this case even if this return were missing.
		return m.Delta{}
	}

	prevOrigin := p.Entity.Rect.Origin
	prevSize := p.Entity.Rect.Size
	targetSize := prevSize.Add(bySize)

	// First grow in minus directions.
	topLeftDelta := bySize.Div(2)
	if topLeftDelta.DX > 0 {
		_, _, trace := p.tryMove(m.Delta{DX: -topLeftDelta.DX, DY: 0}, false)
		if trace != nil {
			p.handleTouchFunc(*trace)
		}
	} else {
		p.Entity.Rect.Origin.X -= topLeftDelta.DX
	}
	p.Entity.Rect.Size.DX += prevOrigin.X - p.Entity.Rect.Origin.X
	if topLeftDelta.DY > 0 {
		_, _, trace := p.tryMove(m.Delta{DX: 0, DY: -topLeftDelta.DY}, false)
		if trace != nil {
			p.handleTouchFunc(*trace)
		}
	} else {
		p.Entity.Rect.Origin.Y -= topLeftDelta.DY
	}
	p.Entity.Rect.Size.DY += prevOrigin.Y - p.Entity.Rect.Origin.Y

	// Then grow in plus directions.
	prevOrigin2 := p.Entity.Rect.Origin
	bottomRightDelta := targetSize.Sub(p.Entity.Rect.Size)
	if bottomRightDelta.DX > 0 {
		_, _, trace := p.tryMove(m.Delta{DX: bottomRightDelta.DX, DY: 0}, false)
		if trace != nil {
			p.handleTouchFunc(*trace)
		}
		p.Entity.Rect.Size.DX += p.Entity.Rect.Origin.X - prevOrigin2.X
		p.Entity.Rect.Origin.X = prevOrigin2.X
	} else {
		p.Entity.Rect.Size.DX += bottomRightDelta.DX
	}
	if bottomRightDelta.DY > 0 {
		_, _, trace := p.tryMove(m.Delta{DX: 0, DY: bottomRightDelta.DY}, false)
		if trace != nil {
			p.handleTouchFunc(*trace)
		}
		p.Entity.Rect.Size.DY += p.Entity.Rect.Origin.Y - prevOrigin2.Y
		p.Entity.Rect.Origin.Y = prevOrigin2.Y
	} else {
		p.Entity.Rect.Size.DY += bottomRightDelta.DY
	}

	// Grow remaining amount in minus directions again.
	prevOrigin3 := p.Entity.Rect.Origin
	topLeftDelta3 := targetSize.Sub(p.Entity.Rect.Size)
	if topLeftDelta3.DX > 0 {
		_, _, trace := p.tryMove(m.Delta{DX: -topLeftDelta3.DX, DY: 0}, false)
		if trace != nil {
			p.handleTouchFunc(*trace)
		}
	} else {
		p.Entity.Rect.Origin.X -= topLeftDelta3.DX
	}
	p.Entity.Rect.Size.DX += prevOrigin3.X - p.Entity.Rect.Origin.X
	if topLeftDelta3.DY > 0 {
		_, _, trace := p.tryMove(m.Delta{DX: 0, DY: -topLeftDelta3.DY}, false)
		if trace != nil {
			p.handleTouchFunc(*trace)
		}
	} else {
		p.Entity.Rect.Origin.Y -= topLeftDelta3.DY
	}
	p.Entity.Rect.Size.DY += prevOrigin3.Y - p.Entity.Rect.Origin.Y

	// Adjust render offset.
	p.Entity.RenderOffset = p.Entity.RenderOffset.Add(topLeftDelta)

	return p.Entity.Rect.Size.Sub(prevSize)
}

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
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/interfaces"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
)

const (
	SubPixelScale = 65536
)

type Physics struct {
	World  *engine.World
	Entity *engine.Entity

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
}

func (p *Physics) Reset() {
	p.OnGround = true
	p.GroundEntity = nil
	p.Velocity = m.Delta{}
	p.SubPixel = m.Delta{}
}

func (p *Physics) Update() {
	oldOrigin := p.Entity.Rect.Origin

	p.SubPixel = p.SubPixel.Add(p.Velocity)
	move := p.SubPixel.Div(SubPixelScale)

	groundChecked := false
	for (move != m.Delta{}) {
		dest := p.Entity.Rect.Origin.Add(move)
		trace := p.World.TraceBox(p.Entity.Rect, dest, engine.TraceOptions{
			Contents:  p.Contents,
			IgnoreEnt: p.IgnoreEnt,
			ForEnt:    p.Entity,
			LoadTiles: true,
		})
		if (trace.HitDelta == m.Delta{}) {
			// Nothing hit. We're done.
			p.SubPixel.DX -= move.DX * SubPixelScale
			p.SubPixel.DY -= move.DY * SubPixelScale
			p.Entity.Rect.Origin = trace.EndPos
			if move.DY != 0 {
				// If move had a Y component, we're flying.
				p.OnGround, p.GroundEntity, groundChecked = false, nil, true
			}
			break
		}
		if trace.HitDelta.DX != 0 {
			// An X hit. Just adjust X subpixel to be as close to the hit as possible.
			if p.SubPixel.DX > SubPixelScale-1 {
				p.SubPixel.DX = SubPixelScale - 1
			} else if p.SubPixel.DX < 0 {
				p.SubPixel.DX = 0
			}
			p.SubPixel.DY -= (trace.EndPos.Y - p.Entity.Rect.Origin.Y) * SubPixelScale
			p.Velocity.DX = 0
			move.DX = 0
			move.DY -= trace.EndPos.Y - p.Entity.Rect.Origin.Y
			p.Entity.Rect.Origin = trace.EndPos

			p.handleTouchFunc(trace)
		} else if trace.HitDelta.DY != 0 {
			// A Y hit. Also update ground status.
			if p.SubPixel.DY > SubPixelScale-1 {
				p.SubPixel.DY = SubPixelScale - 1
			} else if p.SubPixel.DY < 0 {
				p.SubPixel.DY = 0
			}
			p.SubPixel.DX -= (trace.EndPos.X - p.Entity.Rect.Origin.X) * SubPixelScale
			p.Velocity.DY = 0
			move.DX -= trace.EndPos.X - p.Entity.Rect.Origin.X
			move.DY = 0
			p.Entity.Rect.Origin = trace.EndPos

			if trace.HitDelta.DY > 0 {
				p.OnGround, p.GroundEntity, groundChecked = true, trace.HitEntity, true
			} else {
				p.OnGround, p.GroundEntity, groundChecked = false, nil, true
			}

			p.handleTouchFunc(trace)
		}
	}

	if p.OnGround && !groundChecked && p.OnGroundVec != (m.Delta{}) {
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
			p.GroundEntity = trace.HitEntity
			p.handleTouchFunc(trace)
		}
	}

	// Now if I am the ground, push everyone on me.
	delta := p.Entity.Rect.Origin.Delta(oldOrigin)
	if (delta != m.Delta{}) {
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
				if (trace.HitDelta != m.Delta{}) {
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

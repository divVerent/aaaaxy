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
	m "github.com/divVerent/aaaaaa/internal/math"
)

const (
	SubPixelScale = 65536
)

type Physics struct {
	World  *engine.World
	Entity *engine.Entity

	OnGround        bool
	GroundEntity    *engine.Entity
	Velocity        m.Delta // An input to be set changed by caller.
	SubPixel        m.Delta
	IgnoreEnt       *engine.Entity
	handleTouchFunc func(trace engine.TraceResult)
}

func (p *Physics) Init(w *engine.World, e *engine.Entity, handleTouch func(trace engine.TraceResult)) {
	p.World = w
	p.Entity = e
	p.handleTouchFunc = handleTouch
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
			IgnoreEnt: p.IgnoreEnt,
			ForEnt:    p.Entity,
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

	if p.OnGround && !groundChecked {
		trace := p.World.TraceBox(p.Entity.Rect, p.Entity.Rect.Origin.Add(m.Delta{DX: 0, DY: 1}), engine.TraceOptions{
			IgnoreEnt: p.IgnoreEnt,
			ForEnt:    p.Entity,
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
			otherP, ok := other.Impl.(groundEntityer)
			if !ok {
				return
			}
			if otherP.groundEntity() == p.Entity {
				trace := p.World.TraceBox(other.Rect, other.Rect.Origin.Add(delta), engine.TraceOptions{
					IgnoreEnt: p.IgnoreEnt,
					ForEnt:    other,
				})
				other.Rect.Origin = trace.EndPos
				if (trace.HitDelta != m.Delta{}) {
					otherP.handleTouch(trace)
				}
			}
		})
	}
}

func (p *Physics) groundEntity() *engine.Entity {
	return p.GroundEntity
}

func (p *Physics) handleTouch(trace engine.TraceResult) {
	p.handleTouchFunc(trace)
}

type groundEntityer interface {
	engine.EntityImpl
	groundEntity() *engine.Entity
	handleTouch(trace engine.TraceResult)
}

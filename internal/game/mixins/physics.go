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

	OnGround bool
	Velocity m.Delta // An input to be set changed by caller.
	SubPixel m.Delta
}

func (p *Physics) Init(w *engine.World, e *engine.Entity) {
	p.World = w
	p.Entity = e
}

func (p *Physics) Reset() {
	p.OnGround = true
	p.Velocity = m.Delta{}
	p.SubPixel = m.Delta{}
}

func (p *Physics) Update(handleTouch func(delta m.Delta, trace engine.TraceResult)) {
	p.SubPixel = p.SubPixel.Add(p.Velocity)
	move := p.SubPixel.Div(SubPixelScale)

	if move.DX != 0 {
		delta := m.Delta{DX: move.DX, DY: 0}
		dest := p.Entity.Rect.Origin.Add(delta)
		trace := p.World.TraceBox(p.Entity.Rect, dest, engine.TraceOptions{
			IgnoreEnt: p.Entity,
			ForEnt:    p.Entity,
		})
		if trace.EndPos == dest {
			// Nothing hit.
			p.SubPixel.DX -= move.DX * SubPixelScale
			p.Entity.Rect.Origin = trace.EndPos
		} else {
			// Hit something. Move as far as we can in direction of the hit, but not farther than intended.
			if p.SubPixel.DX > SubPixelScale-1 {
				p.SubPixel.DX = SubPixelScale - 1
			} else if p.SubPixel.DX < 0 {
				p.SubPixel.DX = 0
			}
			p.Velocity.DX = 0
			p.Entity.Rect.Origin = trace.EndPos
			handleTouch(delta, trace)
		}
	}

	if move.DY != 0 {
		delta := m.Delta{DX: 0, DY: move.DY}
		dest := p.Entity.Rect.Origin.Add(delta)
		trace := p.World.TraceBox(p.Entity.Rect, dest, engine.TraceOptions{
			IgnoreEnt: p.Entity,
			ForEnt:    p.Entity,
		})
		if trace.EndPos == dest {
			// Nothing hit.
			p.SubPixel.DY -= move.DY * SubPixelScale
			p.Entity.Rect.Origin = trace.EndPos
		} else {
			// Hit something. Move as far as we can in direction of the hit, but not farther than intended.
			if p.SubPixel.DY > SubPixelScale-1 {
				p.SubPixel.DY = SubPixelScale - 1
			} else if p.SubPixel.DY < 0 {
				p.SubPixel.DY = 0
			}
			p.Velocity.DY = 0
			// If moving down, set OnGround flag.
			if move.DY > 0 {
				p.OnGround = true
			}
			p.Entity.Rect.Origin = trace.EndPos
			handleTouch(delta, trace)
		}
	} else if p.OnGround {
		delta := m.Delta{DX: 0, DY: 1}
		trace := p.World.TraceBox(p.Entity.Rect, p.Entity.Rect.Origin.Add(m.Delta{DX: 0, DY: 1}), engine.TraceOptions{
			IgnoreEnt: p.Entity,
			ForEnt:    p.Entity,
		})
		if trace.EndPos != p.Entity.Rect.Origin {
			p.OnGround = false
		} else {
			handleTouch(delta, trace)
		}
	}
}

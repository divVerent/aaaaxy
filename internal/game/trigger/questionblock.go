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

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/interfaces"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/picture"
	"github.com/divVerent/aaaaxy/internal/propmap"
	"github.com/divVerent/aaaaxy/internal/sound"
)

// QuestionBlock is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
type QuestionBlock struct {
	World           *engine.World
	Entity          *engine.Entity
	PersistentState propmap.Map

	Kaizo  bool
	Target mixins.TargetSelection

	Used         bool
	UsedImage    *ebiten.Image
	UseAnimFrame int

	Sound *sound.Sound
}

const (
	UseFramesPerPixel = 2
	UsePixels         = 4
)

func (q *QuestionBlock) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	q.World = w
	q.Entity = e
	q.PersistentState = sp.PersistentState

	var parseErr error
	var err error
	w.SetSolid(e, true)
	w.SetOpaque(e, false)        // These shadows are annoying.
	e.Orientation = m.Identity() // Always show upright.
	q.Kaizo = propmap.ValueOrP(sp.Properties, "kaizo", false, &parseErr)
	q.Target = mixins.ParseTarget(propmap.StringOr(sp.Properties, "target", ""))
	q.Used = propmap.ValueOrP(q.PersistentState, "used", false, &parseErr)
	q.UsedImage, err = picture.Load("sprites", "exclamationblock.png")
	if err != nil {
		return err
	}
	if q.Used {
		e.Image = q.UsedImage
		q.UseAnimFrame = 2 * UseFramesPerPixel * UsePixels
	} else {
		if !q.Kaizo {
			e.Image, err = picture.Load("sprites", "questionblock.png")
			if err != nil {
				return err
			}
		}
	}
	q.Sound, err = sound.Load("questionblock.ogg")
	if err != nil {
		return fmt.Errorf("could not load questionblock sound: %w", err)
	}
	return parseErr
}

func (q *QuestionBlock) Despawn() {}

func (q *QuestionBlock) isAboveFlying(other *engine.Entity) bool {
	onGroundVec := m.Delta{DX: 0, DY: 1}
	onGround := false
	if phys, ok := other.Impl.(interfaces.Physics); ok {
		onGround = phys.ReadOnGround()
		onGroundVec = phys.ReadOnGroundVec()
	}
	return !onGround && q.Entity.Rect.Delta(other.Rect).Dot(onGroundVec) < 0
}

func (q *QuestionBlock) Update() {
	if q.Used {
		if q.UseAnimFrame < UseFramesPerPixel*UsePixels {
			q.UseAnimFrame++
			if q.UseAnimFrame%UseFramesPerPixel == 0 {
				q.Entity.RenderOffset.DY--
			}
		} else if q.UseAnimFrame < 2*UseFramesPerPixel*UsePixels {
			q.UseAnimFrame++
			if q.UseAnimFrame%UseFramesPerPixel == 0 {
				q.Entity.RenderOffset.DY++
			}
		}
		return
	}
	if !q.Kaizo {
		return
	}
	q.World.SetSolid(q.Entity, q.isAboveFlying(q.World.Player))
}

func (q *QuestionBlock) Touch(other *engine.Entity) {
	if other != q.World.Player {
		return
	}
	if !q.isAboveFlying(other) {
		return
	}

	// Send a message. Always do this, even if the block was already used.
	mixins.SetStateOfTarget(q.World, other, q.Entity, q.Target, true)

	if q.Used {
		return
	}
	q.Used = true
	propmap.Set(q.PersistentState, "used", true)
	q.Entity.Image = q.UsedImage
	q.UsedImage = nil
	q.World.SetSolid(q.Entity, true)
	q.Sound.Play()

	// Draw an effect.
	effect := q.Entity.Rect.Add(m.Delta{DX: 0, DY: -12})
	trace := q.World.TraceBox(q.Entity.Rect, effect.Origin, engine.TraceOptions{
		Contents: level.ObjectSolidContents,
		ForEnt:   q.Entity,
	})
	effect.Origin = trace.EndPos
	properties := propmap.New()
	propmap.Set(properties, "animation", "questionblock")
	propmap.Set(properties, "animation_frame_interval", "2")
	propmap.Set(properties, "animation_frames", "8")
	propmap.Set(properties, "animation_group", "hit")
	propmap.Set(properties, "animation_repeat_interval", "16")
	propmap.Set(properties, "fade_despawn", "true")
	propmap.Set(properties, "fade_time", "0s")
	propmap.Set(properties, "invert", "true")
	propmap.Set(properties, "no_transform", "true")
	propmap.Set(properties, "time_to_fade", "0.25s")
	propmap.Set(properties, "velocity", "0 -16") // 4px in 1/4 sec.
	_, err := q.World.SpawnDetached(&level.SpawnableProps{
		EntityType:      "MovingAnimation",
		Orientation:     m.Identity(),
		Properties:      properties,
		PersistentState: propmap.New(),
	}, effect, q.Entity.Orientation, q.Entity)
	if err != nil {
		log.Errorf("could not spawn question block effect: %v", err)
	}
}

func init() {
	engine.RegisterEntityType(&QuestionBlock{})
}

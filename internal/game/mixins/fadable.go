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
	"time"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

const (
	defaultFadeFrames = 16
	solidThreshold    = 8
	opaqueThreshold   = 16
)

// Fadable a mixin to make an object fade in/out when toggled.
// Must be initialized _after_ alpha and contents are set by the entity.
type Fadable struct {
	Settable
	World  *engine.World
	Entity *engine.Entity

	Alpha       float64
	Contents    level.Contents
	FadeFrames  int
	FadeDespawn bool

	AnimDir   int
	AnimFrame int
}

func (f *Fadable) Init(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	f.Settable.Init(sp)

	f.World = w
	f.Entity = e

	// Collect the sprite info.
	f.Alpha = f.Entity.Alpha
	f.Contents = f.Entity.Contents()

	var parseErr error
	animTime := propmap.ValueOrP(sp.Properties, "fade_time", time.Duration(defaultFadeFrames*time.Second/engine.GameTPS), &parseErr)
	f.FadeFrames = int((animTime*engine.GameTPS + (time.Second / 2)) / time.Second)
	if f.FadeFrames < 1 {
		f.FadeFrames = 1
	}
	f.FadeDespawn = propmap.ValueOrP(sp.Properties, "fade_despawn", false, &parseErr)

	// Skip the animation on initial load.
	if f.Settable.State && propmap.ValueOrP(sp.Properties, "fade_skip_animation", true, &parseErr) {
		f.AnimFrame = f.FadeFrames
	} else {
		f.AnimFrame = 0
	}
	f.Update()

	return parseErr
}

func (f *Fadable) SetState(originator, predecessor *engine.Entity, state bool) {
	f.Settable.SetState(originator, predecessor, state)
	if f.FadeDespawn {
		f.World.Detach(f.Entity)
	}
}

func (f *Fadable) Update() {
	if f.Settable.State {
		f.AnimFrame++
	} else {
		f.AnimFrame--
		if f.AnimFrame <= 0 && f.FadeDespawn {
			f.World.Despawn(f.Entity)
			return
		}
	}

	if f.AnimFrame <= 0 {
		f.Entity.Alpha = 0
		f.AnimFrame = 0
	} else if f.AnimFrame >= f.FadeFrames {
		f.Entity.Alpha = f.Alpha
		f.AnimFrame = f.FadeFrames
	} else {
		alpha := float64(f.AnimFrame) / float64(f.FadeFrames)
		f.Entity.Alpha = alpha * f.Alpha
	}

	f.World.MutateContentsBool(f.Entity, f.Contents&level.SolidContents, f.AnimFrame >= solidThreshold)
	f.World.MutateContentsBool(f.Entity, f.Contents&level.OpaqueContents, f.AnimFrame >= opaqueThreshold)
}

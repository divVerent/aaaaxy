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
	"fmt"
	"time"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/level"
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

	Alpha      float64
	Contents   level.Contents
	FadeFrames int

	AnimDir   int
	AnimFrame int
}

func (f *Fadable) Init(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	f.Settable.Init(sp)

	f.World = w
	f.Entity = e

	// Collect the sprite info.
	f.Alpha = f.Entity.Alpha
	f.Contents = f.Entity.Contents()

	fadeString := sp.Properties["fade_time"]
	if fadeString != "" {
		animTime, err := time.ParseDuration(fadeString)
		if err != nil {
			return fmt.Errorf("could not parse fade time: %v", fadeString)
		}
		f.FadeFrames = int((animTime*engine.GameTPS + (time.Second / 2)) / time.Second)
		if f.FadeFrames < 1 {
			f.FadeFrames = 1
		}
	} else {
		f.FadeFrames = defaultFadeFrames
	}

	// Skip the animation on initial load.
	if f.Settable.State {
		f.AnimFrame = f.FadeFrames
	} else {
		f.AnimFrame = 0
	}
	f.Update()

	return nil
}

func (f *Fadable) Update() {
	if f.Settable.State {
		f.AnimFrame++
	} else {
		f.AnimFrame--
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

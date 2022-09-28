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
	"github.com/divVerent/aaaaxy/internal/level"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

// Moving is a physics object with initial velocity. That's all.
type Moving struct {
	Physics
	TouchedSomething bool
}

func (v *Moving) Init(w *engine.World, sp *level.SpawnableProps, e *engine.Entity, contents level.Contents, handleTouch func(engine.TraceResult)) error {
	v.Physics.Init(w, e, contents, handleTouch)
	var parseErr error
	vel := propmap.ValueOrP(sp.Properties, "velocity", m.Delta{}, &parseErr)
	v.Physics.Velocity = e.Transform.Inverse().Apply(
		vel.MulFracFixed(m.NewFixed(constants.SubPixelScale), m.NewFixed(engine.GameTPS)))
	return parseErr
}

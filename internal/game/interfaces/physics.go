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

package interfaces

import (
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
)

type Velocityer interface {
	engine.EntityImpl

	SetVelocity(velocity m.Delta)
	SetVelocityForJump(velocity m.Delta)
	ReadVelocity() m.Delta
	ReadSubPixel() m.Delta
	ReadOnGroundVec() m.Delta
}

type GroundEntityer interface {
	engine.EntityImpl

	ReadGroundEntity() *engine.Entity
}

type HandleToucher interface {
	engine.EntityImpl

	HandleTouch(trace engine.TraceResult)
}

type Contentser interface {
	engine.EntityImpl

	ReadContents() level.Contents
}

type Physics interface {
	Velocityer
	GroundEntityer
	HandleToucher
	Contentser
}

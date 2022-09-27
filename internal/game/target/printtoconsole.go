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

package target

import (
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/fun"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

// PrintToConsoleTarget prints the given text to console when activated.
// Setting state to ON saves the current text, setting state to OFF dumps it.
type PrintToConsoleTarget struct {
	World *engine.World

	Text string

	PrintFrames int
	PrevText    string
}

func (p *PrintToConsoleTarget) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	p.World = w
	var parseErr error
	p.Text = propmap.ValueP(sp.Properties, "text", "", &parseErr)
	return parseErr
}

func (p *PrintToConsoleTarget) Despawn() {
	if p.PrevText != "" {
		log.Infof("%s", p.PrevText)
		p.PrevText = ""
	}
}

func (p *PrintToConsoleTarget) Update() {
	if p.PrintFrames > 0 {
		p.PrintFrames--
		if p.PrintFrames == 0 {
			if p.PrevText != "" {
				log.Infof("%s", p.PrevText)
				p.PrevText = ""
			}
		}
	}
}

func (p *PrintToConsoleTarget) Touch(other *engine.Entity) {}

func (p *PrintToConsoleTarget) SetState(originator, predecessor *engine.Entity, state bool) {
	if state {
		p.PrevText = fun.FormatText(&p.World.PlayerState, p.Text)
		p.PrintFrames = 2
	}
}

func init() {
	engine.RegisterEntityType(&PrintToConsoleTarget{})
}

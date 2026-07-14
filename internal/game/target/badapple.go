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
	"time"

	"github.com/divVerent/aaaaxy/internal/centerprint"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/fun"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/palette"
)

// BadAppleTarget prints the given text to console when activated.
// Setting state to ON saves the current text, setting state to OFF dumps it.
type BadAppleTarget struct {
	World *engine.World

	State bool
}

func (b *BadAppleTarget) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	b.World = w
	return nil
}

func (b *BadAppleTarget) Despawn() {}

func (b *BadAppleTarget) Update() {}

func (b *BadAppleTarget) Touch(other *engine.Entity) {}

func (b *BadAppleTarget) SetState(originator, predecessor *engine.Entity, state bool) {
	if state == b.State {
		return
	}
	b.State = state
	if state {
		b.World.PlayerState.ToggleBadApple()
		err := b.World.Save()
		if err != nil {
			log.Errorf("could not save game: %v", err)
			str := locale.G.Get("Error:\ncould not save game:\n%s", err)
			centerprint.New(fun.FormatText(&b.World.PlayerState, str), centerprint.Important, centerprint.Top, centerprint.NormalFont(), palette.EGA(palette.LightRed, 255), 5*time.Second).SetFadeOut(true)
			return
		}
	}
}

func init() {
	engine.RegisterEntityType(&BadAppleTarget{})
}

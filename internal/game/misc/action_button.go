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

package misc

import (
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/interfaces"
	"github.com/divVerent/aaaaaa/internal/game/mixins"
	"github.com/divVerent/aaaaaa/internal/level"
)

// ActionButton sends a signal along when action button is pressed/released.
type ActionButton struct {
	World  *engine.World
	Entity *engine.Entity

	Target      mixins.TargetSelection
	Invert      bool
	SendPress   bool
	SendRelease bool
	SendHold    bool

	State bool
}

func (g *ActionButton) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	g.World = w
	g.Entity = e
	g.Target = mixins.ParseTarget(sp.Properties["target"])
	g.Invert = sp.Properties["invert"] == "true"            // false by default.
	g.SendPress = sp.Properties["send_press"] == "true"     // false by default.
	g.SendRelease = sp.Properties["send_release"] == "true" // false by default.
	g.SendHold = sp.Properties["send_hold"] == "true"       // false by default.
	return nil
}

func (g *ActionButton) Despawn() {}

func (g *ActionButton) Update() {
	newState := g.World.Player.Impl.(interfaces.ActionPresseder).ActionPressed()
	if newState == g.State && !(g.SendHold && newState) {
		return
	}
	g.State = newState
	if newState && g.SendPress {
		mixins.SetStateOfTarget(g.World, g.World.Player, g.Entity, g.Target, !g.Invert)
	}
	if !newState && g.SendRelease {
		mixins.SetStateOfTarget(g.World, g.World.Player, g.Entity, g.Target, g.Invert)
	}
}

func (g *ActionButton) Touch(other *engine.Entity) {}

func init() {
	engine.RegisterEntityType(&ActionButton{})
}

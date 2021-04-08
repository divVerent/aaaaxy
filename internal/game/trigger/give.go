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

	"github.com/divVerent/aaaaaa/internal/animation"
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/interfaces"
	"github.com/divVerent/aaaaaa/internal/game/mixins"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
)

const (
	giveFadeFrames = 60
)

// Give grants the player an ability when touched.
type Give struct {
	mixins.NonSolidTouchable

	Ability   string
	Text      string
	AnimFrame int

	Anim animation.State
}

func (g *Give) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	g.NonSolidTouchable.Init(w, e)
	g.Ability = s.Properties["ability"]
	g.Text = s.Properties["text"]
	err := g.Anim.Init("can_"+g.Ability, map[string]*animation.Group{
		"default": &animation.Group{
			Frames:        16,
			FrameInterval: 4,
			NextInterval:  4 * 16,
			NextAnim:      "default",
		}}, "default")
	if err != nil {
		return fmt.Errorf("could not initialize give animation: %v", err)
	}
	return nil
}

func (g *Give) Despawn() {}

func (g *Give) Update() {
	g.NonSolidTouchable.Update()
	if g.World.Player.Impl.(interfaces.HasAbilityer).HasAbility(g.Ability) {
		g.AnimFrame--
	} else {
		g.AnimFrame++
	}

	if g.AnimFrame <= 0 {
		g.Entity.Alpha = 0
		g.AnimFrame = 0
	} else if g.AnimFrame >= giveFadeFrames {
		g.Entity.Alpha = 1
		g.AnimFrame = giveFadeFrames
	} else {
		g.Entity.Alpha = float64(g.AnimFrame) / giveFadeFrames
	}

	g.Anim.Update(g.Entity)
}

func (g *Give) Touch(other *engine.Entity) {
	if other != g.World.Player {
		return
	}
	g.World.Player.Impl.(interfaces.HasAbilityer).GiveAbility(g.Ability, g.Text)
}

func (g *Give) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&Give{})
}

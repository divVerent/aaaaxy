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

package checkpoint

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/centerprint"
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/player"
	m "github.com/divVerent/aaaaaa/internal/math"
	"github.com/divVerent/aaaaaa/internal/music"
)

// Checkpoint remembers that it was hit and allows spawning from there again. Also displays a text.
type Checkpoint struct {
	World  *engine.World
	Entity *engine.Entity

	Name  string
	Text  string
	Music string

	PlayerProperty        string
	PlayerPropertyFlipped string
	Inactive              bool
}

func (c *Checkpoint) Spawn(w *engine.World, s *engine.Spawnable, e *engine.Entity) error {
	c.World = w
	c.Entity = e

	// Field contains orientation OF THE PLAYER to make it easier in the map editor.
	// So it is actually a transform as far as this code is concerned.
	requiredTransform, err := m.ParseOrientation(s.Properties["required_orientation"])
	if err != nil {
		return fmt.Errorf("could not parse required orientation: %v", err)
	}

	c.Name = s.Properties["name"]
	c.Text = s.Properties["text"]
	c.Music = s.Properties["music"]

	c.PlayerProperty = "checkpoint_seen." + c.Name
	if c.Entity.Transform == requiredTransform {
		c.PlayerPropertyFlipped = "Identity"
	} else if c.Entity.Transform == requiredTransform.Concat(m.FlipX()) {
		c.PlayerPropertyFlipped = "FlipX"
	} else {
		c.Inactive = true
	}

	return nil
}

func (c *Checkpoint) Despawn() {}

func (c *Checkpoint) Update() {
	if c.Inactive {
		return
	}
	// The "down" direction must match. That way we allow x-flipping and still matching the CP.
	if (c.World.Player.Rect.Delta(c.Entity.Rect) != m.Delta{}) {
		return
	}
	// Checkpoint always sets "mood".
	music.Switch(c.Music)
	player := c.World.Player.Impl.(*player.Player)
	if player.PersistentState["last_checkpoint"] == c.Name && player.PersistentState[c.PlayerProperty] == c.PlayerPropertyFlipped {
		return
	}
	player.PersistentState[c.PlayerProperty] = c.PlayerPropertyFlipped
	player.PersistentState["last_checkpoint"] = c.Name
	centerprint.New(c.Text, centerprint.Important, centerprint.Middle, centerprint.BigFont, color.NRGBA{R: 255, G: 255, B: 255, A: 255}).SetFadeOut(true)
}

func (c *Checkpoint) Touch(other *engine.Entity) {}

func (c *Checkpoint) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&Checkpoint{})
}

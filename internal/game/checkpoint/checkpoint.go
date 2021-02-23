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
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/centerprint"
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/game/mixins"
	"github.com/divVerent/aaaaaa/internal/game/player"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
	"github.com/divVerent/aaaaaa/internal/music"
	"github.com/divVerent/aaaaaa/internal/sound"
)

// Checkpoint remembers that it was hit and allows spawning from there again. Also displays a text.
type Checkpoint struct {
	mixins.NonSolidTouchable
	World  *engine.World
	Entity *engine.Entity

	Text    string
	Music   string
	DeadEnd bool

	FlippedStr string
	Inactive   bool

	Sound *sound.Sound
}

func (c *Checkpoint) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	c.NonSolidTouchable.Init(w, e)
	c.World = w
	c.Entity = e

	// Field contains orientation OF THE PLAYER to make it easier in the map editor.
	// So it is actually a transform as far as this code is concerned.
	requiredTransform, err := m.ParseOrientation(s.Properties["required_orientation"])
	if err != nil {
		return fmt.Errorf("could not parse required orientation: %v", err)
	}

	c.Text = s.Properties["text"]
	c.Music = s.Properties["music"]
	c.DeadEnd = s.Properties["dead_end"] == "true"

	if c.Entity.Transform == requiredTransform {
		c.FlippedStr = "Identity"
	} else if c.Entity.Transform == requiredTransform.Concat(m.FlipX()) {
		c.FlippedStr = "FlipX"
	} else {
		c.Inactive = true
	}

	c.Sound, err = sound.Load("checkpoint.ogg")
	if err != nil {
		return fmt.Errorf("could not load checkpoint sound: %v", err)
	}

	return nil
}

func (c *Checkpoint) Despawn() {}

func (c *Checkpoint) setCheckpoint() bool {
	player := c.World.Player.Impl.(*player.Player)
	changed := false
	cpProperty := "checkpoint_seen." + c.Entity.Name()
	if player.PersistentState[cpProperty] != c.FlippedStr {
		player.PersistentState[cpProperty] = c.FlippedStr
		changed = true
	}
	if player.PersistentState["last_checkpoint"] != c.Entity.Name() {
		edgeProperty := "checkpoints_walked." + player.PersistentState["last_checkpoint"] + "." + c.Entity.Name()
		if player.PersistentState[edgeProperty] != "true" {
			player.PersistentState[edgeProperty] = "true"
			changed = true
		}
		if !c.DeadEnd {
			player.PersistentState["last_checkpoint"] = c.Entity.Name()
			changed = true
		}
	}
	return changed
}

func (c *Checkpoint) Touch(other *engine.Entity) {
	if other != c.World.Player {
		return
	}
	if c.Inactive {
		return
	}
	// All checkpoints set the "mood".
	music.Switch(c.Music)
	if !c.setCheckpoint() {
		return
	}
	err := c.World.Save()
	if err != nil {
		log.Printf("Could not save game: %v", err)
		return
	}
	if !c.DeadEnd {
		centerprint.New(c.Text, centerprint.Important, centerprint.Middle, centerprint.BigFont, color.NRGBA{R: 255, G: 255, B: 255, A: 255}).SetFadeOut(true)
		c.Sound.Play()
	}
}

func (c *Checkpoint) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&Checkpoint{})
}

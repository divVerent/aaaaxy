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
	"github.com/divVerent/aaaaxy/internal/log"
	"image/color"
	"time"

	"github.com/divVerent/aaaaxy/internal/centerprint"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/fun"
	"github.com/divVerent/aaaaxy/internal/game/interfaces"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/level"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/music"
	"github.com/divVerent/aaaaxy/internal/sound"
)

// Checkpoint remembers that it was hit and allows spawning from there again. Also displays a text.
type Checkpoint struct {
	mixins.NonSolidTouchable
	World  *engine.World
	Entity *engine.Entity

	Text  string
	Music string

	Flipped           bool
	Inactive          bool
	VVVVVV            bool
	VVVVVVOnGroundVec m.Delta

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
	c.VVVVVV = s.Properties["vvvvvv"] == "true"
	if onGroundVecStr := s.Properties["vvvvvv_gravity_direction"]; onGroundVecStr != "" {
		_, err := fmt.Sscanf(onGroundVecStr, "%d %d", &c.VVVVVVOnGroundVec.DX, &c.VVVVVVOnGroundVec.DY)
		if err != nil {
			return fmt.Errorf("invalid vvvvvv_gravity_direction: %v", err)
		}
	}

	if c.Entity.Transform == requiredTransform {
		c.Flipped = false
	} else if c.Entity.Transform == requiredTransform.Concat(m.FlipX()) {
		c.Flipped = true
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

func (c *Checkpoint) Touch(other *engine.Entity) {
	if other != c.World.Player {
		return
	}
	if c.Inactive {
		return
	}
	if c.VVVVVV {
		c.World.Player.Impl.(interfaces.VVVVVVer).SetVVVVVV(true, c.VVVVVVOnGroundVec, 1.0)
	}
	c.World.PlayerTouchedCheckpoint(c.Entity)
	// All checkpoints set the "mood".
	music.Switch(c.Music)
	if !c.World.PlayerState.RecordCheckpointEdge(c.Entity.Name(), c.Flipped) {
		return
	}
	err := c.World.Save()
	if err != nil {
		log.Errorf("Could not save game: %v", err)
		return
	}
	if c.Text != "" {
		centerprint.New(fun.FormatText(&c.World.PlayerState, c.Text), centerprint.Important, centerprint.Middle, centerprint.BigFont(), color.NRGBA{R: 255, G: 255, B: 255, A: 255}, time.Second).SetFadeOut(true)
		c.Sound.Play()
	}
}

func init() {
	engine.RegisterEntityType(&Checkpoint{})
}

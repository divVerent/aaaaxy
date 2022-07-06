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
	"time"

	"github.com/divVerent/aaaaxy/internal/centerprint"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/fun"
	"github.com/divVerent/aaaaxy/internal/game/interfaces"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/music"
	"github.com/divVerent/aaaaxy/internal/palette"
	"github.com/divVerent/aaaaxy/internal/sound"
)

// CheckpointTarget remembers that it was activated and allows spawning from there again. Also displays a text.
type CheckpointTarget struct {
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

func (c *CheckpointTarget) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	c.World = w
	c.Entity = e

	// Field contains orientation OF THE PLAYER to make it easier in the map editor.
	// So it is actually a transform as far as this code is concerned.
	requiredTransforms, err := m.ParseOrientations(sp.Properties["required_orientation"])
	if err != nil {
		return fmt.Errorf("could not parse required orientation: %w", err)
	}

	c.Text = sp.Properties["text"]
	c.Music = sp.Properties["music"]
	c.VVVVVV = sp.Properties["vvvvvv"] == "true"
	if onGroundVecStr := sp.Properties["vvvvvv_gravity_direction"]; onGroundVecStr != "" {
		_, err := fmt.Sscanf(onGroundVecStr, "%d %d", &c.VVVVVVOnGroundVec.DX, &c.VVVVVVOnGroundVec.DY)
		if err != nil {
			return fmt.Errorf("invalid vvvvvv_gravity_direction: %w", err)
		}
	}

	c.Inactive = true
	for _, requiredTransform := range requiredTransforms {
		if c.Entity.Transform == requiredTransform {
			c.Flipped = false
			c.Inactive = false
			break
		} else if c.Entity.Transform == requiredTransform.Concat(m.FlipX()) {
			c.Flipped = true
			c.Inactive = false
			break
		}
	}

	c.Sound, err = sound.Load("checkpoint.ogg")
	if err != nil {
		return fmt.Errorf("could not load checkpoint sound: %w", err)
	}

	return nil
}

func (c *CheckpointTarget) Despawn() {}

func (c *CheckpointTarget) Update() {}

func (c *CheckpointTarget) Touch(other *engine.Entity) {}

func (c *CheckpointTarget) SetState(originator, predecessor *engine.Entity, state bool) {
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
		log.Errorf("could not save game: %v", err)
		str := fmt.Sprintf("Error:\ncould not save game:\n%v", err)
		centerprint.New(fun.FormatText(&c.World.PlayerState, str), centerprint.Important, centerprint.Top, centerprint.NormalFont(), palette.NRGBA(255, 85, 85, 255), 5*time.Second).SetFadeOut(true)
		return
	}
	if c.Text != "" {
		centerprint.New(fun.FormatText(&c.World.PlayerState, c.Text), centerprint.Important, centerprint.Middle, centerprint.BigFont(), palette.NRGBA(255, 255, 255, 255), time.Second).SetFadeOut(true)
		c.Sound.Play()
	}
}

func init() {
	engine.RegisterEntityType(&CheckpointTarget{})
}

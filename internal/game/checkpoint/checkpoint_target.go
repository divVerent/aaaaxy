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
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/m"
	"github.com/divVerent/aaaaxy/internal/music"
	"github.com/divVerent/aaaaxy/internal/palette"
	"github.com/divVerent/aaaaxy/internal/propmap"
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

func (c *CheckpointTarget) Precache(sp *level.Spawnable) {
	txtOrig, err := propmap.Value(sp.Properties, "text", "")
	if err != nil {
		log.Warningf("failed to read text on entity %v: %v", sp.ID, err)
		return
	}
	txt, err := fun.TryFormatText(nil, txtOrig)
	if err != nil {
		// Cannot format, requires player state. No bounds checking then.
		return
	}
	bounds := centerprint.BigFont().BoundString(txt)
	if bounds.Size.DX > engine.GameWidth {
		locale.Errorf("text too big: entity %v must fit in width %v but text needs %v: %v",
			sp.ID, engine.GameWidth, bounds.Size, txtOrig)
	}
}

func (c *CheckpointTarget) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	c.World = w
	c.Entity = e

	var parseErr error

	// Field contains orientation OF THE PLAYER to make it easier in the map editor.
	// So it is actually a transform as far as this code is concerned.
	requiredTransforms := propmap.ValueP(sp.Properties, "required_orientation", m.Orientations{}, &parseErr)

	c.Text = propmap.StringOr(sp.Properties, "text", "")
	c.Music = propmap.StringOr(sp.Properties, "music", "")
	c.VVVVVV = propmap.ValueOrP(sp.Properties, "vvvvvv", false, &parseErr)
	c.VVVVVVOnGroundVec = propmap.ValueOrP(sp.Properties, "vvvvvv_gravity_direction", m.Delta{}, &parseErr)

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

	var err error
	c.Sound, err = sound.Load("checkpoint.ogg")
	if err != nil {
		return fmt.Errorf("could not load checkpoint sound: %w", err)
	}

	return parseErr
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
		str := locale.G.Get("Error:\ncould not save game:\n%s", err)
		centerprint.New(fun.FormatText(&c.World.PlayerState, str), centerprint.Important, centerprint.Top, centerprint.NormalFont(), palette.EGA(palette.LightRed, 255), 5*time.Second).SetFadeOut(true)
		return
	}
	if c.Text != "" {
		centerprint.New(fun.FormatText(&c.World.PlayerState, c.Text), centerprint.Important, centerprint.Middle, centerprint.BigFont(), palette.EGA(palette.White, 255), time.Second).SetFadeOut(true)
		c.Sound.Play()
	}
}

func init() {
	engine.RegisterEntityType(&CheckpointTarget{})
}

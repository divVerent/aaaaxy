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
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/centerprint"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/game/constants"
	"github.com/divVerent/aaaaxy/internal/game/mixins"
	"github.com/divVerent/aaaaxy/internal/image"
	"github.com/divVerent/aaaaxy/internal/level"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/sound"
)

// TnihSign just displays a text and remembers that it was hit.
type TnihSign struct {
	mixins.NonSolidTouchable
	World           *engine.World
	Entity          *engine.Entity
	PersistentState map[string]string

	Text      string
	SeenImage *ebiten.Image
	Sound     *sound.Sound
	Target    mixins.TargetSelection

	Touching bool
	Touched  bool

	Centerprint *centerprint.Centerprint
}

const (
	tnihWidth  = 32
	tnihHeight = 32
)

func (t *TnihSign) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	t.NonSolidTouchable.Init(w, e)
	t.NotifyUntouched = true
	t.World = w
	t.Entity = e
	t.PersistentState = s.PersistentState
	var err error
	t.SeenImage, err = image.Load("sprites", "tnihsign_seen.png")
	if err != nil {
		return fmt.Errorf("could not load sign seen sprite: %v", err)
	}
	if s.PersistentState["seen"] == "true" {
		t.Entity.Image = t.SeenImage
	} else {
		t.Entity.Image, err = image.Load("sprites", "tnihsign.png")
		if err != nil {
			return fmt.Errorf("could not load sign sprite: %v", err)
		}
	}
	t.Entity.Orientation = m.Identity()
	w.SetZIndex(t.Entity, constants.TnihSignZ)
	t.Text = strings.ReplaceAll(s.Properties["text"], "  ", "\n")
	t.Sound, err = sound.Load("tnihsign.ogg")
	if err != nil {
		return fmt.Errorf("could not load tnihsign sound: %v", err)
	}
	t.Entity.ResizeImage = s.Properties["resize_image"] == "true"
	if !t.Entity.ResizeImage {
		t.Entity.RenderOffset = t.Entity.Rect.Size.Sub(m.Delta{DX: tnihWidth, DY: tnihHeight}).Div(2)
	}
	t.Target = mixins.ParseTarget(s.Properties["target"])
	return nil
}

func (t *TnihSign) Despawn() {
	if t.Centerprint.Active() {
		t.Centerprint.SetFadeOut(true)
	}
}

func (t *TnihSign) Touch(other *engine.Entity) {
	if other != t.World.Player {
		return
	}
	if t.Centerprint.Active() {
		t.Centerprint.SetFadeOut(false)
	} else {
		importance := centerprint.Important
		if t.PersistentState["seen"] == "true" {
			importance = centerprint.NotImportant
		} else {
			t.Sound.Play()
		}
		t.Centerprint = centerprint.New(t.Text, importance, centerprint.Top, centerprint.NormalFont(), color.NRGBA{R: 255, G: 255, B: 85, A: 255})
		t.PersistentState["seen"] = "true"
		t.Entity.Image = t.SeenImage
		mixins.SetStateOfTarget(t.World, other, t.Entity, t.Target, true)
	}
	t.Touching = true
}

func (t *TnihSign) Update() {
	t.NonSolidTouchable.Update()
	if t.Touched && !t.Touching {
		if t.Centerprint.Active() {
			t.Centerprint.SetFadeOut(true)
		}
	}
	t.Touching, t.Touched = false, t.Touching
}

func init() {
	engine.RegisterEntityType(&TnihSign{})
}

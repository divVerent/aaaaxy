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
	"fmt"
	"image/color"
	"time"

	"github.com/divVerent/aaaaxy/internal/centerprint"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/fun"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/palette"
	"github.com/divVerent/aaaaxy/internal/propmap"
	"github.com/divVerent/aaaaxy/internal/sound"
)

// CenterPrintTarget just displays a text to the screen.
type CenterPrintTarget struct {
	World *engine.World

	Text     string
	Font     *font.Face
	Imp      centerprint.Importance
	Pos      centerprint.InitialPosition
	BGColor  color.Color
	FGColor  color.Color
	FadeTime time.Duration
	Sound    *sound.Sound

	Centerprint *centerprint.Centerprint
}

func (t *CenterPrintTarget) Precache(sp *level.Spawnable) error {
	txtOrig, err := propmap.Value(sp.Properties, "text", "")
	if err != nil {
		log.Warningf("failed to read text on entity %v: %v", sp.ID, err)
		return nil
	}
	txt, err := fun.TryFormatText(nil, txtOrig)
	if err != nil {
		// Cannot format, requires player state. No bounds checking then.
		return nil
	}
	fontName, err := propmap.Value(sp.Properties, "text_font", "Centerprint")
	if err != nil {
		log.Warningf("failed to read font on entity %v: %v", sp.ID, err)
		return nil
	}
	font := font.ByName[fontName]
	if font == nil {
		log.Warningf("failed to find font %q on entity %v", fontName, sp.ID)
		return nil
	}
	bounds := font.BoundString(txt)
	if bounds.Size.DX > engine.GameWidth {
		locale.Errorf("text too big: entity %v must fit in width %v but text needs %v: %v",
			sp.ID, engine.GameWidth, bounds.Size, txtOrig)
	}
	return nil
}

func (t *CenterPrintTarget) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	var parseErr error
	t.World = w
	t.Text = propmap.ValueP(sp.Properties, "text", "", &parseErr)
	fontName := propmap.ValueOrP(sp.Properties, "text_font", "Centerprint", &parseErr)
	t.Font = font.ByName[fontName]
	if t.Font == nil {
		log.Warningf("failed to find font %q", fontName)
		return nil
	}
	t.Imp = propmap.ValueOrP(sp.Properties, "importance", centerprint.Important, &parseErr)
	t.Pos = propmap.ValueOrP(sp.Properties, "initial_position", centerprint.Top, &parseErr)
	t.BGColor = propmap.ValueOrP(sp.Properties, "text_bg", palette.EGA(palette.Black, 255), &parseErr)
	t.FGColor = propmap.ValueOrP(sp.Properties, "text_fg", palette.EGA(palette.White, 255), &parseErr)
	t.FadeTime = propmap.ValueOrP(sp.Properties, "fade_time", 2*time.Second, &parseErr)
	soundName := propmap.ValueP(sp.Properties, "sound", "", &parseErr)
	if soundName != "" {
		var err error
		t.Sound, err = sound.Load(soundName)
		if err != nil {
			return fmt.Errorf("could not load sound %q: %w", soundName, err)
		}
	}
	return parseErr
}

func (t *CenterPrintTarget) Despawn() {
	if t.Centerprint.Active() {
		t.Centerprint.SetFadeOut(true)
	}
}

func (t *CenterPrintTarget) Update() {}

func (t *CenterPrintTarget) Touch(other *engine.Entity) {}

func (t *CenterPrintTarget) SetState(originator, predecessor *engine.Entity, state bool) {
	if state {
		if !t.Centerprint.Active() {
			t.Centerprint = centerprint.NewWithBG(t.Text, t.Imp, t.Pos, t.Font, t.BGColor, t.FGColor, t.FadeTime)
			if t.Sound != nil {
				t.Sound.Play()
			}
		}
	} else {
		if t.Centerprint.Active() {
			t.Centerprint.SetFadeOut(true)
		}
	}
}

func init() {
	engine.RegisterEntityType(&CenterPrintTarget{})
}

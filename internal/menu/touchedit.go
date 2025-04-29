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

package menu

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/m"
	"github.com/divVerent/aaaaxy/internal/palette"
)

type TouchEditScreenItem int

const (
	TouchDone = iota
	TouchReset
	TouchCount
)

type TouchEditScreen struct {
	Controller *Controller
	Item       TouchEditScreenItem
}

func (s *TouchEditScreen) Init(m *Controller) error {
	s.Controller = m
	return nil
}

func touchReset() error {
	input.TouchResetEditor()
	return nil
}

func (s *TouchEditScreen) Update() error {
	clicked := s.Controller.QueryMouseItem(&s.Item, TouchCount)
	if input.Down.JustHit {
		s.Item++
		s.Controller.MoveSound(nil)
	}
	if input.Up.JustHit {
		s.Item--
		s.Controller.MoveSound(nil)
	}
	s.Item = TouchEditScreenItem(m.Mod(int(s.Item), int(TouchCount)))
	if input.Exit.JustHit {
		return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&SettingsScreen{}))
	}
	if input.Jump.JustHit || input.Action.JustHit || clicked != NotClicked {
		switch s.Item {
		case TouchDone:
			return s.Controller.ActivateSound(s.Controller.SaveConfigAndSwitchToScreen(&SettingsScreen{}))
		case TouchReset:
			return s.Controller.ActivateSound(touchReset())
		}
	}
	return nil
}

func (s *TouchEditScreen) Draw(screen *ebiten.Image) {
	input.DrawEditor(screen)
	fgs := palette.EGA(palette.Yellow, 255)
	bgs := palette.EGA(palette.Black, 255)
	fgn := palette.EGA(palette.LightGrey, 255)
	bgn := palette.EGA(palette.DarkGrey, 255)
	font.ByName["MenuBig"].Draw(screen, locale.G.Get("Edit Touch Controls"), m.Pos{X: CenterX, Y: HeaderY}, font.Center, fgs, bgs)
	fg, bg := fgn, bgn
	if s.Item == TouchDone {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("Done"), m.Pos{X: CenterX, Y: ItemBaselineY(TouchDone, TouchCount)}, font.Center, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == TouchReset {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("Reset to Defaults"), m.Pos{X: CenterX, Y: ItemBaselineY(TouchReset, TouchCount)}, font.Center, fg, bg)
}

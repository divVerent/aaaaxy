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
	"github.com/divVerent/aaaaxy/internal/fun"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/locale"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/palette"
)

type MainScreenItem int

const (
	Play = iota
	Settings
	Credits
	Quit
	MainCount
)

type MainScreen struct {
	Controller *Controller
	Item       MainScreenItem
}

func (s *MainScreen) Init(m *Controller) error {
	s.Controller = m
	return nil
}

func (s *MainScreen) Update() error {
	clicked := s.Controller.QueryMouseItem(&s.Item, MainCount)
	if input.Down.JustHit {
		s.Item++
		s.Controller.MoveSound(nil)
	}
	if input.Up.JustHit {
		s.Item--
		s.Controller.MoveSound(nil)
	}
	s.Item = MainScreenItem(m.Mod(int(s.Item), int(MainCount)))

	/*
		Actually not allowed as it could be used for pausebuffering.
		if input.Exit.JustHit {
			return s.Controller.ActivateSound(s.Controller.SwitchToGame())
		}
	*/
	if input.Jump.JustHit || input.Action.JustHit || clicked != NotClicked {
		switch s.Item {
		case Play:
			return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&MapScreen{}))
		case Settings:
			return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&SettingsScreen{}))
		case Credits:
			return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&CreditsScreen{Fancy: false}))
		case Quit:
			return s.Controller.ActivateSound(s.Controller.QuitGame())
		}
	}
	return nil
}

func (s *MainScreen) Draw(screen *ebiten.Image) {
	fgs := palette.EGA(palette.Yellow, 255)
	bgs := palette.EGA(palette.Black, 255)
	fgn := palette.EGA(palette.LightGrey, 255)
	bgn := palette.EGA(palette.DarkGrey, 255)
	font.MenuBig.Draw(screen, "AAAAXY", m.Pos{X: CenterX, Y: HeaderY}, true, fgs, bgs)
	fg, bg := fgn, bgn
	if s.Item == Play {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, locale.G.Get("Play"), m.Pos{X: CenterX, Y: ItemBaselineY(Play, MainCount)}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Settings {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, locale.G.Get("Settings"), m.Pos{X: CenterX, Y: ItemBaselineY(Settings, MainCount)}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Credits {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, locale.G.Get("Credits"), m.Pos{X: CenterX, Y: ItemBaselineY(Credits, MainCount)}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Quit {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, locale.G.Get("Quit"), m.Pos{X: CenterX, Y: ItemBaselineY(Quit, MainCount)}, true, fg, bg)

	// Display stats.
	font.MenuSmall.Draw(screen, fun.FormatText(&s.Controller.World.PlayerState, locale.G.Get("Score: {{Score}}{{SpeedrunCategoriesShort}} | Time: {{GameTime}}")),
		m.Pos{X: CenterX, Y: ItemBaselineY(-2, MainCount)}, true, fgn, bgn)
}

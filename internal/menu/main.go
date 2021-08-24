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
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/input"
	m "github.com/divVerent/aaaaxy/internal/math"
)

var (
	cheatShowFinalCredits = flag.Bool("cheat_show_final_credits", false, "show the final credits screen for testing")
)

type MainScreenItem int

const (
	Play MainScreenItem = iota
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
		Don't need this anymore - we have a quit entry.
		if input.Exit.JustHit {
			return s.Controller.ActivateSound(s.Controller.QuitGame())
		}
	*/
	if input.Jump.JustHit || input.Action.JustHit {
		switch s.Item {
		case Play:
			return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&MapScreen{}))
		case Settings:
			return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&SettingsScreen{}))
		case Credits:
			return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&CreditsScreen{Fancy: *cheatShowFinalCredits}))
		case Quit:
			return s.Controller.ActivateSound(s.Controller.QuitGame())
		}
	}
	return nil
}

func (s *MainScreen) Draw(screen *ebiten.Image) {
	h := engine.GameHeight
	x := engine.GameWidth / 2
	fgs := color.NRGBA{R: 255, G: 255, B: 85, A: 255}
	bgs := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	fgn := color.NRGBA{R: 170, G: 170, B: 170, A: 255}
	bgn := color.NRGBA{R: 85, G: 85, B: 85, A: 255}
	font.MenuBig.Draw(screen, "AAAAXY", m.Pos{X: x, Y: h / 4}, true, fgs, bgs)
	fg, bg := fgn, bgn
	if s.Item == Play {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, "Play", m.Pos{X: x, Y: 21 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Settings {
		fg, bg = fgs, bgs
	}
	// TODO: menu item for signs seen and coins gotten.
	font.Menu.Draw(screen, "Settings", m.Pos{X: x, Y: 23 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Credits {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, "Credits", m.Pos{X: x, Y: 25 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Quit {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, "Quit", m.Pos{X: x, Y: 27 * h / 32}, true, fg, bg)

	// Display stats.
	frames := s.Controller.World.PlayerState.Frames()
	ss, ms := frames/60, (frames%60)*1000/60
	mm, ss := ss/60, ss%60
	hh, mm := mm/60, mm%60
	font.MenuSmall.Draw(screen, fmt.Sprintf("Time: %d:%02d:%02d.%03d", hh, mm, ss, ms),
		m.Pos{X: x, Y: 17 * h / 32}, true, fgn, bgn)
	// TODO: also add
}

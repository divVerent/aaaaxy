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
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/font"
	"github.com/divVerent/aaaaaa/internal/input"
	m "github.com/divVerent/aaaaaa/internal/math"
)

type MainScreenItem int

const (
	Start MainScreenItem = iota
	Reset
	Credits
	Quit
	MainCount
)

type MainScreen struct {
	Menu *Menu
	Item MainScreenItem
}

func (s *MainScreen) Init(m *Menu) error {
	s.Menu = m
	return nil
}

func (s *MainScreen) Update() error {
	if input.Down.JustHit || input.Right.JustHit {
		s.Item++
	}
	if input.Left.JustHit || input.Up.JustHit {
		s.Item--
	}
	s.Item = MainScreenItem(m.Mod(int(s.Item), int(MainCount)))
	if input.Exit.JustHit {
		return s.Menu.QuitGame()
	}
	if input.Jump.JustHit || input.Action.JustHit {
		switch s.Item {
		case Start:
			return s.Menu.SwitchToGame()
		case Reset:
			// TODO
			return s.Menu.SwitchToScreen(&MainScreen{})
		case Credits:
			// TODO
			return s.Menu.SwitchToScreen(&MainScreen{})
		case Quit:
			return s.Menu.QuitGame()
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
	bgn := color.NRGBA{R: 85, G: 85, B: 85, A: 0}
	font.MenuBig.Draw(screen, "AAAAAA", m.Pos{X: x, Y: h / 4}, true, fgs, bgs)
	fg, bg := fgn, bgn
	if s.Item == Start {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, "Start", m.Pos{X: x, Y: 21 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Reset {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, "Reset", m.Pos{X: x, Y: 23 * h / 32}, true, fg, bg)
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
}

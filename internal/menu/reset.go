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
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/m"
	"github.com/divVerent/aaaaxy/internal/palette"
)

type ResetScreenItem int

const (
	ResetNothing = iota
	ResetConfig
	ResetGame
	BackToMain
	ResetCount
)

const resetFrames = 300

type ResetScreen struct {
	Controller                 *Controller
	Item                       ResetScreenItem
	ResetFrame                 int
	WaitForKeyReleaseThenReset bool
}

func (s *ResetScreen) Init(m *Controller) error {
	s.Controller = m
	return nil
}

func (s *ResetScreen) Update() error {
	clicked := s.Controller.QueryMouseItem(&s.Item, ResetCount)
	if input.Down.JustHit {
		s.Item++
		s.Controller.MoveSound(nil)
	}
	if input.Up.JustHit {
		s.Item--
		s.Controller.MoveSound(nil)
	}
	s.Item = ResetScreenItem(m.Mod(int(s.Item), int(ResetCount)))
	if s.Item == ResetGame {
		s.ResetFrame++
	} else {
		s.ResetFrame = 0
		s.WaitForKeyReleaseThenReset = false
	}
	if input.Exit.JustHit {
		return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&SettingsScreen{}))
	}
	if input.Jump.JustHit || input.Action.JustHit || clicked != NotClicked {
		switch s.Item {
		case ResetNothing:
			return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&SettingsScreen{}))
		case ResetConfig:
			flag.ResetToDefaults()
			return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&SettingsScreen{}))
		case ResetGame:
			if s.ResetFrame >= resetFrames {
				s.WaitForKeyReleaseThenReset = true
			}
		case BackToMain:
			return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&MainScreen{}))
		}
	}
	if s.WaitForKeyReleaseThenReset && !input.Jump.Held && !input.Action.Held {
		return s.Controller.ActivateSound(s.Controller.InitGame(resetGame))
	}
	return nil
}

func (s *ResetScreen) Draw(screen *ebiten.Image) {
	fgs := palette.EGA(palette.Yellow, 255)
	bgs := palette.EGA(palette.Black, 255)
	fgn := palette.EGA(palette.LightGrey, 255)
	bgn := palette.EGA(palette.DarkGrey, 255)
	font.ByName["MenuBig"].Draw(screen, locale.G.Get("Reset"), m.Pos{X: CenterX, Y: HeaderY}, font.Center, fgs, bgs)
	fg, bg := fgn, bgn
	if s.Item == ResetNothing {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("Reset Nothing"), m.Pos{X: CenterX, Y: ItemBaselineY(ResetNothing, ResetCount)}, font.Center, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == ResetConfig {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("Reset and Lose Settings"), m.Pos{X: CenterX, Y: ItemBaselineY(ResetConfig, ResetCount)}, font.Center, fg, bg)
	var resetText string
	var dx, dy int
	var save string
	switch *saveState {
	case 0:
		save = "A"
	case 1:
		save = "4"
	case 2:
		save = "X"
	case 3:
		save = "Y"
	default:
		save = fmt.Sprint(*saveState)
	}
	if s.ResetFrame >= resetFrames && s.Item == ResetGame {
		fg, bg = palette.EGA(palette.Red, 255), palette.EGA(palette.Black, 255)
		resetText = locale.G.Get("Reset and Lose SAVE STATE %s", save)
	} else {
		fg, bg = fgn, bgn
		if s.Item == ResetGame {
			fg, bg = palette.EGA(palette.LightRed, 255), palette.EGA(palette.Red, 255)
			if s.WaitForKeyReleaseThenReset {
				resetText = locale.G.Get("Reset and Lose Save State %s (just release buttons)", save)
			} else {
				resetText = locale.G.Get("Reset and Lose Save State %s (think about it for %d sec)", save, (resetFrames-s.ResetFrame+engine.GameTPS-1)/engine.GameTPS)
			}
		} else {
			resetText = locale.G.Get("Reset and Lose Save State %s", save)
		}
		dx = rand.Intn(3) - 1
		dy = rand.Intn(3) - 1
	}
	font.ByName["Menu"].Draw(screen, resetText, m.Pos{X: CenterX + dx, Y: ItemBaselineY(ResetGame, ResetCount) + dy}, font.Center, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == BackToMain {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("Main Menu"), m.Pos{X: CenterX, Y: ItemBaselineY(BackToMain, ResetCount)}, font.Center, fg, bg)
}

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
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/input"
	m "github.com/divVerent/aaaaxy/internal/math"
)

type ResetScreenItem int

const (
	ResetNothing ResetScreenItem = iota
	ResetConfig
	ResetGame
	BackToMain
	ResetCount
)

const resetFrames = 300

type ResetScreen struct {
	Controller *Controller
	Item       ResetScreenItem
	ResetFrame int
}

func (s *ResetScreen) Init(m *Controller) error {
	s.Controller = m
	return nil
}

func (s *ResetScreen) Update() error {
	if s.Item == ResetGame {
		s.ResetFrame++
	} else {
		s.ResetFrame = 0
	}
	if input.Down.JustHit {
		s.Item++
		s.Controller.MoveSound(nil)
	}
	if input.Up.JustHit {
		s.Item--
		s.Controller.MoveSound(nil)
	}
	s.Item = ResetScreenItem(m.Mod(int(s.Item), int(ResetCount)))
	if input.Exit.JustHit {
		return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&SettingsScreen{}))
	}
	if input.Jump.JustHit || input.Action.JustHit {
		switch s.Item {
		case ResetNothing:
			return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&SettingsScreen{}))
		case ResetConfig:
			flag.ResetToDefaults()
			return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&SettingsScreen{}))
		case ResetGame:
			if s.ResetFrame >= resetFrames {
				return s.Controller.ActivateSound(s.Controller.InitGame(resetGame))
			}
		case BackToMain:
			return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&MainScreen{}))
		}
	}
	return nil
}

func (s *ResetScreen) Draw(screen *ebiten.Image) {
	h := engine.GameHeight
	x := engine.GameWidth / 2
	fgs := color.NRGBA{R: 255, G: 255, B: 85, A: 255}
	bgs := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	fgn := color.NRGBA{R: 170, G: 170, B: 170, A: 255}
	bgn := color.NRGBA{R: 85, G: 85, B: 85, A: 255}
	font.MenuBig.Draw(screen, "Reset", m.Pos{X: x, Y: h / 4}, true, fgs, bgs)
	fg, bg := fgn, bgn
	if s.Item == ResetNothing {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, "Reset Nothing", m.Pos{X: x, Y: 23 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == ResetConfig {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, "Reset and Lose Settings", m.Pos{X: x, Y: 25 * h / 32}, true, fg, bg)
	var resetText string
	var dx, dy int
	save := ""
	switch *saveState {
	case 0:
		save = " A"
	case 1:
		save = " B"
	case 2:
		save = " C"
	case 3:
		save = " D"
	}
	if s.ResetFrame >= resetFrames && s.Item == ResetGame {
		fg, bg = color.NRGBA{R: 170, G: 0, B: 0, A: 255}, color.NRGBA{R: 0, G: 0, B: 0, A: 255}
		resetText = fmt.Sprintf("Reset and Lose SAVE STATE%s", save)
	} else {
		fg, bg = fgn, bgn
		if s.Item == ResetGame {
			fg, bg = color.NRGBA{R: 255, G: 85, B: 85, A: 255}, color.NRGBA{R: 170, G: 0, B: 0, A: 255}
		}
		if s.Item == ResetGame {
			resetText = fmt.Sprintf("Reset and Lose Save State%s (think about it for %d sec)", save, (resetFrames-s.ResetFrame+engine.GameTPS-1)/engine.GameTPS)
		} else {
			resetText = fmt.Sprintf("Reset and Lose Save State%s", save)
		}
		dx = rand.Intn(3) - 1
		dy = rand.Intn(3) - 1
	}
	font.Menu.Draw(screen, resetText, m.Pos{X: x + dx, Y: 27*h/32 + dy}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == BackToMain {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, "Main Menu", m.Pos{X: x, Y: 29 * h / 32}, true, fg, bg)
}

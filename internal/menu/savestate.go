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
	"encoding/json"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/fun"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/playerstate"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

type SaveStateScreenItem int

const (
	SaveStateA SaveStateScreenItem = iota
	SaveState4
	SaveStateX
	SaveStateY
	SaveExit
	SaveStateCount
)

type SaveStateScreen struct {
	Controller *Controller
	Item       SaveStateScreenItem
	Text       [4]string
}

func saveStateInfo(initLvl *level.Level, idx int) string {
	saveName := fmt.Sprintf("save-%d.json", idx)
	state, err := vfs.ReadState(vfs.SavedGames, saveName)
	if err != nil {
		return "(empty)"
	}
	save := &level.SaveGame{}
	err = json.Unmarshal(state, save)
	if err != nil {
		return "(empty)"
	}
	err = initLvl.LoadGame(save)
	if err != nil {
		return "(empty)"
	}
	ps := &playerstate.PlayerState{
		Level: initLvl,
	}
	format := "Score: {{Score}}{{SpeedrunCategoriesShort}}"
	if idx == *saveState {
		format += " (current)"
	} else {
		format += " | Time: {{GameTime}}"
	}
	return fun.FormatText(ps, format)
}

func (s *SaveStateScreen) Init(m *Controller) error {
	// TODO: Skip this loading step and have a different way to
	initLvl, err := level.Load("level")
	if err != nil {
		log.Fatalf("could not load level: %v", err)
	}

	s.Controller = m
	s.Text[0] = "A: " + saveStateInfo(initLvl, 0)
	s.Text[1] = "4: " + saveStateInfo(initLvl, 1)
	s.Text[2] = "X: " + saveStateInfo(initLvl, 2)
	s.Text[3] = "Y: " + saveStateInfo(initLvl, 3)
	switch *saveState {
	case 0:
		s.Item = SaveStateA
	case 1:
		s.Item = SaveState4
	case 2:
		s.Item = SaveStateX
	case 3:
		s.Item = SaveStateY
	default:
		s.Item = SaveExit
		return nil
	}
	return nil
}

func (s *SaveStateScreen) Update() error {
	if input.Down.JustHit {
		s.Item++
		s.Controller.MoveSound(nil)
	}
	if input.Up.JustHit {
		s.Item--
		s.Controller.MoveSound(nil)
	}
	s.Item = SaveStateScreenItem(m.Mod(int(s.Item), int(SaveStateCount)))
	if input.Exit.JustHit {
		return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&SettingsScreen{}))
	}
	if input.Jump.JustHit || input.Action.JustHit {
		switch s.Item {
		case SaveStateA:
			return s.Controller.ActivateSound(s.Controller.SwitchSaveState(0))
		case SaveState4:
			return s.Controller.ActivateSound(s.Controller.SwitchSaveState(1))
		case SaveStateX:
			return s.Controller.ActivateSound(s.Controller.SwitchSaveState(2))
		case SaveStateY:
			return s.Controller.ActivateSound(s.Controller.SwitchSaveState(3))
		case SaveExit:
			return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&SettingsScreen{}))
		}
	}
	return nil
}

func (s *SaveStateScreen) Draw(screen *ebiten.Image) {
	h := engine.GameHeight
	x := engine.GameWidth / 2
	fgs := color.NRGBA{R: 255, G: 255, B: 85, A: 255}
	bgs := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	fgn := color.NRGBA{R: 170, G: 170, B: 170, A: 255}
	bgn := color.NRGBA{R: 85, G: 85, B: 85, A: 255}
	font.MenuBig.Draw(screen, "Switch Save State", m.Pos{X: x, Y: h / 4}, true, fgs, bgs)
	fg, bg := fgn, bgn
	if s.Item == SaveStateA {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, s.Text[0], m.Pos{X: x, Y: 21 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == SaveState4 {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, s.Text[1], m.Pos{X: x, Y: 23 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == SaveStateX {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, s.Text[2], m.Pos{X: x, Y: 25 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == SaveStateY {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, s.Text[3], m.Pos{X: x, Y: 27 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == SaveExit {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, "Main Menu", m.Pos{X: x, Y: 29 * h / 32}, true, fg, bg)
}

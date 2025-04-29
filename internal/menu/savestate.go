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

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/fun"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/m"
	"github.com/divVerent/aaaaxy/internal/palette"
	"github.com/divVerent/aaaaxy/internal/playerstate"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

type SaveStateScreenItem int

const (
	SaveStateA = iota
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

func (s *SaveStateScreen) saveStateInfo(initLvl *level.Level, idx int) string {
	var ps *playerstate.PlayerState
	if idx == *saveState {
		ps = &s.Controller.World.PlayerState
	} else {
		saveName := engine.SaveName(idx)
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
		ps = &playerstate.PlayerState{
			Level: initLvl,
		}
	}
	format := locale.G.Get("Score: {{Score}}{{SpeedrunCategoriesShort}} | Time: {{GameTime}}")
	return fun.FormatText(ps, format)
}

func (s *SaveStateScreen) Init(m *Controller) error {
	s.Controller = m

	initLvl := s.Controller.World.Level.Clone()

	s.Text[0] = s.saveStateInfo(initLvl, 0)
	s.Text[1] = s.saveStateInfo(initLvl, 1)
	s.Text[2] = s.saveStateInfo(initLvl, 2)
	s.Text[3] = s.saveStateInfo(initLvl, 3)
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
	clicked := s.Controller.QueryMouseItem(&s.Item, SaveStateCount)

	// Update so one can always see which save state is current.
	if *saveState >= 0 && *saveState < 4 {
		s.Text[*saveState] = s.saveStateInfo(nil, *saveState)
	}

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
	if input.Jump.JustHit || input.Action.JustHit || clicked != NotClicked {
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
	fgs := palette.EGA(palette.Yellow, 255)
	bgs := palette.EGA(palette.Black, 255)
	fgn := palette.EGA(palette.LightGrey, 255)
	bgn := palette.EGA(palette.DarkGrey, 255)
	font.ByName["MenuBig"].Draw(screen, locale.G.Get("Switch Save State"), m.Pos{X: CenterX, Y: HeaderY}, font.Center, fgs, bgs)
	fg, bg := fgn, bgn
	if s.Item == SaveStateA {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("A: %s", s.Text[0]), m.Pos{X: CenterX, Y: ItemBaselineY(SaveStateA, SaveStateCount)}, font.Center, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == SaveState4 {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("4: %s", s.Text[1]), m.Pos{X: CenterX, Y: ItemBaselineY(SaveState4, SaveStateCount)}, font.Center, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == SaveStateX {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("X: %s", s.Text[2]), m.Pos{X: CenterX, Y: ItemBaselineY(SaveStateX, SaveStateCount)}, font.Center, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == SaveStateY {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("Y: %s", s.Text[3]), m.Pos{X: CenterX, Y: ItemBaselineY(SaveStateY, SaveStateCount)}, font.Center, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == SaveExit {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("Main Menu"), m.Pos{X: CenterX, Y: ItemBaselineY(SaveExit, SaveStateCount)}, font.Center, fg, bg)
}

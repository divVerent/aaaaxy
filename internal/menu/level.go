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
	"sort"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/locale"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/palette"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

type LevelScreenItem int

type LevelScreen struct {
	Controller *Controller
	Item       LevelScreenItem
}

var levels []string

func levelInfo(idx int) string {
	s := levels[idx]
	switch s {
	case "level":
		return "AAAAXY"
	default:
		return s
	}
}

func initLevels() error {
	l, err := vfs.ReadDir("maps")
	if err != nil {
		return fmt.Errorf("could not enumerate levels: %w", err)
	}
	for _, level := range l {
		name, isTMX := strings.CutSuffix(level, ".tmx")
		if !isTMX {
			continue
		}
		levels = append(levels, name)
	}
	sort.Slice(levels, func(i, j int) bool {
		return levelInfo(i) < levelInfo(j)
	})
	return nil
}

func (s *LevelScreen) Init(m *Controller) error {
	s.Controller = m

	s.Item = LevelScreenItem(len(levels))
	for i := range levels {
		if levels[i] == engine.LevelName() {
			s.Item = LevelScreenItem(i)
		}
	}

	return nil
}

func (s *LevelScreen) Update() error {
	clicked := s.Controller.QueryMouseItem(&s.Item, len(levels)+1)

	if input.Down.JustHit {
		s.Item++
		s.Controller.MoveSound(nil)
	}
	if input.Up.JustHit {
		s.Item--
		s.Controller.MoveSound(nil)
	}
	s.Item = LevelScreenItem(m.Mod(int(s.Item), len(levels)+1))
	if input.Exit.JustHit {
		return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&MainScreen{}))
	}
	if input.Jump.JustHit || input.Action.JustHit || clicked != NotClicked {
		if s.Item < LevelScreenItem(len(levels)) {
			return s.Controller.ActivateSound(s.Controller.SwitchLevel(levels[s.Item]))
		} else {
			return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&MainScreen{}))
		}
	}
	return nil
}

func (s *LevelScreen) Draw(screen *ebiten.Image) {
	fgs := palette.EGA(palette.Yellow, 255)
	bgs := palette.EGA(palette.Black, 255)
	fgn := palette.EGA(palette.LightGrey, 255)
	bgn := palette.EGA(palette.DarkGrey, 255)
	font.ByName["MenuBig"].Draw(screen, locale.G.Get("Switch World"), m.Pos{X: CenterX, Y: HeaderY}, font.Center, fgs, bgs)

	n := len(levels)

	for i, level := range levels {
		fg, bg := fgn, bgn
		if s.Item == LevelScreenItem(i) {
			fg, bg = fgs, bgs
		}
		font.ByName["Menu"].Draw(screen, locale.G.Get("%s: %s", level, levelInfo(i)), m.Pos{X: CenterX, Y: ItemBaselineY(i, n+1)}, font.Center, fg, bg)
	}

	fg, bg := fgn, bgn
	if s.Item == LevelScreenItem(len(levels)) {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("Main Menu"), m.Pos{X: CenterX, Y: ItemBaselineY(n, n+1)}, font.Center, fg, bg)
}

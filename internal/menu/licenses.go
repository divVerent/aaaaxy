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

	"github.com/divVerent/aaaaxy/internal/credits"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/input"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/palette"
)

const (
	licensesLineHeight = 12
	licensesFrames     = 3
	licensesStep       = 5
)

type LicensesScreen struct {
	Controller *Controller
	Frame      int // Subpixel accumulator.
	ScrollPos  int // Current scroll position.
}

func (s *LicensesScreen) Init(m *Controller) error {
	s.Controller = m
	s.ScrollPos = textScreenStartPos(credits.Licenses, licensesLineHeight)
	return nil
}

func (s *LicensesScreen) Update() error {
	exit := input.Exit.JustHit || input.Left.JustHit || input.Right.JustHit
	up := input.Up.Held
	down := input.Down.Held
	if pos, status := input.Mouse(); status != input.NoMouse {
		if pos.Y < engine.GameHeight/3 {
			up = true
		} else if pos.Y > 2*engine.GameHeight/3 {
			down = true
		} else if status == input.ClickingMouse {
			exit = true
		}
	}
	if exit {
		return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&MainScreen{}))
	}
	if up {
		s.ScrollPos = textScreenAdjustScrollUp(credits.Licenses, s.ScrollPos, licensesStep, licensesLineHeight)
		s.Frame = 0
	}
	if down {
		s.ScrollPos = textScreenAdjustScrollDown(credits.Licenses, s.ScrollPos, licensesStep, licensesLineHeight)
		s.Frame = 0
	}
	s.Frame++
	if s.Frame >= licensesFrames {
		s.ScrollPos = textScreenAdjustScrollDown(credits.Licenses, s.ScrollPos, 1, licensesLineHeight)
		s.Frame = 0
	}
	return nil
}

func (s *LicensesScreen) Draw(screen *ebiten.Image) {
	fg := palette.EGA(palette.LightGrey, 255)
	bg := palette.EGA(palette.Black, 255)
	pos := m.Pos{
		X: 64,
		Y: s.ScrollPos,
	}
	f := font.ByName["MonoSmall"]
	renderTextScreen(screen, f, f, credits.Licenses, pos, font.Left, licensesLineHeight, fg, bg, fg, bg)
}

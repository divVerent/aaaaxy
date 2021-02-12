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

	"github.com/divVerent/aaaaaa/internal/credits"
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/font"
	"github.com/divVerent/aaaaaa/internal/input"
	m "github.com/divVerent/aaaaaa/internal/math"
)

const (
	creditsLineHeight = 24
	creditsFrames     = 3
)

type CreditsScreen struct {
	Menu     *Menu
	Fancy    bool // With music, and constant speed - no scrolling. No exiting. Background image not needed - we use last game screen.
	Frame    int  // Current scroll position.
	MaxFrame int  // Maximum scroll position.
}

func (s *CreditsScreen) Init(m *Menu) error {
	s.Menu = m
	s.Frame = (-engine.GameHeight - creditsLineHeight) * creditsFrames
	s.MaxFrame = (creditsLineHeight*len(credits.Lines) - 3*creditsLineHeight/2 - engine.GameHeight/2) * creditsFrames
	return nil
}

func (s *CreditsScreen) Update() error {
	if s.Fancy {
		if s.Frame >= s.MaxFrame && input.Exit.JustHit {
			return s.Menu.ActivateSound(s.Menu.SwitchToScreen(&MainScreen{}))
		}
	} else {
		if input.Exit.JustHit {
			return s.Menu.ActivateSound(s.Menu.SwitchToScreen(&MainScreen{}))
		}
		if input.Up.Held {
			s.Frame -= 5
		}
		if input.Down.Held {
			s.Frame += 5
		}
	}
	s.Frame++
	if s.Frame > s.MaxFrame {
		s.Frame = s.MaxFrame
	}
	return nil
}

func (s *CreditsScreen) Draw(screen *ebiten.Image) {
	x := engine.GameWidth / 2
	fgs := color.NRGBA{R: 255, G: 255, B: 85, A: 255}
	bgs := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	fgn := color.NRGBA{R: 85, G: 255, B: 255, A: 255}
	bgn := color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	nextIsTitle := true
	for i, line := range credits.Lines {
		if line == "" {
			nextIsTitle = true
			continue
		}
		isTitle := nextIsTitle
		nextIsTitle = false
		y := creditsLineHeight*i - s.Frame/creditsFrames
		if y < 0 || y >= engine.GameHeight+creditsLineHeight {
			continue
		}
		// TODO fade in/out at screen edge
		if isTitle {
			font.MenuBig.Draw(screen, line, m.Pos{X: x, Y: y}, true, fgs, bgs)
		} else {
			font.Menu.Draw(screen, line, m.Pos{X: x, Y: y}, true, fgn, bgn)
		}
	}
}

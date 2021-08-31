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
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/credits"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/fun"
	"github.com/divVerent/aaaaxy/internal/input"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/music"
	"github.com/divVerent/aaaaxy/internal/playerstate"
	"github.com/divVerent/aaaaxy/internal/version"
)

const (
	creditsLineHeight = 24
	creditsFrames     = 3
)

type CreditsScreen struct {
	// Must be set when creating.
	Fancy bool // With music, and constant speed - no scrolling. No exiting. Background image not needed - we use last game screen.

	Controller *Controller
	Lines      []string // Actual lines to display.
	Frame      int      // Current scroll position.
	MaxFrame   int      // Maximum scroll position.
	Exits      int      // How often exit was pressed. Need to press 7 times to leave fancy credits.
}

func (s *CreditsScreen) Init(m *Controller) error {
	s.Controller = m
	s.Lines = append([]string{}, credits.Lines...)
	s.Lines = append(
		s.Lines,
		fmt.Sprintf("Level Version: %d", s.Controller.World.Level.SaveGameVersion),
		"Build: "+version.Revision(),
	)
	if s.Fancy {
		music.Switch(s.Controller.World.Level.CreditsMusic)
		cat := s.Controller.World.PlayerState.SpeedrunCategories()
		timeStr := fun.FormatText(&s.Controller.World.PlayerState, "{{GameTime}}")
		s.Lines = append(
			s.Lines,
			"",
			"Your Time",
			timeStr,
			"",
			"Your Speedrun Categories",
		)
		tryNext := ""
		categories := []string{}
		addCategory := func(cat string, have bool) {
			if have {
				categories = append(categories, cat)
			} else {
				if tryNext == "" {
					tryNext = cat
				}
			}
		}
		if flag.Cheating() {
			addCategory("Cheat%", true)
			addCategory("Without Cheating Of Course", false)
		}
		if cat&playerstate.HundredPercentSpeedrun == 0 {
			addCategory("Any%", cat&playerstate.AnyPercentSpeedrun != 0)
		}
		addCategory("100%", cat&playerstate.HundredPercentSpeedrun != 0)
		addCategory("All Notes", cat&playerstate.AllSignsSpeedrun != 0)
		addCategory("All Paths", cat&playerstate.AllPathsSpeedrun != 0)
		addCategory("All Secrets", cat&playerstate.AllSecretsSpeedrun != 0)
		addCategory("All Flipped", cat&playerstate.AllFlippedSpeedrun != 0)
		noEscape := "No Escape"
		if input.UsingGamepad() {
			noEscape = "No Start"
		}
		addCategory(noEscape, cat&playerstate.NoEscapeSpeedrun != 0)
		l := len(categories)
		switch l {
		case 0:
			s.Lines = append(s.Lines,
				"None")
		case 1:
			s.Lines = append(s.Lines,
				categories[0])
		default:
			s.Lines = append(s.Lines,
				strings.Join(categories[0:l-1], ", ")+" and "+categories[l-1])
		}
		if tryNext != "" {
			s.Lines = append(s.Lines,
				"",
				"Try Next",
				tryNext)
		}
		s.Lines = append(s.Lines,
			"",
			"Thank You!")
	}
	s.Frame = (-engine.GameHeight - creditsLineHeight) * creditsFrames
	s.MaxFrame = (creditsLineHeight*len(s.Lines) - 3*creditsLineHeight/2 - engine.GameHeight*7/8) * creditsFrames
	return nil
}

func (s *CreditsScreen) Update() error {
	if s.Fancy {
		if input.Exit.JustHit {
			s.Exits++
			if s.Frame >= s.MaxFrame {
				return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&MainScreen{}))
			} else if s.Exits >= 6 {
				s.Frame = s.MaxFrame
			}
		}
	} else {
		if input.Exit.JustHit {
			return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&MainScreen{}))
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
	for i, line := range s.Lines {
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

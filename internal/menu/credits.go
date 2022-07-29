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
	"github.com/divVerent/aaaaxy/internal/palette"
	"github.com/divVerent/aaaaxy/internal/playerstate"
	"github.com/divVerent/aaaaxy/internal/version"
)

var (
	cheatShowFinalCredits               = flag.Bool("cheat_show_final_credits", false, "show the final credits screen for testing")
	cheatFinalCreditsSpeedrunCategories = flag.Int("cheat_final_credits_speedrun_categories", -1, "speedrun categories to show for testing")
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
	if *cheatShowFinalCredits {
		s.Fancy = true
	}
	s.Controller = m
	s.Lines = append([]string{}, credits.Lines...)
	s.Lines = append(
		s.Lines,
		fmt.Sprintf("Level Version: %d", s.Controller.World.Level.SaveGameVersion),
		"Build: "+version.Revision(),
	)
	if s.Fancy {
		music.Switch(s.Controller.World.Level.CreditsMusic)
		timeStr := fun.FormatText(&s.Controller.World.PlayerState, "{{GameTime}}")
		cats := s.Controller.World.PlayerState.SpeedrunCategories()
		if *cheatFinalCreditsSpeedrunCategories >= 0 {
			cats = playerstate.SpeedrunCategories(*cheatFinalCreditsSpeedrunCategories)
		}
		categories, tryNext := cats.Describe()
		phrases := strings.Split(categories, ", ")
		categories1, categories2 := categories, ""
		if len(phrases) >= 3 {
			half := len(phrases)/2 + 1 // 3 -> 2|1, 4 -> 3|1, 5 -> 3|2, ...
			categories1 = strings.Join(phrases[:half], ", ") + ", "
			categories2 = strings.Join(phrases[half:], ", ")
		}
		s.Lines = append(
			s.Lines,
			"",
			"Your Time",
			timeStr,
			"",
			"Your Speedrun Categories",
			categories1,
		)
		if categories2 != "" {
			s.Lines = append(s.Lines,
				categories2)
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
	s.MaxFrame = (creditsLineHeight*len(s.Lines) - 0*creditsLineHeight - engine.GameHeight) * creditsFrames
	return nil
}

func (s *CreditsScreen) Update() error {
	exit := input.Exit.JustHit || input.Left.JustHit
	up := input.Up.Held
	down := input.Down.Held
	credits := input.Right.JustHit
	if pos, status := input.Mouse(); status != input.NoMouse {
		if pos.Y < engine.GameHeight/3 {
			up = true
		} else if pos.Y > 2*engine.GameHeight/3 {
			down = true
		} else if status == input.ClickingMouse {
			if pos.X > 2*engine.GameWidth/3 {
				credits = true
			} else {
				exit = true
			}
		}
	}
	if s.Fancy {
		if exit {
			s.Exits++
			if s.Frame >= s.MaxFrame {
				return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&MainScreen{}))
			} else if s.Exits >= 6 {
				s.Frame = s.MaxFrame
			}
		}
	} else {
		if exit {
			return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&MainScreen{}))
		}
		if credits {
			// TODO switch to credits dialog.
			return s.Controller.ActivateSound(s.Controller.SwitchToScreen(&MainScreen{}))
		}
		if up {
			s.Frame -= 5
		}
		if down {
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
	fgs := palette.EGA(palette.Yellow, 255)
	bgs := palette.EGA(palette.Black, 255)
	fgn := palette.EGA(palette.LightCyan, 255)
	bgn := palette.EGA(palette.Black, 0)
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
		if isTitle {
			font.MenuBig.Draw(screen, line, m.Pos{X: x, Y: y}, true, fgs, bgs)
		} else {
			font.Menu.Draw(screen, line, m.Pos{X: x, Y: y}, true, fgn, bgn)
		}
	}
}

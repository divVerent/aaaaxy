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

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/flag"
	"github.com/divVerent/aaaaaa/internal/font"
	"github.com/divVerent/aaaaaa/internal/input"
	m "github.com/divVerent/aaaaaa/internal/math"
)

type SettingsScreenItem int

const (
	Graphics SettingsScreenItem = iota
	Volume
	Reset
	Back
	SettingsCount
)

type SettingsScreen struct {
	Menu *Menu
	Item SettingsScreenItem
}

func (s *SettingsScreen) Init(m *Menu) error {
	s.Menu = m
	return nil
}

type graphicsSetting int

const (
	lowestGraphics graphicsSetting = iota
	lowGraphics
	mediumGraphics
	highGraphics
	maxGraphics
	graphicsSettingCount
)

func (s graphicsSetting) String() string {
	switch s {
	case maxGraphics:
		return "Max"
	case highGraphics:
		return "High"
	case mediumGraphics:
		return "Medium"
	case lowGraphics:
		return "Low"
	case lowestGraphics:
		return "Lowest"
	}
	return "???"
}

func currentGraphics() graphicsSetting {
	if flag.Get("screen_filter").(string) == "linear2xcrt" {
		return maxGraphics
	}
	if flag.Get("draw_outside").(bool) {
		return highGraphics
	}
	if flag.Get("draw_blurs").(bool) {
		return mediumGraphics
	}
	if flag.Get("expand_using_vertices_accurately").(bool) {
		return lowGraphics
	}
	return lowestGraphics
}

func (s graphicsSetting) apply() error {
	switch s {
	case maxGraphics:
		flag.Set("draw_blurs", true)
		flag.Set("draw_outside", true)
		flag.Set("draw_visibility_mask", true)
		flag.Set("expand_using_vertices", true)
		flag.Set("expand_using_vertices_accurately", true)
		flag.Set("screen_filter", "linear2xcrt")
	case highGraphics:
		flag.Set("draw_blurs", true)
		flag.Set("draw_outside", true)
		flag.Set("draw_visibility_mask", true)
		flag.Set("expand_using_vertices", true)
		flag.Set("expand_using_vertices_accurately", true)
		flag.Set("screen_filter", "linear2x")
	case mediumGraphics:
		flag.Set("draw_blurs", true)
		flag.Set("draw_outside", false)
		flag.Set("draw_visibility_mask", true)
		flag.Set("expand_using_vertices", true)
		flag.Set("expand_using_vertices_accurately", true)
		flag.Set("screen_filter", "linear2x")
	case lowGraphics:
		flag.Set("draw_blurs", false)
		flag.Set("draw_outside", false)
		flag.Set("draw_visibility_mask", true)
		flag.Set("expand_using_vertices", true)
		flag.Set("expand_using_vertices_accurately", true)
		flag.Set("screen_filter", "linear2x")
	case lowestGraphics:
		flag.Set("draw_blurs", false)
		flag.Set("draw_outside", false)
		flag.Set("draw_visibility_mask", true)
		flag.Set("expand_using_vertices", true)
		flag.Set("expand_using_vertices_accurately", false)
		flag.Set("screen_filter", "simple")
	}
	return nil
}

func toggleGraphics(delta int) error {
	g := currentGraphics()
	switch delta {
	case 0:
		g++
		if g >= graphicsSettingCount {
			g = 0
		}
	case -1:
		if g > 0 {
			g--
		}
	case +1:
		g++
		if g >= graphicsSettingCount {
			g--
		}
	}
	g.apply()
	return nil
}

func currentVolume() string {
	v := flag.Get("volume").(float64)
	return fmt.Sprintf("%.0f%%", v*100)
}

func toggleVolume(delta int) error {
	v := flag.Get("volume").(float64)
	switch delta {
	case 0:
		v += 0.1
		if v > 1 {
			v = 0
		}
	case -1:
		v -= 0.1
		if v < 0 {
			v = 0
		}
	case +1:
		v += 0.1
		if v > 1 {
			v = 1
		}
	}
	flag.Set("volume", v)
	return nil
}

func (s *SettingsScreen) Update() error {
	if input.Down.JustHit {
		s.Item++
		s.Menu.MoveSound(nil)
	}
	if input.Up.JustHit {
		s.Item--
		s.Menu.MoveSound(nil)
	}
	s.Item = SettingsScreenItem(m.Mod(int(s.Item), int(SettingsCount)))
	if input.Exit.JustHit {
		return s.Menu.ActivateSound(s.Menu.SwitchToScreen(&MainScreen{}))
	}
	if input.Jump.JustHit || input.Action.JustHit {
		switch s.Item {
		case Graphics:
			return s.Menu.ActivateSound(toggleGraphics(0))
		case Volume:
			return s.Menu.ActivateSound(toggleVolume(0))
		case Reset:
			return s.Menu.ActivateSound(s.Menu.SwitchToScreen(&ResetScreen{}))
		case Back:
			return s.Menu.ActivateSound(s.Menu.SwitchToScreen(&MainScreen{}))
		}
	}
	if input.Left.JustHit {
		switch s.Item {
		case Graphics:
			return s.Menu.ActivateSound(toggleGraphics(-1))
		case Volume:
			return s.Menu.ActivateSound(toggleVolume(-1))
		}
	}
	if input.Right.JustHit {
		switch s.Item {
		case Graphics:
			return s.Menu.ActivateSound(toggleGraphics(+1))
		case Volume:
			return s.Menu.ActivateSound(toggleVolume(+1))
		}
	}
	return nil
}

func (s *SettingsScreen) Draw(screen *ebiten.Image) {
	h := engine.GameHeight
	x := engine.GameWidth / 2
	fgs := color.NRGBA{R: 255, G: 255, B: 85, A: 255}
	bgs := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	fgn := color.NRGBA{R: 170, G: 170, B: 170, A: 255}
	bgn := color.NRGBA{R: 85, G: 85, B: 85, A: 255}
	font.MenuBig.Draw(screen, "Settings", m.Pos{X: x, Y: h / 4}, true, fgs, bgs)
	fg, bg := fgn, bgn
	if s.Item == Graphics {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, fmt.Sprintf("Graphics: %v", currentGraphics()), m.Pos{X: x, Y: 21 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Volume {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, fmt.Sprintf("Volume: %v", currentVolume()), m.Pos{X: x, Y: 23 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Reset {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, "Reset", m.Pos{X: x, Y: 25 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Back {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, "Main Menu", m.Pos{X: x, Y: 27 * h / 32}, true, fg, bg)
}

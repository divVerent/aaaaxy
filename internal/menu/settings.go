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

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/input"
	m "github.com/divVerent/aaaaxy/internal/math"
)

type SettingsScreenItem int

const (
	Fullscreen SettingsScreenItem = iota
	Graphics
	Quality
	Volume
	SaveState
	Reset
	Back
	SettingsCount
)

type SettingsScreen struct {
	Controller *Controller
	Item       SettingsScreenItem
}

func (s *SettingsScreen) Init(m *Controller) error {
	s.Controller = m
	return nil
}

type graphicsSetting int

const (
	herculesGraphics graphicsSetting = iota
	cgaGraphics
	egaGraphics
	vgaGraphics
	svgaGraphics
	graphicsSettingCount
)

func (s graphicsSetting) String() string {
	switch s {
	case herculesGraphics:
		return "Hercules" // Mostly Hercules, that is (2 more scan lines; who cares).
	case cgaGraphics:
		return "CGA" // Actually, it takes four CGA cards (or two Tandy) to render 640x360 in 4 colors.
	case egaGraphics:
		return "EGA" // Mostly EGA, that is (10 more scan lines; this probably can be hacked on real EGA though).
	case vgaGraphics:
		return "VGA"
	case svgaGraphics:
		return "SVGA"
	}
	return "???"
}

func currentGraphics() graphicsSetting {
	switch flag.Get("palette").(string) {
	case "mono":
		return herculesGraphics
	case "cga41h":
		return cgaGraphics
	case "ega":
		return egaGraphics
	case "vga":
		return vgaGraphics
	case "none":
		return svgaGraphics
	}
	return svgaGraphics
}

func (s graphicsSetting) apply() error {
	switch s {
	case herculesGraphics:
		flag.Set("palette", "mono")
	case cgaGraphics:
		flag.Set("palette", "cga41h")
	case egaGraphics:
		flag.Set("palette", "ega")
	case vgaGraphics:
		flag.Set("palette", "vga")
	case svgaGraphics:
		flag.Set("palette", "none")
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

type qualitySetting int

const (
	lowestQuality qualitySetting = iota
	lowQuality
	mediumQuality
	highQuality
	maxQuality
	qualitySettingCount
)

func (s qualitySetting) String() string {
	switch s {
	case maxQuality:
		return "Max"
	case highQuality:
		return "High"
	case mediumQuality:
		return "Medium"
	case lowQuality:
		return "Low"
	case lowestQuality:
		return "Lowest"
	}
	return "???"
}

func currentQuality() qualitySetting {
	if flag.Get("screen_filter").(string) == "linear2xcrt" {
		return maxQuality
	}
	if flag.Get("draw_outside").(bool) {
		return highQuality
	}
	if flag.Get("draw_blurs").(bool) {
		return mediumQuality
	}
	if flag.Get("expand_using_vertices_accurately").(bool) {
		return lowQuality
	}
	return lowestQuality
}

func (s qualitySetting) apply() error {
	switch s {
	case maxQuality:
		flag.Set("draw_blurs", true)
		flag.Set("draw_outside", true)
		flag.Set("expand_using_vertices_accurately", true)
		flag.Set("screen_filter", "linear2xcrt") // <-
	case highQuality:
		flag.Set("draw_blurs", true)
		flag.Set("draw_outside", true) // <-
		flag.Set("expand_using_vertices_accurately", true)
		flag.Set("screen_filter", "simple")
	case mediumQuality:
		flag.Set("draw_blurs", true) // <-
		flag.Set("draw_outside", false)
		flag.Set("expand_using_vertices_accurately", true)
		flag.Set("screen_filter", "simple")
	case lowQuality:
		flag.Set("draw_blurs", false)
		flag.Set("draw_outside", false)
		flag.Set("expand_using_vertices_accurately", true) // <-
		flag.Set("screen_filter", "simple")
	case lowestQuality:
		flag.Set("draw_blurs", false)
		flag.Set("draw_outside", false)
		flag.Set("expand_using_vertices_accurately", false)
		flag.Set("screen_filter", "simple")
	}
	return nil
}

func toggleQuality(delta int) error {
	g := currentQuality()
	switch delta {
	case 0:
		g++
		if g >= qualitySettingCount {
			g = 0
		}
	case -1:
		if g > 0 {
			g--
		}
	case +1:
		g++
		if g >= qualitySettingCount {
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
		s.Controller.MoveSound(nil)
	}
	if input.Up.JustHit {
		s.Item--
		s.Controller.MoveSound(nil)
	}
	s.Item = SettingsScreenItem(m.Mod(int(s.Item), int(SettingsCount)))
	if input.Exit.JustHit {
		return s.Controller.ActivateSound(s.Controller.SaveConfigAndSwitchToScreen(&MainScreen{}))
	}
	if input.Jump.JustHit || input.Action.JustHit {
		switch s.Item {
		case Fullscreen:
			return s.Controller.ActivateSound(s.Controller.toggleFullscreen())
		case Graphics:
			return s.Controller.ActivateSound(toggleGraphics(0))
		case Quality:
			return s.Controller.ActivateSound(toggleQuality(0))
		case Volume:
			return s.Controller.ActivateSound(toggleVolume(0))
		case SaveState:
			return s.Controller.ActivateSound(s.Controller.SaveConfigAndSwitchToScreen(&SaveStateScreen{}))
		case Reset:
			return s.Controller.ActivateSound(s.Controller.SaveConfigAndSwitchToScreen(&ResetScreen{}))
		case Back:
			return s.Controller.ActivateSound(s.Controller.SaveConfigAndSwitchToScreen(&MainScreen{}))
		}
	}
	if input.Left.JustHit {
		switch s.Item {
		case Fullscreen:
			return s.Controller.ActivateSound(s.Controller.toggleFullscreen())
		case Graphics:
			return s.Controller.ActivateSound(toggleGraphics(-1))
		case Quality:
			return s.Controller.ActivateSound(toggleQuality(-1))
		case Volume:
			return s.Controller.ActivateSound(toggleVolume(-1))
		}
	}
	if input.Right.JustHit {
		switch s.Item {
		case Fullscreen:
			return s.Controller.ActivateSound(s.Controller.toggleFullscreen())
		case Graphics:
			return s.Controller.ActivateSound(toggleGraphics(+1))
		case Quality:
			return s.Controller.ActivateSound(toggleQuality(+1))
		case Volume:
			return s.Controller.ActivateSound(toggleVolume(+1))
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
	if s.Item == Fullscreen {
		fg, bg = fgs, bgs
	}
	fsText := "Switch to Fullscreen Mode"
	if ebiten.IsFullscreen() {
		fsText = "Switch to Windowed Mode"
	}
	font.Menu.Draw(screen, fsText, m.Pos{X: x, Y: 17 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Graphics {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, fmt.Sprintf("Graphics: %v", currentGraphics()), m.Pos{X: x, Y: 19 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Quality {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, fmt.Sprintf("Quality: %v", currentQuality()), m.Pos{X: x, Y: 21 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Volume {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, fmt.Sprintf("Volume: %v", currentVolume()), m.Pos{X: x, Y: 23 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == SaveState {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, "Switch Save State", m.Pos{X: x, Y: 25 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Reset {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, "Reset", m.Pos{X: x, Y: 27 * h / 32}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Back {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, "Main Menu", m.Pos{X: x, Y: 29 * h / 32}, true, fg, bg)
}

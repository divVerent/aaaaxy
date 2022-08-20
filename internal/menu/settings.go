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

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/palette"
)

var offerFullscreen = flag.SystemDefault(map[string]interface{}{
	"android/*": false,
	"ios/*":     false,
	"*/*":       true,
}).(bool)

type SettingsScreenItem int

const (
	Fullscreen = iota
	Graphics
	Quality
	Volume
	SaveState
	Reset
	Back
	SettingsCount
)

type SettingsScreen struct {
	Controller      *Controller
	Item            SettingsScreenItem
	CurrentGraphics graphicsSetting
}

func (s *SettingsScreen) Init(m *Controller) error {
	s.Controller = m
	s.CurrentGraphics = currentGraphics()
	return nil
}

type graphicsSetting int

type graphicsSettingData struct {
	palette string
	name    string
}

var graphicsSettings = []graphicsSettingData{
	{"mono", "Hercules"},
	{"egamono", "EGA (monochrome)"},
	{"cga41l", "CGA"},
	{"ega", "EGA"},
	{"vga", "VGA"},
	{"none", "SVGA"},
	{"none", "SVGA"},
	{"none", "SVGA"},
	{"none", "SVGA"},
	{"none", "SVGA"},
	{"none", "SVGA"},
	{"none", "SVGA"},
	{"atarist", "Atari ST"},
	{"c64", "C64"},
	{"cga40n", "CGA (NTSC)"},
	{"gb", "Gameboy"},
	{"intellivision", "Intellivision"},
	{"macii", "Mac II"},
	{"msx", "MSX"},
	{"nes", "NES"},
	{"quake", "Quake"},
	{"web", "Web"},
	{"xterm", "XTerm"},
	{"3x3x3", "3x3x3"},
	{"4x4x4", "4x4x4"},
	{"7x7x4", "7x7x4"},
	{"8x8x4", "8x8x4"},
	{"8x8x4", "8x8x4"},
	{"8x8x4", "8x8x4"},
	{"8x8x4", "8x8x4"},
	{"8x8x4", "8x8x4"},
	{"8x8x4", "8x8x4"},
	{"8x8x4", "8x8x4"},
	{"ua3", "U" + "k" + "r" + "a" + "i" + "n" + "i" + "a" + "n"},
}

func init() {
	// Skip intentionally unused palettes.
	intentionallyUnused := map[string]bool{
		// Redundant, do not offer in menu.
		"cga40h": true,
		"cga40l": true,
		"cga41h": true,
		"cga41n": true,
		"cga5h":  true,
		"cga5l":  true,
		"cga6n":  true,
		"div0":   true,
		// Very low quality, do not offer in menu.
		"2x2x2":      true,
		"atarist4":   true,
		"cmyk":       true,
		"egalow":     true,
		"rgb":        true,
		"smb":        true,
		"vgadefault": true,
		// Country flags.
		"de3": true,
		"us4": true,
	}

	used := make(map[string]bool, len(graphicsSettings))
	for _, s := range graphicsSettings {
		if s.palette == "none" {
			continue
		}
		if intentionallyUnused[s.palette] {
			log.Fatalf("used intentionally unused graphics setting: %v", s.palette)
		}
		used[s.palette] = true
		if palette.ByName(s.palette) == nil {
			log.Fatalf("undefined graphics setting: %v", s.palette)
		}
	}

	var unused []string
	for _, name := range palette.Names() {
		if !used[name] && !intentionallyUnused[name] {
			unused = append(unused, name)
		}
	}
	if len(unused) != 0 {
		log.Fatalf("unused palette settings: %v", unused)
	}
}

func (s graphicsSetting) String() string {
	return graphicsSettings[s].name
}

func currentGraphics() graphicsSetting {
	pal := flag.Get("palette").(string)
	for i, s := range graphicsSettings {
		if s.palette == pal {
			return graphicsSetting(i)
		}
	}
	for i, s := range graphicsSettings {
		if s.palette == "none" {
			return graphicsSetting(i)
		}
	}
	return 0
}

func (s graphicsSetting) apply() error {
	flag.Set("palette", graphicsSettings[s].palette)
	return nil
}

func (s *SettingsScreen) toggleGraphics(delta int) error {
	count := graphicsSetting(len(graphicsSettings))
	switch delta {
	case 0:
		s.CurrentGraphics++
		if s.CurrentGraphics >= count {
			s.CurrentGraphics = 0
		}
	case -1:
		if s.CurrentGraphics > 0 {
			s.CurrentGraphics--
			for s.CurrentGraphics > 0 && graphicsSettings[s.CurrentGraphics].palette == graphicsSettings[s.CurrentGraphics-1].palette {
				s.CurrentGraphics--
			}
		}
	case +1:
		s.CurrentGraphics++
		if s.CurrentGraphics >= count {
			s.CurrentGraphics--
		}
	}
	s.CurrentGraphics.apply()
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
	var clicked Direction
	if offerFullscreen {
		clicked = s.Controller.QueryMouseItem(&s.Item, SettingsCount)
	} else {
		s.Item--
		clicked = s.Controller.QueryMouseItem(&s.Item, SettingsCount-1)
		s.Item++
	}
	if input.Down.JustHit {
		s.Item++
		s.Controller.MoveSound(nil)
	}
	if input.Up.JustHit {
		s.Item--
		s.Controller.MoveSound(nil)
	}
	if offerFullscreen {
		s.Item = SettingsScreenItem(m.Mod(int(s.Item), int(SettingsCount)))
	} else {
		s.Item = SettingsScreenItem(m.Mod(int(s.Item-1), int(SettingsCount-1)) + 1)
	}
	if input.Exit.JustHit {
		return s.Controller.ActivateSound(s.Controller.SaveConfigAndSwitchToScreen(&MainScreen{}))
	}
	if input.Jump.JustHit || input.Action.JustHit || clicked == CenterClicked {
		switch s.Item {
		case Fullscreen:
			return s.Controller.ActivateSound(s.Controller.toggleFullscreen())
		case Graphics:
			return s.Controller.ActivateSound(s.toggleGraphics(0))
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
	if input.Left.JustHit || clicked == LeftClicked {
		switch s.Item {
		case Fullscreen:
			return s.Controller.ActivateSound(s.Controller.toggleFullscreen())
		case Graphics:
			return s.Controller.ActivateSound(s.toggleGraphics(-1))
		case Quality:
			return s.Controller.ActivateSound(toggleQuality(-1))
		case Volume:
			return s.Controller.ActivateSound(toggleVolume(-1))
		}
	}
	if input.Right.JustHit || clicked == RightClicked {
		switch s.Item {
		case Fullscreen:
			return s.Controller.ActivateSound(s.Controller.toggleFullscreen())
		case Graphics:
			return s.Controller.ActivateSound(s.toggleGraphics(+1))
		case Quality:
			return s.Controller.ActivateSound(toggleQuality(+1))
		case Volume:
			return s.Controller.ActivateSound(toggleVolume(+1))
		}
	}
	return nil
}

func (s *SettingsScreen) Draw(screen *ebiten.Image) {
	fgs := palette.EGA(palette.Yellow, 255)
	bgs := palette.EGA(palette.Black, 255)
	fgn := palette.EGA(palette.LightGrey, 255)
	bgn := palette.EGA(palette.DarkGrey, 255)
	font.MenuBig.Draw(screen, "Settings", m.Pos{X: CenterX, Y: HeaderY}, true, fgs, bgs)
	if offerFullscreen {
		fg, bg := fgn, bgn
		if s.Item == Fullscreen {
			fg, bg = fgs, bgs
		}
		fsText := "Switch to Fullscreen Mode"
		if ebiten.IsFullscreen() {
			fsText = "Switch to Windowed Mode"
		}
		font.Menu.Draw(screen, fsText, m.Pos{X: CenterX, Y: ItemBaselineY(Fullscreen, SettingsCount)}, true, fg, bg)
	}
	fg, bg := fgn, bgn
	if s.Item == Graphics {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, fmt.Sprintf("Graphics: %v", currentGraphics()), m.Pos{X: CenterX, Y: ItemBaselineY(Graphics, SettingsCount)}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Quality {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, fmt.Sprintf("Quality: %v", currentQuality()), m.Pos{X: CenterX, Y: ItemBaselineY(Quality, SettingsCount)}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Volume {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, fmt.Sprintf("Volume: %v", currentVolume()), m.Pos{X: CenterX, Y: ItemBaselineY(Volume, SettingsCount)}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == SaveState {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, "Switch Save State", m.Pos{X: CenterX, Y: ItemBaselineY(SaveState, SettingsCount)}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Reset {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, "Reset", m.Pos{X: CenterX, Y: ItemBaselineY(Reset, SettingsCount)}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Back {
		fg, bg = fgs, bgs
	}
	font.Menu.Draw(screen, "Main Menu", m.Pos{X: CenterX, Y: ItemBaselineY(Back, SettingsCount)}, true, fg, bg)
}

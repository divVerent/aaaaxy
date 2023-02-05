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

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/game/misc"
	"github.com/divVerent/aaaaxy/internal/image"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/palette"
)

var offerFullscreen = flag.SystemDefault(map[string]bool{
	"android/*": false,
	"ios/*":     false,
	"*/*":       true,
})

type SettingsScreenItem int

const (
	Dynamic1 = iota
	Dynamic2
	Graphics
	Quality
	Volume
	Language
	SaveState
	Reset
	Back
	SettingsCount
)

type SettingsScreen struct {
	Controller      *Controller
	Item            SettingsScreenItem
	CurrentGraphics graphicsSetting
	CurrentLanguage languageSetting
	TopItem         SettingsScreenItem
	EditControls    SettingsScreenItem
	Fullscreen      SettingsScreenItem
}

func (s *SettingsScreen) Init(m *Controller) error {
	s.Controller = m
	s.CurrentGraphics = currentGraphics()
	s.CurrentLanguage.init()
	s.TopItem = Graphics
	if offerFullscreen {
		s.TopItem--
		s.Fullscreen = s.TopItem
	} else {
		s.Fullscreen = SettingsCount
	}
	if input.HaveTouch() {
		s.TopItem--
		s.EditControls = s.TopItem
	} else {
		s.EditControls = SettingsCount
	}
	s.Item = s.TopItem
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
	// Six more identical entries so one has to press right seven times to reach the custom palettes.
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
	{"egagray", "EGA (grayscale)"},
	{"vgagray", "VGA (grayscale)"},
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
	{"ua3", "Ukrainian"},
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
	pal := flag.Get[string]("palette")
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

func (s graphicsSetting) apply(m *Controller) error {
	palName := graphicsSettings[s].palette
	if palName == flag.Get[string]("palette") {
		return nil
	}
	flag.Set("palette", palName)

	pal := palette.ByName(palName)
	if !palette.SetCurrent(pal, flag.Get[bool]("palette_remap_colors")) {
		return nil
	}

	err := image.PaletteChanged()
	if err != nil {
		return fmt.Errorf("could not reapply palette to images: %v", err)
	}
	misc.ClearPrecache()
	err = engine.PaletteChanged()
	if err != nil {
		return fmt.Errorf("could not reapply palette to engine: %v", err)
	}
	err = m.GameChanged()
	if err != nil {
		return fmt.Errorf("could not reapply palette to menu: %v", err)
	}
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
	s.CurrentGraphics.apply(s.Controller)
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
	v := flag.Get[float64]("volume")
	return fmt.Sprintf("%.0f%%", v*100)
}

func toggleVolume(delta int) error {
	v := flag.Get[float64]("volume")
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
	saveItem := s.Item
	clicked := s.Controller.QueryMouseItem(&s.Item, SettingsCount)
	if s.Item < s.TopItem {
		clicked = NotClicked
		s.Item = saveItem
	}
	if input.Down.JustHit {
		s.Item++
		s.Controller.MoveSound(nil)
	}
	if input.Up.JustHit {
		s.Item--
		s.Controller.MoveSound(nil)
	}
	s.Item = SettingsScreenItem(m.Mod(int(s.Item-s.TopItem), int(SettingsCount-s.TopItem))) + s.TopItem
	if input.Exit.JustHit {
		return s.Controller.ActivateSound(s.Controller.SaveConfigAndSwitchToScreen(&MainScreen{}))
	}
	if input.Jump.JustHit || input.Action.JustHit || clicked == CenterClicked {
		switch s.Item {
		case s.Fullscreen:
			return s.Controller.ActivateSound(s.Controller.toggleFullscreen())
		case s.EditControls:
			return s.Controller.ActivateSound(s.Controller.SaveConfigAndSwitchToScreen(&TouchEditScreen{}))
		case Graphics:
			return s.Controller.ActivateSound(s.toggleGraphics(0))
		case Quality:
			return s.Controller.ActivateSound(toggleQuality(0))
		case Volume:
			return s.Controller.ActivateSound(toggleVolume(0))
		case Language:
			return s.Controller.ActivateSound(s.CurrentLanguage.toggle(s.Controller, 0))
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
		case s.Fullscreen:
			return s.Controller.ActivateSound(s.Controller.toggleFullscreen())
		case s.EditControls:
			return s.Controller.ActivateSound(s.Controller.SaveConfigAndSwitchToScreen(&TouchEditScreen{}))
		case Graphics:
			return s.Controller.ActivateSound(s.toggleGraphics(-1))
		case Quality:
			return s.Controller.ActivateSound(toggleQuality(-1))
		case Volume:
			return s.Controller.ActivateSound(toggleVolume(-1))
		case Language:
			return s.Controller.ActivateSound(s.CurrentLanguage.toggle(s.Controller, -1))
		}
	}
	if input.Right.JustHit || clicked == RightClicked {
		switch s.Item {
		case s.Fullscreen:
			return s.Controller.ActivateSound(s.Controller.toggleFullscreen())
		case s.EditControls:
			return s.Controller.ActivateSound(s.Controller.SaveConfigAndSwitchToScreen(&TouchEditScreen{}))
		case Graphics:
			return s.Controller.ActivateSound(s.toggleGraphics(+1))
		case Quality:
			return s.Controller.ActivateSound(toggleQuality(+1))
		case Volume:
			return s.Controller.ActivateSound(toggleVolume(+1))
		case Language:
			return s.Controller.ActivateSound(s.CurrentLanguage.toggle(s.Controller, +1))
		}
	}
	return nil
}

func (s *SettingsScreen) Draw(screen *ebiten.Image) {
	fgs := palette.EGA(palette.Yellow, 255)
	bgs := palette.EGA(palette.Black, 255)
	fgn := palette.EGA(palette.LightGrey, 255)
	bgn := palette.EGA(palette.DarkGrey, 255)
	font.ByName["MenuBig"].Draw(screen, locale.G.Get("Settings"), m.Pos{X: CenterX, Y: HeaderY}, true, fgs, bgs)
	if s.EditControls != SettingsCount {
		fg, bg := fgn, bgn
		if s.Item == s.EditControls {
			fg, bg = fgs, bgs
		}
		font.ByName["Menu"].Draw(screen, locale.G.Get("Edit Touch Controls"), m.Pos{X: CenterX, Y: ItemBaselineY(int(s.EditControls), SettingsCount)}, true, fg, bg)
	}
	if s.Fullscreen != SettingsCount {
		fg, bg := fgn, bgn
		if s.Item == s.Fullscreen {
			fg, bg = fgs, bgs
		}
		fsText := locale.G.Get("Switch to Fullscreen Mode")
		if ebiten.IsFullscreen() {
			fsText = locale.G.Get("Switch to Windowed Mode")
		}
		font.ByName["Menu"].Draw(screen, fsText, m.Pos{X: CenterX, Y: ItemBaselineY(int(s.Fullscreen), SettingsCount)}, true, fg, bg)
	}
	fg, bg := fgn, bgn
	if s.Item == Graphics {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("Graphics: %s", currentGraphics()), m.Pos{X: CenterX, Y: ItemBaselineY(Graphics, SettingsCount)}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Quality {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("Quality: %s", currentQuality()), m.Pos{X: CenterX, Y: ItemBaselineY(Quality, SettingsCount)}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Volume {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("Volume: %s", currentVolume()), m.Pos{X: CenterX, Y: ItemBaselineY(Volume, SettingsCount)}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Language {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("Language: %s", s.CurrentLanguage.name()), m.Pos{X: CenterX, Y: ItemBaselineY(Language, SettingsCount)}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == SaveState {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("Switch Save State"), m.Pos{X: CenterX, Y: ItemBaselineY(SaveState, SettingsCount)}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Reset {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("Reset"), m.Pos{X: CenterX, Y: ItemBaselineY(Reset, SettingsCount)}, true, fg, bg)
	fg, bg = fgn, bgn
	if s.Item == Back {
		fg, bg = fgs, bgs
	}
	font.ByName["Menu"].Draw(screen, locale.G.Get("Main Menu"), m.Pos{X: CenterX, Y: ItemBaselineY(Back, SettingsCount)}, true, fg, bg)
}

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

package aaaaaa

import (
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/audiowrap"
	"github.com/divVerent/aaaaaa/internal/credits"
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/flag"
	"github.com/divVerent/aaaaaa/internal/font"
	"github.com/divVerent/aaaaaa/internal/image"
	"github.com/divVerent/aaaaaa/internal/input"
	m "github.com/divVerent/aaaaaa/internal/math"
	"github.com/divVerent/aaaaaa/internal/noise"
	"github.com/divVerent/aaaaaa/internal/sound"
	"github.com/divVerent/aaaaaa/internal/vfs"
)

var (
	vsync             = flag.Bool("vsync", true, "enable waiting for vertical synchronization")
	fullscreen        = flag.Bool("fullscreen", true, "enable fullscreen mode")
	windowScaleFactor = flag.Float64("window_scale_factor", 0, "window scale factor in device pixels per game pixel (0 means auto integer scaling)")
)

func LoadConfig() (*flag.Config, error) {
	return engine.LoadConfig()
}

func setWindowSize() {
	f := *windowScaleFactor
	log.Printf("Requested window scale factor: %v", f)
	dscale := ebiten.DeviceScaleFactor()
	log.Printf("Device scale factor: %v", dscale)
	if f <= 0 {
		screenw, screenh := ebiten.ScreenSizeInFullscreen()
		log.Printf("Screen size: %vx%v", screenw, screenh)
		// Reserve 128 device independent pixels for system controls.
		maxw, maxh := screenw-128, screenh-128
		log.Printf("Max size: %vx%v", maxw, maxh)
		// Compute max scaling factors.
		maxwf, maxhf := float64(maxw)*dscale/engine.GameWidth, float64(maxh)*dscale/engine.GameHeight
		log.Printf("Max raw scale factor: %v, %v", maxwf, maxhf)
		f = maxwf
		if maxhf < f {
			f = maxhf
		}
		if f < 1 {
			f = 1
		}
		f = math.Floor(f)
		log.Printf("Chosen raw scale factor: %v", f)
	}
	f /= dscale
	w, h := m.Rint(engine.GameWidth*f), m.Rint(engine.GameHeight*f)
	log.Printf("Chosen window size: %vx%v", w, h)
	ebiten.SetWindowSize(w, h)
}

func InitEbiten() error {
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetFullscreen(*fullscreen)
	ebiten.SetInitFocused(true)
	ebiten.SetRunnableOnUnfocused(false)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetScreenTransparent(false)
	ebiten.SetVsyncEnabled(*vsync)
	ebiten.SetWindowDecorated(true)
	ebiten.SetWindowFloating(false)
	ebiten.SetWindowPosition(0, 0)
	ebiten.SetWindowResizable(true)
	setWindowSize()
	ebiten.SetWindowTitle("AAAAAA")

	err := vfs.Init()
	if err != nil {
		return fmt.Errorf("could not initialize VFS: %v", err)
	}
	err = input.Init()
	if err != nil {
		return fmt.Errorf("could not initialize input: %v", err)
	}
	err = font.Init()
	if err != nil {
		return fmt.Errorf("could not initialize fonts: %v", err)
	}
	err = credits.Precache()
	if err != nil {
		return fmt.Errorf("could not precache credits: %v", err)
	}
	err = image.Precache()
	if err != nil {
		return fmt.Errorf("could not precache images: %v", err)
	}
	err = audiowrap.Init()
	if err != nil {
		return fmt.Errorf("could not initialize audio: %v", err)
	}
	err = sound.Precache()
	if err != nil {
		return fmt.Errorf("could not precache sounds: %v", err)
	}
	err = noise.Init()
	if err != nil {
		return fmt.Errorf("could not initialize noise: %v", err)
	}
	err = initDumping()
	if err != nil {
		return fmt.Errorf("could not initialize dumping: %v", err)
	}

	if dumping() || *externalCapture {
		ebiten.SetMaxTPS(ebiten.UncappedTPS)
	} else {
		ebiten.SetMaxTPS(engine.GameTPS)
	}

	return nil
}

func BeforeExit() {
	finishDumping()
}

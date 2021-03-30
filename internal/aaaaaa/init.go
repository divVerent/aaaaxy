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

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"

	"github.com/divVerent/aaaaaa/internal/credits"
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/flag"
	"github.com/divVerent/aaaaaa/internal/font"
	"github.com/divVerent/aaaaaa/internal/image"
	"github.com/divVerent/aaaaaa/internal/input"
	"github.com/divVerent/aaaaaa/internal/noise"
	"github.com/divVerent/aaaaaa/internal/sound"
	"github.com/divVerent/aaaaaa/internal/vfs"
)

var (
	vsync      = flag.Bool("vsync", true, "enable waiting for vertical synchronization")
	fullscreen = flag.Bool("fullscreen", true, "enable fullscreen mode")
)

func LoadConfig() (*flag.Config, error) {
	return engine.LoadConfig()
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
	ebiten.SetWindowSize(engine.GameWidth, engine.GameHeight)
	ebiten.SetWindowTitle("AAAAAA")
	audio.NewContext(48000)

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

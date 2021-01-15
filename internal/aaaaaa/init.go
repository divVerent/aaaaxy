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
	"flag"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"

	"github.com/divVerent/aaaaaa/internal/engine"
)

var (
	vsync = flag.Bool("vsync", true, "enable waiting for vertical synchronization")
)

func InitEbiten() {
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetFullscreen(true)
	ebiten.SetInitFocused(true)
	ebiten.SetMaxTPS(engine.GameTPS)
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
}

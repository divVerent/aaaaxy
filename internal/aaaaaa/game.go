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
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/audiowrap"
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/flag"
	"github.com/divVerent/aaaaaa/internal/image"
	"github.com/divVerent/aaaaaa/internal/input"
	"github.com/divVerent/aaaaaa/internal/menu"
	"github.com/divVerent/aaaaaa/internal/music"
	"github.com/divVerent/aaaaaa/internal/noise"
	"github.com/divVerent/aaaaaa/internal/timing"
)

var (
	dumpVideo       = flag.String("dump_video", "", "filename prefix to dump game frames to")
	externalCapture = flag.Bool("external_dump", false, "assume an external dump application like apitrace is running; makes game run in lock step with rendering")
	loadGame        = flag.String("load_game", "", "filename to load game state from")
	saveGame        = flag.String("save_game", "", "filename to save game state to")
)

type Game struct {
	Menu menu.Menu
}

var _ ebiten.Game = &Game{}

func (g *Game) Update() error {
	timing.ReportRegularly()

	defer timing.Group()()
	timing.Section("update")
	defer timing.Group()()

	timing.Section("input")
	input.Update()

	timing.Section("menu")
	err := g.Menu.Update()
	if err != nil {
		return err
	}

	timing.Section("world")
	err = g.Menu.UpdateWorld()
	if err != nil {
		return err
	}

	// As the world's Update method may change the sound system info,
	// run this part last to reduce sound latency.
	timing.Section("music")
	music.Update()

	timing.Section("noise")
	noise.Update()

	return nil
}

var frameIndex = 0

func (g *Game) Draw(screen *ebiten.Image) {
	defer timing.Group()()
	timing.Section("draw")
	defer timing.Group()()

	timing.Section("world")
	g.Menu.DrawWorld(screen)

	if *dumpVideo != "" || *externalCapture {
		ebiten.SetMaxTPS(ebiten.UncappedTPS)
	}

	timing.Section("menu")
	g.Menu.Draw(screen)

	timing.Section("dump")
	if *dumpVideo != "" {
		image.Save(screen, fmt.Sprintf("%s_%08d.png", *dumpVideo, frameIndex))
	}
	frameIndex++
	audiowrap.Update(time.Duration(frameIndex) * time.Second / engine.GameTPS)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return engine.GameWidth, engine.GameHeight
}

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
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/divVerent/aaaaaa/internal/engine"
	_ "github.com/divVerent/aaaaaa/internal/game" // Load entities.
	"github.com/divVerent/aaaaaa/internal/image"
	"github.com/divVerent/aaaaaa/internal/music"
	"github.com/divVerent/aaaaaa/internal/noise"
	"github.com/divVerent/aaaaaa/internal/timing"
)

var (
	captureVideo    = flag.String("capture_video", "", "filename prefix to capture game frames to")
	externalCapture = flag.Bool("external_capture", false, "assume an external capture application like apitrace is running; makes game run in lock step with rendering")
	showFps         = flag.Bool("show_fps", false, "show fps counter")
	loadGame        = flag.String("load_game", "", "filename to load game state from")
	saveGame        = flag.String("save_game", "", "filename to save game state to")
)

type Game struct {
	World *engine.World
}

var _ ebiten.Game = &Game{}

func (g *Game) Update() error {
	timing.ReportRegularly()

	defer timing.Group()()
	timing.Section("update")
	defer timing.Group()()

	timing.Section("once")
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		if *saveGame != "" {
			file, err := os.Create(*saveGame)
			if err != nil {
				log.Panicf("could not open savegame: %v", err)
			}
			defer file.Close()
			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "\t")
			err = encoder.Encode(g.World.Level.SaveGame())
			if err != nil {
				log.Panicf("could not save game: %v", err)
			}
		}
		return errors.New("esc")
	}
	if g.World == nil {
		g.World = engine.NewWorld()
		if *loadGame != "" {
			file, err := os.Open(*loadGame)
			if err != nil {
				log.Panicf("could not open savegame: %v", err)
			}
			defer file.Close()
			decoder := json.NewDecoder(file)
			save := engine.SaveGame{}
			err = decoder.Decode(&save)
			if err != nil {
				log.Panicf("could not decode savegame: %v", err)
			}
			err = g.World.Level.LoadGame(save)
			if err != nil {
				log.Panicf("could not load savegame: %v", err)
			}
			cpName := g.World.Level.Player.PersistentState["last_checkpoint"]
			cpFlipped := g.World.Level.Player.PersistentState["checkpoint_seen."+cpName] == "FlipX"
			g.World.RespawnPlayer(cpName, cpFlipped)
		}
	}

	timing.Section("music")
	music.Update()

	timing.Section("noise")
	noise.Update()

	timing.Section("world")
	return g.World.Update()
}

var frameIndex = 0

func (g *Game) Draw(screen *ebiten.Image) {
	defer timing.Group()()
	timing.Section("draw")
	defer timing.Group()()

	timing.Section("world")
	g.World.Draw(screen)

	if *captureVideo != "" || *externalCapture {
		ebiten.SetMaxTPS(ebiten.UncappedTPS)
	}

	if *captureVideo != "" {
		timing.Section("capture")
		image.Save(screen, fmt.Sprintf("%s_%08d.png", *captureVideo, frameIndex))
		frameIndex++
	}

	timing.Section("hud")
	// TODO Draw HUD.

	timing.Section("menu")
	// TODO Draw menu.

	if *showFps {
		timing.Section("fps")
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.1f fps, %.1f tps", ebiten.CurrentFPS(), ebiten.CurrentTPS()), 0, engine.GameHeight-16)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return engine.GameWidth, engine.GameHeight
}

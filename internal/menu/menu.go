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
	"github.com/divVerent/aaaaaa/internal/timing"
)

var (
	showFps = flag.Bool("show_fps", false, "show fps counter")
)

type Menu struct{}

func (m *Menu) Update(world *engine.World) error {
	defer timing.Group()()

	timing.Section("once")
	if !world.Initialized() {
		world.Init()
		file, err := os.Open("save")
		if !os.IsNotExist(err) {
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
			err = world.Level.LoadGame(save)
			if err != nil {
				log.Panicf("could not load savegame: %v", err)
			}
			cpName := world.Level.Player.PersistentState["last_checkpoint"]
			cpFlipped := world.Level.Player.PersistentState["checkpoint_seen."+cpName] == "FlipX"
			world.RespawnPlayer(cpName, cpFlipped)
		}
	}

	timing.Section("global_hotkeys")
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		file, err := os.Create("save")
		if err != nil {
			log.Panicf("could not open savegame: %v", err)
		}
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "\t")
		err = encoder.Encode(world.Level.SaveGame())
		if err != nil {
			log.Panicf("could not save game: %v", err)
		}
		return errors.New("esc")
	}

	return nil
}

func (m *Menu) Draw(screen *ebiten.Image) {
	defer timing.Group()()

	timing.Section("global_overlays")
	if *showFps {
		timing.Section("fps")
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.1f fps, %.1f tps", ebiten.CurrentFPS(), ebiten.CurrentTPS()), 0, engine.GameHeight-16)
	}
}

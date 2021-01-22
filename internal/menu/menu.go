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
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/divVerent/aaaaaa/internal/engine"
	_ "github.com/divVerent/aaaaaa/internal/game" // Load entities.
	"github.com/divVerent/aaaaaa/internal/timing"
)

var (
	resetSave = flag.Bool("reset_save", false, "reset the savegame on startup")
	showFps   = flag.Bool("show_fps", false, "show fps counter")
)

type Menu struct{}

func (m *Menu) Update(world *engine.World) error {
	defer timing.Group()()

	timing.Section("once")
	if !world.Initialized() {
		world.Init()
		if !*resetSave {
			err := world.Load()
			if err != nil {
				return err
			}
		}
	}

	timing.Section("global_hotkeys")
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		err := world.Save()
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

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
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/flag"
	_ "github.com/divVerent/aaaaaa/internal/game" // Load entities.
	"github.com/divVerent/aaaaaa/internal/input"
	"github.com/divVerent/aaaaaa/internal/music"
	"github.com/divVerent/aaaaaa/internal/sound"
	"github.com/divVerent/aaaaaa/internal/timing"
)

var RegularTermination = errors.New("exited normally")

var (
	resetSave = flag.Bool("reset_save", false, "reset the savegame on startup")
	showFps   = flag.Bool("show_fps", false, "show fps counter")
)

const (
	blurSize     = 1
	darkenFactor = 0.75
)

type MenuScreen interface {
	Init(m *Menu) error
	Update() error
	Draw(screen *ebiten.Image)
}

type Menu struct {
	initialized   bool
	World         engine.World
	Screen        MenuScreen
	blurImage     *ebiten.Image
	moveSound     *sound.Sound
	activateSound *sound.Sound
}

func (m *Menu) Update() error {
	defer timing.Group()()

	timing.Section("once")
	if !m.initialized {
		err := m.World.Init()
		if err != nil {
			return fmt.Errorf("could not initialize world: %v", err)
		}
		if !*resetSave {
			err := m.World.Load()
			if err != nil {
				return err
			}
		}
		m.blurImage = ebiten.NewImage(engine.GameWidth, engine.GameHeight)
		m.moveSound, err = sound.Load("menu_move.ogg")
		if err != nil {
			return err
		}
		m.activateSound, err = sound.Load("menu_activate.ogg")
		if err != nil {
			return err
		}
		m.initialized = true
	}

	timing.Section("global_hotkeys")
	if input.Exit.JustHit && m.Screen == nil {
		music.Switch("")
		m.World.PlayerState.AddEscape()
		return m.SwitchToScreen(&MainScreen{})
	}
	if input.Fullscreen.JustHit {
		fs := !ebiten.IsFullscreen()
		flag.Set("fullscreen", fs)
		ebiten.SetFullscreen(fs)
	}

	timing.Section("screen")
	if m.Screen != nil {
		err := m.Screen.Update()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Menu) UpdateWorld() error {
	// Increment the frame counter.
	m.World.PlayerState.AddFrame()

	if m.Screen != nil {
		// Game is paused while in menu.
		return nil
	}
	return m.World.Update()
}

func (m *Menu) Draw(screen *ebiten.Image) {
	defer timing.Group()()

	timing.Section("screen")
	if m.Screen != nil {
		m.Screen.Draw(screen)
	}

	timing.Section("global_overlays")
	if *showFps {
		timing.Section("fps")
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.1f fps, %.1f tps", ebiten.CurrentFPS(), ebiten.CurrentTPS()), 0, engine.GameHeight-16)
	}
}

func (m *Menu) DrawWorld(screen *ebiten.Image) {
	m.World.Draw(screen)
	if m.Screen != nil {
		// If a menu screen is active, just draw the previous saved bitmap, but blur it.
		engine.BlurImage(screen, m.blurImage, screen, blurSize, darkenFactor, 0.0)
	}
}

// ResetGame is called by menu screens to reset the game.
func (m *Menu) ResetGame() error {
	m.World.Init()
	m.Screen = nil
	return nil
}

// SwitchToGame switches to a specific checkpoint.
func (m *Menu) SwitchToGame() error {
	m.Screen = nil
	return nil
}

// SwitchToCheckpoint switches to a specific checkpoint.
func (m *Menu) SwitchToCheckpoint(cp string) error {
	err := m.World.RespawnPlayer(cp)
	if err != nil {
		return fmt.Errorf("could not respawn player: %v", err)
	}
	m.Screen = nil
	return nil
}

// SwitchToScreen is called by menu screens to go to a different menu screen.
func (m *Menu) SwitchToScreen(screen MenuScreen) error {
	m.Screen = screen
	return m.Screen.Init(m)
}

// QuitGame is called by menu screens to end the game.
func (m *Menu) QuitGame() error {
	err := m.World.Save()
	if err != nil {
		return fmt.Errorf("could not save game: %v", err)
	}
	err = engine.SaveConfig()
	if err != nil {
		return fmt.Errorf("could not save config: %v", err)
	}
	return RegularTermination
}

// ActivateSound plays the sound effect to activate something.
func (m *Menu) ActivateSound(err error) error {
	if err == nil {
		m.activateSound.Play()
	}
	return err
}

// MoveSound plays the sound effect to activate something.
func (m *Menu) MoveSound(err error) error {
	if err == nil {
		m.moveSound.Play()
	}
	return err
}

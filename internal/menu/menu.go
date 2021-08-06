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

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	_ "github.com/divVerent/aaaaxy/internal/game" // Load entities.
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/music"
	"github.com/divVerent/aaaaxy/internal/sound"
	"github.com/divVerent/aaaaxy/internal/timing"
)

var RegularTermination = errors.New("exited normally")

var (
	saveState = flag.Int("save_state", 0, "number of save state slot")
	showFps   = flag.Bool("show_fps", false, "show fps counter")
)

const (
	blurSize     = 1
	blurFrames   = 32
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
	blurFrame     int
}

func (m *Menu) Update() error {
	defer timing.Group()()

	timing.Section("once")
	if !m.initialized {
		err := m.InitGame(loadGame)
		if err != nil {
			return err
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
	if m.World.ForceCredits {
		m.World.ForceCredits = false
		m.blurFrame = 0
		return m.SwitchToScreen(&CreditsScreen{Fancy: true})
	} else if input.Exit.JustHit && m.Screen == nil {
		m.World.TimerStarted = true
		music.Switch("")
		m.World.PlayerState.AddEscape()
		m.blurFrame = 0
		return m.SwitchToScreen(&MainScreen{})
	}
	if input.Fullscreen.JustHit {
		fs := !ebiten.IsFullscreen()
		flag.Set("fullscreen", fs)
		ebiten.SetFullscreen(fs)
	}

	timing.Section("screen")
	if m.blurFrame < blurFrames {
		m.blurFrame++
	}
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
	// Except when on the credits screen - that time does not count.
	if m.World.TimerStarted && !m.World.TimerStopped {
		m.World.PlayerState.AddFrame()
	}

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
		engine.BlurImage(screen, m.blurImage, screen, blurSize, darkenFactor, 0.0, float64(m.blurFrame)/blurFrames)
	}
}

type resetFlag int

const (
	loadGame resetFlag = iota
	resetGame
)

// InitGame is called by menu screens to load/reset the game.
func (m *Menu) InitGame(f resetFlag) error {
	// Stop the timer.
	m.World.TimerStarted = false

	// Initialize the world.
	err := m.World.Init(*saveState)
	if err != nil {
		return fmt.Errorf("could not initialize world: %v", err)
	}

	// Load the saved state.
	if f == loadGame {
		err := m.World.Load()
		if err != nil {
			return err
		}
	}

	// Go to the game screen.
	m.Screen = nil
	return nil
}

// SwitchSaveState switches to a given save state.
func (m *Menu) SwitchSaveState(state int) error {
	err := m.World.Save()
	if err != nil {
		return fmt.Errorf("could not save game: %v", err)
	}
	*saveState = state
	err = engine.SaveConfig()
	if err != nil {
		return fmt.Errorf("could not save config: %v", err)
	}
	return m.InitGame(loadGame)
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
	m.World.TimerStarted = true
	m.Screen = nil
	return nil
}

// SwitchToScreen is called by menu screens to go to a different menu screen.
func (m *Menu) SwitchToScreen(screen MenuScreen) error {
	m.Screen = screen
	return m.Screen.Init(m)
}

// SaveConfigAndSwitchToScreen is called by menu screens to go to a different menu screen.
func (m *Menu) SaveConfigAndSwitchToScreen(screen MenuScreen) error {
	err := engine.SaveConfig()
	if err != nil {
		return fmt.Errorf("could not save config: %v", err)
	}
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

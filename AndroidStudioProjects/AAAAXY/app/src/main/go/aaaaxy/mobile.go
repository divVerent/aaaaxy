// Copyright 2022 Google LLC
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

//go:build android
// +build android

package aaaaxy

import (
	"errors"
	"fmt"
	"time"

	"github.com/jeandeaual/go-locale"
	"golang.org/x/mobile/app"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/mobile"

	"github.com/divVerent/aaaaxy/internal/aaaaxy"
	"github.com/divVerent/aaaaxy/internal/exitstatus"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

// A Quitter is used to exit the game. The Quit() method will be implemented by MainActivity in Java.
type Quitter interface {
	Quit()
}

type game struct {
	game *aaaaxy.Game

	inited  bool
	drawErr error
}

var (
	g       *game
	quitter Quitter
)

// SetQuitter receives an object that can quit the game.
func SetQuitter(q Quitter) {
	quitter = q
}

func (g *game) Update() (err error) {
	ok := false
	defer func() {
		if !ok {
			err = fmt.Errorf("caught panic during update: %v", recover())
		}
		if err != nil {
			quitter.Quit()
		}
	}()
	if g.drawErr != nil {
		return g.drawErr
	}
	if !g.inited {
		g.inited = true
		locale.SetRunOnJVM(app.RunOnJVM)
		err = g.game.InitEarly()
	}
	if err == nil {
		err = g.game.Update()
		if err != nil {
			errbe := g.game.BeforeExit()
			if !errors.Is(err, exitstatus.RegularTermination) {
				log.Errorf("RunGame exited abnormally: %v", err)
			} else if errbe != nil {
				log.Errorf("BeforeExit exited abnormally: %v", errbe)
			}
		}
	}
	ok = true
	return err
}

func (g *game) Draw(screen *ebiten.Image) {
	if !g.inited {
		return
	}
	ok := false
	defer func() {
		if !ok {
			g.drawErr = fmt.Errorf("caught panic during draw: %v", recover())
		}
	}()
	g.game.Draw(screen)
	ok = true
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.game.Layout(outsideWidth, outsideHeight)
}

func init() {
	log.UsePanic(true)
	g = &game{
		game: aaaaxy.NewGame(),
	}
	mobile.SetGame(g)
}

// SetFilesDir forwards the location of the data files to the app.
func SetFilesDir(dir string) {
	vfs.SetFilesDir(dir)
}

// LoadConfig loads the configuration. To be called after SetFilesDir().
func LoadConfig() {
	// Sorry, some of the stuff SetGame does couldn't use flags then.
	flag.Parse(aaaaxy.LoadConfig)
}

// ForceBenchmarkDemo runs a benchmark demo instead of the game.
// This ignores the config, and should be called instead of LoadConfig() after SetFilesDir().
func ForceBenchmarkDemo() {
	flag.Parse(flag.NoConfig)
	flag.Set("debug_frame_profiling", true)
	flag.Set("debug_profiling", 10*time.Second)
	flag.Set("demo_play", "benchmark.dem")
	flag.Set("demo_timedemo", true)

	/*
	   // Settings for benchmarking:
	   flag.Set("auto_adjust_quality", false)
	   flag.Set("vsync", false)

	   // TEST.
	   flag.Set("pin_fonts_to_cache", true)

	   // Low settings:
	   flag.Set("palette", "none")
	   flag.Set("draw_blurs", false)
	   flag.Set("draw_outside", false)
	   flag.Set("expand_using_vertices_accurately", false)
	   flag.Set("screen_filter", "nearest")

	   // Highest settings:
	   flag.Set("palette", "vga")
	   flag.Set("draw_blurs", true)
	   flag.Set("draw_outside", true)
	   flag.Set("expand_using_vertices_accurately", true)
	   flag.Set("screen_filter", "linear2xcrt")
	*/
}

// BackPressed notifies the game that the back button has been pressed.
func BackPressed() {
	input.ExitPressed()
}

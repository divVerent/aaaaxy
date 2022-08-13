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
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/mobile"

	"github.com/divVerent/aaaaxy/internal/aaaaxy"
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
		err = g.game.InitEarly()
	}
	if err == nil {
		err = g.game.Update()
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
	// Only now we can actually load the config.
	// Sorry, some of the stuff SetGame does couldn't use flags then.
	flag.Parse(aaaaxy.LoadConfig)
}

// BackPressed notifies the game that the back button has been pressed.
func BackPressed() {
	input.ExitPressed()
}

// Copyright 2023 Google LLC
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

//go:build ios
// +build ios

package aaaaxy

import (
	"errors"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/mobile"

	"github.com/divVerent/aaaaxy/internal/aaaaxy"
	"github.com/divVerent/aaaaxy/internal/exitstatus"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
)

type game struct {
	game *aaaaxy.Game

	inited  bool
	drawErr error
}

var (
	g *game
)

func (g *game) Update() (err error) {
	ok := false
	defer func() {
		if !ok {
			err = fmt.Errorf("caught panic during update: %v", recover())
		}
		if err != nil {
			// Do We need to notify the ObjC side here? Android does:
			// quitter.Quit()
		}
	}()
	if g.drawErr != nil {
		return g.drawErr
	}
	if !g.inited {
		g.inited = true
		flag.Parse(aaaaxy.LoadConfig)
		err = g.game.InitEarly()
	}
	if err == nil {
		err = g.game.Update()
		if err != nil {
			errbe := g.game.BeforeExit()
			if !errors.Is(err, exitstatus.ErrRegularTermination) {
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

func (g *game) DrawFinalScreen(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
	if !g.inited {
		return
	}
	ok := false
	defer func() {
		if !ok {
			g.drawErr = fmt.Errorf("caught panic during final screen draw: %v", recover())
		}
	}()
	g.game.DrawFinalScreen(screen, offscreen, geoM)
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

// Dummy is an exported name to make ebitenmobile happy. It does nothing.
func Dummy() {
}

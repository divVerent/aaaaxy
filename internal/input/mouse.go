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

package input

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
	m "github.com/divVerent/aaaaxy/internal/math"
)

var (
	mouse = flag.Bool("mouse", true, "enable mouse input")
)

const (
	mouseHoverFrames = 5 * 60
	mouseBlockFrames = 30
)

var (
	mousePos        m.Pos
	mousePrevPos    m.Pos
	mouseHoverFrame int
	mouseBlockFrame int
	mouseClicking   bool
	mouseVisible    bool = true
	mouseWantClicks bool
)

func mouseUpdate(screenWidth, screenHeight, gameWidth, gameHeight int, crtK1, crtK2 float64) {
	wantVisible := *mouse && mouseWantClicks && mouseHoverFrame > 0
	if wantVisible != mouseVisible {
		mouseVisible = wantVisible
		if wantVisible {
			ebiten.SetCursorMode(ebiten.CursorModeVisible)
		} else {
			ebiten.SetCursorMode(ebiten.CursorModeHidden)
		}
	}

	if !*mouse {
		return
	}

	x, y := ebiten.CursorPosition()
	mousePos = pointerCoords(screenWidth, screenHeight, gameWidth, gameHeight, crtK1, crtK2, x, y)

	if mousePos != mousePrevPos {
		mouseHoverFrame = mouseHoverFrames
	}
	mousePrevPos = mousePos

	if mouseBlockFrame > 0 {
		mouseBlockFrame--
		mouseHoverFrame = 0
		mouseClicking = false
		return
	}

	if mouseHoverFrame > 0 {
		mouseHoverFrame--
		hoverPos = &mousePos
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseClicking = true
	} else if mouseClicking {
		// Click on release.
		mouseClicking = false
		clickPos = &mousePrevPos
	}
}

func mouseSetWantClicks(want bool) {
	mouseWantClicks = want
}

func mouseCancel() {
	mouseHoverFrame = 0
	mouseBlockFrame = mouseBlockFrames
}

func (i *impulse) mousePressed() InputMap {
	if !i.mouseControl {
		return 0
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		return AnyInput
	}
	return 0
}

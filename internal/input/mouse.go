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

	m "github.com/divVerent/aaaaxy/internal/math"
)

const (
	mouseHoverFrames = 5 * 60
)

var (
	mousePrevPos    m.Pos
	mouseHoverFrame int
	mouseClicking   bool
	mouseWantClicks bool = true
)

func mouseUpdate() {
	x, y := ebiten.CursorPosition()
	pos := m.Pos{X: x, Y: y}
	if pos != mousePrevPos {
		mouseHoverFrame = mouseHoverFrames
	}
	mousePrevPos = pos

	if mouseHoverFrame > 0 {
		mouseHoverFrame--
		hoverPos = &pos
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseClicking = true
	} else if mouseClicking {
		mouseClicking = false
		clickPos = &pos
	}
}

func mouseSetWantClicks(want bool) {
	if want == mouseWantClicks {
		return
	}
	mouseWantClicks = want
	if want {
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
	} else {
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
	}
}

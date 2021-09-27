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
)

var (
	leftKeys = []ebiten.Key{
		ebiten.KeyLeft, // DOS, NES.
		ebiten.KeyA,    // FPS.
		ebiten.KeyH,    // Vi.
	}
	rightKeys = []ebiten.Key{
		ebiten.KeyRight, // DOS, NES.
		ebiten.KeyD,     // FPS.
		ebiten.KeyL,     // Vi.
	}
	upKeys = []ebiten.Key{
		ebiten.KeyUp, // DOS, NES.
		ebiten.KeyW,  // FPS.
		ebiten.KeyK,  // Vi.
	}
	downKeys = []ebiten.Key{
		ebiten.KeyDown, // DOS, NES.
		ebiten.KeyS,    // FPS.
		ebiten.KeyJ,    // Vi.
	}
	jumpKeys = []ebiten.Key{
		ebiten.KeyControl, // DOS.
		ebiten.KeySpace,   // DOS, FPS, Vi.
		ebiten.KeyX,       // NES.
	}
	actionKeys = []ebiten.Key{
		ebiten.KeyAlt,   // DOS.
		ebiten.KeyShift, // DOS, FPS, Vi.
		ebiten.KeyZ,     // NES.
		ebiten.KeyTab,   // FPS, Vi.
		ebiten.KeyEnter, // Vi.
	}
	exitKeys = []ebiten.Key{
		ebiten.KeyEscape,    // DOS, NES, FPS, Vi.
		ebiten.KeyBackspace, // Vi.
	}
	fullscreenKeys = []ebiten.Key{
		ebiten.KeyF11, // Common.
		ebiten.KeyF,   // Common.
	}
)

func (i *impulse) keyboardPressed() bool {
	for _, k := range i.keys {
		if ebiten.IsKeyPressed(k) {
			return true
		}
	}
	return false
}

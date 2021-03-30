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
	leftKeys       = []ebiten.Key{ebiten.KeyLeft, ebiten.KeyA}
	rightKeys      = []ebiten.Key{ebiten.KeyRight, ebiten.KeyD}
	upKeys         = []ebiten.Key{ebiten.KeyUp, ebiten.KeyW}
	downKeys       = []ebiten.Key{ebiten.KeyDown, ebiten.KeyS}
	jumpKeys       = []ebiten.Key{ebiten.KeyControl, ebiten.KeySpace, ebiten.KeyX}
	actionKeys     = []ebiten.Key{ebiten.KeyAlt, ebiten.KeyE, ebiten.KeyZ, ebiten.KeyEnter}
	exitKeys       = []ebiten.Key{ebiten.KeyEscape}
	fullscreenKeys = []ebiten.Key{ebiten.KeyF11, ebiten.KeyF}
)

func (i *impulse) keyboardPressed() bool {
	for _, k := range i.keys {
		if ebiten.IsKeyPressed(k) {
			return true
		}
	}
	return false
}

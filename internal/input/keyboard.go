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
	leftKeys = map[ebiten.Key]InputMap{
		ebiten.KeyLeft: DOSKeyboard | NESKeyboard,
		ebiten.KeyA:    FPSKeyboard,
		ebiten.KeyH:    ViKeyboard,
	}
	rightKeys = map[ebiten.Key]InputMap{
		ebiten.KeyRight: DOSKeyboard | NESKeyboard,
		ebiten.KeyD:     FPSKeyboard,
		ebiten.KeyL:     ViKeyboard,
	}
	upKeys = map[ebiten.Key]InputMap{
		ebiten.KeyUp: DOSKeyboard | NESKeyboard,
		ebiten.KeyW:  FPSKeyboard,
		ebiten.KeyK:  ViKeyboard,
	}
	downKeys = map[ebiten.Key]InputMap{
		ebiten.KeyDown: DOSKeyboard | NESKeyboard,
		ebiten.KeyS:    FPSKeyboard,
		ebiten.KeyJ:    ViKeyboard,
	}
	jumpKeys = map[ebiten.Key]InputMap{
		ebiten.KeyControl: DOSKeyboard,
		ebiten.KeySpace:   DOSKeyboard | FPSKeyboard | ViKeyboard,
		ebiten.KeyX:       NESKeyboard,
	}
	actionKeys = map[ebiten.Key]InputMap{
		ebiten.KeyAlt:   DOSKeyboard,
		ebiten.KeyShift: DOSKeyboard | FPSKeyboard | ViKeyboard,
		ebiten.KeyE:     FPSKeyboard,
		ebiten.KeyZ:     NESKeyboard,
		ebiten.KeyTab:   FPSKeyboard | ViKeyboard,
		ebiten.KeyEnter: DOSKeyboard | ViKeyboard,
	}
	exitKeys = map[ebiten.Key]InputMap{
		ebiten.KeyEscape:    AnyKeyboardWithEscape,
		ebiten.KeyBackspace: AnyKeyboardWithBackspace,
	}
	fullscreenKeys = map[ebiten.Key]InputMap{
		ebiten.KeyF11: AnyInput,
		ebiten.KeyF:   AnyInput,
	}
)

func (i *impulse) keyboardPressed() InputMap {
	for k, m := range i.keys {
		if ebiten.IsKeyPressed(k) {
			return m
		}
	}
	return NoInput
}

func keyboardEasterEggKeyState() easterEggKeyState {
	var state easterEggKeyState
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		state |= easterEggA
	}
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		state |= easterEggX
	}
	if ebiten.IsKeyPressed(ebiten.KeyY) {
		state |= easterEggY
	}
	return state
}

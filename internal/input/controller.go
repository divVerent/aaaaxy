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

type impulse struct {
	Name    string
	Held    bool
	JustHit bool

	keys        []ebiten.Key
	padControls padControls
}

var (
	Left       = (&impulse{Name: "Left", keys: leftKeys, padControls: leftPad}).register()
	Right      = (&impulse{Name: "Right", keys: rightKeys, padControls: rightPad}).register()
	Up         = (&impulse{Name: "Up", keys: upKeys, padControls: upPad}).register()
	Down       = (&impulse{Name: "Down", keys: downKeys, padControls: downPad}).register()
	Jump       = (&impulse{Name: "Jump", keys: jumpKeys, padControls: jumpPad}).register()
	Action     = (&impulse{Name: "Action", keys: actionKeys, padControls: actionPad}).register()
	Exit       = (&impulse{Name: "Exit", keys: exitKeys, padControls: exitPad}).register()
	Fullscreen = (&impulse{Name: "Fullscreen", keys: fullscreenKeys /* no padControls */}).register()

	impulses = []*impulse{}

	usingGamepad bool

	// Wait for first frame to detect initial gamepad situation.
	firstUpdate = true
)

func (i *impulse) register() *impulse {
	impulses = append(impulses, i)
	return i
}

func (i *impulse) update() {
	keyboardHeld := i.keyboardPressed()
	gamepadHeld := i.gamepadPressed()
	held := keyboardHeld || gamepadHeld
	if held && !i.Held {
		i.JustHit = true
		// Whenever a new key is pressed, update the flag whether we're actually
		// _using_ the gamepad. Used for some in-game text messages.
		if keyboardHeld != gamepadHeld {
			usingGamepad = gamepadHeld
		}
	} else {
		i.JustHit = false
	}
	i.Held = held
}

func Init() error {
	gamepadInit()
	return nil
}

func Update() {
	gamepadScan()
	if firstUpdate {
		// At first, assume gamepad whenever one is present.
		usingGamepad = len(gamepads) > 0
		firstUpdate = false
	}
	for _, i := range impulses {
		i.update()
	}
}

func UsingGamepad() bool {
	return usingGamepad
}

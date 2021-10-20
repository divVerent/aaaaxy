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

type ImpulseState struct {
	Held    bool `json:",omitempty"`
	JustHit bool `json:",omitempty"`
}

type InputMap int

func (i InputMap) ContainsAny(o InputMap) bool {
	return i&o != 0
}

type impulse struct {
	ImpulseState
	Name string

	keys        map[ebiten.Key]InputMap
	padControls padControls
}

const (
	NoInput     InputMap = 0
	DOSKeyboard InputMap = 1
	NESKeyboard InputMap = 2
	FPSKeyboard InputMap = 4
	ViKeyboard  InputMap = 8
	AnyKeyboard InputMap = DOSKeyboard | NESKeyboard | FPSKeyboard | ViKeyboard
	Gamepad     InputMap = 16
	AnyInput    InputMap = AnyKeyboard | Gamepad
)

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

	inputMap InputMap

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
	held := keyboardHeld | gamepadHeld
	if held != 0 && !i.Held {
		i.JustHit = true
		// Whenever a new key is pressed, update the flag whether we're actually
		// _using_ the gamepad. Used for some in-game text messages.
		inputMap &= held
		if inputMap == 0 {
			inputMap = held
		}
	} else {
		i.JustHit = false
	}
	i.Held = held != 0
}

func Init() error {
	gamepadInit()
	return nil
}

func Update() {
	gamepadScan()
	if firstUpdate {
		// At first, assume gamepad whenever one is present.
		if len(gamepads) > 0 {
			inputMap = Gamepad
		} else {
			inputMap = DOSKeyboard
		}
		firstUpdate = false
	}
	for _, i := range impulses {
		i.update()
	}
	easterEggUpdate()
}

func EasterEggJustHit() bool {
	return easterEggJustHit
}

func Map() InputMap {
	return inputMap
}

// Demo code.

type DemoState struct {
	InputMap         InputMap
	Left             ImpulseState
	Right            ImpulseState
	Up               ImpulseState
	Down             ImpulseState
	Jump             ImpulseState
	Action           ImpulseState
	Exit             ImpulseState
	EasterEggJustHit bool
}

func LoadFromDemo(state *DemoState) {
	if state == nil {
		return
	}
	inputMap = state.InputMap
	Left.ImpulseState = state.Left
	Right.ImpulseState = state.Right
	Up.ImpulseState = state.Up
	Down.ImpulseState = state.Down
	Jump.ImpulseState = state.Jump
	Action.ImpulseState = state.Action
	Exit.ImpulseState = state.Exit
	easterEggJustHit = state.EasterEggJustHit
}

func SaveToDemo() *DemoState {
	return &DemoState{
		InputMap:         inputMap,
		Left:             Left.ImpulseState,
		Right:            Right.ImpulseState,
		Up:               Up.ImpulseState,
		Down:             Down.ImpulseState,
		Jump:             Jump.ImpulseState,
		Action:           Action.ImpulseState,
		Exit:             Exit.ImpulseState,
		EasterEggJustHit: easterEggJustHit,
	}
}

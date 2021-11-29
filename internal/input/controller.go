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
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImpulseState struct {
	Held    bool `json:",omitempty"`
	JustHit bool `json:",omitempty"`
}

func (i *ImpulseState) Empty() bool {
	return !i.Held && !i.JustHit
}

func (i *ImpulseState) OrEmpty() ImpulseState {
	if i == nil {
		return ImpulseState{}
	}
	return *i
}

func (i *ImpulseState) UnlessEmpty() *ImpulseState {
	if i.Empty() {
		return nil
	}
	return i
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

type ExitButtonID int

var exitKey = Escape

const (
	Escape ExitButtonID = iota
	Backspace
	Start
)

func ExitButton() ExitButtonID {
	if inputMap.ContainsAny(Gamepad) {
		return Start
	}
	if runtime.GOOS == "js" {
		// On JS, the Esc key is kinda "reserved" for leaving fullsreeen.
		return Backspace
	}
	return exitKey
}

// Demo code.

type DemoState struct {
	InputMap         InputMap      `json:",omitempty"`
	ExitKey          ExitButtonID  `json:",omitempty"`
	Left             *ImpulseState `json:",omitempty"`
	Right            *ImpulseState `json:",omitempty"`
	Up               *ImpulseState `json:",omitempty"`
	Down             *ImpulseState `json:",omitempty"`
	Jump             *ImpulseState `json:",omitempty"`
	Action           *ImpulseState `json:",omitempty"`
	Exit             *ImpulseState `json:",omitempty"`
	EasterEggJustHit bool          `json:",omitempty"`
}

func LoadFromDemo(state *DemoState) {
	if state == nil {
		state = &DemoState{}
	}
	inputMap = state.InputMap
	exitKey = state.ExitKey
	Left.ImpulseState = state.Left.OrEmpty()
	Right.ImpulseState = state.Right.OrEmpty()
	Up.ImpulseState = state.Up.OrEmpty()
	Down.ImpulseState = state.Down.OrEmpty()
	Jump.ImpulseState = state.Jump.OrEmpty()
	Action.ImpulseState = state.Action.OrEmpty()
	Exit.ImpulseState = state.Exit.OrEmpty()
	easterEggJustHit = state.EasterEggJustHit
}

func SaveToDemo() *DemoState {
	return &DemoState{
		InputMap:         inputMap,
		ExitKey:          exitKey,
		Left:             Left.ImpulseState.UnlessEmpty(),
		Right:            Right.ImpulseState.UnlessEmpty(),
		Up:               Up.ImpulseState.UnlessEmpty(),
		Down:             Down.ImpulseState.UnlessEmpty(),
		Jump:             Jump.ImpulseState.UnlessEmpty(),
		Action:           Action.ImpulseState.UnlessEmpty(),
		Exit:             Exit.ImpulseState.UnlessEmpty(),
		EasterEggJustHit: easterEggJustHit,
	}
}

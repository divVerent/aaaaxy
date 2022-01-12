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
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

var (
	gamepadAxisOnThreshold  = flag.Float64("gamepad_axis_on_threshold", 0.6, "Minimum amount to push the game pad for registering an action. Can be zero to accept any movement.")
	gamepadAxisOffThreshold = flag.Float64("gamepad_axis_off_threshold", 0.4, "Maximum amount to push the game pad for unregistering an action. Can be zero to accept any movement.")
	gamepadOverride         = flag.String("gamepad_override", "", "Entries in SDL_GameControllerDB format to add/override gamepad support. Multiple entries are permitted and can be separated by newlines or semicolons. Can also be provided via $SDL_GAMECONTROLLERCONFIG environment variable.")
)

type (
	padControls struct {
		buttons       []ebiten.StandardGamepadButton
		axes          []ebiten.StandardGamepadAxis
		axisDirection float64
	}
)

var (
	leftPad = padControls{
		buttons: []ebiten.StandardGamepadButton{
			ebiten.StandardGamepadButtonLeftLeft,
		},
		axes: []ebiten.StandardGamepadAxis{
			ebiten.StandardGamepadAxisLeftStickHorizontal,
			ebiten.StandardGamepadAxisRightStickHorizontal,
		},
		axisDirection: -1,
	}
	rightPad = padControls{
		buttons: []ebiten.StandardGamepadButton{
			ebiten.StandardGamepadButtonLeftRight,
		},
		axes: []ebiten.StandardGamepadAxis{
			ebiten.StandardGamepadAxisLeftStickHorizontal,
			ebiten.StandardGamepadAxisRightStickHorizontal,
		},
		axisDirection: +1,
	}
	upPad = padControls{
		buttons: []ebiten.StandardGamepadButton{
			ebiten.StandardGamepadButtonLeftTop,
		},
		axes: []ebiten.StandardGamepadAxis{
			ebiten.StandardGamepadAxisLeftStickVertical,
			ebiten.StandardGamepadAxisRightStickVertical,
		},
		axisDirection: -1,
	}
	downPad = padControls{
		buttons: []ebiten.StandardGamepadButton{
			ebiten.StandardGamepadButtonLeftBottom,
		},
		axes: []ebiten.StandardGamepadAxis{
			ebiten.StandardGamepadAxisLeftStickVertical,
			ebiten.StandardGamepadAxisRightStickVertical,
		},
		axisDirection: +1,
	}
	jumpPad = padControls{
		buttons: []ebiten.StandardGamepadButton{
			ebiten.StandardGamepadButtonRightBottom,
			ebiten.StandardGamepadButtonRightTop,
			ebiten.StandardGamepadButtonFrontBottomRight,
		},
	}
	actionPad = padControls{
		buttons: []ebiten.StandardGamepadButton{
			ebiten.StandardGamepadButtonRightLeft,
			ebiten.StandardGamepadButtonRightRight,
			ebiten.StandardGamepadButtonFrontBottomLeft,
		},
	}
	exitPad = padControls{
		buttons: []ebiten.StandardGamepadButton{
			ebiten.StandardGamepadButtonFrontTopLeft,
			ebiten.StandardGamepadButtonFrontTopRight,
			ebiten.StandardGamepadButtonCenterLeft,
			ebiten.StandardGamepadButtonCenterRight,
			ebiten.StandardGamepadButtonCenterCenter,
		},
	}

// Ignore ebiten.StandardGamepadButtonLeftStick.
// Ignore ebiten.StandardGamepadButtonRightStick.
)

var (
	// gamepadInvAxisOnThreshold is 1.0 divided by the variable gamepadAxisOnThreshold. Done to save a division for every axis test.
	gamepadInvAxisOnThreshold float64
	// gamepadInvAxisOffThreshold is 1.0 divided by the variable gamepadAxisOffThreshold. Done to save a division for every axis test.
	gamepadInvAxisOffThreshold float64
	// gamepads is the set of currently active gamepads. The boolean value should always be true, except during rescanning, where it's set to false temporarily to detect removed gamepads.
	gamepads = map[ebiten.GamepadID]struct{}{}
	// allGamepads is the set of all gamepads, even unsupported ones.
	allGamepads = map[ebiten.GamepadID]bool{}
	// allGamepadsList is the list of all gamepads. Global to reduce allocation.
	allGamepadsList []ebiten.GamepadID
)

func (i *impulse) gamepadPressed() InputMap {
	t := *gamepadAxisOnThreshold
	if i.Held {
		t = *gamepadAxisOffThreshold
	}
	for p := range gamepads {
		for _, b := range i.padControls.buttons {
			if ebiten.IsStandardGamepadButtonPressed(p, b) {
				return Gamepad
			}
		}
		for _, a := range i.padControls.axes {
			if ebiten.StandardGamepadAxisValue(p, a)*i.padControls.axisDirection >= t {
				return Gamepad
			}
		}
	}
	return NoInput
}

func gamepadScan() {
	// List new gamepads.
	allGamepadsList = ebiten.AppendGamepadIDs(allGamepadsList[:0])
	// Detect added/removed gamepads.
	for p := range allGamepads {
		allGamepads[p] = false
	}
	for _, p := range allGamepadsList {
		_, alreadyThere := allGamepads[p]
		if alreadyThere {
			allGamepads[p] = true
			continue
		}
		log.Infof("gamepad %v added", ebiten.GamepadName(p))
		allGamepads[p] = true
		if !ebiten.IsStandardGamepadLayoutAvailable(p) {
			log.Errorf("gamepad %v has no standard layout - cannot use", ebiten.GamepadName(p))
			continue
		}
		// A good gamepad! Add it.
		gamepads[p] = struct{}{}
	}
	for p, stillThere := range allGamepads {
		if stillThere {
			continue
		}
		log.Infof("gamepad %v removed", ebiten.GamepadName(p))
		delete(allGamepads, p)
		delete(gamepads, p)
	}
}

func applyAndLogGameControllerDb(config string, err error, name string) {
	if os.IsNotExist(err) {
		log.Infof("%v not provided - OK", name)
		return
	}
	if err != nil {
		log.Warningf("could not load %v: %v", name, err)
		return
	}
	if config == "" {
		log.Infof("%v not provided - OK", name)
		return
	}
	anyApplied := false
	anyError := false
	for line, entry := range strings.Split(config, "\n") {
		applied, err := ebiten.UpdateStandardGamepadLayoutMappings(entry)
		if err != nil {
			log.Errorf("could not add %v line %v: %v", name, line+1, err)
			anyError = true
		}
		if applied {
			anyApplied = true
		}
	}
	if anyError {
		// Already logged something, no need to log more.
	} else if anyApplied {
		log.Infof("%v applied", name)
	} else {
		log.Infof("%v exist but are not used on this platform", name)
	}
}

func readBuiltinGamepadMappings() (string, error) {
	configHandle, err := vfs.Load("input", "gamecontrollerdb.txt")
	if err != nil {
		return "", fmt.Errorf("open: %v", err)
	}
	defer configHandle.Close()
	configBytes, err := ioutil.ReadAll(configHandle)
	if err != nil {
		return "", fmt.Errorf("read: %v", err)
	}
	return string(configBytes), nil
}

func gamepadInit() {
	// Note: we're also stripping spaces before/after a semicolon
	// as a user might be putting some given they're usual in English,
	// yet they're technically invalid in SDL_GameControllerDB format.
	semiRE := regexp.MustCompile(`\s*;\s*`)

	// Support an included gamecontrollerdb.txt override.
	// Doing this because Ebiten's lags behind.
	mappings, err := readBuiltinGamepadMappings()
	applyAndLogGameControllerDb(mappings, err, "included gamepad mappings")

	// Support ~/.config/AAAAXY/gamecontrollerdb.txt.
	configBytes, err := vfs.ReadState(vfs.Config, "gamecontrollerdb.txt")
	applyAndLogGameControllerDb(string(configBytes), err, "gamepad mappings from gamecontrollerdb.txt")

	// Support the environment variable.
	config := semiRE.ReplaceAllString(os.Getenv("SDL_GAMECONTROLLERCONFIG"), "\n")
	applyAndLogGameControllerDb(config, nil, "gamepad mappings from $SDL_GAMECONTROLLERCONFIG")

	// Also support the flag. Note that the flag value is saved.
	config = semiRE.ReplaceAllString(*gamepadOverride, "\n")
	applyAndLogGameControllerDb(config, nil, "gamepad mappings from --gamepad_override")
}

func gamepadEasterEggKeyState() easterEggKeyState {
	var state easterEggKeyState
	for p := range gamepads {
		if ebiten.IsStandardGamepadButtonPressed(p, ebiten.StandardGamepadButtonRightBottom) {
			state |= easterEggA
		}
		if ebiten.IsStandardGamepadButtonPressed(p, ebiten.StandardGamepadButtonRightLeft) {
			state |= easterEggX
		}
		if ebiten.IsStandardGamepadButtonPressed(p, ebiten.StandardGamepadButtonRightTop) {
			state |= easterEggY
		}
	}
	return state
}

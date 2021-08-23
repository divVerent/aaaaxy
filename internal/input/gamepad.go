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
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	gamepadAxisOnThreshold  = flag.Float64("gamepad_axis_on_threshold", 0.6, "Minimum amount to push the game pad for registering an action. Can be zero to accept any movement.")
	gamepadAxisOffThreshold = flag.Float64("gamepad_axis_off_threshold", 0.4, "Maximum amount to push the game pad for unregistering an action. Can be zero to accept any movement.")
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
			ebiten.StandardGamepadButtonRightLeft,
			ebiten.StandardGamepadButtonRightBottom,
			ebiten.StandardGamepadButtonFrontBottomRight,
		},
	}
	actionPad = padControls{
		buttons: []ebiten.StandardGamepadButton{
			ebiten.StandardGamepadButtonRightTop,
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

func (i *impulse) gamepadPressed() bool {
	t := *gamepadAxisOnThreshold
	if i.Held {
		t = *gamepadAxisOffThreshold
	}
	for p := range gamepads {
		for _, b := range i.padControls.buttons {
			if ebiten.IsStandardGamepadButtonPressed(p, b) {
				return true
			}
		}
		for _, a := range i.padControls.axes {
			if ebiten.StandardGamepadAxisValue(p, a)*i.padControls.axisDirection >= t {
				return true
			}
		}
	}
	return false
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
			continue
		}
		log.Infof("Gamepad %v added.", ebiten.GamepadName(p))
		allGamepads[p] = true
		if !ebiten.IsStandardGamepadLayoutAvailable(p) {
			log.Errorf("Gamepad %v has no standard layout - cannot use.", ebiten.GamepadName(p))
		}
		// A good gamepad! Add it.
		gamepads[p] = struct{}{}
	}
	for p, stillThere := range allGamepads {
		if stillThere {
			continue
		}
		log.Infof("Gamepad %v removed.", ebiten.GamepadName(p))
		delete(allGamepads, p)
		delete(gamepads, p)
	}
}

func gamepadInit() {
	config := os.Getenv("SDL_GAMECONTROLLERCONFIG")
	if config != "" {
		applied, err := ebiten.UpdateStandardGamepadLayoutMappings(config)
		if err != nil {
			log.Errorf("Could not add SDL_GAMECONTROLLERCONFIG mappings: %v", err)
		} else if applied {
			log.Infof("SDL_GAMECONTROLLERCONFIG applied.")
		} else {
			log.Warningf("SDL_GAMECONTROLLERCONFIG set but not used on this platform.")
		}
	}
}

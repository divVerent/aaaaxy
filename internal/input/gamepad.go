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
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

var (
	gamepad                 = flag.Bool("gamepad", true, "enable gamepad input")
	gamepadAxisOnThreshold  = flag.Float64("gamepad_axis_on_threshold", 0.6, "minimum amount to push the game pad for registering an action; can be zero to accept any movement")
	gamepadAxisOffThreshold = flag.Float64("gamepad_axis_off_threshold", 0.4, "maximum amount to push the game pad for unregistering an action; can be zero to accept any movement")
	gamepadOverride         = flag.String("gamepad_override", "", "entries in SDL_GameControllerDB format to add/override gamepad support; multiple entries are permitted and can be separated by newlines or semicolons; can also be provided via $SDL_GAMECONTROLLERCONFIG environment variable")
	debugGamepadLogging     = flag.Bool("debug_gamepad_logging", false, "log all gamepad states (spammy)")
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
			if ignoredGamepadButtons[b] {
				continue
			}
			if ebiten.IsStandardGamepadButtonPressed(p, b) {
				return Gamepad
			}
		}
		for _, a := range i.padControls.axes {
			if ignoredGamepadAxes[a] {
				continue
			}
			if ebiten.StandardGamepadAxisValue(p, a)*i.padControls.axisDirection >= t {
				return Gamepad
			}
		}
	}
	return NoInput
}

func encodeAxis[K comparable](f float64, m map[K]string, i K) {
	if f < -0.333 {
		m[i] = "-"
	}
	if f > 0.333 {
		m[i] = "+"
	}
}

func standardAxisName(a int) string {
	switch ebiten.StandardGamepadAxis(a) {
	case ebiten.StandardGamepadAxisLeftStickHorizontal:
		return "LX"
	case ebiten.StandardGamepadAxisLeftStickVertical:
		return "LY"
	case ebiten.StandardGamepadAxisRightStickHorizontal:
		return "RX"
	case ebiten.StandardGamepadAxisRightStickVertical:
		return "RY"
	}
	return "?"
}

func standardButtonName(b int) string {
	switch ebiten.StandardGamepadButton(b) {
	case ebiten.StandardGamepadButtonRightBottom:
		return "RB"
	case ebiten.StandardGamepadButtonRightRight:
		return "RR"
	case ebiten.StandardGamepadButtonRightLeft:
		return "RL"
	case ebiten.StandardGamepadButtonRightTop:
		return "RT"
	case ebiten.StandardGamepadButtonFrontTopLeft:
		return "FTL"
	case ebiten.StandardGamepadButtonFrontTopRight:
		return "FTR"
	case ebiten.StandardGamepadButtonFrontBottomLeft:
		return "FBL"
	case ebiten.StandardGamepadButtonFrontBottomRight:
		return "FBR"
	case ebiten.StandardGamepadButtonCenterLeft:
		return "CL"
	case ebiten.StandardGamepadButtonCenterRight:
		return "CR"
	case ebiten.StandardGamepadButtonLeftStick:
		return "LS"
	case ebiten.StandardGamepadButtonRightStick:
		return "RS"
	case ebiten.StandardGamepadButtonLeftTop:
		return "LT"
	case ebiten.StandardGamepadButtonLeftBottom:
		return "LB"
	case ebiten.StandardGamepadButtonLeftLeft:
		return "LL"
	case ebiten.StandardGamepadButtonLeftRight:
		return "LR"
	case ebiten.StandardGamepadButtonCenterCenter:
		return "CC"
	}
	return "?"
}

func gamepadLog() {
	if !*debugGamepadLogging {
		return
	}
	type state struct {
		Name           string
		SDLID          string
		AxisCount      int
		ButtonCount    int
		HasStandard    bool
		Axis           map[int]string
		Button         []int
		StandardAxis   map[string]string
		StandardButton []string
	}
	var states []state
	for _, p := range allGamepadsList {
		s := state{
			Name:         ebiten.GamepadName(p),
			SDLID:        ebiten.GamepadSDLID(p),
			AxisCount:    ebiten.GamepadAxisCount(p),
			ButtonCount:  ebiten.GamepadButtonCount(p),
			HasStandard:  ebiten.IsStandardGamepadLayoutAvailable(p),
			Axis:         map[int]string{},
			StandardAxis: map[string]string{},
		}
		for i := 0; i < ebiten.GamepadAxisCount(p); i++ {
			encodeAxis(ebiten.GamepadAxisValue(p, i), s.Axis, i)
		}
		for i := 0; i < ebiten.GamepadButtonCount(p); i++ {
			if ebiten.IsGamepadButtonPressed(p, ebiten.GamepadButton(i)) {
				s.Button = append(s.Button, i)
			}
		}
		for i := 0; i <= int(ebiten.StandardGamepadAxisMax); i++ {
			encodeAxis(ebiten.StandardGamepadAxisValue(p, ebiten.StandardGamepadAxis(i)), s.StandardAxis, standardAxisName(i))
		}
		for i := 0; i <= int(ebiten.StandardGamepadButtonMax); i++ {
			if ebiten.IsStandardGamepadButtonPressed(p, ebiten.StandardGamepadButton(i)) {
				s.StandardButton = append(s.StandardButton, standardButtonName(i))
			}
		}
		states = append(states, s)
	}
	log.Infof("gamepad states: %+v", states)
}

func gamepadScan() {
	if !*gamepad {
		for p := range gamepads {
			delete(gamepads, p)
		}
		return
	}

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
		log.Infof("gamepad %v (%v) added", ebiten.GamepadName(p), ebiten.GamepadSDLID(p))
		allGamepads[p] = true
		if !ebiten.IsStandardGamepadLayoutAvailable(p) {
			log.Errorf("gamepad %v (%v) has no standard layout - cannot use", ebiten.GamepadName(p), ebiten.GamepadSDLID(p))
			continue
		}
		// A good gamepad! Add it.
		gamepads[p] = struct{}{}
	}
	for p, stillThere := range allGamepads {
		if stillThere {
			continue
		}
		log.Infof("gamepad %v (%v) removed", ebiten.GamepadName(p), ebiten.GamepadSDLID(p))
		delete(allGamepads, p)
		delete(gamepads, p)
	}

	gamepadLog()
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
		return "", fmt.Errorf("open: %w", err)
	}
	defer configHandle.Close()
	configBytes, err := io.ReadAll(configHandle)
	if err != nil {
		return "", fmt.Errorf("read: %w", err)
	}
	return string(configBytes), nil
}

func gamepadInit() {
	// Note: we're also stripping spaces before/after a semicolon
	// as a user might be putting some given they're usual in English,
	// yet they're technically invalid in SDL_GameControllerDB format.
	semiRE := regexp.MustCompile(`\s*;\s*`)

	// Support an included gamecontrollerdb.txt override.
	// Doing this because Ebitengine's lags behind.
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

func gamepadEasterEggKeyState() int {
	state := 0
	for p := range gamepads {
		if ebiten.IsStandardGamepadButtonPressed(p, ebiten.StandardGamepadButtonRightBottom) {
			state |= easterEggA
		}
		if ebiten.IsStandardGamepadButtonPressed(p, ebiten.StandardGamepadButtonRightRight) {
			state |= easterEggB
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

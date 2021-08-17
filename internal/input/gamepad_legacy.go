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

// +build !newgamepad

package input

import (
	"bufio"
	"fmt"
	"github.com/divVerent/aaaaxy/internal/log"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

var (
	gamepadAxisOnThreshold  = flag.Float64("gamepad_axis_on_threshold", 0.6, "Minimum amount to push the game pad for registering an action. Can be zero to accept any movement.")
	gamepadAxisOffThreshold = flag.Float64("gamepad_axis_off_threshold", 0.4, "Maximum amount to push the game pad for unregistering an action. Can be zero to accept any movement.")
)

type (
	padControl struct {
		pad           ebiten.GamepadID
		axisDirection float64 // 0 if button, -1 or 1 if axis.
		axis          int
		button        ebiten.GamepadButton
	}
)

var (
	// gamepadPlatform is the current platform name in the same form as gamecontrollerdb.txt uses.
	gamepadPlatform string
	// gamepadDatabase is the currently loaded gamepad database, keyed by GUID, each entry being a list of all gamepad definitions for that GUID.
	gamepadDatabase map[string][][]string
	// gamepadInvAxisOnThreshold is 1.0 divided by the variable gamepadAxisOnThreshold. Done to save a division for every axis test.
	gamepadInvAxisOnThreshold float64
	// gamepadInvAxisOffThreshold is 1.0 divided by the variable gamepadAxisOffThreshold. Done to save a division for every axis test.
	gamepadInvAxisOffThreshold float64
	// gamepads is the set of currently active gamepads. The boolean value should always be true, except during rescanning, where it's set to false temporarily to detect removed gamepads.
	gamepads = map[ebiten.GamepadID]bool{}
	// defRE is a regular expression to match a gamecontrollerdb.txt assignment. See also the def* constants that match parts of this RE.
	defRE = regexp.MustCompile(`^(?:(\w+):(?:([+-]?)a(\d+)(~?)|b(\d+)|h(\d+)\.(\d+))|platform:(\S+))$`)
)

const (
	defAssignment = 1
	defHalfAxis   = 2
	defAxis       = 3
	defInvertAxis = 4
	defButton     = 5
	defHat        = 6
	defHatBit     = 7
	defPlatform   = 8
)

func (i *impulse) gamepadPressed() bool {
	for _, c := range i.padControls {
		if c.axisDirection == 0 {
			if ebiten.IsGamepadButtonPressed(c.pad, c.button) {
				return true
			}
		} else {
			t := gamepadInvAxisOnThreshold
			if i.Held {
				t = gamepadInvAxisOffThreshold
			}
			if ebiten.GamepadAxis(c.pad, c.axis)*c.axisDirection*t >= 1 {
				return true
			}
		}
	}
	return false
}

func gamepadUpdate() {
	gamepadInvAxisOnThreshold = 1.0 / *gamepadAxisOnThreshold
	gamepadInvAxisOffThreshold = 1.0 / *gamepadAxisOffThreshold
	for pad := range gamepads {
		gamepads[pad] = false
	}
	for _, pad := range ebiten.GamepadIDs() {
		if _, found := gamepads[pad]; !found {
			gamepadAdd(pad)
		}
		gamepads[pad] = true
	}
	for pad, stillThere := range gamepads {
		if !stillThere {
			gamepadRemove(pad)
		}
	}
}

func gamepadRemove(pad ebiten.GamepadID) {
	log.Infof("Removing gamepad %v", ebiten.GamepadName(pad))
	for _, i := range impulses {
		l := len(i.padControls)
		for j := 0; j < l; j++ {
			if i.padControls[j].pad == pad {
				l--
				i.padControls[j] = i.padControls[l]
				i.padControls = i.padControls[:l]
				j--
			}
		}
	}
	delete(gamepads, pad)
}

func gamepadInit() {
	// We need to filter by platform.
	switch runtime.GOOS {
	case "android":
		gamepadPlatform = "Android"
	case "windows":
		gamepadPlatform = "Windows"
	case "darwin":
		if runtime.GOARCH == "arm64" {
			gamepadPlatform = "iOS"
		} else {
			gamepadPlatform = "Mac OS X"
		}
	default: // Include the BSDs too.
		gamepadPlatform = "Linux"
	}

	gamepadLoadDatabase()
	gamepadAddOverrides()
}

func gamepadLoadDatabase() {
	gamepadDatabase = map[string][][]string{}
	f, err := vfs.Load("input", "gamecontrollerdb.txt")
	if err != nil {
		log.Errorf("could not open game controller database: %v", err)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		gamepadAddLine(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Errorf("could not read game controller database: %v", err)
	}
}

func gamepadAddOverrides() {
	for _, line := range strings.Split(os.Getenv("SDL_GAMECONTROLLERCONFIG"), "\n") {
		gamepadAddLine(line)
	}
}

func gamepadAddLine(line string) {
	if line == "" || line[0] == '#' {
		return
	}
	data := strings.Split(line, ",")
	id := data[0]
	gamepadDatabase[id] = append(gamepadDatabase[id], data[1:])
}

func gamepadAdd(pad ebiten.GamepadID) {
	log.Infof("Adding gamepad %v", ebiten.GamepadName(pad))
	gamepads[pad] = true // Don't try again.
	type gamepadError struct {
		name string
		err  error
	}
	var errors []gamepadError
	for _, config := range gamepadDatabase[ebiten.GamepadSDLID(pad)] {
		name, err := gamepadAddWithConfig(pad, config)
		if err != nil {
			errors = append(errors, gamepadError{name, err})
			continue
		}
		return
	}
	for _, ge := range errors {
		log.Errorf("Error adding gamepad %v as %v: %v", ebiten.GamepadName(pad), ge.name, ge.err)
	}
	if len(errors) == 0 {
		log.Errorf("Error adding gamepad %v: no suitable entry found for %v", ebiten.GamepadName(pad), ebiten.GamepadSDLID(pad))
	}
}

func gamepadAddWithConfig(pad ebiten.GamepadID, config []string) (string, error) {
	axes := ebiten.GamepadAxisNum(pad)
	buttons := ebiten.GamepadButtonNum(pad)
	name := config[0]
	controls := map[*impulse][]padControl{
		Jump:   nil,
		Action: nil,
		Left:   nil,
		Right:  nil,
		Up:     nil,
		Down:   nil,
	}
	for _, def := range config[1:] {
		match := defRE.FindStringSubmatch(def)
		if match == nil {
			log.Errorf("Unmatched game pad definition directive: %v", def)
			continue
		}
		if platform := match[defPlatform]; platform != "" {
			if platform != gamepadPlatform {
				return name, fmt.Errorf("wrong platform: got %v, want %v", gamepadPlatform, platform)
			}
			continue
		}
		var addTo, addInverseTo *impulse
		switch match[defAssignment] {
		case "a", "righttrigger", "x":
			addTo = Jump
		case "b", "lefttrigger", "y":
			addTo = Action
		case "back", "guide", "leftshoulder", "rightshoulder", "start":
			addTo = Exit
		case "dpleft":
			addTo = Left
		case "dpright":
			addTo = Right
		case "dpdown":
			addTo = Down
		case "dpup":
			addTo = Up
		case "leftx", "rightx":
			addTo, addInverseTo = Right, Left
		case "lefty", "righty":
			addTo, addInverseTo = Down, Up
		case "leftstick", "misc1", "paddle1", "paddle2", "paddle3", "paddle4", "rightstick", "touchpad":
		// Ignore.
		// Where to put fullscreen?
		default:
			log.Warningf("Unknown assignment in game pad definition directive: %v", def)
			continue
		}
		switch match[defHalfAxis] {
		case "+":
			addInverseTo = nil
		case "-":
			addTo = nil
		}
		if match[defInvertAxis] != "" {
			addTo, addInverseTo = addInverseTo, addTo
		}
		if a := match[defAxis]; a != "" {
			ax, err := strconv.Atoi(a)
			if err != nil {
				log.Warningf("Could not parse axis %v: %v", a, err)
				continue
			}
			if ax < 0 || ax >= axes {
				log.Warningf("Invalid axis : got %v, want 0 <= ax < %d", ax, axes)
				continue
			}
			if addTo != nil {
				controls[addTo] = append(controls[addTo], padControl{
					pad:           pad,
					axisDirection: 1.0,
					axis:          ax,
				})
			}
			if addInverseTo != nil {
				controls[addInverseTo] = append(controls[addInverseTo], padControl{
					pad:           pad,
					axisDirection: -1.0,
					axis:          ax,
				})
			}
		}
		if b := match[defButton]; b != "" {
			bt, err := strconv.Atoi(b)
			if err != nil {
				log.Warningf("Could not parse button %v: %v", b, err)
				continue
			}
			if bt < 0 || bt >= buttons {
				log.Warningf("Invalid button: got %v, want 0 <= bt < %d", bt, buttons)
				continue
			}
			if addTo != nil {
				controls[addTo] = append(controls[addTo], padControl{
					pad:    pad,
					button: ebiten.GamepadButton(bt),
				})
			}
		}
		if h, b := match[defHat], match[defHatBit]; h != "" {
			ht, err := strconv.Atoi(h)
			if err != nil {
				log.Warningf("Could not parse hat %v: %v", b, err)
				continue
			}
			bt, err := strconv.Atoi(b)
			if err != nil {
				log.Warningf("Could not parse hat bit %v: %v", b, err)
				continue
			}
			// Note: ebiten currently doesn't support "hats" in GLFW properly.
			// However, hat 0 always occupies the last four buttons.
			if ht != 0 {
				log.Warningf("Sorry, non-zero hat numbers are not supported right now")
				continue
			}
			vbt := buttons - map[int]int{
				1: 4, 2: 3, 4: 2, 8: 1,
			}[bt]
			if vbt < 0 || vbt >= buttons {
				log.Warningf("Invalid hat button: got %v, want 0 <= vbt < %d", vbt, buttons)
				continue
			}
			if addTo != nil {
				controls[addTo] = append(controls[addTo], padControl{
					pad:    pad,
					button: ebiten.GamepadButton(vbt),
				})
			}
		}
	}
	for i, c := range controls {
		if c == nil {
			return name, fmt.Errorf("missing assignment for %v", i.Name)
		}
	}
	// If we get here, the pad is configured fully.
	for i, c := range controls {
		i.padControls = append(i.padControls, c...)
	}
	log.Infof("Gamepad configured and found: %v (configured: %v)", ebiten.GamepadName(pad), name)
	return name, nil
}

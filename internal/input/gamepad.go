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
	"log"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/flag"
)

var (
	debugGamepadString   = flag.String("debug_gamepad_string", "03000000d62000000228000001010000,PowerA Pro Ex,a:b0,b:b1,back:b6,dpdown:h0.4,dpleft:h0.8,dpright:h0.2,dpup:h0.1,guide:b8,leftshoulder:b4,leftstick:b9,lefttrigger:a2,leftx:a0,lefty:a1,rightshoulder:b5,rightstick:b10,righttrigger:a5,rightx:a3,righty:a4,start:b7,x:b2,y:b3,platform:Linux", "SDL gamepad definition")
	gamepadAxisThreshold = flag.Float64("gamepad_axis_threshold", 0.5, "Minimum amount to push the game pad for registering an action. Can be zero to accept any movement.")
)

type (
	padControl struct {
		pad              ebiten.GamepadID
		invAxisThreshold float64 // 0 if button.
		axis             int
		button           ebiten.GamepadButton
	}
)

func (i *impulse) gamepadPressed() bool {
	for _, c := range i.padControls {
		if c.invAxisThreshold == 0 {
			if ebiten.IsGamepadButtonPressed(c.pad, c.button) {
				return true
			}
		} else {
			if ebiten.GamepadAxis(c.pad, c.axis)*c.invAxisThreshold >= 1 {
				return true
			}
		}
	}
	return false
}

var (
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

var (
	gamepads = map[ebiten.GamepadID]bool{}
)

func gamepadUpdate() {
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
	log.Printf("Removing gamepad %v", ebiten.GamepadName(pad))
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

var gamepadPlatform string

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
}

func gamepadAdd(pad ebiten.GamepadID) error {
	log.Printf("Adding gamepad %v", ebiten.GamepadName(pad))
	gamepads[pad] = true // Don't try again.
	if *debugGamepadString == "" {
		return fmt.Errorf("no config found")
	}
	// TODO: support the entire database and environment override.
	return gamepadAddWithConfig(pad, *debugGamepadString)
}

func gamepadAddWithConfig(pad ebiten.GamepadID, config string) error {
	data := strings.Split(config, ",")
	id := data[0]
	name := data[1]
	padID := ebiten.GamepadSDLID(pad)
	if id != padID {
		return fmt.Errorf("different ID: got %v, want %v", padID, id)
	}
	controls := map[*impulse][]padControl{
		Jump:   nil,
		Action: nil,
		Left:   nil,
		Right:  nil,
		Up:     nil,
		Down:   nil,
	}
	for _, def := range data[2:] {
		match := defRE.FindStringSubmatch(def)
		if match == nil {
			log.Printf("Unmatched game pad definition directive: %v", def)
			continue
		}
		if platform := match[defPlatform]; platform != "" {
			if platform != gamepadPlatform {
				return fmt.Errorf("wrong platform: got %v, want %v", gamepadPlatform, platform)
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
			log.Printf("Unknown assignment in game pad definition directive: %v", def)
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
				log.Printf("Could not parse axis %v: %v", a, err)
				continue
			}
			if addTo != nil {
				controls[addTo] = append(controls[addTo], padControl{
					pad:              pad,
					invAxisThreshold: 1.0 / *gamepadAxisThreshold,
					axis:             ax,
				})
			}
			if addInverseTo != nil {
				controls[addInverseTo] = append(controls[addInverseTo], padControl{
					pad:              pad,
					invAxisThreshold: -1.0 / *gamepadAxisThreshold,
					axis:             ax,
				})
			}
		}
		if b := match[defButton]; b != "" {
			bt, err := strconv.Atoi(b)
			if err != nil {
				log.Printf("Could not parse button %v: %v", b, err)
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
				log.Printf("Could not parse hat %v: %v", b, err)
				continue
			}
			bt, err := strconv.Atoi(b)
			if err != nil {
				log.Printf("Could not parse hat bit %v: %v", b, err)
				continue
			}
			// Note: ebiten currently doesn't support "hats" in GLFW properly.
			// However, hat 0 always occupies the last four buttons.
			if ht != 0 {
				log.Printf("Sorry, non-zero hat numbers are not supported right now")
				continue
			}
			vbt := ebiten.GamepadButtonNum(pad) - map[int]int{
				1: 3, 2: 2, 4: 1, 8: 0,
			}[bt]
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
			return fmt.Errorf("missing assignment for %v", i.Name)
		}
	}
	// If we get here, the pad is configured fully.
	for i, c := range controls {
		i.padControls = append(i.padControls, c...)
	}
	log.Printf("Gamepad configured and found: %v (configured: %v)", ebiten.GamepadName(pad), name)
	return nil
}

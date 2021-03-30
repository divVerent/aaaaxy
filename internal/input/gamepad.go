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
	debugGamepadOverridePlatform = flag.String("debug_gamepad_override_platform", "", "the platform name to look for in the game pad definition")
	debugGamepadString           = flag.String("debug_gamepad_string", "03000000d62000006dca000000000000,PowerA Pro Ex,a:b1,b:b2,back:b8,dpdown:h0.4,dpleft:h0.8,dpright:h0.2,dpup:h0.1,guide:b12,leftshoulder:b4,leftstick:b10,lefttrigger:b6,leftx:a0,lefty:a1,rightshoulder:b5,rightstick:b11,righttrigger:b7,rightx:a2,righty:a3,start:b9,x:b0,y:b3,platform:Windows,", "SDL gamepad definition")
	gamepadAxisThreshold         = flag.Float64("gamepad_axis_threshold", 0.1, "Minimum amount to push the game pad for registering an action. Can be zero to accept any movement.")
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
	defRE = regexp.MustCompile(`^(?:(\w+):(?:a([+-]?)(\d+)(~?)|b(\d+)|h(\d+)\.(\d+))|platform:(\S+))$`)
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

func gamepadInit() error {
	// First clear all existing gamepad mappings.
	for _, i := range impulses {
		i.padControls = nil
	}
	// We need to filter by platform.
	wantPlatform := *debugGamepadOverridePlatform
	if wantPlatform == "" {
		switch runtime.GOOS {
		case "android":
			wantPlatform = "Android"
		case "windows":
			wantPlatform = "Windows"
		case "darwin":
			if runtime.GOARCH == "arm64" {
				wantPlatform = "iOS"
			} else {
				wantPlatform = "Mac OS X"
			}
		default: // Include the BSDs too.
			wantPlatform = "Linux"
		}
	}
	// Now configure the gamepad.
	if *debugGamepadString == "" {
		return nil
	}
	// TODO: support the entire database and environment override.
	data := strings.Split(*debugGamepadString, ",")
	id := data[0]
	name := data[1]
	havePad := false
	ids := []string{}
nextPad:
	for _, pad := range ebiten.GamepadIDs() {
		padID := ebiten.GamepadSDLID(pad)
		ids = append(ids, padID)
		if id != ebiten.GamepadSDLID(pad) {
			continue
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
			if match[defPlatform] != "" {
				if match[defPlatform] != wantPlatform {
					continue nextPad
				}
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
				addTo, addInverseTo = Up, Down
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
				}
				if addTo != nil {
					controls[addTo] = append(controls[addTo], padControl{
						pad:              pad,
						invAxisThreshold: 1.0 / *gamepadAxisThreshold,
						axis:             ax,
					})
				}
				if addInverseTo != nil {
					controls[addTo] = append(controls[addTo], padControl{
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
				}
				if addTo != nil {
					controls[addTo] = append(controls[addTo], padControl{
						pad:    pad,
						button: ebiten.GamepadButton(bt),
					})
				}
			}
		}
		for i, c := range controls {
			if c == nil {
				log.Printf("Skipping gamepad %v (missing assignment for %v)", name, i.Name)
				continue nextPad
			}
		}
		for i, c := range controls {
			i.padControls = append(i.padControls, c...)
		}
		log.Printf("Gamepad configured and found: %v (configured: %v)", ebiten.GamepadName(pad), name)
		havePad = true
	}
	if !havePad {
		return fmt.Errorf("Gamepad configured but not present: got %v, want %v", ids, id)
	}
	return nil
}

// Copyright 2022 Google LLC
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

//go:build android
// +build android

package input

import (
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// For some Android gamepads, AXIS_RX/AXIS_RY and AXIS_Z/AXIS_RZ are confused.
	// RX means "rotate X", not "right X"; a right stick should according to docs be mapped to Z/RZ.
	// Linux OTOH defines ABS_RX/ABS_RY as right X and Y axis...
	//
	// HOWEVER, due to the SDL-inherited way of how Ebitengine maps axes,
	// this should not matter anymore. So for now keeping a separate file
	// just in case I need to add ignores back, and enabling all axes.
	ignoredGamepadButtons = map[ebiten.StandardGamepadButton]bool(nil)
	ignoredGamepadAxes    = map[ebiten.StandardGamepadAxis]bool(nil)
)

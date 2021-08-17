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
)

func (i *impulse) register() *impulse {
	impulses = append(impulses, i)
	return i
}

func (i *impulse) update() {
	held := i.keyboardPressed() || i.gamepadPressed()
	i.JustHit = held && !i.Held
	i.Held = held
}

func Init() error {
	gamepadInit()
	return nil
}

func Update() {
	gamepadUpdate()
	for _, i := range impulses {
		i.update()
	}
}

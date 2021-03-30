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
	Held    bool
	JustHit bool

	keys []ebiten.Key
}

var (
	Left       = (&impulse{keys: leftKeys}).register()
	Right      = (&impulse{keys: rightKeys}).register()
	Up         = (&impulse{keys: upKeys}).register()
	Down       = (&impulse{keys: downKeys}).register()
	Jump       = (&impulse{keys: jumpKeys}).register()
	Action     = (&impulse{keys: actionKeys}).register()
	Exit       = (&impulse{keys: exitKeys}).register()
	Fullscreen = (&impulse{keys: fullscreenKeys}).register()

	impulses = []*impulse{}
)

func (i *impulse) register() *impulse {
	impulses = append(impulses, i)
	return i
}

func (i *impulse) update() {
	held := i.keyboardPressed()
	i.JustHit = held && !i.Held
	i.Held = held
}

func Init() {}

func Update() {
	for _, i := range impulses {
		i.update()
	}
}

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

type Impulse struct {
	Held    bool
	JustHit bool

	keys []ebiten.Key
}

var (
	Left   = (&Impulse{keys: []ebiten.Key{ebiten.KeyLeft, ebiten.KeyA}}).register()
	Right  = (&Impulse{keys: []ebiten.Key{ebiten.KeyRight, ebiten.KeyD}}).register()
	Up     = (&Impulse{keys: []ebiten.Key{ebiten.KeyUp, ebiten.KeyW}}).register()
	Down   = (&Impulse{keys: []ebiten.Key{ebiten.KeyDown, ebiten.KeyS}}).register()
	Jump   = (&Impulse{keys: []ebiten.Key{ebiten.KeyControl, ebiten.KeySpace, ebiten.KeyX}}).register()
	Action = (&Impulse{keys: []ebiten.Key{ebiten.KeyAlt, ebiten.KeyE, ebiten.KeyZ, ebiten.KeyEnter}}).register()
	Exit   = (&Impulse{keys: []ebiten.Key{ebiten.KeyEscape}}).register()

	impulses = []*Impulse{}
)

func (i *Impulse) register() *Impulse {
	impulses = append(impulses, i)
	return i
}

func (i *Impulse) update() {
	held := false
	for _, k := range i.keys {
		if ebiten.IsKeyPressed(k) {
			held = true
			break
		}
	}
	i.JustHit = held && !i.Held
	i.Held = held
}

func Update() {
	for _, i := range impulses {
		i.update()
	}
}

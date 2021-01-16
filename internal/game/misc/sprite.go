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

package misc

import (
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/image"
	m "github.com/divVerent/aaaaaa/internal/math"
)

// Sprite is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
type Sprite struct{}

func (s *Sprite) Spawn(w *engine.World, sp *engine.Spawnable, e *engine.Entity) error {
	var err error
	directory := sp.Properties["image_dir"]
	if directory == "" {
		directory = "sprites"
	}
	e.Image, err = image.Load(directory, sp.Properties["image"])
	if err != nil {
		return err
	}
	e.ResizeImage = true
	e.Solid = sp.Properties["solid"] != "false"
	e.Opaque = sp.Properties["opaque"] != "false"
	if sp.Properties["alpha"] != "" {
		e.Alpha, err = strconv.ParseFloat(sp.Properties["alpha"], 64)
		if err != nil {
			return fmt.Errorf("could not decode alpha %q: %v", sp.Properties["alpha"], err)
		}
	}
	if sp.Properties["z_index"] != "" {
		e.ZIndex, err = strconv.Atoi(sp.Properties["z_index"])
		if err != nil {
			return fmt.Errorf("could not decode z index %q: %v", sp.Properties["z_index"], err)
		}
	}
	return nil
}

func (s *Sprite) Despawn() {}

func (s *Sprite) Update() {}

func (s *Sprite) Touch(other *engine.Entity) {}

func (s *Sprite) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&Sprite{})
}

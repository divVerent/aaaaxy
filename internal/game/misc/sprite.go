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
	go_image "image"
	"image/color"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/animation"
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/font"
	"github.com/divVerent/aaaaaa/internal/game/constants"
	"github.com/divVerent/aaaaaa/internal/image"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
)

// Sprite is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
type Sprite struct {
	SpriteBase
}

func (s *Sprite) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	directory := sp.Properties["image_dir"]
	if directory == "" {
		directory = "sprites"
	}
	var err error
	e.Image, err = image.Load(directory, sp.Properties["image"])
	if err != nil {
		return err
	}
	e.ResizeImage = true
	subX, subY := 0, 0
	subW, subH := e.Image.Size()
	regionString := sp.Properties["image_region"]
	if regionString != "" {
		if _, err := fmt.Sscanf(regionString, "%d %d %d %d", &subX, &subY, &subW, &subH); err != nil {
			return fmt.Errorf("could not decode image region %q: %v", regionString, err)
		}
		e.Image = e.Image.SubImage(go_image.Rectangle{
			Min: go_image.Point{
				X: subX,
				Y: subY,
			},
			Max: go_image.Point{
				X: subX + subW,
				Y: subY + subH,
			},
		}).(*ebiten.Image)
	}
	return s.SpriteBase.Spawn(w, sp, e)
}

func init() {
	engine.RegisterEntityType(&Sprite{})
}

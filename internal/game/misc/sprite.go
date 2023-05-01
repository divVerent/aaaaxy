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
	go_image "image"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/image"
	"github.com/divVerent/aaaaxy/internal/level"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

var (
	checkSprites = flag.Bool("check_sprites", false, "check that all sprites exist at startup (NOT recommended)")
)

// Sprite is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
type Sprite struct {
	SpriteBase
}

var _ engine.Precacher = &Sprite{}

func (s *Sprite) Precache(sp *level.Spawnable) error {
	if !*checkSprites {
		return nil
	}
	var parseErr error
	directory := propmap.StringOr(sp.Properties, "image_dir", "sprites")
	imgSrc := propmap.ValueP(sp.Properties, "image", "", &parseErr)
	_, err := image.Load(directory, imgSrc)
	if err != nil {
		return err
	}
	imgSrcByOrientation, err := level.ParseImageSrcByOrientation(imgSrc, sp.Properties)
	if err != nil {
		return err
	}
	for _, thisSrc := range imgSrcByOrientation {
		if thisSrc == "" {
			continue
		}
		_, err := image.Load(directory, thisSrc)
		if err != nil {
			return err
		}
	}
	return parseErr
}

func (s *Sprite) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	var parseErr error
	directory := propmap.StringOr(sp.Properties, "image_dir", "sprites")
	imgSrc := propmap.ValueP(sp.Properties, "image", "", &parseErr)
	imgSrcByOrientation, err := level.ParseImageSrcByOrientation(imgSrc, sp.Properties)
	if err != nil {
		return err
	}
	imgSrc, e.Orientation = level.ResolveImage(e.Transform, e.Orientation, imgSrc, imgSrcByOrientation)
	e.Image, err = image.Load(directory, imgSrc)
	if err != nil {
		return err
	}
	e.RenderOffset = propmap.ValueOrP(sp.Properties, "render_offset", m.Delta{}, &parseErr)
	if e.RenderOffset.IsZero() {
		e.ResizeImage = true
	}
	region := propmap.ValueOrP(sp.Properties, "image_region", m.Rect{}, &parseErr)
	if !region.Size.IsZero() {
		e.Image = e.Image.SubImage(go_image.Rectangle{
			Min: go_image.Point{
				X: region.Origin.X,
				Y: region.Origin.Y,
			},
			Max: go_image.Point{
				X: region.Origin.X + region.Size.DX,
				Y: region.Origin.Y + region.Size.DY,
			},
		}).(*ebiten.Image)
	}
	e.BorderPixels = propmap.ValueOrP(sp.Properties, "border_pixels", 0, &parseErr)
	err = s.SpriteBase.Spawn(w, sp, e)
	if err != nil {
		return err
	}
	return parseErr
}

func init() {
	engine.RegisterEntityType(&Sprite{})
}

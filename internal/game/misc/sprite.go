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

func (s *Sprite) Precache(id level.EntityID, sp *level.SpawnableProps) error {
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
	subX, subY := 0, 0
	subW, subH := e.Image.Size()
	// This pattern is very specific - just have it here and not clutter the generic parser with it.
	regionString := propmap.StringOr(sp.Properties, "image_region", "")
	if regionString != "" {
		if _, err := fmt.Sscanf(regionString, "%d %d %d %d", &subX, &subY, &subW, &subH); err != nil {
			return fmt.Errorf("could not decode image region %q: %w", regionString, err)
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

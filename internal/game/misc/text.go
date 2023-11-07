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
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/fun"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/playerstate"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

var (
	precacheText            = flag.Bool("precache_text", true, "preload all text objects at startup (VERY recommended)")
	memImagesForStaticText  = flag.Bool("mem_images_for_static_text", true, "use in-memory images for static text objects (faster startup)")
	memImagesForDynamicText = flag.Bool("mem_images_for_dynamic_text", false, "use in-memory images for dynamic text objects (seems to update slower in-game)")
)

// Text is a simple entity type that renders text.
type Text struct {
	SpriteBase

	World  *engine.World
	Entity *engine.Entity

	Key     textCacheKey
	MyImage bool
}

var _ engine.Precacher = &Text{}

type textCacheKey struct {
	font   string
	fg, bg color.NRGBA
	text   string
}

var textCache = map[textCacheKey]*ebiten.Image{}

func cacheKey(sp *level.SpawnableProps) textCacheKey {
	return textCacheKey{
		font: propmap.ValueP(sp.Properties, "text_font", "", nil),
		fg:   propmap.ValueP(sp.Properties, "text_fg", color.NRGBA{}, nil),
		bg:   propmap.ValueP(sp.Properties, "text_bg", color.NRGBA{}, nil),
		text: propmap.ValueP(sp.Properties, "text", "", nil),
	}
}

func (key textCacheKey) load(ps *playerstate.PlayerState) (*ebiten.Image, error) {
	fnt := font.ByName[key.font]
	if fnt.Face == nil {
		return nil, fmt.Errorf("could not find font %q", key.font)
	}
	txt, err := fun.TryFormatText(ps, key.text)
	if err != nil {
		if ps == nil {
			// On template execution failure, we do not fail precaching.
			// However later rendering may fail too then.
			return nil, nil
		}
		return nil, err
	}
	bounds := fnt.BoundString(txt)
	useMemImages := *memImagesForStaticText
	if ps != nil {
		useMemImages = *memImagesForDynamicText
	}
	if useMemImages {
		img := image.NewRGBA( // image.RGBA is Ebitengine's fast path.
			image.Rectangle{
				Min: image.Point{
					X: 0,
					Y: 0,
				},
				Max: image.Point{
					X: bounds.Size.DX,
					Y: bounds.Size.DY,
				},
			})
		fnt.Draw(img, txt, bounds.Origin.Mul(-1), font.AsBounds, key.fg, key.bg)
		img2 := ebiten.NewImageFromImage(img)
		return img2, nil
	} else {
		img := ebiten.NewImage(bounds.Size.DX, bounds.Size.DY)
		fnt.Draw(img, txt, bounds.Origin.Mul(-1), font.AsBounds, key.fg, key.bg)
		return img, nil
	}
}

func ClearPrecache() {
	/*
		// Can't dispose them right away, as the images are still referenced.
		// Instead, the GC will eventually do this.
		for _, img := range textCache {
			img.Dispose()
		}
	*/
	textCache = map[textCacheKey]*ebiten.Image{}
}

func (t *Text) Precache(sp *level.Spawnable) error {
	if !*precacheText {
		return nil
	}
	log.Debugf("precaching text for entity %v", sp.ID)
	key := cacheKey(&sp.SpawnableProps)
	if textCache[key] != nil {
		return nil
	}
	img, err := key.load(nil)
	if err != nil {
		return fmt.Errorf("could not precache text image for entity %v: %w", sp.ID, err)
	}
	textCache[key] = img
	if img == nil {
		// Nothing precached.
		return nil
	}
	rx, ry := sp.RectInTile.Size.DX, sp.RectInTile.Size.DY
	if sp.Orientation.Right.DX == 0 {
		rx, ry = ry, rx
	}
	sz := img.Bounds().Size()
	if sz.X > rx+2 || (locale.Active.AuditHeight() && sz.Y > ry+3) {
		// Tolerate 3 extra pixels for diacritics or borders.
		locale.Errorf("text too big: entity %v has size %v but text needs %v: %v",
			sp.ID, m.Delta{DX: rx, DY: ry}, sz, key.text)
	}
	return nil
}

func (t *Text) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	propmap.SetDefault(sp.Properties, "no_flip", "x")

	t.World = w
	t.Entity = e

	t.Key = cacheKey(sp)
	err := t.updateText()
	if err != nil {
		return err
	}

	e.ResizeImage = false

	return t.SpriteBase.Spawn(w, sp, e)
}

func (t *Text) updateText() error {
	if t.MyImage {
		t.Entity.Image.Dispose()
	}
	t.Entity.Image = nil
	if *precacheText {
		var found bool
		t.Entity.Image, found = textCache[t.Key]
		t.MyImage = false
		if !found {
			return fmt.Errorf("could not find precached text image for entity %v", t.Key)
		}
	}
	if t.Entity.Image == nil {
		var err error
		t.Entity.Image, err = t.Key.load(&t.World.PlayerState)
		t.MyImage = true
		if err != nil {
			return fmt.Errorf("could not render text image for entity %v: %w", t.Key, err)
		}
	}
	sz := t.Entity.Image.Bounds().Size()
	if t.Entity.Orientation.Right.DX == 0 {
		sz.X, sz.Y = sz.Y, sz.X
	}
	centerOffset := t.Entity.Rect.Size.Sub(m.Delta{DX: sz.X, DY: sz.Y}).Div(2)
	t.Entity.RenderOffset = centerOffset
	return nil
}

func (t *Text) Despawn() {
	if t.MyImage {
		t.Entity.Image.Dispose()
	}
	t.Entity.Image = nil
}

func init() {
	engine.RegisterEntityType(&Text{})
}

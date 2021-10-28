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
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/fun"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/playerstate"
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
	fg, bg string
	text   string
}

var textCache = map[textCacheKey]*ebiten.Image{}

func cacheKey(sp *level.Spawnable) textCacheKey {
	return textCacheKey{
		font: sp.Properties["text_font"],
		fg:   sp.Properties["text_fg"],
		bg:   sp.Properties["text_bg"],
		text: sp.Properties["text"],
	}
}

func (key textCacheKey) load(ps *playerstate.PlayerState) (*ebiten.Image, error) {
	fnt := font.ByName[key.font]
	if fnt.Face == nil {
		return nil, fmt.Errorf("could not find font %q", key.font)
	}
	var fg, bg color.NRGBA
	if _, err := fmt.Sscanf(key.fg, "#%02x%02x%02x%02x", &fg.A, &fg.R, &fg.G, &fg.B); err != nil {
		return nil, fmt.Errorf("could not decode color %q: %v", key.fg, err)
	}
	if _, err := fmt.Sscanf(key.bg, "#%02x%02x%02x%02x", &bg.A, &bg.R, &bg.G, &bg.B); err != nil {
		return nil, fmt.Errorf("could not decode color %q: %v", key.bg, err)
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
	txt = strings.ReplaceAll(txt, "  ", "\n")
	bounds := fnt.BoundString(txt)
	useMemImages := *memImagesForStaticText
	if ps != nil {
		useMemImages = *memImagesForDynamicText
	}
	if useMemImages {
		img := image.NewRGBA( // image.RGBA is ebiten's fast path.
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
		fnt.Draw(img, txt, bounds.Origin.Mul(-1), false, fg, bg)
		img2 := ebiten.NewImageFromImage(img)
		return img2, nil
	} else {
		img := ebiten.NewImage(bounds.Size.DX, bounds.Size.DY)
		fnt.Draw(img, txt, bounds.Origin.Mul(-1), false, fg, bg)
		return img, nil
	}
}

func (t *Text) Precache(sp *level.Spawnable) error {
	if !*precacheText {
		return nil
	}
	log.Debugf("precaching text for entity %v", sp.ID)
	key := cacheKey(sp)
	if textCache[key] != nil {
		return nil
	}
	img, err := key.load(nil)
	if err != nil {
		return fmt.Errorf("could not precache text image for entity %v: %v", sp, err)
	}
	textCache[key] = img
	return nil
}

func (t *Text) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	if sp.Properties["no_flip"] == "" {
		sp.Properties["no_flip"] = "x"
	}

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
			return fmt.Errorf("could not render text image for entity %v: %v", t.Key, err)
		}
	}
	dx, dy := t.Entity.Image.Size()
	if t.Entity.Orientation.Right.DX == 0 {
		dx, dy = dy, dx
	}
	centerOffset := t.Entity.Rect.Size.Sub(m.Delta{DX: dx, DY: dy}).Div(2)
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

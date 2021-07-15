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
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/level"
	m "github.com/divVerent/aaaaxy/internal/math"
)

var (
	precacheText = flag.Bool("precache_text", true, "preload all text objects at startup (VERY recommended)")
)

// Text is a simple entity type that renders text.
type Text struct {
	SpriteBase
	Entity *engine.Entity
}

var _ engine.Precacher = &Text{}

type textCacheKey struct {
	font   string
	fg, bg string
	text   string
}

var textCache = map[textCacheKey]*ebiten.Image{}

func cacheKey(s *level.Spawnable) textCacheKey {
	return textCacheKey{
		font: s.Properties["text_font"],
		fg:   s.Properties["text_fg"],
		bg:   s.Properties["text_bg"],
		text: s.Properties["text"],
	}
}

func (key textCacheKey) load() (*ebiten.Image, error) {
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
	txt := strings.ReplaceAll(key.text, "  ", "\n")
	bounds := fnt.BoundString(txt)
	img := ebiten.NewImage(bounds.Size.DX, bounds.Size.DY)
	fnt.Draw(img, txt, bounds.Origin.Mul(-1), false, fg, bg)
	// NewImageFromImage forces the text to actually be written to the atlas.
	// Sadly we can only do that once actually initialized, as it reads from an *ebiten.Image.
	// If only we could render the text to an image.Image... TBD.
	// TODO(divVerent): Fix that, and move level precaching into engine.Precache.
	img2 := ebiten.NewImageFromImage(img)
	img.Dispose()
	return img2, nil
}

func (t *Text) Precache(s *level.Spawnable) error {
	if !*precacheText {
		return nil
	}
	log.Printf("precaching text for entity %v", s.ID)
	key := cacheKey(s)
	if textCache[key] != nil {
		return nil
	}
	img, err := key.load()
	if err != nil {
		return fmt.Errorf("could not precache text image for entity %v: %v", s, err)
	}
	textCache[key] = img
	return nil
}

func (t *Text) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	if s.Properties["no_flip"] == "" {
		s.Properties["no_flip"] = "x"
	}

	key := cacheKey(s)
	if *precacheText {
		e.Image = textCache[key]
		if e.Image == nil {
			return fmt.Errorf("could not find precached text image for entity %v", s)
		}
	} else {
		var err error
		e.Image, err = key.load()
		if err != nil {
			return fmt.Errorf("could not render text image for entity %v: %v", s, err)
		}
	}

	t.Entity = e
	e.ResizeImage = false
	dx, dy := e.Image.Size()
	if e.Orientation.Right.DX == 0 {
		dx, dy = dy, dx
	}
	centerOffset := e.Rect.Size.Sub(m.Delta{DX: dx, DY: dy}).Div(2)
	e.RenderOffset = e.RenderOffset.Add(centerOffset)
	return t.SpriteBase.Spawn(w, s, e)
}

func (t *Text) Despawn() {
	if *precacheText {
		return
	}
	t.Entity.Image.Dispose()
	t.Entity.Image = nil
}

func init() {
	engine.RegisterEntityType(&Text{})
}

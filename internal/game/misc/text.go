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
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/level"
)

// Text is a simple entity type that renders text.
type Text struct {
	SpriteBase
	Entity *engine.Entity
}

func (t *Text) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	if s.Properties["no_flip"] == "" {
		s.Properties["no_flip"] = "x"
	}
	t.Entity = e
	fntString := s.Properties["text_font"]
	fnt := font.ByName[fntString]
	if fnt.Face == nil {
		return fmt.Errorf("could not find font %q", fntString)
	}
	var fg, bg color.NRGBA
	fgString := s.Properties["text_fg"]
	if _, err := fmt.Sscanf(fgString, "#%02x%02x%02x%02x", &fg.A, &fg.R, &fg.G, &fg.B); err != nil {
		return fmt.Errorf("could not decode color %q: %v", fgString, err)
	}
	bgString := s.Properties["text_bg"]
	if _, err := fmt.Sscanf(bgString, "#%02x%02x%02x%02x", &bg.A, &bg.R, &bg.G, &bg.B); err != nil {
		return fmt.Errorf("could not decode color %q: %v", bgString, err)
	}
	txt := strings.ReplaceAll(s.Properties["text"], "  ", "\n")
	bounds := fnt.BoundString(txt)
	e.Image = ebiten.NewImage(bounds.Size.DX, bounds.Size.DY)
	fnt.Draw(e.Image, txt, bounds.Origin.Mul(-1), false, fg, bg)
	e.ResizeImage = false
	if e.Orientation.Right.DX == 0 {
		bounds.Size.DX, bounds.Size.DY = bounds.Size.DY, bounds.Size.DX
	}
	centerOffset := e.Rect.Size.Sub(bounds.Size).Div(2)
	e.RenderOffset = e.RenderOffset.Add(centerOffset)
	return t.SpriteBase.Spawn(w, s, e)
}

func (t *Text) Despawn() {
	t.Entity.Image.Dispose()
	t.Entity.Image = nil
}

func init() {
	engine.RegisterEntityType(&Text{})
}

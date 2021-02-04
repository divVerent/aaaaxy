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
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/animation"
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/font"
	"github.com/divVerent/aaaaaa/internal/image"
	m "github.com/divVerent/aaaaaa/internal/math"
)

// Sprite is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
type Sprite struct {
	Entity  *engine.Entity
	MyImage bool
	Anim    animation.State
}

func (s *Sprite) Spawn(w *engine.World, sp *engine.Spawnable, e *engine.Entity) error {
	s.Entity = e
	var err error
	directory := sp.Properties["image_dir"]
	if directory == "" {
		directory = "sprites"
	}
	if sp.Properties["image"] == "" && sp.Properties["text"] != "" && sp.Properties["animation"] != "" {
		fntString := sp.Properties["text_font"]
		fnt := font.ByName[fntString]
		if fnt.Face == nil {
			return fmt.Errorf("could not find font %q", fntString)
		}
		var fg, bg color.NRGBA
		fgString := sp.Properties["text_fg"]
		if _, err := fmt.Sscanf(fgString, "#%02x%02x%02x%02x", &fg.A, &fg.R, &fg.G, &fg.B); err != nil {
			return fmt.Errorf("could not decode color %q: %v", fgString, err)
		}
		bgString := sp.Properties["text_bg"]
		if _, err := fmt.Sscanf(bgString, "#%02x%02x%02x%02x", &bg.A, &bg.R, &bg.G, &bg.B); err != nil {
			return fmt.Errorf("could not decode color %q: %v", bgString, err)
		}
		txt := strings.ReplaceAll(sp.Properties["text"], "  ", "\n")
		bounds := fnt.BoundString(txt)
		e.Image = ebiten.NewImage(bounds.Size.DX, bounds.Size.DY)
		fnt.Draw(e.Image, txt, bounds.Origin.Mul(-1), false, fg, bg)
		e.ResizeImage = false
		centerOffset := e.Rect.Size.Sub(bounds.Size).Div(2)
		e.RenderOffset = e.RenderOffset.Add(centerOffset)
		s.MyImage = true
	} else if sp.Properties["text"] == "" && sp.Properties["image"] != "" && sp.Properties["animation"] != "" {
		e.Image, err = image.Load(directory, sp.Properties["image"])
		if err != nil {
			return err
		}
		e.ResizeImage = true
	} else if sp.Properties["animation"] != "" {
		prefix := sp.Properties["animation"]
		group := &animation.Group{
			NextAnim: "default",
		}
		framesString := sp.Properties["animation_frames"]
		if _, err := fmt.Sscanf(framesString, "%d", &group.Frames); err != nil {
			return fmt.Errorf("could not decode animation_frames %q: %v", framesString, err)
		}
		frameIntervalString := sp.Properties["animation_frame_interval"]
		if _, err := fmt.Sscanf(frameIntervalString, "%d", &group.FrameInterval); err != nil {
			return fmt.Errorf("could not decode animation_frame_interval %q: %v", frameIntervalString, err)
		}
		repeatIntervalString := sp.Properties["animation_repeat_interval"]
		if _, err := fmt.Sscanf(repeatIntervalString, "%d", &group.NextInterval); err != nil {
			return fmt.Errorf("could not decode animation_repeat_interval %q: %v", repeatIntervalString, err)
		}
		syncToMusicOffsetString := sp.Properties["animation_sync_to_music_offset"]
		if group.SyncToMusicOffset, err = time.ParseDuration(syncToMusicOffsetString); err != nil {
			return fmt.Errorf("could not decode animation_sync_to_music_offset %q: %v", syncToMusicOffsetString, err)
		}
		s.Anim.Init(prefix, map[string]*animation.Group{"default": group}, "default")
	} else {
		return fmt.Errorf("Sprite entity requires exactly one of image, text and animation to be set")
	}
	e.Solid = sp.Properties["solid"] == "true"
	e.Opaque = sp.Properties["opaque"] == "true"
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
	if sp.Properties["no_transform"] == "true" {
		// Undo transform of orientation by tile.
		e.Orientation = sp.Orientation
	}
	if e.Transform.Determinant() < 0 {
		// e.Orientation: in-editor transform. Applied first.
		// Normally the formula is e.Transform.Inverse().Concat(e.Orientation).
		// Add an FlipX() between the two to "undo" any sense difference in the editor.
		// This flips the view on the _level editor_ X axis.
		switch sp.Properties["no_flip"] {
		case "x":
			e.Orientation = e.Transform.Inverse().Concat(m.FlipX()).Concat(sp.Orientation)
		case "y":
			e.Orientation = e.Transform.Inverse().Concat(m.FlipY()).Concat(sp.Orientation)
		}
	}
	return nil
}

func (s *Sprite) Despawn() {
	if s.MyImage {
		s.Entity.Image.Dispose()
		s.MyImage = false
	}
	if s.Anim.Groups != nil {
		s.Anim.Dispose()
	}
}

func (s *Sprite) Update() {
	if s.Anim.Groups != nil {
		s.Anim.Update(s.Entity)
	}
}

func (s *Sprite) Touch(other *engine.Entity) {}

func (s *Sprite) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&Sprite{})
}

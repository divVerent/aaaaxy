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
	"time"

	"github.com/divVerent/aaaaxy/internal/animation"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/level"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

// Animation is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
type Animation struct {
	SpriteBase
	Entity *engine.Entity
	Anim   animation.State
}

func (a *Animation) Spawn(w *engine.World, sp *level.SpawnableProps, e *engine.Entity) error {
	a.Entity = e
	var parseErr error
	prefix := propmap.ValueP(sp.Properties, "animation", "", &parseErr)
	groupName := propmap.ValueP(sp.Properties, "animation_group", "", &parseErr)
	group := &animation.Group{
		NextAnim: groupName,
	}
	group.Frames = propmap.ValueP(sp.Properties, "animation_frames", 0, &parseErr)
	group.Symmetric = propmap.ValueOrP(sp.Properties, "animation_symmetric", false, &parseErr)
	group.FrameInterval = propmap.ValueP(sp.Properties, "animation_frame_interval", 0, &parseErr)
	group.NextInterval = propmap.ValueP(sp.Properties, "animation_repeat_interval", 0, &parseErr)
	group.SyncToMusicOffset = propmap.ValueOrP(sp.Properties, "animation_sync_to_music_offset", time.Duration(0), &parseErr)
	e.RenderOffset = propmap.ValueOrP(sp.Properties, "render_offset", m.Delta{}, &parseErr)
	if e.RenderOffset.IsZero() {
		e.ResizeImage = true
	}
	e.BorderPixels = propmap.ValueOrP(sp.Properties, "border_pixels", 0, &parseErr)
	err := a.Anim.Init(prefix, map[string]*animation.Group{groupName: group}, groupName)
	if err != nil {
		return fmt.Errorf("could not initialize animation %v: %w", prefix, err)
	}
	err = a.SpriteBase.Spawn(w, sp, e)
	if err != nil {
		return err
	}
	return parseErr
}

func (a *Animation) Update() {
	a.Anim.Update(a.Entity)
}

func init() {
	engine.RegisterEntityType(&Animation{})
}

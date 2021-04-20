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

// Animation is a simple entity type that renders a static sprite. It can be optionally solid and/or opaque.
type Animation struct {
	SpriteBase
	Entity *engine.Entity
	Anim   animation.State
}

func (a *Animation) Spawn(w *engine.World, s *level.Spawnable, e *engine.Entity) error {
	a.Entity = e
	prefix := s.Properties["animation"]
	groupName := s.Properties["animation_group"]
	group := &animation.Group{
		NextAnim: groupName,
	}
	framesString := s.Properties["animation_frames"]
	if _, err := fmt.Sscanf(framesString, "%d", &group.Frames); err != nil {
		return fmt.Errorf("could not decode animation_frames %q: %v", framesString, err)
	}
	frameIntervalString := s.Properties["animation_frame_interval"]
	if _, err := fmt.Sscanf(frameIntervalString, "%d", &group.FrameInterval); err != nil {
		return fmt.Errorf("could not decode animation_frame_interval %q: %v", frameIntervalString, err)
	}
	repeatIntervalString := s.Properties["animation_repeat_interval"]
	if _, err := fmt.Sscanf(repeatIntervalString, "%d", &group.NextInterval); err != nil {
		return fmt.Errorf("could not decode animation_repeat_interval %q: %v", repeatIntervalString, err)
	}
	syncToMusicOffsetString := s.Properties["animation_sync_to_music_offset"]
	if group.SyncToMusicOffset, err = time.ParseDuration(syncToMusicOffsetString); err != nil {
		return fmt.Errorf("could not decode animation_sync_to_music_offset %q: %v", syncToMusicOffsetString, err)
	}
	err := a.Anim.Init(prefix, map[string]*animation.Group{groupName: group}, groupName)
	if err != nil {
		return fmt.Errorf("could not initialize animation %v: %v", prefix, err)
	}
	return a.SpriteBase.Spawn(w, s, e)
}

func (a *Animation) Update() {
	a.Anim.Update(a.Entity)
}

func init() {
	engine.RegisterEntityType(&Animation{})
}

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

package animation

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/image"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/music"
)

type Group struct {
	Frames            int           // Number of frames of anim.
	Symmetric         bool          // If true, repeat reversed frames 1..(Frames-2) at end of sequence.
	FrameInterval     int           // Time till next frame.
	NextInterval      int           // Time till NextAnim.
	WaitFinish        bool          // Set if this anim shouldn't be interrupted.
	NextAnim          string        // Name of next animation.
	SyncToMusicOffset time.Duration // Time in music to sync to frame 0.

	// These will be filled in by Init.
	Images    []*ebiten.Image // One image per frame.
	NextGroup *Group          // Pointer to same.
}

type State struct {
	// Global state.
	Groups map[string]*Group

	// Current status.
	Group     *Group
	Frame     int
	WantNext  bool
	NextGroup *Group
}

func (s *State) Init(spritePrefix string, groups map[string]*Group, initialGroup string) error {
	for name, group := range groups {
		if group.NextAnim == "" {
			group.NextGroup = nil
		} else {
			group.NextGroup = groups[group.NextAnim]
			if group.NextGroup == nil {
				return fmt.Errorf("animation group %q references nonexisting frame group %q", name, group.NextAnim)
			}
		}
		images := group.Frames
		if group.Symmetric {
			// Even count: 0 1 2 3 -> 0 1 2 3 2 1, so 6 -> 4
			// Odd count: 0 1 2 3 -> 0 1 2 3 3 2 1, so 7 -> 4
			images = images/2 + 1
		}
		group.Images = make([]*ebiten.Image, images)
		for i := range group.Images {
			var spriteName string
			if images > 1 {
				spriteName = fmt.Sprintf("%s_%s_%d.png", spritePrefix, name, i)
			} else {
				spriteName = fmt.Sprintf("%s_%s.png", spritePrefix, name)
			}
			var err error
			group.Images[i], err = image.Load("sprites", spriteName)
			if err != nil {
				return fmt.Errorf("could not load image %v for group %q: %w", spriteName, name, err)
			}
		}
	}
	s.Groups = groups
	s.ForceGroup(initialGroup)
	s.Group = s.NextGroup // Don't crash on SetGroup calls.
	return nil
}

func (s *State) ForceGroup(group string) {
	requested := s.Groups[group]
	if requested == nil {
		fmt.Printf("Trying to switch to nonexisting frame group: %q.", group)
	}
	s.WantNext = true
	s.NextGroup = requested
}

func (s *State) SetGroup(group string) {
	requested := s.Groups[group]
	if requested == nil {
		fmt.Printf("Trying to switch to nonexisting frame group: %q.", group)
		return
	}
	// Moving to same non-WaitFinish group does nothing.
	if requested == s.Group && !s.Group.WaitFinish {
		return
	}
	// If there's already a scheduled next group, prefer those with WaitFinish.
	if s.NextGroup != nil && s.NextGroup.WaitFinish {
		return
	}
	// Immediately switch over if current group isn't WaitFinish.
	s.WantNext = !s.Group.WaitFinish
	s.NextGroup = requested
}

func (s *State) Update(e *engine.Entity) {
	s.Frame += 1
	if s.NextGroup != nil && (s.WantNext || s.Frame >= s.Group.NextInterval) {
		s.Frame = 0
		s.Group = s.NextGroup
		s.WantNext = false
		s.NextGroup = s.Group.NextGroup
	}
	frame := 0
	if s.Group.FrameInterval != 0 {
		frame = s.Frame / s.Group.FrameInterval
	}
	if frame >= s.Group.Frames {
		frame = s.Group.Frames - 1
	}
	if s.Group.SyncToMusicOffset != 0 {
		absFrame := int((music.Now() - s.Group.SyncToMusicOffset) * engine.GameTPS / (time.Second * time.Duration(s.Group.FrameInterval)))
		frame = m.Mod(absFrame, s.Group.Frames)
	}
	image := frame
	if s.Group.Symmetric {
		other := s.Group.Frames - image
		if other < image {
			image = other
		}
	}
	e.Image = s.Group.Images[image]
}

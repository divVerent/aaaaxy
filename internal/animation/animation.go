package animation

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
)

type Group struct {
	Frames        int    // Number of frames of anim.
	FrameInterval int    // Time till next frame.
	NextInterval  int    // Time till NextAnim.
	WaitFinish    bool   // Set if this anim shouldn't be interrupted.
	NextAnim      string // Name of next animation.

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

func (s *State) Init(spritePrefix string, groups map[string]*Group, initialGroup string) {
	for name, group := range groups {
		group.NextGroup = groups[group.NextAnim]
		if group.NextGroup == nil {
			log.Panicf("Animation group %q references nonexisting frame group %q", group, group.NextGroup)
		}
		group.Images = make([]*ebiten.Image, group.Frames)
		for i := range group.Images {
			var spriteName string
			if group.Frames > 1 {
				spriteName = fmt.Sprintf("%s_%s_%d.png", spritePrefix, name, i)
			} else {
				spriteName = fmt.Sprintf("%s_%s.png", spritePrefix, name)
			}
			var err error
			group.Images[i], err = engine.LoadImage("sprites", spriteName)
			if err != nil {
				log.Panicf("Could not load image %v for group %q: %v", spriteName, name, err)
			}
		}
	}
	s.Groups = groups
	s.ForceGroup(initialGroup)
}

func (s *State) ForceGroup(group string) {
	requested := s.Groups[group]
	if requested == nil {
		log.Panicf("Trying to switch to nonexisting frame group: %q.", group)
	}
	s.WantNext = true
	s.NextGroup = requested
}

func (s *State) SetGroup(group string) {
	requested := s.Groups[group]
	if requested == nil {
		log.Panicf("Trying to switch to nonexisting frame group: %q.", group)
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
	s.WantNext = s.Group == nil || !s.Group.WaitFinish
	s.NextGroup = requested
}

func (s *State) Update(e *engine.Entity) {
	s.Frame += 1
	if s.WantNext || s.Frame >= s.Group.NextInterval {
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
	e.Image = s.Group.Images[frame]
}

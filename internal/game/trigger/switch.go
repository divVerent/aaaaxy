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

package trigger

import (
	"github.com/divVerent/aaaaxy/internal/animation"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/level"
)

// Switch overrides the boolean state of a warpzone or entity.
type Switch struct {
	SetState
	Entity    *engine.Entity
	Anim      animation.State
	AnimState bool
}

func (s *Switch) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	s.Entity = e
	// Switches ALWAYS send every frame.
	// TODO: replace by making SetState's default this and all other uses explicit.
	sp.Properties["send_every_frame"] = "true"
	err := s.SetState.Spawn(w, sp, e)
	if err != nil {
		return err
	}
	err = s.Anim.Init("switch", map[string]*animation.Group{
		"switchon": {
			Frames:        10,
			FrameInterval: 2,
			NextInterval:  2 * 10,
			NextAnim:      "on",
		},
		"on": {
			Frames: 1,
		},
		"switchoff": {
			Frames:        10,
			FrameInterval: 2,
			NextInterval:  2 * 10,
			NextAnim:      "off",
		},
		"off": {
			Frames: 1,
		},
	}, "off")
	if err != nil {
		return err
	}
	return nil
}

func (s *Switch) Update() {
	s.SetState.Update()
	if s.State != s.AnimState {
		if s.State {
			s.Anim.SetGroup("switchon")
			s.AnimState = true
		} else if s.SendUntouch {
			s.Anim.SetGroup("switchoff")
			s.AnimState = false
		}
	}
	s.Anim.Update(s.Entity)
}

func init() {
	engine.RegisterEntityType(&Switch{})
}

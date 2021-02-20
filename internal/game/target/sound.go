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

package target

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/audiowrap"
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
	"github.com/divVerent/aaaaaa/internal/sound"
)

// SoundTarget just changes the music track to the given one.
type SoundTarget struct {
	Sound  *sound.Sound
	Player *audiowrap.Player
}

func (s *SoundTarget) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	var err error
	s.Sound, err = sound.Load(sp.Properties["sound"])
	if err != nil {
		return fmt.Errorf("could not load sound: %v", err)
	}
	return nil
}

func (s *SoundTarget) Despawn() {}

func (s *SoundTarget) Update() {}

func (s *SoundTarget) Touch(other *engine.Entity) {}

func (s *SoundTarget) SetState(state bool) {
	if state {
		if s.Player == nil || s.Player.IsPlaying() {
			s.Player = s.Sound.Play()
		}
	} else {
		if s.Player != nil {
			s.Player.Close()
			s.Player = nil
		}
	}
}

func (s *SoundTarget) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&SoundTarget{})
}

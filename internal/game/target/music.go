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
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/level"
	m "github.com/divVerent/aaaaaa/internal/math"
	"github.com/divVerent/aaaaaa/internal/music"
)

// SwitchMusicTarget just changes the music track to the given one.
type SwitchMusicTarget struct {
	Music string
}

func (s *SwitchMusicTarget) Spawn(w *engine.World, sp *level.Spawnable, e *engine.Entity) error {
	s.Music = sp.Properties["music"]
	return nil
}

func (s *SwitchMusicTarget) Despawn() {}

func (s *SwitchMusicTarget) Update() {}

func (s *SwitchMusicTarget) Touch(other *engine.Entity) {}

func (s *SwitchMusicTarget) SetState(state bool) {
	if state {
		music.Switch(s.Music)
	}
}

func (s *SwitchMusicTarget) DrawOverlay(screen *ebiten.Image, scrollDelta m.Delta) {}

func init() {
	engine.RegisterEntityType(&SwitchMusicTarget{})
}

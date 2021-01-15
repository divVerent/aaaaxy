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

package sound

import (
	"fmt"
	"io/ioutil"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"

	"github.com/divVerent/aaaaaa/internal/vfs"
)

// Sound represents a sound effect.
type Sound struct {
	sound   []byte
	players []*audio.Player
}

// Sounds are preloaded as byte streams.
var cache = map[string]*Sound{}

// Load loads a sound effect.
// Multiple Load calls to the same sound effect return the same cached instance.
func Load(name string) (*Sound, error) {
	if sound, found := cache[name]; found {
		return sound, nil
	}
	data, err := vfs.Load("sounds", name)
	if err != nil {
		return nil, fmt.Errorf("could not load: %v", err)
	}
	defer data.Close()
	stream, err := vorbis.Decode(audio.CurrentContext(), data)
	if err != nil {
		return nil, fmt.Errorf("could not start decoding: %v", err)
	}
	decoded, err := ioutil.ReadAll(stream)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %v", err)
	}
	sound := &Sound{sound: decoded}
	cache[name] = sound
	return sound, nil
}

// Play plays the given sound effect.
func (s *Sound) Play() {
	audio.NewPlayerFromBytes(audio.CurrentContext(), s.sound).Play()
}

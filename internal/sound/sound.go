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
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"

	"github.com/divVerent/aaaaaa/internal/audiowrap"
	"github.com/divVerent/aaaaaa/internal/flag"
	"github.com/divVerent/aaaaaa/internal/vfs"
)

var (
	precacheSounds = flag.Bool("precache_sounds", true, "preload all sounds at startup (VERY recommended)")
	soundVolume    = flag.Float64("sound_volume", 0.5, "sound volume (0..1)")
)

// Sound represents a sound effect.
type Sound struct {
	sound        []byte
	players      []*audiowrap.Player
	volumeAdjust float64
}

// Sounds are preloaded as byte streams.
var (
	cache       = map[string]*Sound{}
	cacheFrozen bool
)

type soundJson struct {
	VolumeAdjust float64 `json:"volume_adjust"`
}

// Load loads a sound effect.
// Multiple Load calls to the same sound effect return the same cached instance.
func Load(name string) (*Sound, error) {
	cacheName := vfs.Canonical("sounds", name)
	if sound, found := cache[cacheName]; found {
		return sound, nil
	}
	if cacheFrozen {
		return nil, fmt.Errorf("sound %v was not precached", name)
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
	config := soundJson{
		VolumeAdjust: 1,
	}
	j, err := vfs.Load("sounds", name+".json")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("could not load sound json config file for %q: %v", name, err)
	}
	if j != nil {
		defer j.Close()
		err = json.NewDecoder(j).Decode(&config)
		if err != nil {
			return nil, fmt.Errorf("could not decode sound json config file for %q: %v", name, err)
		}
	}
	sound := &Sound{
		sound:        decoded,
		volumeAdjust: config.VolumeAdjust,
	}
	cache[cacheName] = sound
	return sound, nil
}

// Play plays the given sound effect.
func (s *Sound) Play() *audiowrap.Player {
	player := audiowrap.NewPlayerFromBytes(s.sound)
	player.SetVolume(s.volumeAdjust * *soundVolume)
	player.Play()
	return player
}

func Precache() error {
	if !*precacheSounds {
		return nil
	}
	names, err := vfs.ReadDir("sounds")
	if err != nil {
		return fmt.Errorf("could not enumerate sounds: %v", err)
	}
	for _, name := range names {
		if !strings.HasSuffix(name, ".ogg") {
			continue
		}
		_, err := Load(name)
		if err != nil {
			return fmt.Errorf("could not precache %v: %v", name, err)
		}
	}
	cacheFrozen = true
	return nil
}

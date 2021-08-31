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

package music

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"

	"github.com/divVerent/aaaaxy/internal/audiowrap"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

var (
	musicVolume   = flag.Float64("music_volume", 0.5, "music volume (0..1)")
	musicFadeTime = flag.Duration("music_fade_time", time.Second, "music fade time")
)

const (
	bytesPerSample = 4
)

type musicJson struct {
	PlayStart  int64   `json:"play_start"`
	ReplayGain float64 `json:"replay_gain"`
	LoopStart  int64   `json:"loop_start"`
	LoopEnd    int64   `json:"loop_end"`
}

type sampleCutter struct {
	base   io.ReadSeeker
	closer io.Closer

	offset int64
}

var _ io.ReadSeeker = &sampleCutter{}

func (c *sampleCutter) Read(b []byte) (int, error) {
	return c.base.Read(b)
}

func (c *sampleCutter) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekStart {
		offset += c.offset
	}
	o, err := c.base.Seek(offset, whence)
	return o - c.offset, err
}

func (c *sampleCutter) Close() error {
	return c.closer.Close()
}

func newSampleCutter(base io.ReadSeeker, offset int64, closer io.Closer) (*sampleCutter, error) {
	f := &sampleCutter{
		base:   base,
		offset: offset,
		closer: closer,
	}
	_, err := f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("could not build sample cutter: %v", err)
	}
	return f, nil
}

var (
	currentName string
	player      *audiowrap.Player
	active      bool
)

func Enable() {
	if !active && player != nil {
		player.Play()
	}
	active = true
}

// Now returns the current music playback time.
func Now() time.Duration {
	if player != nil && player.IsPlaying() {
		return player.Current()
	}
	return 0
}

// Switch switches from the currently playing music to the given track.
// Passing an empty string means fading to silence.
func Switch(name string) {
	cacheName := vfs.Canonical("music", name)
	if cacheName == currentName {
		return
	}

	// Fade out the current music.
	if player != nil {
		player.FadeOutIn(*musicFadeTime)
		player = nil
	}

	// Switch to it.
	currentName = cacheName

	// If we're playing silence, we're done.
	if cacheName == "" {
		return
	}

	// Now load the new track.
	config := musicJson{
		PlayStart:  0,
		LoopStart:  0,
		LoopEnd:    -1,
		ReplayGain: 1,
	}
	j, err := vfs.Load("music", name+".json")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Errorf("could not load music json config file for %q: %v", name, err)
		return
	}
	if j != nil {
		defer j.Close()
		err = json.NewDecoder(j).Decode(&config)
		if err != nil {
			log.Errorf("could not decode music json config file for %q: %v", name, err)
			return
		}
	}
	player, err = audiowrap.NewPlayer(func() (io.ReadCloser, error) {
		handle, err := vfs.Load("music", name)
		if err != nil {
			return nil, fmt.Errorf("could not load music %q: %v", name, err)
		}
		data, err := vorbis.DecodeWithSampleRate(audiowrap.SampleRate(), handle)
		if err != nil {
			return nil, fmt.Errorf("could not start decoding music %q: %v", name, err)
		}
		loopEnd := data.Length()
		if config.LoopEnd >= 0 {
			loopEnd = config.LoopEnd * bytesPerSample
		}
		return newSampleCutter(audio.NewInfiniteLoopWithIntro(data, config.LoopStart*bytesPerSample, loopEnd), config.PlayStart*bytesPerSample, handle)
	})
	if err != nil {
		log.Errorf("could not start playing music %q: %v", name, err)
		return
	}

	// We have a valid player.
	player.SetVolume(*musicVolume * config.ReplayGain)
	if active {
		player.Play()
	}
}

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
	"github.com/divVerent/aaaaxy/internal/log"
	"io"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"

	"github.com/divVerent/aaaaxy/internal/audiowrap"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

var (
	musicVolume   = flag.Float64("music_volume", 0.5, "music volume (0..1)")
	musicFadeTime = flag.Duration("music_fade_time", time.Second, "music fade time")
)

const (
	bytesPerSample = 4
)

type track struct {
	name       string
	valid      bool
	replayGain float64
	player     *audiowrap.Player
}

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

func (t *track) open(cacheName, name string) error {
	t.stop()
	t.valid = true
	if cacheName == "" {
		return nil
	}
	t.name = cacheName

	config := musicJson{
		PlayStart:  0,
		LoopStart:  0,
		LoopEnd:    -1,
		ReplayGain: 1,
	}
	j, err := vfs.Load("music", name+".json")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.stop()
		return fmt.Errorf("could not load music json config file for %q: %v", name, err)
	}
	if j != nil {
		defer j.Close()
		err = json.NewDecoder(j).Decode(&config)
		if err != nil {
			t.stop()
			return fmt.Errorf("could not decode music json config file for %q: %v", name, err)
		}
	}
	t.replayGain = config.ReplayGain

	t.player, err = audiowrap.NewPlayer(func() (io.ReadCloser, error) {
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
		t.stop()
		return fmt.Errorf("could not start playing music %q: %v", name, err)
	}

	return nil
}

func (t *track) play() {
	if t.player != nil {
		t.player.SetVolume(*musicVolume * t.replayGain)
		t.player.Play()
	}
}

func (t *track) stop() {
	if t.player != nil {
		t.player.FadeOutIn(*musicFadeTime)
	}
	t.player = nil
	t.valid = false
	t.name = ""
}

var (
	current, next track
	active        bool
)

func Enable() {
	active = true
}

func Update() {
	if !active {
		return
	}
	// Advance track.
	if next.valid {
		if current.valid {
			current.stop()
		}
		next.play()
		current, next = next, current
	}
}

// Now returns the current music playback time.
func Now() time.Duration {
	if current.valid && current.player != nil && current.player.IsPlaying() {
		return current.player.Current()
	}
	return 0
}

// Switch switches from the currently playing music to the given track.
// Passing an empty string means fading to silence.
func Switch(name string) {
	cacheName := vfs.Canonical("music", name)
	if next.valid {
		if next.name == cacheName {
			return
		}
	} else if current.valid {
		if current.name == cacheName {
			return
		}
	}
	if next.player != nil {
		next.player.CloseInstantly()
		next.player = nil
	}
	err := next.open(cacheName, name)
	if err != nil {
		log.Errorf("could not open music %q: %v", name, err)
	}
}

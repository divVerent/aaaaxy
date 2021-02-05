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
	"flag"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"

	"github.com/divVerent/aaaaaa/internal/audiowrap"
	"github.com/divVerent/aaaaaa/internal/vfs"
)

var (
	musicVolume = flag.Float64("music_volume", 0.5, "music volume (0..1)")
)

const (
	xFadeFrameOut  = 40
	xFadeFrameIn   = 20
	xFadeFrameEnd  = 60
	bytesPerSample = 4
)

type track struct {
	name   string
	valid  bool
	handle vfs.ReadSeekCloser
	data   *vorbis.Stream
	player *audiowrap.Player
}

type musicJson struct {
	LoopStart int64 `json:"loop_start"`
	LoopEnd   int64 `json:"loop_end"`
}

func (t *track) open(name string) {
	t.stop()
	t.valid = true
	if name == "" {
		return
	}
	t.name = name
	var err error
	t.handle, err = vfs.Load("music", name)
	if err != nil {
		log.Panicf("Could not load music %q: %v", name, err)
	}
	t.data, err = vorbis.Decode(audio.CurrentContext(), t.handle)
	if err != nil {
		log.Panicf("Could not start decoding music %q: %v", name, err)
	}
	config := musicJson{
		LoopStart: 0,
		LoopEnd:   t.data.Length() / bytesPerSample,
	}
	j, err := vfs.Load("music", name+".json")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Panicf("Could not load music json config file for %q: %v", name, err)
	}
	if j != nil {
		defer j.Close()
		err = json.NewDecoder(j).Decode(&config)
		if err != nil {
			log.Panicf("Could not decode music json config file for %q: %v", name, err)
		}
	}
	loop := audio.NewInfiniteLoopWithIntro(t.data, config.LoopStart*bytesPerSample, config.LoopEnd*bytesPerSample)
	t.player, err = audiowrap.NewPlayer(loop)
	if err != nil {
		log.Panicf("Could not start playing music %q: %v", name, err)
	}
}

func (t *track) play() {
	if t.player != nil {
		t.player.Play()
	}
}

func (t *track) setVolume(vol float64) {
	if t.player != nil {
		t.player.SetVolume(vol * *musicVolume)
	}
}

func (t *track) stop() {
	if t.player != nil {
		t.player.Close()
	}
	t.player = nil
	t.data = nil
	if t.handle != nil {
		t.handle.Close()
	}
	t.handle = nil
	t.valid = false
	t.name = ""
}

var (
	current, fadeTo, next track
	haveNext              bool
	xFadeFrame            int
	idleNotifier          chan<- struct{}
)

func Update() {
	// If idle, initiate fading.
	if !fadeTo.valid && next.valid {
		fadeTo, next = next, fadeTo
		xFadeFrame = 0
		// Skip right to fade-in if we are currently playing silence.
		if current.player == nil {
			xFadeFrame = xFadeFrameIn
		}
	}
	// Nothing to fade?
	if !fadeTo.valid {
		if idleNotifier != nil {
			close(idleNotifier)
			idleNotifier = nil
		}
		return
	}
	// Advance.
	xFadeFrame++
	if current.valid {
		if xFadeFrame == xFadeFrameOut {
			current.stop()
		} else {
			current.setVolume(float64(xFadeFrameOut-xFadeFrame) / float64(xFadeFrameOut))
		}
	}
	if fadeTo.valid && xFadeFrame > xFadeFrameIn {
		fadeTo.setVolume(float64(xFadeFrame-xFadeFrameIn) / float64(xFadeFrameEnd-xFadeFrameIn))
		if xFadeFrame == xFadeFrameIn+1 {
			fadeTo.play()
		}
	}
	if xFadeFrame == xFadeFrameEnd {
		current, fadeTo = fadeTo, current
		xFadeFrame = 0
	}
}

// Now returns the current music playback time.
func Now() time.Duration {
	if fadeTo.valid && fadeTo.player != nil && fadeTo.player.IsPlaying() {
		return fadeTo.player.Current()
	}
	if current.valid && current.player != nil && current.player.IsPlaying() {
		return current.player.Current()
	}
	return 0
}

// Switch switches from the currently playing music to the given track.
// Passing an empty string means fading to silence.
func Switch(name string) {
	name = vfs.Canonical(name)
	if next.valid {
		if next.name == name {
			return
		}
	} else if fadeTo.valid {
		if fadeTo.name == name {
			return
		}
	} else if current.valid {
		if current.name == name {
			return
		}
	}
	if next.player != nil {
		next.player.Close()
		next.player = nil
		next.data = nil
		next.handle.Close()
		next.handle = nil
	}
	next.open(name)
}

// End ends all music playback, then notifies the given channel.
func End() <-chan struct{} {
	Switch("")
	ch := make(chan struct{})
	idleNotifier = ch
	return ch
}

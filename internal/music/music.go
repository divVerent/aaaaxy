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
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"

	"github.com/divVerent/aaaaaa/internal/vfs"
)

const (
	xFadeFrameOut = 120
	xFadeFrameIn  = 60
	xFadeFrameEnd = 180
)

type track struct {
	name   string
	valid  bool
	handle vfs.ReadSeekCloser
	data   *vorbis.Stream
	player *audio.Player
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
		log.Panicf("Could not load music %q: %v", name)
	}
	t.data, err = vorbis.Decode(audio.CurrentContext(), t.handle)
	if err != nil {
		log.Panicf("Could not start decoding music %q: %v", name)
	}
	loop := audio.NewInfiniteLoop(t.data, t.data.Length())
	t.player, err = audio.NewPlayer(audio.CurrentContext(), loop)
	if err != nil {
		log.Panicf("Could not start playing music %q: %v", name)
	}
}

func (t *track) play() {
	if t.player != nil {
		t.player.Play()
	}
}

func (t *track) setVolume(vol float64) {
	if t.player != nil {
		t.player.SetVolume(vol)
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

// Switch switches from the currently playing music to the given track.
// Passing an empty string means fading to silence.
func Switch(name string) {
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

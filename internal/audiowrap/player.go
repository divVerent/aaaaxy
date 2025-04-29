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

package audiowrap

import (
	"bytes"
	"io"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"

	"github.com/divVerent/aaaaxy/internal/dontgc"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	audioFlag     = flag.Bool("audio", true, "enable audio")
	audioRate     = flag.Int("audio_rate", 44100, "preferred audio sample rate")
	volume        = flag.Float64("volume", 0.5, "global volume (0..1)")
	soundFadeTime = flag.Duration("sound_fade_time", time.Second, "default sound fade time")
)

type Player struct {
	ebi       *audio.Player
	ebiCloser io.Closer
	dmp       *dumper

	// These fields are only really used when -audio=false.
	accumulatedTime time.Duration
	playTime        time.Time

	// Debug info to print if this were to be GC'd while still playing.
	dontGCState dontgc.State

	// State for fading out.
	volume     float64
	fadeFrames int
	fadeFrame  int
}

func NoPlayer() *Player {
	p := &Player{}
	p.dontGCState = dontgc.SetUp(p)
	return p
}

type FadeHandle struct {
	player *Player
}

var (
	fadingOutPlayers = map[*Player]struct{}{}
	fadingInPlayers  = map[*Player]struct{}{}
)

func Rate() int {
	return *audioRate
}

func Init() error {
	if *audioFlag {
		audio.NewContext(*audioRate)

		// Workaround: for some reason playing the first sound can incur significant delay.
		// So let's do this at the start.
		audio.CurrentContext().NewPlayerFromBytes([]byte{}).Play()
	}
	return nil
}

func SampleRate() int {
	if *audioFlag {
		return audio.CurrentContext().SampleRate()
	}
	return *audioRate
}

func Update() {
	for p := range fadingOutPlayers {
		p.fadeFrame--
		if p.fadeFrame == 0 {
			p.CloseInstantly()
			delete(fadingOutPlayers, p)
		}
		v := p.volume * float64(p.fadeFrame) / float64(p.fadeFrames)
		p.setVolume(v)
	}
	for p := range fadingInPlayers {
		p.fadeFrame++
		if p.fadeFrame == p.fadeFrames {
			delete(fadingInPlayers, p)
		}
		v := p.volume * float64(p.fadeFrame) / float64(p.fadeFrames)
		p.setVolume(v)
	}
}

func ebiPlayer(src io.Reader) (*audio.Player, error) {
	if !*audioFlag {
		return nil, nil
	}
	return audio.CurrentContext().NewPlayer(src)
}

func NewPlayer(src func() (io.ReadCloser, error)) (*Player, error) {
	dmp, err := newDumper(src)
	if err != nil {
		return nil, err
	}
	srcReader, err := src()
	if err != nil {
		return nil, err
	}
	ebi, err := ebiPlayer(srcReader)
	if err != nil {
		return nil, err
	}
	p := &Player{
		ebi:       ebi,
		ebiCloser: srcReader,
		dmp:       dmp,
	}
	p.dontGCState = dontgc.SetUp(p)
	return p, nil
}

func (p *Player) CheckGC() dontgc.State {
	if !p.IsPlaying() {
		return nil
	}
	p.CloseInstantly()
	return p.dontGCState
}

func ebiPlayerFromBytes(src []byte) *audio.Player {
	if !*audioFlag {
		return nil
	}
	return audio.CurrentContext().NewPlayerFromBytes(src)
}

func NewPlayerFromBytes(src []byte) (*Player, error) {
	dmp, err := newDumper(func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(src)), nil
	})
	if err != nil {
		return nil, err
	}
	ebi := ebiPlayerFromBytes(src)
	return &Player{
		ebi: ebi,
		dmp: dmp,
	}, nil
}

func (p *Player) CloseInstantly() error {
	p.playTime = time.Time{}
	if p.dmp != nil {
		p.dmp.Close()
	}
	var err error = nil
	if p.ebi != nil {
		err2 := p.ebi.Close()
		if err == nil {
			err = err2
		}
	}
	if p.ebiCloser != nil {
		err2 := p.ebiCloser.Close()
		if err == nil {
			err = err2
		}
	}
	return err
}

func (p *Player) Close() error {
	if p.volume == 0 || !p.IsPlaying() {
		p.CloseInstantly()
	} else {
		p.FadeOutIn(*soundFadeTime)
	}
	return nil
}

func toFrames(d time.Duration) int {
	frames := int((d*engine.GameTPS + (time.Second / 2)) / time.Second)
	if frames < 0 {
		frames = 0
	}
	return frames
}

func (p *Player) FadeOutIn(d time.Duration) *FadeHandle {
	frames := toFrames(d)
	if _, found := fadingInPlayers[p]; found {
		// ceil-convert the frame number. Then next frame has the lowest possible reduction.
		p.fadeFrame = (p.fadeFrame*frames + p.fadeFrames - 1) / p.fadeFrames
		// Need at least frame 1 so next frame is at least 0.
		// Note: this can only happen if p.fadeFrame was previously 0, which is... odd.
		// However, it ccan happen when RestoreIn was called the same frame and resulted in a frame number of zero.
		if p.fadeFrame < 1 {
			p.fadeFrame = 1
		}
		delete(fadingInPlayers, p)
	} else {
		p.fadeFrame = frames
	}
	p.fadeFrames = frames
	fadingOutPlayers[p] = struct{}{}
	return &FadeHandle{
		player: p,
	}
}

func (f *FadeHandle) RestoreIn(d time.Duration) *Player {
	frames := toFrames(d)
	p := f.player
	if _, found := fadingOutPlayers[p]; !found {
		return nil
	}
	delete(fadingOutPlayers, p)
	// floor-convert the frame number. Then next frame has the lowest possible increase.
	p.fadeFrame = p.fadeFrame * frames / p.fadeFrames
	// Need at most frame frames-1 so next frame is at most frames.
	if p.fadeFrame > frames-1 {
		p.fadeFrame = frames - 1
	}
	p.fadeFrames = frames
	fadingInPlayers[p] = struct{}{}
	return p
}

func (p *Player) Position() time.Duration {
	if p.dmp != nil {
		return p.dmp.Position()
	}
	if p.ebi != nil {
		// Type switch to get around deprecation warning on Current() in Ebitengine v2.6.
		// TODO(divVerent): Remove when requiring Ebitengine 2.6.
		switch ebi := (interface{})(p.ebi).(type) {
		case interface{ Position() time.Duration }:
			return ebi.Position()
		case interface{ Current() time.Duration }:
			return ebi.Current()
		default:
			log.Fatalf("Ebitengine player of unknown type: %#v", ebi)
		}
	}
	t := p.accumulatedTime
	if !p.playTime.IsZero() {
		t += time.Since(p.playTime)
	}
	return t
}

func (p *Player) IsPlaying() bool {
	if p.dmp != nil {
		return p.dmp.IsPlaying()
	}
	if p.ebi != nil {
		return p.ebi.IsPlaying()
	}
	return !p.playTime.IsZero()
}

func (p *Player) Pause() {
	if p.dmp != nil {
		p.dmp.Pause()
	}
	if p.ebi != nil {
		p.ebi.Pause()
	}
	if !p.playTime.IsZero() {
		p.accumulatedTime += time.Since(p.playTime)
		p.playTime = time.Time{}
	}
}

func (p *Player) Play() {
	if p.dmp != nil {
		p.dmp.Play()
	}
	if p.ebi != nil {
		p.ebi.Play()
	}
	if p.playTime.IsZero() {
		p.playTime = time.Now()
	}
}

func (p *Player) SetVolume(vol float64) {
	p.volume = vol // For fading.
	p.setVolume(vol)
	delete(fadingInPlayers, p)
}

func (p *Player) setVolume(vol float64) {
	if p.dmp != nil {
		p.dmp.SetVolume(vol * *volume)
	}
	if p.ebi != nil {
		p.ebi.SetVolume(vol * *volume)
	}
}

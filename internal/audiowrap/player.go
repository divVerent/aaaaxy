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

	ebiaudio "github.com/hajimehoshi/ebiten/v2/audio"

	"github.com/divVerent/aaaaxy/internal/dontgc"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/go117"
	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	audio         = flag.Bool("audio", true, "enable audio")
	audioRate     = flag.Int("audio_rate", 44100, "preferred audio sample rate")
	volume        = flag.Float64("volume", 0.5, "global volume (0..1)")
	soundFadeTime = flag.Duration("sound_fade_time", time.Second, "default sound fade time")
	// TODO: add a way to simulate audio and write to disk, syncing with the frame clock (i.e. each frame renders exactly 1/60 sec of audio).
	// Also a way to don't actually render audio (but still advance clock) would be nice.
)

type Player struct {
	ebi       *ebiaudio.Player
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
	if *audio {
		ebiaudio.NewContext(*audioRate)
	}
	return nil
}

func SampleRate() int {
	if *audio {
		return ebiaudio.CurrentContext().SampleRate()
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

func ebiPlayer(src io.Reader) (*ebiaudio.Player, error) {
	if !*audio {
		return nil, nil
	}
	return ebiaudio.NewPlayer(ebiaudio.CurrentContext(), src)
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

func ebiPlayerFromBytes(src []byte) *ebiaudio.Player {
	if !*audio {
		return nil
	}
	return ebiaudio.NewPlayerFromBytes(ebiaudio.CurrentContext(), src)
}

func NewPlayerFromBytes(src []byte) *Player {
	dmp, err := newDumper(func() (io.ReadCloser, error) {
		return go117.NopCloser(bytes.NewReader(src)), nil
	})
	if err != nil {
		log.Fatalf("UNREACHABLE CODE: newDumper returned an error despite passed an always-succeed function: %v", err)
		return nil
	}
	ebi := ebiPlayerFromBytes(src)
	return &Player{
		ebi: ebi,
		dmp: dmp,
	}
}

func (p *Player) CloseInstantly() error {
	p.playTime = time.Time{}
	if p.dmp != nil {
		p.dmp.Close()
	}
	if p.ebi != nil {
		return p.ebi.Close()
	}
	if p.ebiCloser != nil {
		p.ebiCloser.Close()
	}
	return nil
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

func (p *Player) Current() time.Duration {
	if p.dmp != nil {
		return p.dmp.Current()
	}
	if p.ebi != nil {
		return p.ebi.Current()
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

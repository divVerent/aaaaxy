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

	"github.com/divVerent/aaaaxy/internal/dontgc"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
	ebiaudio "github.com/hajimehoshi/ebiten/v2/audio"
)

var (
	audio     = flag.Bool("audio", true, "enable audio")
	audioRate = flag.Int("audio_rate", 44100, "preferred audio sample rate")
	volume    = flag.Float64("volume", 1.0, "global volume (0..1)")
	// TODO: add a way to simulate audio and write to disk, syncing with the frame clock (i.e. each frame renders exactly 1/60 sec of audio).
	// Also a way to don't actually render audio (but still advance clock) would be nice.
)

type Player struct {
	ebi *ebiaudio.Player
	dmp *dumper

	// These fields are only really used when -audio=false.
	accumulatedTime time.Duration
	playTime        time.Time

	// Debug info to print if this were to be GC'd while still playing.
	dontGCState dontgc.State
}

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

func ebiPlayer(src io.Reader) (*ebiaudio.Player, error) {
	if !*audio {
		return nil, nil
	}
	return ebiaudio.NewPlayer(ebiaudio.CurrentContext(), src)
}

func NewPlayer(src func() (io.Reader, error)) (*Player, error) {
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
		ebi: ebi,
		dmp: dmp,
	}
	p.dontGCState = dontgc.SetUp(p)
	return p, nil
}

func (p *Player) CheckGC() dontgc.State {
	if !p.IsPlaying() {
		return nil
	}
	p.Close()
	return p.dontGCState
}

func ebiPlayerFromBytes(src []byte) *ebiaudio.Player {
	if !*audio {
		return nil
	}
	return ebiaudio.NewPlayerFromBytes(ebiaudio.CurrentContext(), src)
}

func NewPlayerFromBytes(src []byte) *Player {
	dmp, err := newDumper(func() (io.Reader, error) {
		return bytes.NewReader(src), nil
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

func (p *Player) Close() error {
	if p.dmp != nil {
		p.dmp.Close()
	}
	if p.ebi != nil {
		return p.ebi.Close()
	}
	p.playTime = time.Time{}
	return nil
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
	if p.dmp != nil {
		p.dmp.SetVolume(vol * *volume)
	}
	if p.ebi != nil {
		p.ebi.SetVolume(vol * *volume)
	}
}

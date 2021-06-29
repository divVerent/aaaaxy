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
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	ebiaudio "github.com/hajimehoshi/ebiten/v2/audio"
)

type dumper struct {
	reader  io.Reader
	volume  float64
	playing bool
	played  int
}

var (
	dumping       bool
	currentSounds []*dumper
	sampleIndex   int
)

func InitDumping() {
	dumping = true
}

func DumpFrame(dumpFile io.Writer, toTime time.Duration) error {
	if !dumping {
		return errors.New("DumpFrame called when not dumping")
	}
	toSample := int(toTime * time.Duration(ebiaudio.CurrentContext().SampleRate()) / time.Second)
	samples := toSample - sampleIndex
	sampleIndex = toSample
	return dumpSamples(dumpFile, samples)
}

func dumpSamples(dumpFile io.Writer, samples int) error {
	buf := make([]int16, 2*samples)
	toClose := []*dumper{}
	for _, dmp := range currentSounds {
		if dmp.addTo(buf) != nil {
			toClose = append(toClose, dmp)
		}
	}
	for _, dmp := range toClose {
		dmp.Close()
	}
	err := binary.Write(dumpFile, binary.LittleEndian, buf)
	if err != nil {
		dumping = false
		return fmt.Errorf("cannot dump audio frame: %v", err)
	}
	return nil
}

func newDumper(src func() (io.Reader, error)) (*dumper, error) {
	if !dumping {
		return nil, nil
	}
	srcReader, err := src()
	if err != nil {
		return nil, err
	}
	dmp := &dumper{
		reader:  srcReader,
		volume:  0.0,
		playing: false,
	}
	currentSounds = append(currentSounds, dmp)
	return dmp, nil
}

func (d *dumper) Close() {
	d.playing = false
	for i, snd := range currentSounds {
		if snd == d {
			currentSounds = append(currentSounds[:i], currentSounds[(i+1):]...)
			return
		}
	}
}

func (d *dumper) Current() time.Duration {
	return time.Duration(d.played) * time.Second / time.Duration(ebiaudio.CurrentContext().SampleRate())
}

func (d *dumper) IsPlaying() bool {
	return d.playing
}

func (d *dumper) Pause() {
	d.playing = false
}

func (d *dumper) Play() {
	d.playing = true
}

func (d *dumper) SetVolume(vol float64) {
	d.volume = vol
}

func (d *dumper) addTo(buf []int16) error {
	if !d.playing {
		return nil
	}
	addBuf := make([]int16, len(buf))
	err := binary.Read(d.reader, binary.LittleEndian, addBuf)
	for i, s := range addBuf {
		buf[i] += int16(d.volume * float64(s))
	}
	d.played += len(buf) / 2
	return err
}

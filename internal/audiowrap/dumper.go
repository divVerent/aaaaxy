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
	"io"
	"log"
	"os"
	"time"

	ebiaudio "github.com/hajimehoshi/ebiten/v2/audio"

	"github.com/divVerent/aaaaaa/internal/flag"
)

var (
	dumpAudio = flag.String("dump_audio", "", "filename to dump audio to")
)

type dumper struct {
	reader  io.Reader
	volume  float64
	playing bool
	played  int
}

var (
	dumpFile      io.WriteCloser
	currentSounds []*dumper
	sampleIndex   int
)

func Update(toTime time.Duration) {
	if *dumpAudio == "" {
		return
	}
	if dumpFile == nil {
		var err error
		dumpFile, err = os.Create(*dumpAudio)
		if err != nil {
			log.Printf("cannot create audio dump file: %v", err)
			*dumpAudio = ""
		}
	}
	toSample := int(toTime * time.Duration(ebiaudio.CurrentContext().SampleRate()) / time.Second)
	samples := toSample - sampleIndex
	dumpSamples(samples)
}

func dumpSamples(samples int) {
	buf := make([]int16, 2*samples)
	for _, dmp := range currentSounds {
		dmp.addTo(buf)
	}
	err := binary.Write(dumpFile, binary.LittleEndian, buf)
	if err != nil {
		log.Printf("cannot dump audio frame: %v", err)
		dumpFile.Close()
		dumpFile = nil
		*dumpAudio = ""
	}
}

func newDumper(src io.Reader) *dumper {
	if *dumpAudio == "" {
		return nil
	}
	dmp := &dumper{
		reader:  src,
		volume:  0.0,
		playing: false,
	}
	currentSounds = append(currentSounds, dmp)
	return dmp
}

func newDumperWithTee(src io.Reader) (*dumper, io.Reader) {
	if *dumpAudio == "" {
		return nil, src
	}
	pipeRd, pipeWr := io.Pipe()
	teeRd := io.TeeReader(src, pipeWr)
	return newDumper(teeRd), pipeRd
}

func (d *dumper) Close() {
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

func (d *dumper) addTo(buf []int16) {
	if !d.playing {
		return
	}
	addBuf := make([]int16, len(buf))
	binary.Read(d.reader, binary.LittleEndian, addBuf)
	for i, s := range addBuf {
		buf[i] += int16(d.volume * float64(s))
	}
	d.played += len(buf) / 2
}

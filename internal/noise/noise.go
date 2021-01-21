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

package noise

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"

	"github.com/divVerent/aaaaaa/internal/vfs"
)

var (
	noiseVolume = flag.Float64("noise_volume", 0.5, "noise volume (0..1)")
)

const (
	shrinkagePerFrame = 0.05
)

var (
	amount float64 = 0.0
	noise  *audio.Player
)

func Init() {
	data, err := vfs.Load("sounds", "stereonoise.ogg")
	if err != nil {
		log.Panicf("Could not load stereonoise: %v")
	}
	defer data.Close()
	stream, err := vorbis.Decode(audio.CurrentContext(), data)
	if err != nil {
		log.Panicf("Could not start decoding stereonosie: %v", err)
	}
	decoded, err := ioutil.ReadAll(stream)
	if err != nil {
		log.Panicf("Could not decode stereonoise: %v", err)
	}
	loop := audio.NewInfiniteLoop(bytes.NewReader(decoded), int64(len(decoded)))
	noise, err = audio.NewPlayer(audio.CurrentContext(), loop)
	if err != nil {
		log.Panicf("could not start playing noise: %v", err)
	}
}

func Update() {
	if amount > 0 {
		noise.SetVolume(amount * *noiseVolume)
		noise.Play()
	} else {
		noise.Pause()
	}
	amount -= shrinkagePerFrame
}

func Set(noise float64) {
	if noise > 1 {
		noise = 1
	}
	if noise > amount {
		amount = noise
	}
}

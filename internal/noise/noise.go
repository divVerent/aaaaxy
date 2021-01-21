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
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/audio"

	m "github.com/divVerent/aaaaaa/internal/math"
)

var (
	noiseVolume = flag.Float64("noise_volume", 0.5, "noise volume (0..1)")
)

const (
	shrinkagePerFrame = 0.05
	noiseSize         = 65536
	lowpass           = 256
	stereoOffset      = 8
	volumeFactor      = 0.5
)

var (
	amount float64 = 0.0
	noise  *audio.Player
)

func Init() {
	lowpassedData := make([]int16, noiseSize)
	randData := make([]int, noiseSize)
	for i := 0; i < noiseSize; i++ {
		randData[i] = rand.Intn(65536) - 32768
	}
	volumeAdj := volumeFactor / math.Sqrt(lowpass)
	for i := 0; i < noiseSize; i++ {
		sum := 0
		for j := 0; j < lowpass; j++ {
			sum += randData[m.Mod(i+j, noiseSize)]
		}
		lowpassed := int(math.Floor(float64(sum)*volumeAdj + 0.5))
		if lowpassed < -32768 {
			lowpassed = -32768
		}
		if lowpassed > 32767 {
			lowpassed = 32767
		}
		lowpassedData[i] = int16(lowpassed)
	}
	leData := make([]byte, 4*noiseSize)
	for i := 0; i < noiseSize; i++ {
		left := lowpassedData[i]
		right := lowpassedData[m.Mod(i+stereoOffset, noiseSize)]
		leData[4*i] = byte(left & 0xFF)
		leData[4*i+1] = byte(left >> 8)
		leData[4*i+2] = byte(right & 0xFF)
		leData[4*i+3] = byte(right >> 8)
	}
	randBuf := bytes.NewReader(leData)
	randLoop := audio.NewInfiniteLoop(randBuf, int64(len(leData)))
	var err error
	noise, err = audio.NewPlayer(audio.CurrentContext(), randLoop)
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

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

package aaaaaa

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/audiowrap"
	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/flag"
)

var (
	dumpVideo = flag.String("dump_video", "", "filename prefix to dump game frames to")
	dumpAudio = flag.String("dump_audio", "", "filename to dump game audio to")
)

var (
	dumpFrameCount = 0
	dumpVideoFile  *os.File
	dumpAudioFile  *os.File
	pngEncoder     = png.Encoder{
		CompressionLevel: png.BestSpeed,
	}
)

func initDumping() error {
	if *dumpAudio != "" {
		var err error
		dumpAudioFile, err = os.Create(*dumpAudio)
		if err != nil {
			return fmt.Errorf("could not initialize audio dump: %v", err)
		}
		audiowrap.InitDumping()
	}

	if *dumpVideo != "" {
		var err error
		dumpVideoFile, err = os.Create(*dumpVideo)
		if err != nil {
			return fmt.Errorf("could not initialize video dump: %v", err)
		}
	}

	return nil
}

func dumping() bool {
	return dumpAudioFile != nil || dumpVideoFile != nil
}

func dumpFrame(screen *ebiten.Image) {
	if !dumping() {
		return
	}
	dumpFrameCount++
	if dumpVideoFile != nil {
		err := pngEncoder.Encode(dumpVideoFile, screen)
		if err != nil {
			log.Printf("Failed to encode video - expect corruption: %v", err)
			dumpVideoFile.Close()
			dumpVideoFile = nil
		}
	}
	if dumpAudioFile != nil {
		err := audiowrap.DumpFrame(dumpAudioFile, time.Duration(dumpFrameCount)*time.Second/engine.GameTPS)
		if err != nil {
			log.Printf("Failed to encode audio - expect corruption: %v", err)
			dumpAudioFile.Close()
			dumpAudioFile = nil
		}
	}
}

func finishDumping() {
	if !dumping() {
		return
	}
	if dumpAudioFile != nil {
		err := dumpAudioFile.Close()
		if err != nil {
			log.Printf("Failed to close audio - expect corruption: %v", err)
		}
		dumpAudioFile = nil
	}
	if dumpVideoFile != nil {
		err := dumpVideoFile.Close()
		if err != nil {
			log.Printf("Failed to close video - expect corruption: %v", err)
		}
		dumpVideoFile = nil
	}
	log.Print("Media has been dumped.")
	log.Print("To convert to something uploadable, run:")
	inputs := []string{}
	settings := []string{}
	if *dumpAudio != "" {
		inputs = append(inputs, fmt.Sprintf("-f s16le -ac 2 -ar %d  -i '%s'", audiowrap.Rate(), strings.ReplaceAll(*dumpAudio, "'", "'\\''")))
		settings = append(settings, "-codec:a aac -b:a 128k")
	}
	if *dumpVideo != "" {
		inputs = append(inputs, fmt.Sprintf("-f png_pipe -r %d -i '%s'", engine.GameTPS, strings.ReplaceAll(*dumpVideo, "'", "'\\''")))
		settings = append(settings, "-codec:v libx264 -profile:v high444 -preset:v fast -crf:v 10 -preset:v fast -vf premultiply=inplace=1,scale=1280:720:flags=neighbor,scale=1920:1080")
	}
	log.Printf("ffmpeg %s %s -vsync vfr video.mp4", strings.Join(inputs, " "), strings.Join(settings, " "))
}

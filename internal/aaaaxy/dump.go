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

package aaaaxy

import (
	"fmt"
	"github.com/divVerent/aaaaxy/internal/log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/audiowrap"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	m "github.com/divVerent/aaaaxy/internal/math"
)

var (
	dumpVideo            = flag.String("dump_video", "", "filename prefix to dump game frames to")
	dumpVideoFpsDivisor  = flag.Int("dump_video_fps_divisor", 1, "frame rate divisor (try 2 for faster dumping)")
	dumpAudio            = flag.String("dump_audio", "", "filename to dump game audio to")
	cheatDumpSlowAndGood = flag.Bool("cheat_dump_slow_and_good", false, "non-realtime video dumping (slows down the game, thus considered a cheat))")
)

var (
	dumpFrameCount = int64(0)
	dumpVideoFile  *os.File
	dumpAudioFile  *os.File
)

const (
	dumpVideoFrameSize = engine.GameWidth * engine.GameHeight * 4
)

var (
	dumpVideoFrame int64 = 0
	dumpVideoWg    sync.WaitGroup
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

func slowDumping() bool {
	return dumping() && *cheatDumpSlowAndGood
}

func dumpFrameThenReturnTo(screen *ebiten.Image, to chan *ebiten.Image, frames int) {
	if !dumping() || frames == 0 {
		to <- screen
		return
	}
	if dumpVideoFile != nil {
		dumpVideoFrameBegin := dumpFrameCount / int64(*dumpVideoFpsDivisor)
		dumpFrameCount += int64(frames)
		dumpVideoFrameEnd := dumpFrameCount / int64(*dumpVideoFpsDivisor)
		cnt := dumpVideoFrameEnd - dumpVideoFrameBegin
		if cnt > 0 {
			if cnt > 1 {
				log.Infof("video dump: %v frames dropped", cnt-1)
			}
			dumpVideoWg.Add(1)
			dumpPixelsRGBA(screen, func(pix []byte, err error) {
				to <- screen
				if err == nil {
					for i := dumpVideoFrameBegin; i <= dumpVideoFrameEnd; i++ {
						_, err = dumpVideoFile.WriteAt(pix, i*dumpVideoFrameSize)
						if err != nil {
							break
						}
					}
				}
				if err != nil {
					log.Errorf("Failed to encode video - expect corruption: %v", err)
					// dumpVideoFile.Close()
					// dumpVideoFile = nil
				}
				dumpVideoWg.Done()
			})
		} else {
			// log.Infof("video dump: frame skipped")
			to <- screen
		}
	} else {
		to <- screen
	}
	if dumpAudioFile != nil {
		err := audiowrap.DumpFrame(dumpAudioFile, time.Duration(dumpFrameCount)*time.Second/engine.GameTPS)
		if err != nil {
			log.Errorf("Failed to encode audio - expect corruption: %v", err)
			dumpAudioFile.Close()
			dumpAudioFile = nil
		}
	}
}

func ffmpegCommand(audio, video, output, screenFilter string) string {
	var pre string
	inputs := []string{}
	settings := []string{}
	// Video first, so we can refer to the video stream as [0:v] for sure.
	if video != "" {
		fps := float64(engine.GameTPS) / float64(*dumpVideoFpsDivisor)
		inputs = append(inputs, fmt.Sprintf("-f rawvideo -pixel_format rgba -video_size %dx%d -r %v -i '%s'", engine.GameWidth, engine.GameHeight, fps, strings.ReplaceAll(video, "'", "'\\''")))
		var maybeFilterComplex string
		switch screenFilter {
		case "linear":
			maybeFilterComplex = " -filter_complex '[0:v]premultiply=inplace=1,scale=1920:1080'"
		case "linear2x":
			// Note: the two step upscale simulates the effect of the linear2xcrt shader.
			maybeFilterComplex = " -filter_complex '[0:v]premultiply=inplace=1,scale=1280:720:flags=neighbor,scale=1920:1080'"
		case "linear2xcrt":
			// For 3x scale, pattern is: 1 (1-2/3*f) 1.
			// darkened := m.Rint(255 * (1.0 - 2.0/3.0**screenFilterScanLines))
			// pre = fmt.Sprintf("echo 'P2 1 3 255 %d 255 %d' | convert -size 1920x1080 TILE:PNM:- scanlines.png; ", darkened, darkened)
			// Then second scale is to 1920:1080.
			// But for the lens correction, we gotta do better.
			// For 6x scale, pattern is: (1-5/6*f) (1-3/6*f) (1-1/6*f) (1-1/6*f) (1-3/6*f) (1-5/6*f).
			pre = fmt.Sprintf("echo 'P2 1 6 255 %d %d %d %d %d %d' | convert -size 3840:2160 TILE:PNM:- scanlines.png; ",
				m.Rint(255*(1.0-5.0/6.0**screenFilterScanLines)),
				m.Rint(255*(1.0-3.0/6.0**screenFilterScanLines)),
				m.Rint(255*(1.0-1.0/6.0**screenFilterScanLines)),
				m.Rint(255*(1.0-1.0/6.0**screenFilterScanLines)),
				m.Rint(255*(1.0-3.0/6.0**screenFilterScanLines)),
				m.Rint(255*(1.0-5.0/6.0**screenFilterScanLines)))
			maybeFilterComplex = fmt.Sprintf(" -filter_complex '[0:v]premultiply=inplace=1,scale=1280:720:flags=neighbor,scale=3840:2160,format=gbrp[scaled]; movie=scanlines.png,format=gbrp[scanlines]; [scaled][scanlines]blend=all_mode=multiply,lenscorrection=i=bilinear:k1=%f:k2=%f'", crtK1(), crtK2())
		case "simple", "nearest":
			maybeFilterComplex = " -filter_complex '[0:v]premultiply=inplace=1,scale=1920:1080:flags=neighbor'"
		case "":
			maybeFilterComplex = ""
		}
		// Note: using high quality, fast settings and many keyframes
		// as the assumption is that the output file will be further edited.
		// Note: disabling 8x8 DCT here as some older FFmpeg versions -
		// or even newer versions with decoding options changed for compatibility,
		// if the video file has also been losslessly cut -
		// have trouble decoding that.
		settings = append(settings, "-codec:v libx264 -profile:v high444 -preset:v fast -crf:v 10 -8x8dct:v 0 -keyint_min 10 -g 60"+maybeFilterComplex)
	}
	if audio != "" {
		inputs = append(inputs, fmt.Sprintf("-f s16le -ac 2 -ar %d  -i '%s'", audiowrap.Rate(), strings.ReplaceAll(audio, "'", "'\\''")))
		settings = append(settings, "-codec:a aac -b:a 128k")
	}
	return fmt.Sprintf("%sffmpeg %s %s -vsync vfr %s", pre, strings.Join(inputs, " "), strings.Join(settings, " "), strings.ReplaceAll(output, "'", "'\\''"))
}

func finishDumping() {
	if !dumping() {
		return
	}
	if dumpAudioFile != nil {
		err := dumpAudioFile.Close()
		if err != nil {
			log.Errorf("Failed to close audio - expect corruption: %v", err)
		}
		dumpAudioFile = nil
	}
	if dumpVideoFile != nil {
		dumpVideoWg.Wait()
		err := dumpVideoFile.Close()
		if err != nil {
			log.Errorf("Failed to close video - expect corruption: %v", err)
		}
		dumpVideoFile = nil
	}
	log.Infof("Media has been dumped.")
	log.Infof("To create a preview file (DO NOT UPLOAD):")
	log.Infof("  " + ffmpegCommand(*dumpAudio, *dumpVideo, "video-preview.mp4", ""))
	if *dumpVideo != "" {
		if *screenFilter != "linear2xcrt" {
			log.Infof("With current settings (1080p, MEDIUM QUALITY):")
			log.Infof("  " + ffmpegCommand(*dumpAudio, *dumpVideo, "video-medium.mp4", *screenFilter))
		}
		log.Infof("Preferred for uploading (4K, GOOD QUALITY):")
		log.Infof("  " + ffmpegCommand(*dumpAudio, *dumpVideo, "video-high.mp4", "linear2xcrt"))
	}
}

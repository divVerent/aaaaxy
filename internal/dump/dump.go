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

package dump

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/atexit"
	"github.com/divVerent/aaaaxy/internal/audiowrap"
	"github.com/divVerent/aaaaxy/internal/demo"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/namedpipe"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

var (
	dumpVideo               = flag.String("dump_video", "", "filename prefix to dump game frames to")
	dumpVideoFpsDivisor     = flag.Int("dump_video_fps_divisor", 1, "frame rate divisor (try 2 for faster dumping)")
	dumpAudio               = flag.String("dump_audio", "", "filename to dump game audio to")
	dumpMedia               = flag.String("dump_media", "", "filename to dump game media to; exclusive with dump_video and dump_audio; when not changing any dump_*_settings, this should have a .mkv, .mov, .avi or .nut extension")
	dumpVideoCodecSettings  = flag.String("dump_video_codec_settings", "-codec:v mjpeg -q:v 4", "FFmpeg settings for video encoding; set to \"\" to disable the video stream for -dump_media")
	dumpAudioCodecSettings  = flag.String("dump_audio_codec_settings", "-codec:a pcm_s16le", "FFmpeg settings for audio encoding; set to \"\" to disable the audio stream for -dump_media")
	dumpMediaFormatSettings = flag.String("dump_media_format_settings", "-vsync vfr", "FFmpeg flags for muxing")
	cheatDumpSlowAndGood    = flag.Bool("cheat_dump_slow_and_good", false, "non-realtime video dumping (slows down the game, thus considered a cheat))")
	dumpMediaFrameTimeout   = flag.Duration("dump_media_frame_timeout", 300*time.Second, "maximum processing time per frame; after this time it is assumed that ffmpeg died and dumping ends")
)

type Params struct {
	FPSDivisor            int
	ScreenFilter          string
	ScreenFilterScanLines float64
	CRTK1                 float64
	CRTK2                 float64
}

type WriteCloserAt interface {
	io.Writer
	io.WriterAt
	io.Closer
}

var (
	frameCount   = int64(0)
	videoWriter  WriteCloserAt
	audioWriter  WriteCloserAt
	videoPipe    *namedpipe.Fifo
	audioPipe    *namedpipe.Fifo
	mediaCmd     *exec.Cmd
	mediaCmdDone chan struct{}
	params       Params
)

const (
	dumpVideoFrameSize = engine.GameWidth * engine.GameHeight * 4
)

var (
	dumpVideoWg sync.WaitGroup
)

func InitEarly(p Params) error {
	params = p

	if *dumpMedia != "" {
		if *dumpVideo != "" || *dumpAudio != "" {
			return errors.New("-dump_media is mutually exclusive with -dump_video/-dump_audio")
		}
		if *dumpAudioCodecSettings == "" && *dumpVideoCodecSettings == "" {
			return errors.New("not both of -dump_audio_codec_settings and -dump_video_codec_settings may be empty - we need at least one stream")
		}
		var err error
		if *dumpAudioCodecSettings != "" {
			audioPipe, err = namedpipe.New("aaaaxy-audio", 120, 4*96000, *dumpMediaFrameTimeout)
			if err != nil {
				return fmt.Errorf("could not create audio pipe: %w", err)
			}
			audioWriter = namedpipe.NewWriteCloserAt(audioPipe)
			audiowrap.InitDumping()
		}
		if *dumpVideoCodecSettings != "" {
			videoPipe, err = namedpipe.New("aaaaxy-video", 120, dumpVideoFrameSize, *dumpMediaFrameTimeout)
			if err != nil {
				return fmt.Errorf("could not create video pipe: %w", err)
			}
			videoWriter = namedpipe.NewWriteCloserAt(videoPipe)
		}
	}

	if *dumpAudio != "" {
		var err error
		audioWriter, err = vfs.OSCreate(vfs.WorkDir, *dumpAudio)
		if err != nil {
			return fmt.Errorf("could not initialize audio dump: %w", err)
		}
		audiowrap.InitDumping()
	}

	if *dumpVideo != "" {
		var err error
		videoWriter, err = vfs.OSCreate(vfs.WorkDir, *dumpVideo)
		if err != nil {
			return fmt.Errorf("could not initialize video dump: %w", err)
		}
	}

	return nil
}

func InitLate() error {
	if *dumpMedia != "" {
		audioPath := ""
		if audioPipe != nil {
			audioPath = audioPipe.Path()
		}
		videoPath := ""
		if videoPipe != nil {
			videoPath = videoPipe.Path()
		}
		cmdLine, _, err := ffmpegCommand(audioPath, videoPath, *dumpMedia, params.ScreenFilter)
		if err != nil {
			return err
		}
		mediaCmd := exec.Command(cmdLine[0], cmdLine[1:]...)
		mediaCmd.Stdout = os.Stdout
		mediaCmd.Stderr = os.Stderr
		err = mediaCmd.Start()
		if err != nil {
			return fmt.Errorf("could not launch FFmpeg: %w", err)
		}
		mediaCmdDone = make(chan struct{})
		go func() {
			err := mediaCmd.Wait()
			if err != nil {
				log.Fatalf("FFmpeg died: %v", err)
			}
			close(mediaCmdDone)
		}()
	}

	return nil
}

func Active() bool {
	return audioWriter != nil || videoWriter != nil
}

func Slow() bool {
	return Active() && (*cheatDumpSlowAndGood || demo.Playing())
}

func ProcessFrameThenReturnTo(screen *ebiten.Image, to chan *ebiten.Image, frames int) {
	if !Active() || frames == 0 {
		to <- screen
		return
	}
	prevFrameCount := frameCount
	frameCount += int64(frames)
	if videoWriter != nil {
		dumpVideoFrameBegin := prevFrameCount / int64(*dumpVideoFpsDivisor)
		dumpVideoFrameEnd := frameCount / int64(*dumpVideoFpsDivisor)
		cnt := dumpVideoFrameEnd - dumpVideoFrameBegin
		if cnt > 0 {
			if cnt > 1 {
				log.Infof("video dump: %v frames dropped", cnt-1)
			}
			dumpVideoWg.Add(1)
			dumpPixelsRGBA(screen, func(pix []byte, err error) {
				to <- screen
				if err == nil {
					for i := dumpVideoFrameBegin; i < dumpVideoFrameEnd; i++ {
						_, err = videoWriter.WriteAt(pix, i*dumpVideoFrameSize)
						if err != nil {
							break
						}
					}
				}
				if err != nil {
					log.Errorf("failed to encode video - expect corruption: %v", err)
					// videoWriter.Close()
					// videoWriter = nil
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
	if audioWriter != nil {
		err := audiowrap.DumpFrame(audioWriter, time.Duration(frameCount)*time.Second/engine.GameTPS)
		if err != nil {
			log.Errorf("failed to encode audio - expect corruption: %v", err)
			audioWriter.Close()
			audioWriter = nil
		}
	}
}

func ffmpegCommand(audio, video, output, screenFilter string) ([]string, string, error) {
	precmd := ""
	inputs := []string{}
	settings := []string{"-y"}
	// Video first, so we can refer to the video stream as [0:v] for sure.
	if video != "" {
		fps := float64(engine.GameTPS) / (float64(params.FPSDivisor) * float64(*dumpVideoFpsDivisor))
		inputs = append(inputs, "-f", "rawvideo", "-pixel_format", "rgba", "-video_size", fmt.Sprintf("%dx%d", engine.GameWidth, engine.GameHeight), "-r", fmt.Sprint(fps), "-i", video)
		filterComplex := "[0:v]premultiply=inplace=1,format=gbrp[lowres]; "
		switch screenFilter {
		case "linear":
			filterComplex += "[lowres]scale=1920:1080"
		case "linear2x":
			// Note: the two step upscale simulates the effect of the linear2xcrt shader.
			// "simple" does the same as "linear2x" if the screen res is exactly 1080p.
			filterComplex += "[lowres]scale=1280:720:flags=neighbor,scale=1920:1080"
		case "linear2xcrt":
			// For 3x scale, pattern is: 1 (1-2/3*f) 1.
			// darkened := m.Rint(255 * (1.0 - 2.0/3.0**screenFilterScanLines))
			// pnm := fmt.Sprintf("P2 1 3 255 %d 255 %d", darkened, darkened)
			// Then second scale is to 1920:1080.
			// But for the lens correction, we gotta do better.
			// For 6x scale, pattern is: (1-5/6*f) (1-3/6*f) (1-1/6*f) (1-1/6*f) (1-3/6*f) (1-5/6*f).
			pnmHeader1 := []byte("P2\n")
			pnmHeader2 := []byte("1 2160 255\n")
			pnmLine := []byte(fmt.Sprintf("%d %d %d %d %d %d\n",
				m.Rint(255*(1.0-5.0/6.0*params.ScreenFilterScanLines)),
				m.Rint(255*(1.0-3.0/6.0*params.ScreenFilterScanLines)),
				m.Rint(255*(1.0-1.0/6.0*params.ScreenFilterScanLines)),
				m.Rint(255*(1.0-1.0/6.0*params.ScreenFilterScanLines)),
				m.Rint(255*(1.0-3.0/6.0*params.ScreenFilterScanLines)),
				m.Rint(255*(1.0-5.0/6.0*params.ScreenFilterScanLines))))
			tempFile, err := os.CreateTemp("", "aaaaxy-*")
			if err != nil {
				return nil, "", err
			}
			atexit.Delete(tempFile.Name())
			_, err = tempFile.Write(pnmHeader1)
			if err != nil {
				return nil, "", err
			}
			_, err = tempFile.Write(pnmHeader2)
			if err != nil {
				return nil, "", err
			}
			for range make([]struct{}, 360) {
				_, err = tempFile.Write(pnmLine)
				if err != nil {
					return nil, "", err
				}
			}
			err = tempFile.Close()
			if err != nil {
				return nil, "", err
			}
			precmd = fmt.Sprintf("{ echo '%s'; echo '%s'; for i in `seq 1 360`; do echo '%s'; done } > '%s'; ", pnmHeader1[:len(pnmHeader1)-1], pnmHeader2[:len(pnmHeader2)-1], pnmLine[:len(pnmLine)-1], tempFile.Name())
			inputs = append(inputs, "-f", "pgm_pipe", "-i", tempFile.Name())
			filterComplex += fmt.Sprintf("[lowres]scale=1280:720:flags=neighbor,scale=3840:2160[scaled]; [1:v]scale=3840:2160:flags=neighbor,format=gbrp[scanlines]; [scaled][scanlines]blend=all_mode=multiply,lenscorrection=i=bilinear:k1=%f:k2=%f", params.CRTK1, params.CRTK2)
		case "nearest":
			filterComplex += "[lowres]scale=1920:1080:flags=neighbor"
		case "":
			filterComplex += "[lowres]copy"
		}
		// Note: using high quality, fast settings and many keyframes
		// as the assumption is that the output file will be further edited.
		// Note: disabling 8x8 DCT here as some older FFmpeg versions -
		// or even newer versions with decoding options changed for compatibility,
		// if the video file has also been losslessly cut -
		// have trouble decoding that.
		if *dumpVideoCodecSettings != "" {
			settings = append(settings, strings.Split(*dumpVideoCodecSettings, " ")...)
		}
		settings = append(settings, "-filter_complex", filterComplex)
	}
	if audio != "" {
		inputs = append(inputs, "-f", "s16le", "-ac", "2", "-ar", fmt.Sprint(audiowrap.SampleRate()), "-i", audio)
		if *dumpAudioCodecSettings != "" {
			settings = append(settings, strings.Split(*dumpAudioCodecSettings, " ")...)
		}
	}
	if *dumpMediaFormatSettings != "" {
		settings = append(settings, strings.Split(*dumpMediaFormatSettings, " ")...)
	}
	cmd := []string{"ffmpeg"}
	cmd = append(cmd, inputs...)
	cmd = append(cmd, settings...)
	cmd = append(cmd, output)
	return cmd, precmd, nil
}

func printCommand(cmd []string) string {
	r := []string{}
	for _, arg := range cmd {
		r = append(r, "'"+strings.ReplaceAll(arg, "'", "'\\''")+"'")
	}
	return strings.Join(r, " ")
}

func Finish() error {
	if !Active() {
		return nil
	}
	if videoWriter != nil {
		dumpVideoWg.Wait()
	}
	// Closing audio and video file concurrently, which helps in case they're pipes, as it's unclear in which state FFmpeg tries to read them.
	var wg sync.WaitGroup
	var videoErr, audioErr error
	if audioWriter != nil {
		wg.Add(1)
		go func() {
			audioErr = audioWriter.Close()
			audioWriter = nil
			wg.Done()
		}()
	}
	if videoWriter != nil {
		wg.Add(1)
		go func() {
			videoErr = videoWriter.Close()
			videoWriter = nil
			wg.Done()
		}()
	}
	wg.Wait()
	if audioErr != nil {
		return fmt.Errorf("failed to close audio - expect corruption: %w", audioErr)
	}
	if videoErr != nil {
		return fmt.Errorf("failed to close video - expect corruption: %w", videoErr)
	}
	if mediaCmd != nil {
		log.Infof("waiting for FFmpeg to exit...")
		<-mediaCmdDone
		mediaCmdDone = nil
	}
	log.Infof("media has been dumped")
	if *dumpAudio != "" || *dumpVideo != "" {
		log.Infof("to create a preview file (DO NOT UPLOAD):")
		cmd, precmd, err := ffmpegCommand(*dumpAudio, *dumpVideo, "video-preview.mkv", "")
		if err != nil {
			return err
		}
		log.Infof("  %v%v", precmd, printCommand(cmd))
		if params.ScreenFilter != "linear2xcrt" {
			log.Infof("with current settings (1080p, MEDIUM QUALITY):")
			cmd, precmd, err := ffmpegCommand(*dumpAudio, *dumpVideo, "video-medium.mkv", params.ScreenFilter)
			if err != nil {
				return err
			}
			log.Infof("  %v%v", precmd, printCommand(cmd))
		}
		log.Infof("preferred for uploading (4K, GOOD QUALITY):")
		cmd, precmd, err = ffmpegCommand(*dumpAudio, *dumpVideo, "video-high.mkv", "linear2xcrt")
		if err != nil {
			return err
		}
		log.Infof("  %v%v", precmd, printCommand(cmd))
	}
	return nil
}

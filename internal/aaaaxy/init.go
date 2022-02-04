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
	"encoding/json"
	"fmt"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/audiowrap"
	"github.com/divVerent/aaaaxy/internal/credits"
	"github.com/divVerent/aaaaxy/internal/demo"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/image"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/noise"
	"github.com/divVerent/aaaaxy/internal/sound"
	"github.com/divVerent/aaaaxy/internal/splash"
	"github.com/divVerent/aaaaxy/internal/timing"
	"github.com/divVerent/aaaaxy/internal/version"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

var (
	vsync                 = flag.Bool("vsync", true, "enable waiting for vertical synchronization")
	fullscreen            = flag.Bool("fullscreen", true, "enable fullscreen mode")
	windowScaleFactor     = flag.Float64("window_scale_factor", 0, "window scale factor in device pixels per game pixel (0 means auto integer scaling)")
	runnableWhenUnfocused = flag.Bool("runnable_when_unfocused", flag.SystemDefault(map[string]interface{}{"js/*": true, "*/*": false}).(bool), "keep running the game even when not focused")
	dumpLoadingFractions  = flag.String("dump_loading_fractions", "", "file name to dump actual loading fractions to")
	debugJustInit         = flag.Bool("debug_just_init", false, "just init everything, then quit right away")
	fpsDivisor            = flag.Int("fps_divisor", 1, "framerate divisor (use on very low systems, but this may make the game unwinnable or harder as it restricts input; must be a divisor of "+fmt.Sprint(engine.GameTPS))
)

func LoadConfig() (*flag.Config, error) {
	return engine.LoadConfig()
}

func setWindowSize() {
	logicalF := *windowScaleFactor
	log.Infof("requested logical scale factor: %v", logicalF)
	dscale := ebiten.DeviceScaleFactor()
	log.Infof("device scale factor: %v", dscale)
	var physicalF float64
	if logicalF <= 0 {
		screenw, screenh := ebiten.ScreenSizeInFullscreen()
		log.Infof("screen size: %vx%v", screenw, screenh)
		// Reserve 128 device independent pixels for system controls.
		maxw, maxh := screenw-128, screenh-128
		log.Infof("max size: %vx%v", maxw, maxh)
		// Compute max scaling factors.
		maxwf, maxhf := float64(maxw)*dscale/engine.GameWidth, float64(maxh)*dscale/engine.GameHeight
		log.Infof("max physical scale factors: %v, %v", maxwf, maxhf)
		physicalF = math.Min(maxwf, maxhf)
	} else {
		physicalF = logicalF * dscale
	}
	log.Infof("requested physical scale factor: %v", physicalF)
	// Make output pixels an integer multiple of input pixels (looks better).
	physicalF = math.Floor(physicalF)
	if physicalF < 1 {
		physicalF = 1
	}
	log.Infof("chosen physical scale factor: %v", physicalF)
	// Convert back to logical scale factor as ebiten needs that.
	logicalF = physicalF / dscale
	log.Infof("chosen logical pixel scale factor: %v", logicalF)
	w, h := m.Rint(engine.GameWidth*logicalF), m.Rint(engine.GameHeight*logicalF)
	log.Infof("chosen window size: %vx%v", w, h)
	ebiten.SetWindowSize(w, h)
}

func (g *Game) InitEbiten() error {
	// Ensure fps divisor is valid. We can only do integer TPS.
	if *fpsDivisor < 1 || engine.GameTPS%*fpsDivisor != 0 {
		*fpsDivisor = 1
	}

	ebiten.SetFullscreen(*fullscreen)
	ebiten.SetInitFocused(true)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetScreenTransparent(false)
	ebiten.SetVsyncEnabled(*vsync)
	ebiten.SetWindowDecorated(true)
	ebiten.SetWindowResizable(true)
	setWindowSize()
	ebiten.SetWindowTitle("AAAAXY")

	// Initialize some stuff that is needed early.
	err := vfs.Init()
	if err != nil {
		return fmt.Errorf("could not initialize VFS: %v", err)
	}
	err = version.Init()
	if err != nil {
		return fmt.Errorf("could not initialize version: %v", err)
	}
	err = demo.Init()
	if err != nil {
		return fmt.Errorf("could not initialize demo: %v", err)
	}
	err = initDumping()
	if err != nil {
		return fmt.Errorf("could not initialize dumping: %v", err)
	}
	err = font.Init()
	if err != nil {
		return fmt.Errorf("could not initialize fonts: %v", err)
	}

	// When dumping video or benchmarking, do precisely one render frame per update.
	if slowDumping() || demo.Timedemo() {
		ebiten.SetMaxTPS(ebiten.UncappedTPS)
	} else {
		ebiten.SetMaxTPS(engine.GameTPS / *fpsDivisor)
	}

	// Pause when unfocused, except when recording demos.
	ebiten.SetRunnableOnUnfocused(*runnableWhenUnfocused || (demo.Playing() && dumping()))

	return nil
}

type initState struct {
	splash.State
	started bool
	done    bool
}

func (g *Game) provideLoadingFractions() error {
	j, err := vfs.Load("splash", "loading_fractions.json")
	if err != nil {
		return err
	}
	defer j.Close()
	var m map[string]float64
	err = json.NewDecoder(j).Decode(&m)
	if err != nil {
		return err
	}
	g.init.ProvideFractions(m)
	return nil
}

func (g *Game) InitStep() error {
	if !g.init.started {
		g.init.started = true
		err := g.provideLoadingFractions()
		if err != nil {
			log.Errorf("could not provide loading fractions: %v", err)
		}
	}
	status, err := g.init.Enter("precaching credits", "could not precache credits", splash.Single(credits.Precache))
	if status != splash.Continue {
		return err
	}
	status, err = g.init.Enter("initializing audio", "could not initialize audio", splash.Single(audiowrap.Init))
	if status != splash.Continue {
		return err
	}
	status, err = g.init.Enter("initializing noise", "could not initialize noise", splash.Single(noise.Init))
	if status != splash.Continue {
		return err
	}
	status, err = g.init.Enter("precaching sounds", "could not precache sounds", sound.Precache)
	if status != splash.Continue {
		return err
	}
	status, err = g.init.Enter("precaching images", "could not precache images", splash.Single(image.Precache))
	if status != splash.Continue {
		return err
	}
	status, err = g.init.Enter("initializing input", "could not initialize input", splash.Single(input.Init))
	if status != splash.Continue {
		return err
	}
	status, err = g.init.Enter("precaching engine", "could not precache engine", engine.Precache)
	if status != splash.Continue {
		return err
	}
	if *dumpLoadingFractions != "" {
		f, err := os.Create(*dumpLoadingFractions)
		if err != nil {
			return fmt.Errorf("could not open loading fractions file: %v", err)
		}
		j := json.NewEncoder(f)
		j.SetIndent("", "\t")
		err = j.Encode(g.init.ToFractions())
		if err != nil {
			return fmt.Errorf("could not encode to loading fractions file: %v", err)
		}
		err = f.Close()
		if err != nil {
			return fmt.Errorf("could not close loading fractions file: %v", err)
		}
	}
	if *debugJustInit {
		log.Errorf("requested early termination via --debug_just_init")
		return RegularTermination
	}
	g.init.done = true
	return nil
}

func (g *Game) InitFull() error {
	for !g.init.done {
		err := g.InitStep()
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) BeforeExit() error {
	timing.PrintReport()
	err := finishDumping()
	if err != nil {
		return fmt.Errorf("could not finish dumping: %v", err)
	}
	err = demo.BeforeExit()
	if err != nil {
		return fmt.Errorf("could not finalize demo: %v", err)
	}
	return nil
}

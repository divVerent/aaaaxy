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
	"github.com/divVerent/aaaaxy/internal/dump"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/exitstatus"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/image"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/locale/initlocale"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/noise"
	"github.com/divVerent/aaaaxy/internal/palette"
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
	runnableWhenUnfocused = flag.Bool("runnable_when_unfocused", flag.SystemDefault(map[string]bool{
		// Focus didn't quite work well on JS. TODO: try testing again later.
		"js/*": true,
		"*/*":  false,
	}), "keep running the game even when not focused")
	dumpLoadingFractions = flag.String("dump_loading_fractions", "", "file name to dump actual loading fractions to")
	debugJustInit        = flag.Bool("debug_just_init", false, "just init everything, then quit right away")
	fpsDivisor           = flag.Int("fps_divisor", 1, "framerate divisor (use on very low systems, but this may make the game unwinnable or harder as it restricts input; must be a divisor of "+fmt.Sprint(engine.GameTPS))
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
	// Convert back to logical scale factor as Ebitengine needs that.
	logicalF = physicalF / dscale
	log.Infof("chosen logical pixel scale factor: %v", logicalF)
	w, h := m.Rint(engine.GameWidth*logicalF), m.Rint(engine.GameHeight*logicalF)
	log.Infof("chosen window size: %vx%v", w, h)
	ebiten.SetWindowSize(w, h)
}

// NOTE: This function only runs on desktop systems.
// On mobile, we instead run InitEarly only.
func (g *Game) InitEbitengine() error {
	ebiten.SetInitFocused(true)
	ebiten.SetScreenTransparent(false)
	ebiten.SetWindowDecorated(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	setWindowSize()
	return g.InitEarly()
}

// InitEarly is the beginning of our initialization,
// and may take place in the first frame
// if there is no way to run this before the main loop (e.g. on mobile).
func (g *Game) InitEarly() error {
	ebiten.SetFullscreen(*fullscreen)
	ebiten.SetScreenClearedEveryFrame(false)
	if *vsync {
		ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	} else {
		ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	}
	ebiten.SetWindowTitle("AAAAXY")

	// Ensure fps divisor is valid. We can only do integer TPS.
	if *fpsDivisor < 1 || engine.GameTPS%*fpsDivisor != 0 {
		*fpsDivisor = 1
	}

	// Initialize some stuff that is needed early.
	err := vfs.Init()
	if err != nil {
		return fmt.Errorf("could not initialize VFS: %w", err)
	}
	err = initlocale.Init()
	if err != nil {
		return fmt.Errorf("could not initialize locale: %w", err)
	}
	err = version.Init()
	if err != nil {
		return fmt.Errorf("could not initialize version: %w", err)
	}
	err = demo.Init()
	if err != nil {
		return fmt.Errorf("could not initialize demo: %w", err)
	}
	err = dump.InitEarly(dump.Params{
		FPSDivisor:            *fpsDivisor,
		ScreenFilter:          *screenFilter,
		ScreenFilterScanLines: *screenFilterScanLines,
		CRTK1:                 crtK1(),
		CRTK2:                 crtK2(),
	})
	if err != nil {
		return fmt.Errorf("could not preinitialize dumping: %w", err)
	}

	// Load images with the right palette from the start.
	palette.SetCurrent(palette.ByName(*paletteFlag), *paletteRemapColors)

	// When dumping video or benchmarking, do precisely one render frame per update.
	if dump.Slow() || demo.Timedemo() {
		ebiten.SetMaxTPS(ebiten.SyncWithFPS)
	} else {
		ebiten.SetMaxTPS(engine.GameTPS / *fpsDivisor)
	}

	// Pause when unfocused, except when recording demos.
	ebiten.SetRunnableOnUnfocused(*runnableWhenUnfocused || (demo.Playing() && dump.Active()))

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
	status, err := g.init.Enter("precaching fonts", locale.G.Get("precaching fonts"), "could not precache fonts", splash.Single(font.KeepInCache))
	if status != splash.Continue {
		return err
	}
	status, err = g.init.Enter("precaching credits", locale.G.Get("precaching credits"), "could not precache credits", splash.Single(credits.Precache))
	if status != splash.Continue {
		return err
	}
	status, err = g.init.Enter("initializing audio", locale.G.Get("initializing audio"), "could not initialize audio", splash.Single(audiowrap.Init))
	if status != splash.Continue {
		return err
	}
	status, err = g.init.Enter("initializing noise", locale.G.Get("initializing noise"), "could not initialize noise", splash.Single(noise.Init))
	if status != splash.Continue {
		return err
	}
	status, err = g.init.Enter("precaching sounds", locale.G.Get("precaching sounds"), "could not precache sounds", sound.Precache)
	if status != splash.Continue {
		return err
	}
	status, err = g.init.Enter("precaching images", locale.G.Get("precaching images"), "could not precache images", splash.Single(image.Precache))
	if status != splash.Continue {
		return err
	}
	status, err = g.init.Enter("initializing input", locale.G.Get("initializing input"), "could not initialize input", splash.Single(input.Init))
	if status != splash.Continue {
		return err
	}
	status, err = g.init.Enter("precaching engine", locale.G.Get("precaching engine"), "could not precache engine", engine.Precache)
	if status != splash.Continue {
		return err
	}
	status, err = g.init.Enter("initializing dumping", locale.G.Get("initializing dumping"), "could not initialize dumping", splash.Single(dump.InitLate))
	if status != splash.Continue {
		return err
	}
	if *dumpLoadingFractions != "" {
		f, err := os.Create(*dumpLoadingFractions)
		if err != nil {
			return fmt.Errorf("could not open loading fractions file: %w", err)
		}
		j := json.NewEncoder(f)
		j.SetIndent("", "\t")
		err = j.Encode(g.init.ToFractions())
		if err != nil {
			return fmt.Errorf("could not encode to loading fractions file: %w", err)
		}
		err = f.Close()
		if err != nil {
			return fmt.Errorf("could not close loading fractions file: %w", err)
		}
	}
	if *debugJustInit {
		log.Errorf("requested early termination via --debug_just_init")
		return exitstatus.RegularTermination
	}
	log.Infof("game started")
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
	err := dump.Finish()
	if err != nil {
		return fmt.Errorf("could not finish dumping: %w", err)
	}
	err = demo.BeforeExit()
	if err != nil {
		return fmt.Errorf("could not finalize demo: %w", err)
	}
	return nil
}

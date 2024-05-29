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

package main

import (
	"errors"
	"runtime"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/aaaaxy"
	"github.com/divVerent/aaaaxy/internal/atexit"
	"github.com/divVerent/aaaaxy/internal/exitstatus"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

var (
	debugCpuprofile           = flag.String("debug_cpuprofile", "", "write CPU profile to file")
	debugLoadingCpuprofile    = flag.String("debug_loading_cpuprofile", "", "write CPU profile of loading to file")
	debugProfile              = flag.StringMap[string]("debug_profile", map[string]string{}, "key=value map to write profile indicated by key to given file path; possible keys include heap, allocs, threadcreate, block and mutex.")
	debugMemprofileRate       = flag.Int("debug_memprofile_rate", runtime.MemProfileRate, "fraction of bytes to be included in -debug_profile=heap=... and -debug_profile=allocs=...")
	debugBlockprofileRate     = flag.Int("debug_blockprofile_rate", 1, "resolution in nanoseconds for recording the block profile")
	debugMutexprofileFraction = flag.Int("debug_mutexprofile_fraction", max(runtime.SetMutexProfileFraction(-1), 1), "resolution in nanoseconds for recording the block profile")
	debugLogFile              = flag.String("debug_log_file", "", "log file to write all messages to (may be slow)")
)

func setProfileRates() {
	// Set the profile rates as soon as possible.
	if _, found := (*debugProfile)["heap"]; found {
		runtime.MemProfileRate = *debugMemprofileRate
	}
	if _, found := (*debugProfile)["allocs"]; found {
		runtime.MemProfileRate = *debugMemprofileRate
	}
	if _, found := (*debugProfile)["block"]; found {
		runtime.SetBlockProfileRate(*debugBlockprofileRate)
	}
	if _, found := (*debugProfile)["mutex"]; found {
		runtime.SetMutexProfileFraction(*debugMutexprofileFraction)
	}
}

func runGame(game *aaaaxy.Game) error {
	if *debugLoadingCpuprofile != "" {
		f, err := vfs.OSCreate(vfs.WorkDir, *debugLoadingCpuprofile)
		if err != nil {
			log.Fatalf("could not create loading CPU profile: %v", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatalf("could not start CPU profile: %v", err)
		}
	}
	err := game.InitFull()
	if *debugLoadingCpuprofile != "" {
		pprof.StopCPUProfile()
	}

	if err != nil {
		log.Fatalf("could not initialize game: %v", err)
	}

	if *debugCpuprofile != "" {
		f, err := vfs.OSCreate(vfs.WorkDir, *debugCpuprofile)
		if err != nil {
			log.Fatalf("could not create CPU profile: %v", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatalf("could not start CPU profile: %v", err)
		}
	}
	err = ebiten.RunGame(game)
	if *debugCpuprofile != "" {
		pprof.StopCPUProfile()
	}

	if len(*debugProfile) != 0 {
		runtime.GC() // Ensure up to date memory profiles.
	}

	for profile, path := range *debugProfile {
		f, err := vfs.OSCreate(vfs.WorkDir, path)
		if err != nil {
			log.Fatalf("could not create %s profile: %v", profile, err)
		}
		defer f.Close()
		if err := pprof.Lookup(profile).WriteTo(f, 0); err != nil {
			log.Fatalf("could not write %s profile: %v", profile, err)
		}
	}

	return err
}

func main() {
	defer atexit.Finish()

	// Turn all panics into Fatalf for uniform exception handling.
	ok := false
	defer func() {
		if !ok {
			log.Fatalf("got panic: %v", recover())
		}
	}()

	flag.Parse(aaaaxy.LoadConfig)

	setProfileRates()

	if *debugLogFile != "" {
		log.AddLogFile(*debugLogFile)
	}
	defer log.CloseLogFile()

	game := aaaaxy.NewGame()
	err := game.InitEbitengine()
	if err != nil {
		if errors.Is(err, exitstatus.ErrRegularTermination) {
			ok = true
			return
		}
		log.Fatalf("could not initialize game: %v", err)
	}
	err = runGame(game)
	errbe := game.BeforeExit()
	// From here on, nothing can panic.
	ok = true
	if err != nil && !errors.Is(err, exitstatus.ErrRegularTermination) {
		log.Fatalf("RunGame exited abnormally: %v", err)
	}
	if errbe != nil {
		log.Fatalf("BeforeExit exited abnormally: %v", errbe)
	}
}

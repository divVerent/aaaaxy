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
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/aaaaxy"
	"github.com/divVerent/aaaaxy/internal/atexit"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	debugCpuprofile     = flag.String("debug_cpuprofile", "", "write CPU profile to file")
	debugMemprofile     = flag.String("debug_memprofile", "", "write memory profile to file")
	debugMemprofileRate = flag.Int("debug_memprofile_rate", runtime.MemProfileRate, "fraction of bytes to be included in -debug_memprofile")
	debugLogFile        = flag.String("debug_log_file", "", "log file to write all messages to (may be slow)")
)

func runGame(game *aaaaxy.Game) error {
	if *debugCpuprofile != "" {
		err := game.InitFull()
		if err != nil {
			log.Fatalf("could not initialize game: %v", err)
		}
		f, err := os.Create(*debugCpuprofile)
		if err != nil {
			log.Fatalf("could not create CPU profile: %v", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatalf("could not start CPU profile: %v", err)
		}
		defer pprof.StopCPUProfile()
	}
	err := ebiten.RunGame(game)
	if *debugMemprofile != "" {
		f, err := os.Create(*debugMemprofile)
		if err != nil {
			log.Fatalf("could not create memory profile: %v", err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatalf("could not write memory profile: %v", err)
		}
	}
	return err
}

func main() {
	// Turn all panics into Fatalf for uniform exception handling.
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("got panic: %v", r)
		}
	}()

	flag.Parse(aaaaxy.LoadConfig)

	if *debugMemprofile != "" {
		// Set the memory profile rate as soon as possible.
		runtime.MemProfileRate = *debugMemprofileRate
	}

	if *debugLogFile != "" {
		log.AddLogFile(*debugLogFile)
	}
	game := &aaaaxy.Game{}
	err := game.InitEbiten()
	if err != nil {
		log.Fatalf("could not initialize game: %v", err)
	}
	err = runGame(game)
	errbe := game.BeforeExit()
	if err != nil && err != aaaaxy.RegularTermination {
		log.Fatalf("RunGame exited abnormally: %v", err)
	}
	if errbe != nil && errbe != aaaaxy.RegularTermination {
		log.Fatalf("BeforeExit exited abnormally: %v", errbe)
	}
	log.CloseLogFile()
	atexit.Finish()
}

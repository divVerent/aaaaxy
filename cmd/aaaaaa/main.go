package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/aaaaaa"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write CPU profile to file")
)

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	aaaaaa.InitEbiten()
	game := &aaaaaa.Game{}
	if err := ebiten.RunGame(game); err != nil {
		log.Print(err)
	}
}

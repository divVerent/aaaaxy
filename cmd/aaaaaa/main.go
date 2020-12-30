package main

import (
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/aaaaaa"
)

func main() {
	flag.Parse()
	aaaaaa.InitEbiten()
	game := &aaaaaa.Game{}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

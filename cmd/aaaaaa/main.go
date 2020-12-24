package main

import (
	"github.com/divVerent/aaaaaa/internal/aaaaaa"
	"github.com/hajimehoshi/ebiten/v2"
	"log"

	_ "github.com/divVerent/aaaaaa/internal/assets/statik"
)

func main() {
	aaaaaa.InitEbiten()
	game := &aaaaaa.Game{}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

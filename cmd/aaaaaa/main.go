package main

import (
	"github.com/divVerent/aaaaaa/internal/aaaaaa"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	game := &aaaaaa.Game{}
	ebiten.SetWindowSize(640, 360)
	ebiten.SetWindowTitle("AAAAAA")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

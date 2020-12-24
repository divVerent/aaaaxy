package aaaaaa

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

var _ ebiten.Game = &Game{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Maybe enlarge to contain a largest 16:9 subrect?
	return 640, 360
}

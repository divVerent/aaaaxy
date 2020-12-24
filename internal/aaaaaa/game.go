package aaaaaa

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

var _ ebiten.Game = &Game{}

func (g *Game) Update() error {
	// Update game.
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw world.
	// Draw HUD.
	// Draw menu.
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}

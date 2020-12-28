package aaaaaa

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	World *World
}

var _ ebiten.Game = &Game{}

func (g *Game) Update() error {
	if g.World == nil {
		g.World = NewWorld()
	}
	return g.World.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.World.Draw(screen)
	// Draw HUD.
	// Draw menu.
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}

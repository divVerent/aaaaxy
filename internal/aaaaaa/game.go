package aaaaaa

import (
	"errors"
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	World *World
}

var _ ebiten.Game = &Game{}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("esc")
	}
	if g.World == nil {
		g.World = NewWorld()
	}
	return g.World.Update()
}

var frameIndex = 0

func (g *Game) Draw(screen *ebiten.Image) {
	g.World.Draw(screen)

	if os.Getenv("CAPTUREVIDEO") != "" {
		SaveImage(screen, fmt.Sprintf("frame_%08d.png", frameIndex))
		frameIndex++
	}
	// Draw HUD.
	// Draw menu.
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}

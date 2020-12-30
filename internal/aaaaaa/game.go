package aaaaaa

import (
	"errors"
	"flag"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	captureVideo = flag.String("capture_video", "", "filename prefix to capture game frames to")
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
	g.World.Draw(screen)

	if *captureVideo != "" {
		ebiten.SetMaxTPS(ebiten.UncappedTPS)
		SaveImage(screen, fmt.Sprintf("%s_%08d.png", *captureVideo, frameIndex))
		frameIndex++
	}

	// Draw HUD.
	// Draw menu.
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}

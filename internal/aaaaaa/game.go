package aaaaaa

import (
	"errors"
	"flag"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	captureVideo = flag.String("capture_video", "", "filename prefix to capture game frames to")
	showFps      = flag.Bool("show_fps", false, "show fps counter")
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

	if *showFps {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.1f fps, %.1f tps", ebiten.CurrentFPS(), ebiten.CurrentTPS()), 0, GameHeight-16)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}

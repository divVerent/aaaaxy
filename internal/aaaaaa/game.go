package aaaaaa

import (
	"errors"
	"flag"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/divVerent/aaaaaa/internal/engine"
	_ "github.com/divVerent/aaaaaa/internal/game" // Load entities.
	"github.com/divVerent/aaaaaa/internal/timing"
)

var (
	captureVideo    = flag.String("capture_video", "", "filename prefix to capture game frames to")
	externalCapture = flag.Bool("external_capture", false, "assume an external capture application like apitrace is running; makes game run in lock step with rendering")
	showFps         = flag.Bool("show_fps", false, "show fps counter")
)

type Game struct {
	World *engine.World
}

var _ ebiten.Game = &Game{}

func (g *Game) Update() error {
	timing.ReportRegularly()

	defer timing.Group()()
	timing.Section("update")
	defer timing.Group()()

	timing.Section("once")
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("esc")
	}
	if g.World == nil {
		g.World = engine.NewWorld()
	}

	timing.Section("world")
	return g.World.Update()
}

var frameIndex = 0

func (g *Game) Draw(screen *ebiten.Image) {
	defer timing.Group()()
	timing.Section("draw")
	defer timing.Group()()

	timing.Section("world")
	g.World.Draw(screen)

	if *captureVideo != "" || *externalCapture {
		ebiten.SetMaxTPS(ebiten.UncappedTPS)
	}

	if *captureVideo != "" {
		timing.Section("capture")
		saveImage(screen, fmt.Sprintf("%s_%08d.png", *captureVideo, frameIndex))
		frameIndex++
	}

	timing.Section("hud")
	// TODO Draw HUD.

	timing.Section("menu")
	// TODO Draw menu.

	if *showFps {
		timing.Section("fps")
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%.1f fps, %.1f tps", ebiten.CurrentFPS(), ebiten.CurrentTPS()), 0, engine.GameHeight-16)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return engine.GameWidth, engine.GameHeight
}

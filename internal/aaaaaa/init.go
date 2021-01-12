package aaaaaa

import (
	"flag"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
)

var (
	vsync = flag.Bool("vsync", true, "enable waiting for vertical synchronization")
)

func InitEbiten() {
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetFullscreen(true)
	ebiten.SetInitFocused(true)
	ebiten.SetMaxTPS(engine.GameTPS)
	ebiten.SetRunnableOnUnfocused(false)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetScreenTransparent(false)
	ebiten.SetVsyncEnabled(*vsync)
	ebiten.SetWindowDecorated(true)
	ebiten.SetWindowFloating(false)
	ebiten.SetWindowPosition(0, 0)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(engine.GameWidth, engine.GameHeight)
	ebiten.SetWindowTitle("AAAAAA")
}

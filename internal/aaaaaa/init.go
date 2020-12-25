package aaaaaa

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func InitEbiten() {
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetFullscreen(true)
	ebiten.SetInitFocused(true)
	ebiten.SetMaxTPS(GameTPS)
	ebiten.SetRunnableOnUnfocused(false)
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetScreenTransparent(false)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowDecorated(true)
	ebiten.SetWindowFloating(false)
	ebiten.SetWindowPosition(0, 0)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(GameWidth, GameHeight)
	ebiten.SetWindowTitle("AAAAAA")
}

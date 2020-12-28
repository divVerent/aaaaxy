package aaaaaa

import (
	"fmt"
	"image"
	_ "image/png"

	"github.com/divVerent/aaaaaa/internal/vfs"

	"github.com/hajimehoshi/ebiten/v2"
)

func LoadImage(purpose, name string) (*ebiten.Image, error) {
	data, err := vfs.Load(purpose, name)
	if err != nil {
		return nil, fmt.Errorf("could not load: %v", err)
	}
	defer data.Close()
	img, _, err := image.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %v", err)
	}
	return ebiten.NewImageFromImage(img), nil
}

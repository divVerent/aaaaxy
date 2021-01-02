package engine

import (
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/vfs"
)

type imagePath = struct {
	Purpose string
	Name    string
}

var loadedImages = map[imagePath]*ebiten.Image{}

func LoadImage(purpose, name string) (*ebiten.Image, error) {
	ip := imagePath{purpose, name}
	if img, found := loadedImages[ip]; found {
		return img, nil
	}
	data, err := vfs.Load(purpose, name)
	if err != nil {
		return nil, fmt.Errorf("could not load: %v", err)
	}
	defer data.Close()
	img, _, err := image.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %v", err)
	}
	eImg := ebiten.NewImageFromImage(img)
	loadedImages[ip] = eImg
	return eImg, nil
}

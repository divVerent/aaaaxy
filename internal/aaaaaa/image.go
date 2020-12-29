package aaaaaa

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/vfs"
)

var pngEncoder = png.Encoder{
	CompressionLevel: png.BestSpeed,
}

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

func SaveImage(img *ebiten.Image, name string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	err = pngEncoder.Encode(file, img)
	if err != nil {
		file.Close()
		return err
	}
	return file.Close()
}

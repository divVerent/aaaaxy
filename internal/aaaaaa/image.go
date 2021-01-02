package aaaaaa

import (
	"image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

var pngEncoder = png.Encoder{
	CompressionLevel: png.BestSpeed,
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

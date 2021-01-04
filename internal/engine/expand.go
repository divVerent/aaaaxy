package engine

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func expandImage(img, tmp *ebiten.Image, size int, weight float64) {
	opts := ebiten.DrawImageOptions{
		CompositeMode: ebiten.CompositeModeLighter,
		Filter:        ebiten.FilterNearest,
	}
	opts.ColorM.Scale(weight, weight, weight, 1)
	for size > 1 {
		size /= 2
		tmp.Fill(color.Gray{0})
		opts.GeoM.Reset()
		opts.GeoM.Translate(-float64(size), 0)
		tmp.DrawImage(img, &opts)
		opts.GeoM.Reset()
		opts.GeoM.Translate(float64(size), 0)
		tmp.DrawImage(img, &opts)
		img.Fill(color.Gray{0})
		opts.GeoM.Reset()
		opts.GeoM.Translate(0, -float64(size))
		img.DrawImage(tmp, &opts)
		opts.GeoM.Reset()
		opts.GeoM.Translate(0, float64(size))
		img.DrawImage(tmp, &opts)
	}
}

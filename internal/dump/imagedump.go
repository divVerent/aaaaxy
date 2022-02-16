// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dump

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func getPixelsRGBA(img *ebiten.Image) ([]byte, error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y
	pix := make([]byte, 4*width*height)
	p := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y).(color.RGBA)
			pix[p] = c.R
			pix[p+1] = c.G
			pix[p+2] = c.B
			pix[p+3] = c.A
			p += 4
		}
	}
	return pix, nil
}

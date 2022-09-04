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

package input

import (
	"math"

	m "github.com/divVerent/aaaaxy/internal/math"
)

func pointerCoords(screenWidth, screenHeight, gameWidth, gameHeight int, crtK1, crtK2 float64, x, y int) m.Pos {
	inX := float64(x)*float64(gameWidth)/float64(screenWidth) + 0.5
	inY := float64(y)*float64(gameHeight)/float64(screenHeight) + 0.5

	// Straight ported from linear2xcrt.kage.tmpl.
	// Assume srcImageSize is 1:1 -> "square pixels".
	srcImageSizeSrcSizeLen := math.Hypot(float64(gameWidth), float64(gameHeight))
	mapVecX := 1 / (0.5 * srcImageSizeSrcSizeLen)
	mapVecY := 1 / (0.5 * srcImageSizeSrcSizeLen)
	srcMidX := float64(gameWidth) / 2
	srcMidY := float64(gameHeight) / 2
	inRelX := (inX - srcMidX) * mapVecX
	inRelY := (inY - srcMidY) * mapVecY
	inLen := math.Hypot(inRelX, inRelY)
	inLen2 := inLen * inLen
	outFac := 1.0 + inLen2*(crtK1+inLen2*crtK2)
	outRelX := inRelX * outFac
	outRelY := inRelY * outFac
	outX := srcMidX + outRelX/mapVecX
	outY := srcMidY + outRelY/mapVecY
	return m.Pos{
		X: int(math.Floor(outX)),
		Y: int(math.Floor(outY)),
	}
}

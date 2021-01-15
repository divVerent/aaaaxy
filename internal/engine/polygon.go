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

package engine

import (
	"github.com/hajimehoshi/ebiten/v2"

	m "github.com/divVerent/aaaaaa/internal/math"
)

func makeVertex(geoM *ebiten.GeoM, p m.Pos) ebiten.Vertex {
	x, y := geoM.Apply(float64(p.X), float64(p.Y))
	return ebiten.Vertex{
		DstX:   float32(x),
		DstY:   float32(y),
		SrcX:   float32(x),
		SrcY:   float32(y),
		ColorR: 1,
		ColorG: 1,
		ColorB: 1,
		ColorA: 1,
	}
}

func drawPolygonAround(dst *ebiten.Image, center m.Pos, vertices []m.Pos, src *ebiten.Image, geoM ebiten.GeoM, options *ebiten.DrawTrianglesOptions) {
	eVerts := make([]ebiten.Vertex, len(vertices)+1)
	eIndices := make([]uint16, 3*len(vertices))
	eVerts[0] = makeVertex(&geoM, center)
	for i, vert := range vertices {
		eVerts[i+1] = makeVertex(&geoM, vert)
		eIndices[3*i] = 0
		if i == 0 {
			eIndices[3*i+1] = uint16(len(vertices))
		} else {
			eIndices[3*i+1] = uint16(i)
		}
		eIndices[3*i+2] = uint16(i + 1)
	}
	dst.DrawTriangles(eVerts, eIndices, src, options)
}

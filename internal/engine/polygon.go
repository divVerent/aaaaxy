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
	"image/color"

	m "github.com/divVerent/aaaaaa/internal/math"
)

func makeVertex(geoM, texM *ebiten.GeoM, p m.Pos, r, g, b, a float32) ebiten.Vertex {
	x, y := geoM.Apply(float64(p.X), float64(p.Y))
	tx, ty := texM.Apply(x, y)
	return ebiten.Vertex{
		DstX:   float32(x),
		DstY:   float32(y),
		SrcX:   float32(tx),
		SrcY:   float32(ty),
		ColorR: r,
		ColorG: g,
		ColorB: b,
		ColorA: a,
	}
}

func drawPolygonAround(dst *ebiten.Image, center m.Pos, vertices []m.Pos, src *ebiten.Image, color color.Color, geoM, texM ebiten.GeoM, options *ebiten.DrawTrianglesOptions) {
	rI, gI, bI, aI := color.RGBA()
	r, g, b, a := float32(rI)/65535.0, float32(gI)/65535.0, float32(bI)/65535.0, float32(aI)/65535.0
	eVerts := make([]ebiten.Vertex, len(vertices)+1)
	eIndices := make([]uint16, 3*len(vertices))
	eVerts[0] = makeVertex(&geoM, &texM, center, r, g, b, a)
	for i, vert := range vertices {
		eVerts[i+1] = makeVertex(&geoM, &texM, vert, r, g, b, a)
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

func drawAntiPolygonAround(dst *ebiten.Image, center m.Pos, vertices []m.Pos, src *ebiten.Image, color color.Color, geoM, texM ebiten.GeoM, options *ebiten.DrawTrianglesOptions) {
	rI, gI, bI, aI := color.RGBA()
	r, g, b, a := float32(rI)/65535.0, float32(gI)/65535.0, float32(bI)/65535.0, float32(aI)/65535.0
	eVerts := make([]ebiten.Vertex, len(vertices)*2)
	eIndices := make([]uint16, 6*len(vertices))
	c := makeVertex(&geoM, &texM, center, r, g, b, a)
	for i, vert := range vertices {
		v := makeVertex(&geoM, &texM, vert, r, g, b, a)
		eVerts[2*i] = v
		// Now project v coordinates to the outside.
		d2x := v.DstX - c.DstX
		d2y := v.DstY - c.DstY
		fL := -d2x / c.DstX
		fU := -d2y / c.DstY
		fR := d2x / (GameWidth - c.DstX)
		fD := d2y / (GameHeight - c.DstY)
		f := fL
		if f < fU {
			f = fU
		}
		if f < fR {
			f = fR
		}
		if f < fD {
			f = fD
		}
		v.DstX = d2x/f + c.DstX
		v.DstY = d2y/f + c.DstY
		tx, ty := texM.Apply(float64(v.DstX), float64(v.DstY))
		v.SrcX, v.SrcY = float32(tx), float32(ty)
		eVerts[2*i+1] = v
		if i == 0 {
			eIndices[6*i] = uint16(2*len(vertices) - 2)
			eIndices[6*i+1] = uint16(2 * i)
			eIndices[6*i+2] = uint16(2*len(vertices) - 1)
			eIndices[6*i+3] = uint16(2 * i)
			eIndices[6*i+4] = uint16(2*len(vertices) - 1)
			eIndices[6*i+5] = uint16(2*i + 1)
		} else {
			eIndices[6*i] = uint16(2*i - 2)
			eIndices[6*i+1] = uint16(2 * i)
			eIndices[6*i+2] = uint16(2*i - 1)
			eIndices[6*i+3] = uint16(2 * i)
			eIndices[6*i+4] = uint16(2*i - 1)
			eIndices[6*i+5] = uint16(2*i + 1)
		}
	}
	dst.DrawTriangles(eVerts, eIndices, src, options)
}

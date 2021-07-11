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
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	m "github.com/divVerent/aaaaxy/internal/math"
)

// DrawPolyLine draws a polygon line for the given points.
// The screen coords are mapped to texcoords using texM (NOT VICE VERSA).
func DrawPolyLine(src *ebiten.Image, thickness float64, points []m.Pos, img *ebiten.Image, color color.Color, texM *ebiten.GeoM, options *ebiten.DrawTrianglesOptions) {
	rI, gI, bI, aI := color.RGBA()
	r, g, b, a := float32(rI)/65535.0, float32(gI)/65535.0, float32(bI)/65535.0, float32(aI)/65535.0
	vertices := make([]ebiten.Vertex, 2*len(points))
	indices := make([]uint16, 6*(len(points)-1))
	for i, p := range points {
		// Add vertices for this point.
		tX, tY := 0.0, 0.0
		if i > 0 {
			vX := float64(p.X - points[i-1].X)
			vY := float64(p.Y - points[i-1].Y)
			l := math.Hypot(vX, vY)
			tX += vX / l
			tY += vY / l
		}
		if i < len(points)-1 {
			vX := float64(points[i+1].X - p.X)
			vY := float64(points[i+1].Y - p.Y)
			l := math.Hypot(vX, vY)
			tX += vX / l
			tY += vY / l
		}
		l := math.Hypot(tX, tY)
		dX, dY := tX*thickness/l*0.5, tY*thickness/l*0.5
		dX0, dY0, dX1, dY1 := dY, -dX, -dY, dX
		if i == 0 {
			dX0 -= dX
			dY0 -= dY
			dX1 -= dX
			dY1 -= dY
		}
		if i == len(points)-1 {
			dX0 += dX
			dY0 += dY
			dX1 += dX
			dY1 += dY
		}
		dstX0, dstY0 := float64(p.X)+dX0+0.5, float64(p.Y)+dY0+0.5
		dstX1, dstY1 := float64(p.X)+dX1+0.5, float64(p.Y)+dY1+0.5
		srcX0, srcY0 := texM.Apply(dstX0, dstY0)
		srcX1, srcY1 := texM.Apply(dstX1, dstY1)
		vertices[2*i] = ebiten.Vertex{
			DstX:   float32(dstX0),
			DstY:   float32(dstY0),
			SrcX:   float32(srcX0),
			SrcY:   float32(srcY0),
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		}
		vertices[2*i+1] = ebiten.Vertex{
			DstX:   float32(dstX1),
			DstY:   float32(dstY1),
			SrcX:   float32(srcX1),
			SrcY:   float32(srcY1),
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		}
		// Add indices for the triangles.
		if i > 0 {
			indices[6*i-6] = uint16(2*i - 2)
			indices[6*i-5] = uint16(2*i - 1)
			indices[6*i-4] = uint16(2 * i)
			indices[6*i-3] = uint16(2*i - 1)
			indices[6*i-2] = uint16(2 * i)
			indices[6*i-1] = uint16(2*i + 1)
		}
	}
	src.DrawTriangles(vertices, indices, img, options)
}

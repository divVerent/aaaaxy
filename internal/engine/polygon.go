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

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/m"
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

var (
	polyVerts   []ebiten.Vertex
	polyIndices []uint16
)

func allocVerts(verts int) []ebiten.Vertex {
	if cap(polyVerts) < verts {
		polyVerts = make([]ebiten.Vertex, 2*verts)
	}
	return polyVerts[:verts]
}

func allocIndices(indices int) []uint16 {
	if cap(polyIndices) < indices {
		polyIndices = make([]uint16, 2*indices)
	}
	return polyIndices[:indices]
}

// drawPolygonAround draws a filled polygon.
func drawPolygonAround(dst *ebiten.Image, center m.Pos, vertices []m.Pos, src *ebiten.Image, color color.Color, geoM, texM ebiten.GeoM, options *ebiten.DrawTrianglesOptions) {
	rI, gI, bI, aI := color.RGBA()
	r, g, b, a := float32(rI)/65535.0, float32(gI)/65535.0, float32(bI)/65535.0, float32(aI)/65535.0
	eVerts := allocVerts(len(vertices) + 1)
	eIndices := allocIndices(3 * len(vertices))
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

// drawAntiPolygonAround draws all pixels except for the ones covered by the polygon.
// The polygon must go exactly clockwise or counterclockwise from center.
// Minkowski expanded polygons do NOT fulfill this right now, as they can contain self intersections! TODO fix this?
func drawAntiPolygonAround(dst *ebiten.Image, center m.Pos, vertices []m.Pos, src *ebiten.Image, color color.Color, geoM, texM ebiten.GeoM, options *ebiten.DrawTrianglesOptions) {
	rI, gI, bI, aI := color.RGBA()
	r, g, b, a := float32(rI)/65535.0, float32(gI)/65535.0, float32(bI)/65535.0, float32(aI)/65535.0
	eVerts := allocVerts(len(vertices) * 2)
	eIndices := allocIndices(6 * len(vertices))
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

// expandSimpleFrac is the fraction of pixels to move. Should be odd so m.Div() is symmetric.
const expandSimpleFrac = 5

func expandSimpleCoord(x, x0, shift int) int {
	d := m.Div(x-x0, expandSimpleFrac)
	if d > shift {
		return x + shift
	}
	if d < -shift {
		return x - shift
	}
	return x + d
}

// expandSimple expands the given polygon IN PLACE (i.e. clobbers polygon).
func expandSimple(center m.Pos, polygon []m.Pos, shift int) []m.Pos {
	for i, v1 := range polygon {
		// Rather approximate polygon expanding: just push each vertex shift away from the center.
		// Unlike correct polygon expansion perpendicular to sides,
		// this way ensures that we never include more than distance shift from the polugon.
		// However this is just approximate and causes artifacts when close to a wall.
		// Also, it is not _quite_ away from the center (doing this approximate method for better symmetry),
		// so the resulting polygon may not be QUITE correct for drawAntiPolygonAround,
		// however the error is small enough to not matter.
		polygon[i].X = expandSimpleCoord(v1.X, center.X, shift)
		polygon[i].Y = expandSimpleCoord(v1.Y, center.Y, shift)
		// More accurate but asymmetric and bad farther from the player:
		// d := v1.Delta(center)
		// polygon[i] = v1.Add(d.WithLengthFixed(m.NewFixed(shift)))
	}
	return polygon
}

// intersection returns the intersection of the lines a..b and c..d.
func intersection(a, b, c, d m.Pos) m.Pos {
	dab := b.Delta(a)
	dcd := d.Delta(c)
	den := dab.DX*dcd.DY - dab.DY*dcd.DX
	if den == 0 {
		// Parallel. Return any common point. In our concrete scenario the midpoint of BC is best.
		return b.Add(c.Delta(m.Pos{})).Div(2)
	}
	cdxy := c.X*d.Y - c.Y*d.X
	abxy := a.X*b.Y - a.Y*b.X
	nx := dab.DX*cdxy - dcd.DX*abxy
	ny := dab.DY*cdxy - dcd.DY*abxy
	return m.Pos{
		X: (2*nx + den) / (2 * den),
		Y: (2*ny + den) / (2 * den),
	}
}

func collinear(a, b, c m.Pos) bool {
	return (b.X-a.X)*(c.Y-b.Y)-(b.Y-a.Y)*(c.X-b.X) == 0
}

// Global buffers to lose thread safety and reduce allocations.
var (
	minkowskiNoSame      []m.Pos
	minkowskiNoCollinear []m.Pos
	minkowskiEdgeCorner  []m.Delta
	minkowskiOut         []m.Pos
)

// expandMinkowski expands a given polygon to its Minkowski sum with a box from -boxSize,-boxSize to boxSize,boxSize.
func expandMinkowski(polygon []m.Pos, boxSize int) []m.Pos {
	// First simplify the polygon. We can't have any duplicate or collinear vertices.
	// Sadly we need to remove dupes first and collinearities second,
	// or a dupe at a corner causes us to lose an entire vertex.
	minkowskiNoSame = minkowskiNoSame[:0]
	for i, a := range polygon {
		b := polygon[m.Mod(i+1, len(polygon))]
		if a != b {
			minkowskiNoSame = append(minkowskiNoSame, b)
		}
	}
	minkowskiNoCollinear = minkowskiNoCollinear[:0]
	for i, b := range minkowskiNoSame {
		a := minkowskiNoSame[m.Mod(i-1, len(minkowskiNoSame))]
		c := minkowskiNoSame[m.Mod(i+1, len(minkowskiNoSame))]
		if !collinear(a, b, c) {
			minkowskiNoCollinear = append(minkowskiNoCollinear, b)
		}
	}
	// Iterate over all edges.
	minkowskiEdgeCorner = minkowskiEdgeCorner[:0]
	for i, a := range polygon {
		b := polygon[m.Mod(i+1, len(polygon))]
		dab := b.Delta(a)
		nab := m.Left().Apply(dab)
		corner := m.Delta{DX: boxSize, DY: boxSize}
		if nab.DX < 0 {
			corner.DX = -corner.DX
		}
		if nab.DY < 0 {
			corner.DY = -corner.DY
		}
		minkowskiEdgeCorner = append(minkowskiEdgeCorner, corner)
	}
	// Iterate over all edge pairs.
	minkowskiOut = minkowskiOut[:0]
	for i, b := range polygon {
		a := polygon[m.Mod(i-1, len(polygon))]
		c := polygon[m.Mod(i+1, len(polygon))]
		cab := minkowskiEdgeCorner[m.Mod(i-1, len(polygon))]
		cbc := minkowskiEdgeCorner[i]
		dab := b.Delta(a)
		dbc := c.Delta(b)
		nab := m.Left().Apply(dab)
		isConcave := nab.Dot(dbc) > 0
		if isConcave {
			// The "concave" case. None of the corners remains.
			minkowskiOut = append(minkowskiOut, intersection(
				a.Add(cab), b.Add(cab),
				b.Add(cbc), c.Add(cbc),
			))
		} else {
			// The "convex" case. Add all corners on the path. Often just one.
			corners := make([]m.Delta, 0, 4)
			corner := cab
			for corner != cbc {
				corners = append(corners, corner)
				corner = m.Right().Apply(corner)
			}
			corners = append(corners, corner)
			for _, corner := range corners {
				minkowskiOut = append(minkowskiOut, b.Add(corner))
			}
		}
	}
	return minkowskiOut
}

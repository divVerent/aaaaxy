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

package math

import (
	"bytes"
	"fmt"
)

// Rect represents a rectangle.
type Rect struct {
	// Origin is the origin of the rectangle, typically the top left corner.
	Origin Pos
	// Size is the size of the rectangle, typically positive.
	Size Delta
}

// 3 l 2 = {3, 4}
// 3 l 1 = {3}
// 3 l 0 = {}
// 3 l -1 = {2}
// 3 l -2 = {2, 1}

// Normalized returns a rectangle such that its size is nonnegative.
func (r Rect) Normalized() Rect {
	if r.Size.DX < 0 {
		r.Origin.X += r.Size.DX
		r.Size.DX = -r.Size.DX
	}
	if r.Size.DY < 0 {
		r.Origin.Y += r.Size.DY
		r.Size.DY = -r.Size.DY
	}
	return r
}

func RectFromPoints(a, b Pos) Rect {
	var r Rect
	if a.X < b.X {
		r.Origin.X = a.X
		r.Size.DX = b.Y - a.X + 1
	} else {
		r.Origin.X = b.X
		r.Size.DX = a.Y - b.X + 1
	}
	if a.Y < b.Y {
		r.Origin.Y = a.Y
		r.Size.DY = b.Y - a.Y + 1
	} else {
		r.Origin.Y = b.Y
		r.Size.DY = a.Y - b.Y + 1
	}
	return r
}

// Add creates a new rectangle moved by the given delta.
func (r Rect) Add(d Delta) Rect {
	return Rect{
		Origin: r.Origin.Add(d),
		Size:   r.Size,
	}
}

// Grow creates a new rectangle grown by the given delta.
func (r Rect) Grow(d Delta) Rect {
	return Rect{
		Origin: r.Origin.Sub(d),
		Size:   r.Size.Add(d.Mul(2)),
	}
}

// OppositeCorner returns the coordinate of the opposite corner of the rectangle. Only correct on normalized rectangles.
func (r Rect) OppositeCorner() Pos {
	return r.Origin.Add(r.Size).Sub(Delta{DX: 1, DY: 1})
}

// Center returns the coordinate in the middle of the rectangle.
func (r Rect) Center() Pos {
	return r.Origin.Add(r.Size.Div(2))
}

// Foot returns the coordinate in the bottom middle of the rectangle.
func (r Rect) Foot() Pos {
	return r.Origin.Add(Delta{DX: r.Size.DX / 2, DY: r.Size.DY - 1})
}

func intervalDistance(a0, a1, b0, b1 int) int {
	// If intervals are separated, compute separation amount.
	if b0 > a1 {
		return a1 - b0 // -1 when touching.
	}
	if a0 > b1 {
		return a0 - b1 // 1 when touching.
	}
	// Otherwise, we have b0 <= a1 && a0 <= b1. They overlap.
	return 0
}

// delta returns the vector between the closest points of two rectangles.
func (r Rect) delta(c10, c11 Pos) Delta {
	c00 := r.Origin
	c01 := r.OppositeCorner()
	xDist := intervalDistance(c00.X, c01.X, c10.X, c11.X)
	yDist := intervalDistance(c00.Y, c01.Y, c10.Y, c11.Y)
	return Delta{DX: xDist, DY: yDist}
}

// Delta returns the vector between the closest points of two rectangles.
func (r Rect) Delta(other Rect) Delta {
	return r.delta(other.Origin, other.OppositeCorner())
}

// DeltaPos returns the vector between the closest points of a rectangle and a point.
func (r Rect) DeltaPos(other Pos) Delta {
	return r.delta(other, other)
}

// GridPos converts coordinates of p within the rect into grid coords.
func (r Rect) GridPos(p Pos, nx int, ny int) (int, int) {
	dx := p.X - r.Origin.X
	dy := p.Y - r.Origin.Y
	ix, _ := mulFracModUint64(uint64(dx), uint64(nx), uint64(r.Size.DX))
	iy, _ := mulFracModUint64(uint64(dy), uint64(ny), uint64(r.Size.DY))
	return int(ix), int(iy)
}

func (r Rect) String() string {
	return fmt.Sprintf("%d %d %d %d", r.Origin.X, r.Origin.Y, r.Size.DX, r.Size.DY)
}

func (r Rect) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}

func (r *Rect) UnmarshalText(text []byte) error {
	_, err := fmt.Fscanf(bytes.NewReader(text), "%d %d %d %d", &r.Origin.X, &r.Origin.Y, &r.Size.DX, &r.Size.DY)
	return err
}

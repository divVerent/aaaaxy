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

// OppositeCorner returns the coordinate of the opposite corner of the rectangle. Only correct on normalized rectangles.
func (r Rect) OppositeCorner() Pos {
	return r.Origin.Add(r.Size).Sub(Delta{DX: 1, DY: 1})
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

// Diff returns the vector between the closest points of two rectangles.
func (r Rect) Delta(other Rect) Delta {
	c00 := r.Origin
	c01 := r.OppositeCorner()
	c10 := other.Origin
	c11 := other.OppositeCorner()
	xDist := intervalDistance(c00.X, c01.X, c10.X, c11.X)
	yDist := intervalDistance(c00.Y, c01.Y, c10.Y, c11.Y)
	return Delta{DX: xDist, DY: yDist}
}

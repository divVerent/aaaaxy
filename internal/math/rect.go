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

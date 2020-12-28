package math

// Pos represents a pixel position, where X points right and Y points down.
type Pos struct {
	X, Y int
}

// Add applies a delta to a position.
func (p Pos) Add(d Delta) Pos {
	return Pos{p.X + d.DX, p.Y + d.DY}
}

func (p Pos) Sub(d Delta) Pos {
	return Pos{p.X - d.DX, p.Y - d.DY}
}

func (p Pos) Delta(p2 Pos) Delta {
	return Delta{p.X - p2.X, p.Y - p2.Y}
}

func (p Pos) Scale(n int, m int) Pos {
	return Pos{p.X * n / m, p.Y * n / m}
}

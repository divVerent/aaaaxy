package aaaaaa

// Pos represents a pixel position.
type Pos struct {
	X, Y int
}

// Delta represents a move between two pixel positions.
type Delta struct {
	DX, DY int
}

// Orientation represents a transformation matrix.
type Orientation struct {
	Right Delta
	Up    Delta
}

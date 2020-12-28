package math

// Delta represents a move between two pixel positions.
type Delta struct {
	DX, DY int
}

func (d Delta) Scale(n int, m int) Delta {
	return Delta{d.DX * n / m, d.DY * n / m}
}

func North() Delta {
	return Delta{DX: 0, DY: -1}
}
func East() Delta {
	return Delta{DX: 1, DY: 0}
}
func South() Delta {
	return Delta{DX: 0, DY: 1}
}
func West() Delta {
	return Delta{DX: -1, DY: 0}
}

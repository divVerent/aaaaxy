package math

import (
	"math"
)

// Delta represents a move between two pixel positions.
type Delta struct {
	DX, DY int
}

func (d Delta) Norm1() int {
	norm := 0
	if d.DX >= 0 {
		norm += d.DX
	} else {
		norm -= d.DX
	}
	if d.DY >= 0 {
		norm += d.DY
	} else {
		norm -= d.DY
	}
	return norm
}

func (d Delta) Length2() int {
	return d.DX*d.DX + d.DY*d.DY
}

func (d Delta) Length() float64 {
	return math.Sqrt(float64(d.Length2()))
}

func (d Delta) Add(d2 Delta) Delta {
	return Delta{DX: d.DX + d2.DX, DY: d.DY + d2.DY}
}

func (d Delta) Sub(d2 Delta) Delta {
	return Delta{DX: d.DX - d2.DX, DY: d.DY - d2.DY}
}

func (d Delta) Mul(n int) Delta {
	return Delta{DX: d.DX * n, DY: d.DY * n}
}

func (d Delta) Div(m int) Delta {
	return Delta{DX: Div(d.DX, m), DY: Div(d.DY, m)}
}

func (d Delta) MulFloat(f float64) Delta {
	return Delta{DX: int(float64(d.DX)*f + 0.5), DY: int(float64(d.DY)*f + 0.5)}
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

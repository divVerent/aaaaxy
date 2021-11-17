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
	"math"
)

// Delta represents a move between two pixel positions.
type Delta struct {
	DX, DY int
}

func (d Delta) Norm0() int {
	norm := 0
	if d.DX > norm {
		norm = d.DX
	} else if -d.DX > norm {
		norm = -d.DX
	}
	if d.DY > norm {
		norm = d.DY
	} else if -d.DY > norm {
		norm = -d.DY
	}
	return norm
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

func (d Delta) Length2() int64 {
	return int64(d.DX)*int64(d.DX) + int64(d.DY)*int64(d.DY)
}

func (d Delta) Length() float64 {
	return math.Sqrt(float64(d.Length2()))
}

func (d Delta) LengthFixed() Fixed {
	return NewFixedInt64(d.Length2()).Sqrt()
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

func (d Delta) Mul2(mx, my int) Delta {
	return Delta{DX: d.DX * mx, DY: d.DY * my}
}

func (d Delta) Div(m int) Delta {
	return Delta{DX: Div(d.DX, m), DY: Div(d.DY, m)}
}

func (d Delta) MulFixed(f Fixed) Delta {
	return Delta{DX: NewFixed(d.DX).Mul(f).Rint(), DY: NewFixed(d.DY).Mul(f).Rint()}
}

func (d Delta) WithLengthFixed(f Fixed) Delta {
	n := d.LengthFixed()
	if n == 0 {
		return d
	}
	return d.MulFixed(f.Div(n))
}

func (d Delta) WithMaxLengthFixed(f Fixed) Delta {
	n := d.LengthFixed()
	if n <= f {
		return d
	}
	return d.MulFixed(f.Div(n))
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
func (d Delta) Dot(d2 Delta) int {
	return d.DX*d2.DX + d.DY*d2.DY
}

func (d Delta) IsZero() bool {
	return d.DX == 0 && d.DY == 0
}

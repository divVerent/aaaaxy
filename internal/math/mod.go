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
	"math/bits"
)

func Mod(a, b int) int {
	if a < 0 {
		return b - 1 - ^a%b
	}
	return a % b
}

func Div(a, b int) int {
	if a < 0 {
		return ^(^a / b)
	}
	return a / b
}

func Rint(f float64) int {
	return int(math.RoundToEven(f))
}

// MulFracInt64 returns a*b/d rounded to even.
func MulFracInt64(a, b, d int64) int64 {
	sign := int64(1)
	if a < 0 {
		sign, a = -sign, -a
	}
	if b < 0 {
		sign, b = -sign, -b
	}
	if d < 0 {
		sign, d = -sign, -d
	}
	du := uint64(d)
	ch, cl := bits.Mul64(uint64(a), uint64(b))
	q, r := bits.Div64(ch, cl, du)
	rcut := du / 2
	if q%2 == 0 && du%2 == 0 {
		// Round to even logic: if result is even and we're at exactly half, don't increment.
		rcut++
	}
	if r >= rcut {
		q++
	}
	return int64(q) * sign
}

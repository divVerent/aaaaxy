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
	"fmt"
	"math"

	"golang.org/x/image/math/fixed"
)

type Fixed fixed.Int52_12

func NewFixed(i int) Fixed {
	return Fixed(i) << 12
}

func NewFixedInt64(i int64) Fixed {
	return Fixed(i) << 12
}

func NewFixedFloat64(f float64) Fixed {
	return Fixed(math.RoundToEven(f * 4096.0))
}

func (f Fixed) Mul(g Fixed) Fixed {
	return Fixed(fixed.Int52_12(f).Mul(fixed.Int52_12(g)))
}

func (f Fixed) MulFrac(n, d Fixed) Fixed {
	// TODO:
	// f64 := int64(f)
	// n64 := int64(n)
	// d64 := int64(d)
	// Compute f * n / d without accuracy loss.
	return f.Mul(n.Div(d))
}

func (f Fixed) Div(g Fixed) Fixed {
	return Fixed((f<<12 + g>>1) / g)
}

func (f Fixed) Rint() int {
	return int((f + 2048) >> 12)
}

func (f Fixed) Float64() float64 {
	return float64(f) * (1.0 / 4096.0)
}

func (f Fixed) String() string {
	return fmt.Sprintf("%d.0x%03x", int64(f)>>12, int64(f)&0xFFF)
}

func (f Fixed) Sqrt() Fixed {
	if f == 0 {
		return 0
	}

	// Compute a wild guess using the FPU.
	guess := NewFixedFloat64(math.Sqrt(f.Float64()))

	// Want unique number s so that, where delta=0.5:
	//   s-delta <= 4096*sqrt(f/4096) < s+delta
	// Square everything; assumes s-delta >= 0. Thus the check above.
	//   s^2 - s <= 4096 * f - 1/4 < s^2 + s
	//   s^2 - s < 4096 * f <= s^2 + s

	// In practice these loops tend to execute only once.

	goal := f << 12
	// fixes := 0
	s := guess
	for s*s-s >= goal {
		s--
		// fixes++
	}
	for s*s+s < goal {
		s++
		// fixes++
	}
	// if fixes > 16 {
	// log.Fatalf("too many fixes for Sqrt(%v): %v fixes, guess was %v, result is %v", f, fixes, guess, s)
	// }

	return s
}

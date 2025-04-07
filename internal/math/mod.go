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

func Mod64(a, b int64) int64 {
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

func Div64(a, b int64) int64 {
	if a < 0 {
		return ^(^a / b)
	}
	return a / b
}

func Rint(f float64) int {
	return int(math.RoundToEven(f))
}

func MulFrac(a, b, d int) int {
	return int(mulFracInt64(int64(a), int64(b), int64(d)))
}

// mulFracInt64 returns a*b/d rounded to even.
func mulFracInt64(a, b, d int64) int64 {
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
	q, r := mulFracModUint64(uint64(a), uint64(b), uint64(d))
	q = roundToEvenUint64(uint64(d), uint64(q), uint64(r))
	return int64(q) * sign
}

// mulFracModUint64 returns a*b/d with remainder like in C (may be negative).
func mulFracModUint64(a, b, d uint64) (uint64, uint64) {
	ch, cl := bits.Mul64(a, b)
	return bits.Div64(ch, cl, d)
}

// roundToEvenUint64 rounds q with remainder r at divisor d towards nearest, and towards even on a tie.
func roundToEvenUint64(d, q, r uint64) uint64 {
	// Cutoff:
	// d odd: always d/2+1.
	// d even: d/2+1 if q even, d/2 if q odd.
	rcut := d / 2
	if d%2 != 0 || q%2 == 0 {
		rcut++
	}
	if r >= rcut {
		q++
	}
	return q
}

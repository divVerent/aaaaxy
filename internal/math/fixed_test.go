package math

import (
	"fmt"
	"testing"
)

func TestNewFixed(t *testing.T) {
	got := NewFixed(42)
	const want Fixed = 42 * 0x1000
	if got != want {
		t.Errorf("NewFixed(42): got %v, want %v", got, want)
	}
}

func TestNewFixedInt64(t *testing.T) {
	got := NewFixedInt64(42)
	const want Fixed = 42 * 0x1000
	if got != want {
		t.Errorf("NewFixedInt64(42): got %v, want %v", got, want)
	}
}

func TestNewFixedFloat64(t *testing.T) {
	for _, tc := range []struct {
		In   float64
		Want Fixed
	}{
		{In: -1.5 / 4096, Want: -2},
		{In: -0.5 / 4096, Want: 0},
		{In: 2.5 / 4096, Want: 2},
		{In: 3.5 / 4096, Want: 4},
		{In: 4.5 / 4096, Want: 4},
	} {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			got := NewFixedFloat64(tc.In)
			if got != tc.Want {
				t.Errorf("NewFixedFloat64(In): got %v, want %v", got, tc.Want)
			}
		})
	}
}

func TestFixedMul(t *testing.T) {
	for _, tc := range []struct {
		A, B Fixed
		Want Fixed
	}{
		{A: NewFixedFloat64(2), B: NewFixedFloat64(4), Want: NewFixedFloat64(8)},
		{A: NewFixedFloat64(0.5), B: NewFixedFloat64(0.25), Want: NewFixedFloat64(0.125)},
		{A: NewFixedFloat64(1.0 / 64), B: NewFixedFloat64(1.0 / 64), Want: NewFixedFloat64(1.0 / 4096)},
		{A: NewFixedFloat64(-7.0 / 128), B: NewFixedFloat64(1.0 / 64), Want: NewFixedFloat64(-4.0 / 4096)},
		{A: NewFixedFloat64(-13.0 / 256), B: NewFixedFloat64(1.0 / 64), Want: NewFixedFloat64(-3.0 / 4096)},
		{A: NewFixedFloat64(-5.0 / 128), B: NewFixedFloat64(1.0 / 64), Want: NewFixedFloat64(-2.0 / 4096)},
		{A: NewFixedFloat64(1.0 / 128), B: NewFixedFloat64(1.0 / 64), Want: NewFixedFloat64(0)},
		{A: NewFixedFloat64(3.0 / 128), B: NewFixedFloat64(1.0 / 64), Want: NewFixedFloat64(2.0 / 4096)},
		{A: NewFixedFloat64(5.0 / 128), B: NewFixedFloat64(1.0 / 64), Want: NewFixedFloat64(2.0 / 4096)},
	} {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			got := tc.A.Mul(tc.B)
			if got != tc.Want {
				t.Errorf("A.Mul(B): got %v, want %v", got, tc.Want)
			}
		})
	}
}

func TestFixedMulFrac(t *testing.T) {
	for _, tc := range []struct {
		A, B, C Fixed
		Want    Fixed
	}{
		{A: NewFixedFloat64(2), B: NewFixedFloat64(8), C: NewFixedFloat64(2), Want: NewFixedFloat64(8)},
		{A: NewFixedFloat64(0.5), B: NewFixedFloat64(0.125), C: NewFixedFloat64(0.5), Want: NewFixedFloat64(0.125)},
		{A: NewFixedFloat64(1.0 / 64), B: NewFixedFloat64(1.0 / 64), C: NewFixedFloat64(1), Want: NewFixedFloat64(1.0 / 4096)},
		{A: NewFixedFloat64(-7.0 / 128), B: NewFixedFloat64(1), C: NewFixedFloat64(64), Want: NewFixedFloat64(-4.0 / 4096)},
		{A: NewFixedFloat64(-13.0 / 256), B: NewFixedFloat64(1 << 40), C: NewFixedFloat64(64 << 40), Want: NewFixedFloat64(-3.0 / 4096)},
		{A: NewFixedFloat64(-5.0 / 128), B: NewFixedFloat64(-1), C: NewFixedFloat64(-64), Want: NewFixedFloat64(-2.0 / 4096)},
		{A: NewFixedFloat64(1.0 / 128), B: NewFixedFloat64(1 << 44), C: NewFixedFloat64(64 << 44), Want: NewFixedFloat64(0)},
		{A: NewFixedFloat64(3.0 / 128), B: NewFixedFloat64(-1 << 45), C: NewFixedFloat64(-64 << 45), Want: NewFixedFloat64(2.0 / 4096)},
		{A: NewFixedFloat64(5.0 / 128), B: NewFixedFloat64(0.5), C: NewFixedFloat64(32), Want: NewFixedFloat64(2.0 / 4096)},
	} {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			got := tc.A.MulFrac(tc.B, tc.C)
			if got != tc.Want {
				t.Errorf("A.MulFrac(B, C): got %v, want %v", got, tc.Want)
			}
		})
	}
}

func TestFixedDiv(t *testing.T) {
	for _, tc := range []struct {
		A, B Fixed
		Want Fixed
	}{
		{A: NewFixedFloat64(2), B: NewFixedFloat64(0.25), Want: NewFixedFloat64(8)},
		{A: NewFixedFloat64(0.5), B: NewFixedFloat64(4), Want: NewFixedFloat64(0.125)},
		{A: NewFixedFloat64(1.0 / 64), B: NewFixedFloat64(64), Want: NewFixedFloat64(1.0 / 4096)},
		{A: NewFixedFloat64(-7.0 / 128), B: NewFixedFloat64(64), Want: NewFixedFloat64(-4.0 / 4096)},
		{A: NewFixedFloat64(-13.0 / 256), B: NewFixedFloat64(64), Want: NewFixedFloat64(-3.0 / 4096)},
		{A: NewFixedFloat64(-5.0 / 128), B: NewFixedFloat64(64), Want: NewFixedFloat64(-2.0 / 4096)},
		{A: NewFixedFloat64(1.0 / 128), B: NewFixedFloat64(64), Want: NewFixedFloat64(0)},
		{A: NewFixedFloat64(3.0 / 128), B: NewFixedFloat64(64), Want: NewFixedFloat64(2.0 / 4096)},
		{A: NewFixedFloat64(5.0 / 128), B: NewFixedFloat64(64), Want: NewFixedFloat64(2.0 / 4096)},
	} {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			got := tc.A.Div(tc.B)
			if got != tc.Want {
				t.Errorf("A.Div(B): got %v, want %v", got, tc.Want)
			}
		})
	}
}

func TestFixedRint(t *testing.T) {
	for _, tc := range []struct {
		A Fixed
		Want int
	}{
		{A: NewFixedFloat64(1<<30), Want: 1<<30},
		{A: NewFixedFloat64(0.5), Want: 0},
		{A: NewFixedFloat64(1.25), Want: 1},
		{A: NewFixedFloat64(1.5), Want: 2},
		{A: NewFixedFloat64(1.75), Want: 2},
		{A: NewFixedFloat64(2.25), Want: 2},
		{A: NewFixedFloat64(2.5), Want: 2},
		{A: NewFixedFloat64(2.75), Want: 3},
		{A: NewFixedFloat64(3.5), Want: 4},
		{A: NewFixedFloat64(-0.5), Want: 0},
		{A: NewFixedFloat64(-1.25), Want: -1},
		{A: NewFixedFloat64(-1.5), Want: -2},
		{A: NewFixedFloat64(-1.75), Want: -2},
		{A: NewFixedFloat64(-2.25), Want: -2},
		{A: NewFixedFloat64(-2.5), Want: -2},
		{A: NewFixedFloat64(-2.75), Want: -3},
		{A: NewFixedFloat64(-3.5), Want: -4},
	} {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			got := tc.A.Rint()
			if got != tc.Want {
				t.Errorf("A.Rint(): got %v, want %v", got, tc.Want)
			}
		})
	}
}

func TestFixedFloat64(t *testing.T) {
	for _, tc := range []struct {
		A Fixed
		Want float64
	}{
		{A: NewFixedFloat64(1<<30), Want: 1<<30},
		{A: NewFixedFloat64(0.5), Want: 0.5},
		{A: NewFixedFloat64(1.25), Want: 1.25},
		{A: NewFixedFloat64(1.5), Want: 1.5},
	} {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			got := tc.A.Float64()
			if got != tc.Want {
				t.Errorf("A.Float64(): got %v, want %v", got, tc.Want)
			}
		})
	}
}

func TestFixedSqrt(t *testing.T) {
	for _, tc := range []struct {
		A Fixed
		Want Fixed
	}{
		{A: NewFixedFloat64(0), Want: NewFixedFloat64(0)},
		{A: NewFixedFloat64(1.0 / 4096), Want: NewFixedFloat64(1.0/64)},
		{A: NewFixedFloat64(1), Want: NewFixedFloat64(1)},
	} {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			got := tc.A.Sqrt()
			if got != tc.Want {
				t.Errorf("A.Sqrt(): got %v, want %v", got, tc.Want)
			}
		})
	}
}

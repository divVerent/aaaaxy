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
	"strings"
)

// Orientation represents a transformation matrix, written as a right and a down vector.
type Orientation struct {
	Right Delta
	Down  Delta
}

// Concat returns the orientation o * o2 so that o.Concat(o2).Apply(d) == o.Apply(o2.Apply(d)).
func (o Orientation) Concat(o2 Orientation) Orientation {
	return Orientation{
		Right: o.Apply(o2.Right),
		Down:  o.Apply(o2.Down),
	}
}

// Apply rotates a delta by an orientation.
func (o Orientation) Apply(d Delta) Delta {
	return Delta{
		DX: o.Right.DX*d.DX + o.Down.DX*d.DY,
		DY: o.Right.DY*d.DX + o.Down.DY*d.DY,
	}
}

// ApplyToRect2 rotates a rectangle by the given orientation, mapping the pivot to itself.
// The pivot is given in doubled coordinates to support half-pixel pivots.
// Note: odd numbers are pixel centers, even numbers are pixel corners!
func (o Orientation) ApplyToRect2(pivot2 Pos, r Rect) Rect {
	return Rect{
		Origin: pivot2.Add(o.Apply(r.Origin.Mul(2).Delta(pivot2))).Div(2),
		Size:   o.Apply(r.Size),
	}.Normalized()
}

// Inverse returns an orientation so that o.Concat(o.Invert()) == Identity().
func (o Orientation) Inverse() Orientation {
	// There is probably a more efficient way, but all our orientations are identity when applied four times.
	return o.Concat(o).Concat(o)
}

// Identity yields the default orientation.
func Identity() Orientation {
	return Orientation{Right: East(), Down: South()}
}

// FlipX yields an orientation where X is flipped.
func FlipX() Orientation {
	return Orientation{Right: West(), Down: South()}
}

// FlipY yields an orientation where Y is flipped.
func FlipY() Orientation {
	return Orientation{Right: East(), Down: North()}
}

// FlipD yields an orientation where X/Y are swapped.
func FlipD() Orientation {
	return Orientation{Right: South(), Down: East()}
}

// Left yields an orientation that turns left.
func Left() Orientation {
	return Orientation{Right: North(), Down: East()}
}

// Right yields an orientation that turns right.
func Right() Orientation {
	return Orientation{Right: South(), Down: West()}
}

// Left yields an orientation that turns left.
func TurnAround() Orientation {
	return Orientation{Right: West(), Down: North()}
}

// ParseOrientation parses an orientation from a string. It is given by the right and down directions in that order.
func ParseOrientation(s string) (Orientation, error) {
	switch s {
	case "EN":
		return Orientation{Right: East(), Down: North()}, nil
	case "ES":
		return Orientation{Right: East(), Down: South()}, nil
	case "NE":
		return Orientation{Right: North(), Down: East()}, nil
	case "NW":
		return Orientation{Right: North(), Down: West()}, nil
	case "SE":
		return Orientation{Right: South(), Down: East()}, nil
	case "SW":
		return Orientation{Right: South(), Down: West()}, nil
	case "WN":
		return Orientation{Right: West(), Down: North()}, nil
	case "WS":
		return Orientation{Right: West(), Down: South()}, nil
	default:
		return Orientation{}, fmt.Errorf("unsupported orientation %q; want <right><down> direction like ES", s)
	}
}

func ParseOrientations(s string) ([]Orientation, error) {
	orientations := strings.Split(s, " ")
	if len(orientations) == 0 {
		return nil, fmt.Errorf("unsupported orientation list: empty")
	}
	out := make([]Orientation, len(orientations))
	for i, orientation := range orientations {
		var err error
		out[i], err = ParseOrientation(orientation)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func (o Orientation) Determinant() int {
	return o.Right.DX*o.Down.DY - o.Right.DY*o.Down.DX
}

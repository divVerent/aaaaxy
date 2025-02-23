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
	"errors"
	"fmt"
	"strings"
)

// Orientation represents a transformation matrix, written as a right and a down vector.
//
// The zero value is not valid here.
type Orientation struct {
	Right Delta
	Down  Delta
}

// IsZero returns whether o is the zero value. The zero value is not valid to use.
func (o Orientation) IsZero() bool {
	return o.Right.IsZero() && o.Down.IsZero()
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

// Apply rotates a pos by an orientation, mapping the pivot to itself.
// The pivot is given in doubled coordinates to support half-pixel pivots.
// Note: odd numbers are pixel centers, even numbers are pixel corners, and p is considered a pixel corner!
// Alternatively, to consider p a pixel center, just subtract 1 from the pivot coordinates.
func (o Orientation) Apply2(pivot2, p Pos) Pos {
	return pivot2.Add(o.Apply(p.Mul(2).Delta(pivot2))).Div(2)
}

// ApplyToRect2 rotates a rectangle by the given orientation, mapping the pivot to itself.
// The pivot is given in doubled coordinates to support half-pixel pivots.
// Note: odd numbers are pixel centers, even numbers are pixel corners!
func (o Orientation) ApplyToRect2(pivot2 Pos, r Rect) Rect {
	return Rect{
		Origin: o.Apply2(pivot2, r.Origin),
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

type Orientations []Orientation

var AllOrientations = Orientations{
	Orientation{Right: East(), Down: North()},
	Orientation{Right: East(), Down: South()},
	Orientation{Right: North(), Down: East()},
	Orientation{Right: North(), Down: West()},
	Orientation{Right: South(), Down: East()},
	Orientation{Right: South(), Down: West()},
	Orientation{Right: West(), Down: North()},
	Orientation{Right: West(), Down: South()},
}

func ParseOrientations(s string) (Orientations, error) {
	if s == "" {
		return Orientations{}, nil
	}
	orientations := strings.Split(s, " ")
	if len(orientations) == 0 {
		return nil, errors.New("unsupported orientation list: empty")
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

func (o Orientation) MarshalText() ([]byte, error) {
	switch o {
	case Orientation{Right: East(), Down: North()}:
		return []byte("EN"), nil
	case Orientation{Right: East(), Down: South()}:
		return []byte("ES"), nil
	case Orientation{Right: North(), Down: East()}:
		return []byte("NE"), nil
	case Orientation{Right: North(), Down: West()}:
		return []byte("NW"), nil
	case Orientation{Right: South(), Down: East()}:
		return []byte("SE"), nil
	case Orientation{Right: South(), Down: West()}:
		return []byte("SW"), nil
	case Orientation{Right: West(), Down: North()}:
		return []byte("WN"), nil
	case Orientation{Right: West(), Down: South()}:
		return []byte("WS"), nil
	case Orientation{}:
		// Used on some optional fields, otherwise should not happen.
		return []byte(""), nil
	default:
		return nil, fmt.Errorf("unsupported Orientation{Right: %v, Down: %v}", o.Right, o.Down)
	}
}

func (o Orientation) String() string {
	text, err := o.MarshalText()
	if err != nil {
		return err.Error()
	}
	return string(text)
}

func (o *Orientation) UnmarshalText(text []byte) error {
	var err error
	*o, err = ParseOrientation(string(text))
	return err
}

func (o Orientations) MarshalText() ([]byte, error) {
	if len(o) == 0 {
		return []byte{}, nil
	}
	out := make([]byte, 0, 3*len(o)-1)
	for _, o1 := range o {
		out1, err := o1.MarshalText()
		if err != nil {
			return nil, err
		}
		if len(out) != 0 {
			out = append(out, ' ')
		}
		out = append(out, out1...)
	}
	return out, nil
}

func (o *Orientations) UnmarshalText(text []byte) error {
	var err error
	*o, err = ParseOrientations(string(text))
	return err
}

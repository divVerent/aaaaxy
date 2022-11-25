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
	"bytes"
	"fmt"
)

// Pos represents a pixel position, where X points right and Y points down.
type Pos struct {
	X, Y int
}

// Add applies a delta to a position.
func (p Pos) Add(d Delta) Pos {
	return Pos{p.X + d.DX, p.Y + d.DY}
}

func (p Pos) Sub(d Delta) Pos {
	return Pos{p.X - d.DX, p.Y - d.DY}
}

func (p Pos) Delta(p2 Pos) Delta {
	return Delta{p.X - p2.X, p.Y - p2.Y}
}

func (p Pos) Mul(n int) Pos {
	return Pos{X: p.X * n, Y: p.Y * n}
}

func (p Pos) Div(m int) Pos {
	return Pos{X: Div(p.X, m), Y: Div(p.Y, m)}
}

func (p Pos) FromRectToRect(a Rect, b Rect) Pos {
	return Pos{
		X: b.Origin.X + Div((p.X-a.Origin.X)*b.Size.DX, a.Size.DX),
		Y: b.Origin.Y + Div((p.Y-a.Origin.Y)*b.Size.DY, a.Size.DY),
	}
}

func (p Pos) String() string {
	return fmt.Sprintf("%d %d", p.X, p.Y)
}

func (p Pos) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *Pos) UnmarshalText(text []byte) error {
	_, err := fmt.Fscanf(bytes.NewReader(text), "%d %d", &p.X, &p.Y)
	return err
}

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

package menu

import (
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/m"
)

const CenterX = engine.GameWidth / 2
const HeaderY = engine.GameHeight / 4

type Direction int

const (
	NotClicked Direction = iota
	LeftClicked
	CenterClicked
	RightClicked
)

func ItemBaselineY(i, n int) int {
	return engine.GameHeight * (31 - 2*(n-i)) / 32
}

func ItemClicked(pos m.Pos, n int) (int, Direction) {
	// Clicked far at side?
	if pos.X < engine.GameWidth/8 || pos.X > 7*engine.GameWidth/8 {
		return -1, NotClicked
	}

	// Adjust for baseline.
	y := pos.Y + engine.GameHeight/64

	// Map to index.
	i := n - (31-y*32/engine.GameHeight)/2
	if i >= 0 && i < n {
		dir := CenterClicked
		if pos.X < engine.GameWidth/3 {
			dir = LeftClicked
		} else if pos.X > 2*engine.GameWidth/3 {
			dir = RightClicked
		}
		return i, dir
	}

	// Outside.
	return -1, NotClicked
}

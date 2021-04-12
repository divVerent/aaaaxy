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

package level

import (
	m "github.com/divVerent/aaaaaa/internal/math"
)

// Contents indicates what kind of tiles/objects we want to hit.
type Contents int

const (
	NoContents          Contents = 0
	OpaqueContents      Contents = 1
	PlayerSolidContents Contents = 2
	ObjectSolidContents Contents = 4
	SolidContents       Contents = PlayerSolidContents | ObjectSolidContents
)

func (c Contents) Empty() bool {
	return c == NoContents
}

func (c Contents) Opaque() bool {
	return c&OpaqueContents != 0
}

func (c Contents) PlayerSolid() bool {
	return c&PlayerSolidContents != 0
}

func (c Contents) ObjectSolid() bool {
	return c&ObjectSolidContents != 0
}

type VisibilityFlags int

const (
	FrameVis  VisibilityFlags = 1
	TracedVis VisibilityFlags = 2
)

// A Tile is a single game tile.
type Tile struct {
	// Info needed for gameplay.
	Contents   Contents
	Spawnables []*Spawnable // NOTE: not adjusted for transform!

	// Info needed for loading more tiles.
	LevelPos        m.Pos
	Transform       m.Orientation
	VisibilityFlags VisibilityFlags

	// Info needed for rendering.
	Orientation m.Orientation
	ImageSrc    string

	// If provided, these are used instead of image for "nicer" rotation (e.g. for shadow effects).
	// Because Orientation is also set, looking these up is tricky; we want things to show up as in the editor but potentially rotated.
	// We know:
	// - Transform * Orientation = orientationInEditor
	// - If we pick tile I and render at orientation O, we actually render at full orientation O * I.
	// - BUT lighting direction orientation is just O.
	// - we want O = orientationInEditor.
	// - Solve: Orientation = orientationInEditor * I
	// - Orientation = (Transform * Orientation) * I
	// - O = Transform Orientation
	// - I = O^-1 Orientation
	ImageSrcByOrientation map[m.Orientation]string

	// Debug info.
	LoadedFromNeighbor m.Pos
}

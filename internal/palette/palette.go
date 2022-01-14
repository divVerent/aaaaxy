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

package palette

import (
	"sort"
	"strings"
)

// Palette encapsulates a color palette.
type Palette struct {
	// size is the number of colors this palette has. Is > 0 for any valid palette.
	size int

	// bayerScale is a recommended Bayer pattern scale in the range from 0 to 1.
	// It should be close to a typical distance between two close colors.
	bayerScale float64

	// colors are the palette colors.
	colors []uint32
}

func (p *Palette) BayerScale() float64 {
	return p.bayerScale
}

func newPalette(bayerScale float64, c []uint32) *Palette {
	return &Palette{
		size:       len(c),
		colors:     c,
		bayerScale: bayerScale,
	}
}

// Names returns the names of all palettes, in quoted comma separated for, for inclusion in a flag description.
func Names() string {
	l := make([]string, 0, len(data))
	for p := range data {
		l = append(l, p)
	}
	sort.Strings(l)
	return "'" + strings.Join(l, "', '") + "'"
}

// ByName returns the PalData for the given palette. Do not modify the returned object.
func ByName(name string) *Palette {
	return data[name]
}

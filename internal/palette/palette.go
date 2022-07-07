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
	// name is the name of the palette.
	name string

	// size is the number of colors this palette has. Is > 0 for any valid palette.
	size int

	// egaIndices is a list with one entry per protected color whose value is the EGA color index.
	egaIndices []int

	// protected is the number of protected colors.
	// This is also used to compute the Bayer pattern size.
	protected int

	// colors are the palette colors.
	colors []uint32

	// remap is the color remapping.
	remap map[uint32]uint32

	// ega is the set of EGA colors after remapping.
	ega [EGACount]uint32
}

var current *Palette

func newPalette(egaIndices []int, c0 []uint32) *Palette {
	// Keep only unique colors beyond egaIndices.
	h := make(map[uint32]struct{}, len(c0))
	c := make([]uint32, 0, len(c0))
	for i, u := range c0 {
		if _, found := h[u]; found && i >= len(egaIndices) {
			continue
		}
		h[u] = struct{}{}
		c = append(c, u)
	}

	protected := len(egaIndices)
	if protected == 0 {
		protected = len(c)
	}
	ega := egaColors
	remap := map[uint32]uint32{}
	for thisIdx, egaIdx := range egaIndices {
		from := toRGB(egaColors[egaIdx]).toUint32()
		to := toRGB(c[thisIdx]).toUint32()
		if from != to {
			remap[from] = to
			ega[egaIdx] = to
		}
	}
	if len(remap) == 0 {
		remap = nil
	}
	pal := &Palette{
		size:       len(c),
		egaIndices: egaIndices,
		protected:  protected,
		colors:     c,
		remap:      remap,
		ega:        ega,
	}
	return pal
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

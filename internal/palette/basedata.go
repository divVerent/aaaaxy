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

type EGAIndex int

const (
	Black EGAIndex = iota
	Blue
	Green
	Cyan
	Red
	Magenta
	Brown
	LightGrey
	DarkGrey
	LightBlue
	LightGreen
	LightCyan
	LightRed
	LightMagenta
	Yellow
	White
	EGACount
)

// egaColors is the set of reference colors.
var egaColors = [EGACount]uint32{
	0x000000,
	0x0000AA,
	0x00AA00,
	0x00AAAA,
	0xAA0000,
	0xAA00AA,
	0xAA5500,
	0xAAAAAA,
	0x555555,
	0x5555FF,
	0x55FF55,
	0x55FFFF,
	0xFF5555,
	0xFF55FF,
	0xFFFF55,
	0xFFFFFF,
}

var egaColorsSet = map[uint32]bool{}

func init() {
	for _, c := range egaColors {
		egaColorsSet[c] = true
	}
}

var nesColors = [64]uint32{
	// 0x00
	0x666666,
	0x002A88,
	0x1412A8,
	0x3B00A4,
	0x5C007E,
	0x6E0040,
	0x6C0700,
	0x571D00,
	0x343500,
	0x0C4900,
	0x005200,
	0x004F08,
	0x00404E,
	0x000000,
	0x000000,
	0x000000,
	// 0x10
	0xAEAEAE,
	0x155FDA,
	0x4240FE,
	0x7627FF,
	0xA11BCD,
	0xB81E7C,
	0xB53220,
	0x994F00,
	0x6C6E00,
	0x388700,
	0x0D9400,
	0x009032,
	0x007C8E,
	0x000000,
	0x000000,
	0x000000,
	// 0x20
	0xFEFEFE,
	0x64B0FE,
	0x9390FE,
	0xC777FE,
	0xF36AFE,
	0xFE6ECD,
	0xFE8270,
	0xEB9F23,
	0xBDBF00,
	0x89D900,
	0x5DE530,
	0x45E182,
	0x48CEDF,
	0x4F4F4F,
	0x000000,
	0x000000,
	// 0x30
	0xFEFEFE,
	0xC1E0FE,
	0xD4D3FE,
	0xE9C8FE,
	0xFBC3FE,
	0xFEC5EB,
	0xFECDC6,
	0xF7D9A6,
	0xE5E695,
	0xD0F097,
	0xBEF5AB,
	0xB4F3CD,
	0xB5ECF3,
	0xB8B8B8,
	0x000000,
	0x000000,
}

func grays(n int) []uint32 {
	m := n - 1
	l := make([]uint32, 0, m)
	for i := 0; i <= m; i++ {
		l = append(l,
			0x10101*((2*255*uint32(i)+uint32(m))/(2*uint32(m))))
	}
	return l
}

func colorCube(nr, ng, nb int) []uint32 {
	mr, mg, mb := nr-1, ng-1, nb-1
	l := make([]uint32, 0, nr*ng*nb)
	for r := 0; r <= mr; r++ {
		for g := 0; g <= mg; g++ {
			for b := 0; b <= mb; b++ {
				l = append(l,
					0x10000*((2*255*uint32(r)+uint32(mr))/(2*uint32(mr)))+
						0x100*((2*255*uint32(g)+uint32(mg))/(2*uint32(mg)))+
						((2*255*uint32(b)+uint32(mb))/(2*uint32(mb))))
			}
		}
	}
	return l
}

func midpoints(p1, p2, points []uint32, count uint32) []uint32 {
	out := make([]uint32, 0, len(p1)*len(p2))
	for _, ci := range p1 {
		ri := ci >> 16
		gi := (ci >> 8) & 0xFF
		bi := ci & 0xFF
		for _, cj := range p2 {
			rj := cj >> 16
			gj := (cj >> 8) & 0xFF
			bj := cj & 0xFF
			for _, i := range points {
				r := (2*(ri*i+rj*(count-i)) + count) / (2 * count)
				g := (2*(gi*i+gj*(count-i)) + count) / (2 * count)
				b := (2*(bi*i+bj*(count-i)) + count) / (2 * count)
				out = append(out, (r<<16)|(g<<8)|b)
			}
		}
	}
	return out
}

func rounded(p []uint32, n uint32) []uint32 {
	out := make([]uint32, 0, len(p))
	for _, c := range p {
		r := c >> 16
		g := (c >> 8) & 0xFF
		b := c & 0xFF
		r = (n*r + 127) / 255
		g = (n*g + 127) / 255
		b = (n*b + 127) / 255
		r = (2*255*r + n) / (2 * n)
		g = (2*255*g + n) / (2 * n)
		b = (2*255*b + n) / (2 * n)
		out = append(out, (r<<16)|(g<<8)|b)
	}
	return out
}

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

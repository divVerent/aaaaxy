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

package font

import (
	"log"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gomonobold"
)

// Face is an alias to font.Face so users do not need to import the font package.
type Face = font.Face

var (
	Centerprint    Face
	CenterprintBig Face
	DebugSmall     Face
)

func init() {
	// Load the fonts.
	italic, err := truetype.Parse(goitalic.TTF)
	if err != nil {
		log.Panicf("Could not load goitalic font: %v", err)
	}
	mono, err := truetype.Parse(gomono.TTF)
	if err != nil {
		log.Panicf("Could not load gomono font: %v", err)
	}
	monobold, err := truetype.Parse(gomonobold.TTF)
	if err != nil {
		log.Panicf("Could not load gomonobold font: %v", err)
	}

	Centerprint = truetype.NewFace(italic, &truetype.Options{
		Size:       16,
		Hinting:    font.HintingFull,
		SubPixelsX: 1,
		SubPixelsY: 1,
	})
	CenterprintBig = truetype.NewFace(monobold, &truetype.Options{
		Size:       24,
		Hinting:    font.HintingFull,
		SubPixelsX: 1,
		SubPixelsY: 1,
	})
	DebugSmall = truetype.NewFace(mono, &truetype.Options{
		Size:       5,
		Hinting:    font.HintingFull,
		SubPixelsX: 1,
		SubPixelsY: 1,
	})
}

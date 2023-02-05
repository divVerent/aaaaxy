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
	"fmt"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/gomedium"
	"golang.org/x/image/font/gofont/gomediumitalic"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gosmallcaps"
	"golang.org/x/image/font/opentype"

	m "github.com/divVerent/aaaaxy/internal/math"
)

var (
	ByName         = map[string]Face{}
	Centerprint    Face
	CenterprintBig Face
	DebugSmall     Face
	MenuBig        Face
	Menu           Face
	MenuSmall      Face
)

func makeGoFontFace(fnt *opentype.Font, options *opentype.FaceOptions) (Face, error) {
	f, err := opentype.NewFace(fnt, options)
	if err != nil {
		return Face{}, err
	}
	return makeFace(f, m.Rint(options.Size*options.DPI/72))
}

func InitGoFont() error {
	// Load the fonts.
	regular, err := opentype.Parse(gomedium.TTF)
	if err != nil {
		return fmt.Errorf("could not load gomedium font: %w", err)
	}
	italic, err := opentype.Parse(gomediumitalic.TTF)
	if err != nil {
		return fmt.Errorf("could not load gomediumitalic font: %w", err)
	}
	bold, err := opentype.Parse(gobold.TTF)
	if err != nil {
		return fmt.Errorf("could not load gosmallcaps font: %w", err)
	}
	mono, err := opentype.Parse(gomono.TTF)
	if err != nil {
		return fmt.Errorf("could not load gomono font: %w", err)
	}
	smallcaps, err := opentype.Parse(gosmallcaps.TTF)
	if err != nil {
		return fmt.Errorf("could not load gosmallcaps font: %w", err)
	}

	ByName["Small"], err = makeGoFontFace(regular, &opentype.FaceOptions{
		Size:    10,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	ByName["Regular"], err = makeGoFontFace(regular, &opentype.FaceOptions{
		Size:    14,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	ByName["Italic"], err = makeGoFontFace(italic, &opentype.FaceOptions{
		Size:    14,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	ByName["Bold"], err = makeGoFontFace(bold, &opentype.FaceOptions{
		Size:    14,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	ByName["Mono"], err = makeGoFontFace(mono, &opentype.FaceOptions{
		Size:    14,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	ByName["MonoSmall"], err = makeGoFontFace(mono, &opentype.FaceOptions{
		Size:    10,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	ByName["SmallCaps"], err = makeGoFontFace(smallcaps, &opentype.FaceOptions{
		Size:    14,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	Centerprint, err = makeGoFontFace(italic, &opentype.FaceOptions{
		Size:    14,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	CenterprintBig, err = makeGoFontFace(smallcaps, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	DebugSmall, err = makeGoFontFace(regular, &opentype.FaceOptions{
		Size:    9,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	Menu, err = makeGoFontFace(smallcaps, &opentype.FaceOptions{
		Size:    18,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	MenuBig, err = makeGoFontFace(smallcaps, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	MenuSmall, err = makeGoFontFace(smallcaps, &opentype.FaceOptions{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}

	return nil
}

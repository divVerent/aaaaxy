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
)

func makeGoFontFace(fnt *opentype.Font, size int) (*Face, error) {
	f, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}
	return makeFace(f, size), nil
}

func initGoFont() error {
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

	ByName["Small"], err = makeGoFontFace(regular, 10)
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	ByName["Regular"], err = makeGoFontFace(regular, 14)
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	ByName["Italic"], err = makeGoFontFace(italic, 14)
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	ByName["Bold"], err = makeGoFontFace(bold, 14)
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	ByName["Mono"], err = makeGoFontFace(mono, 14)
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	ByName["MonoSmall"], err = makeGoFontFace(mono, 10)
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	ByName["SmallCaps"], err = makeGoFontFace(smallcaps, 14)
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	ByName["Centerprint"] = ByName["Italic"]
	ByName["CenterprintBig"], err = makeGoFontFace(smallcaps, 24)
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	ByName["DebugSmall"], err = makeGoFontFace(regular, 9)
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	ByName["Menu"], err = makeGoFontFace(smallcaps, 18)
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}
	ByName["MenuBig"] = ByName["CenterprintBig"]
	ByName["MenuSmall"], err = makeGoFontFace(smallcaps, 12)
	if err != nil {
		return fmt.Errorf("could not create face: %w", err)
	}

	return nil
}

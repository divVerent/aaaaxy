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
	"compress/gzip"
	"fmt"
	"io"

	"github.com/zachomedia/go-bdf"

	"github.com/divVerent/aaaaxy/internal/vfs"
)

func initUnifont() error {
	unifontRawHandle, err := vfs.Load("fonts", "unifont-15.0.01.bdf.gz")
	if err != nil {
		return fmt.Errorf("could not open unifont: %w", err)
	}
	defer unifontRawHandle.Close()
	unifontHandle, err := gzip.NewReader(unifontRawHandle)
	if err != nil {
		return fmt.Errorf("could not decompress unifont: %w", err)
	}
	defer unifontHandle.Close()
	unifontData, err := io.ReadAll(unifontHandle)
	if err != nil {
		return fmt.Errorf("could not read unifont: %w", err)
	}
	unifont, err := bdf.Parse(unifontData)
	if err != nil {
		return fmt.Errorf("could not parse unifont: %w", err)
	}
	// 14, which is 16 when adding back the outline.
	face := makeFace(unifont.NewFace(), 14)

	ByName["Small"] = face
	ByName["Regular"] = face
	ByName["Italic"] = face
	ByName["Bold"] = face
	ByName["Mono"] = face
	ByName["MonoSmall"] = face
	ByName["SmallCaps"] = face
	ByName["Centerprint"] = face
	ByName["CenterprintBig"] = face
	ByName["Menu"] = face
	ByName["MenuBig"] = face
	ByName["MenuSmall"] = face

	return nil
}

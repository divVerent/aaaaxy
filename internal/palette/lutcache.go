// Copyright 2022 Google LLC
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
	"fmt"
	"image"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/screenshot"
)

var (
	dumpPaletteLUTs = flag.String("dump_palette_luts", "", "file name prefix to dump all palette LUT textures to; the game will then exit")
)

func Init(w, h int) error {
	if *dumpPaletteLUTs != "" {
		for name, p := range data {
			for numLUTs := 1; numLUTs <= 2; numLUTs++ {
				bounds := image.Rectangle{
					Min: image.Point{},
					Max: image.Point{X: w, Y: h},
				}
				img, size, perRow, width := p.computeLUT(bounds, numLUTs)
				name := fmt.Sprintf("%s.%s_%d.png", *dumpPaletteLUTs, name, numLUTs)
				log.Infof("saving %s (size=%d perRow=%d width=%d)...", name, size, perRow, width)
				err := screenshot.Write(img, name)
				if err != nil {
					return err
				}
			}
		}
		return fmt.Errorf("palette LUTs dumped")
	}
	return nil
}

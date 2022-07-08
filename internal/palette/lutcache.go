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
	"encoding/json"
	"fmt"
	"image"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/exitstatus"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/screenshot"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

var (
	dumpPaletteLUTsPrefix = flag.String("dump_palette_luts_prefix", "", "file name prefix to dump all palette LUT textures to; the game will then exit")
	paletteMaxCycles      = flag.Float64("palette_max_cycles", 640*360*256*4, "maximum number of cycles to spend on palette generation; only applies if there is no cached palette file in the game")
)

type lutMeta struct {
	Size   int `json:"size"`
	PerRow int `json:"per_row"`
	Width  int `json:"width"`
}

func Init(w, h int) error {
	if *dumpPaletteLUTsPrefix != "" {
		for name, p := range data {
			for numLUTs := 1; numLUTs <= 2; numLUTs++ {
				bounds := image.Rectangle{
					Min: image.Point{},
					Max: image.Point{X: w, Y: h},
				}
				img, size, perRow, width := p.computeLUT(bounds, numLUTs, math.Inf(+1))
				name := fmt.Sprintf("%s%s_%d.png", *dumpPaletteLUTsPrefix, name, numLUTs)
				log.Infof("saving %s (size=%d perRow=%d width=%d)...", name, size, perRow, width)
				err := screenshot.Write(img, name)
				if err != nil {
					return err
				}
				meta := lutMeta{
					Size:   size,
					PerRow: perRow,
					Width:  width,
				}
				metaName := name + ".json"
				f, err := os.Create(metaName)
				if err != nil {
					return err
				}
				j := json.NewEncoder(f)
				j.SetIndent("", "\t")
				err = j.Encode(meta)
				if err != nil {
					return err
				}
				err = f.Close()
				if err != nil {
					return err
				}
			}
		}
		log.Errorf("requested a palette LUT dump - not running the game")
		return exitstatus.RegularTermination
	}
	return nil
}

func (p *Palette) loadLUT(numLUTs int) (image.Image, int, int, int, error) {
	name := fmt.Sprintf("lut_%s_%d.png", p.name, numLUTs)
	data, err := vfs.Load("palettes", name)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("could not open %v: %w", name, err)
	}
	defer data.Close()
	img, _, err := image.Decode(data)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("could not decode %v: %w", name, err)
	}
	var meta lutMeta
	metaName := name + ".json"
	j, err := vfs.Load("palettes", metaName)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("could not open %v: %w", metaName, err)
	}
	defer j.Close()
	err = json.NewDecoder(j).Decode(&meta)
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("could not decode palette LUT json config file %v: %w", metaName, err)
	}
	return img, meta.Size, meta.PerRow, meta.Width, nil
}

func (p *Palette) ToLUT(numLUTs int, img *ebiten.Image) (int, int, int) {
	lut, lutSize, perRow, lutWidth, err := p.loadLUT(numLUTs)
	if err != nil {
		log.Warningf("cached palette data not found, generating at runtime: %v", err)
		lut, lutSize, perRow, lutWidth = p.computeLUT(img.Bounds(), numLUTs, *paletteMaxCycles)
	}
	if nrgba, ok := lut.(*image.NRGBA); ok {
		img.SubImage(nrgba.Rect).(*ebiten.Image).ReplacePixels(nrgba.Pix)
	} else {
		log.Fatalf("expecting a NRGBA LUT, got %T", img)
	}
	return lutSize, perRow, lutWidth
}

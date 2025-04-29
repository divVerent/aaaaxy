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

package picture

import (
	"bufio"
	"fmt"
	"image"
	_ "image/png"
	"path"
	"regexp"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/palette"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

var (
	precachePictures = flag.Bool("precache_pictures", true, "preload all pictures at startup (VERY recommended)")
)

type picturePath = struct {
	Purpose string
	Name    string
}

var (
	cache       = map[picturePath]*ebiten.Image{}
	cacheFrozen bool

	// This should be in sync with exclusions in scripts/audit-pictures.sh.
	noPaletteSprites = regexp.MustCompile(`^(?:warpzone|clock|gradient|magic)_.*`)
)

func load(purpose, name string, force bool) (*ebiten.Image, error) {
	ip := picturePath{purpose, name}
	cachedImg, found := cache[ip]
	if found && !force {
		return cachedImg, nil
	}
	if cacheFrozen && !found {
		return nil, fmt.Errorf("picture %v was not precached", ip)
	}
	data, err := vfs.Load(purpose, name)
	if err != nil {
		return nil, fmt.Errorf("could not load: %w", err)
	}
	defer data.Close()
	img, _, err := image.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}
	usePalette := true
	if purpose == "sprites" {
		if noPaletteSprites.MatchString(name) {
			usePalette = false
		}
	}
	if usePalette {
		img = palette.Current().ApplyToImage(img, name)
	}
	eImg := ebiten.NewImageFromImage(img)
	if eImg.Bounds().Min != (image.Point{}) {
		return nil, fmt.Errorf("could not get zero origin: %v", eImg.Bounds())
	}
	cache[ip] = eImg
	return eImg, nil
}

func Load(purpose, name string) (*ebiten.Image, error) {
	return load(purpose, name, false)
}

func Precache() error {
	if !*precachePictures {
		return nil
	}
	toLoad := map[picturePath]struct{}{}
	for _, purpose := range []string{"tiles", "sprites"} {
		names, err := vfs.ReadDir(purpose)
		if err != nil {
			return fmt.Errorf("could not enumerate files in %v: %w", purpose, err)
		}
		for _, name := range names {
			if !strings.HasSuffix(name, ".png") {
				continue
			}
			toLoad[picturePath{Purpose: purpose, Name: name}] = struct{}{}
		}
	}
	listFile, err := vfs.Load("generated", "picture_load_order.txt")
	if err != nil {
		return fmt.Errorf("could query load order: %w", err)
	}
	listScanner := bufio.NewScanner(listFile)
	for listScanner.Scan() {
		line := listScanner.Text()
		purpose := path.Dir(line)
		name := path.Base(line)
		item := picturePath{Purpose: purpose, Name: name}
		if _, found := toLoad[item]; found {
			_, err := Load(item.Purpose, item.Name)
			if err != nil {
				return fmt.Errorf("could not precache %v: %w", item, err)
			}
			delete(toLoad, item)
		} else {
			return fmt.Errorf("could not find file for precache item %v", item)
		}
	}
	for item := range toLoad {
		return fmt.Errorf("could not find precache item for file %v", item)
	}
	cacheFrozen = true
	return nil
}

func PaletteChanged() error {
	for ip := range cache {
		_, err := load(ip.Purpose, ip.Name, true)
		if err != nil {
			return err
		}
	}
	return nil
}

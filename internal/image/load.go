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

package image

import (
	"bufio"
	"fmt"
	"image"
	_ "image/png"
	"path"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

var (
	precacheImages = flag.Bool("precache_images", true, "preload all images at startup (VERY recommended)")
)

type imagePath = struct {
	Purpose string
	Name    string
}

var (
	cache       = map[imagePath]*ebiten.Image{}
	cacheFrozen bool
)

func Load(purpose, name string) (*ebiten.Image, error) {
	ip := imagePath{purpose, name}
	if img, found := cache[ip]; found {
		return img, nil
	}
	if cacheFrozen {
		return nil, fmt.Errorf("image %v was not precached", ip)
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
	eImg := ebiten.NewImageFromImage(img)
	cache[ip] = eImg
	return eImg, nil
}

func Precache() error {
	if !*precacheImages {
		return nil
	}
	toLoad := map[imagePath]struct{}{}
	for _, purpose := range []string{"tiles", "sprites"} {
		names, err := vfs.ReadDir(purpose)
		if err != nil {
			return fmt.Errorf("could not enumerate files in %v: %w", purpose, err)
		}
		for _, name := range names {
			if !strings.HasSuffix(name, ".png") {
				continue
			}
			toLoad[imagePath{Purpose: purpose, Name: name}] = struct{}{}
		}
	}
	listFile, err := vfs.Load("generated", "image_load_order.txt")
	if err != nil {
		return fmt.Errorf("could query load order: %w", err)
	}
	listScanner := bufio.NewScanner(listFile)
	for listScanner.Scan() {
		line := listScanner.Text()
		purpose := path.Dir(line)
		name := path.Base(line)
		item := imagePath{Purpose: purpose, Name: name}
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

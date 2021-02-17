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
	"fmt"
	"image"
	_ "image/png"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/flag"
	"github.com/divVerent/aaaaaa/internal/vfs"
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
	name = vfs.Canonical(name)
	ip := imagePath{purpose, name}
	if img, found := cache[ip]; found {
		return img, nil
	}
	if cacheFrozen {
		return nil, fmt.Errorf("image %v was not precached", ip)
	}
	data, err := vfs.Load(purpose, name)
	if err != nil {
		return nil, fmt.Errorf("could not load: %v", err)
	}
	defer data.Close()
	img, _, err := image.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %v", err)
	}
	eImg := ebiten.NewImageFromImage(img)
	cache[ip] = eImg
	return eImg, nil
}

func Precache() {
	if !*precacheImages {
		return
	}
	for _, purpose := range []string{"tiles", "sprites"} {
		names, err := vfs.ReadDir(purpose)
		if err != nil {
			log.Panicf("could not enumerate files in %v: %v", purpose, err)
		}
		for _, name := range names {
			if !strings.HasSuffix(name, ".png") {
				continue
			}
			_, err := Load(purpose, name)
			if err != nil {
				log.Panicf("could not precache %v in %v: %v", name, purpose, err)
			}
		}
	}
	cacheFrozen = true
}

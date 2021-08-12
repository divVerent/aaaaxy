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

// +build !ebitensinglethread

package aaaaxy

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
)

var (
	debugThreadedImageDumping = flag.Bool("debug_threaded_image_dumping", false, "do image dumping in a background thread (should be faster, further boosted using -num_offscreen_images)")
)

func dumpPixelsRGBA(img *ebiten.Image, cb func(pix []byte, err error)) {
	if *debugThreadedImageDumping {
		go func() {
			pix, err := getPixelsRGBA(img)
			cb(pix, err)
		}()
	} else {
		pix, err := getPixelsRGBA(img)
		cb(pix, err)
	}
}

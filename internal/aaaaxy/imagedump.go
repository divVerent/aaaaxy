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

package aaaaxy

import (
	"fmt"
	"image/color"
	"log"
	"reflect"
	"time"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
)

var (
	debugDirtyImageDumping    = flag.Bool("debug_dirty_image_dumping", false, "use a really dirty hack (ebiten internals) to dump image pixels faster")
	debugThreadedImageDumping = flag.Bool("debug_threaded_image_dumping", false, "use a really dirty hack (ebiten internals) to dump image pixels faster")
)

var (
	dumpPixelsBufferSize = 8
	dumpPixelsBuffer     chan *ebiten.Image
)

func getPixelsRGBA(img *ebiten.Image) (pix []byte, err error) {
	bounds := img.Bounds()
	width, height := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y
	defer func(t0 time.Time) {
		log.Printf("image frame dump took %v", time.Since(t0))
	}(time.Now())
	if *debugDirtyImageDumping {
		// FASTER but dirty. May, or rather, will break in ebiten releases.
		mipMap := reflect.ValueOf(img).Elem().FieldByName("mipmap")
		// This is unexported, so we can't use it yet. Let's hack it unsafely.
		flagField, found := reflect.TypeOf(&mipMap).Elem().FieldByName("flag")
		if !found {
			return nil, fmt.Errorf("could not hack a reflect.Value exported: field 'flag' not found")
		}
		flagPtr := unsafe.Pointer(uintptr(unsafe.Pointer(&mipMap)) + flagField.Offset)
		flag := (*int)(flagPtr)
		*flag &^= (1 << 5) | (1 << 6) // Remove flagRO.
		// Now we can!
		pix, err := mipMap.Interface().(interface {
			Pixels(x, y, width, height int) ([]byte, error)
		}).Pixels(
			bounds.Min.X,
			bounds.Min.Y,
			width,
			height)
		if err != nil {
			return nil, fmt.Errorf("could not get image pixels: %v", err)
		}
		return pix, nil
	} else {
		pix := make([]byte, 4*width*height)
		p := 0
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				c := img.At(x, y).(color.RGBA)
				pix[p] = c.R
				pix[p+1] = c.G
				pix[p+2] = c.B
				pix[p+3] = c.A
				p += 4
			}
		}
		return pix, nil
	}
}

func startDumpPixelsRGBA() {
	if !*debugThreadedImageDumping {
		return
	}
	dumpPixelsBuffer = make(chan *ebiten.Image, dumpPixelsBufferSize)
	for i := 0; i < dumpPixelsBufferSize; i++ {
		dumpPixelsBuffer <- ebiten.NewImage(engine.GameWidth, engine.GameHeight)
	}
}

func dumpPixelsRGBA(img *ebiten.Image, cb func(pix []byte, err error)) {
	if *debugThreadedImageDumping {
		dup := <-dumpPixelsBuffer
		dup.DrawImage(img, &ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeCopy,
		})
		go func() {
			pix, err := getPixelsRGBA(dup)
			dumpPixelsBuffer <- dup
			cb(pix, err)
		}()
	} else {
		pix, err := getPixelsRGBA(img)
		cb(pix, err)
	}
}

func stopDumpPixelsRGBA() {
	if !*debugThreadedImageDumping {
		return
	}
	for i := 0; i < dumpPixelsBufferSize; i++ {
		<-dumpPixelsBuffer
	}
}

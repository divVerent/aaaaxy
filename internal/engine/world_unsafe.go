// Copyright 2023 Google LLC
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

package engine

import (
	"unsafe"

	"github.com/divVerent/aaaaxy/internal/level"
)

func (w *World) forEachTile(f func(i int, t *level.Tile)) {
	// Same as for i, t := range w.tiles[:], but 15% faster.
	first := &w.tiles[0]
	d := unsafe.Pointer(first)
	s := unsafe.Sizeof(first)
	for i := uintptr(0); i < tileWindowWidth*tileWindowHeight; i++ {
		t := *(**level.Tile)(unsafe.Pointer(uintptr(d) + s*i))
		if t == nil {
			continue
		}
		f(int(i), t)
	}
}

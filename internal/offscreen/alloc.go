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

package offscreen

import (
	"reflect"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	debugOffscreen           = flag.Bool("debug_offscreen", false, "log offscreen allocation summary")
	debugOffscreenEvents     = flag.Bool("debug_offscreen_events", false, "log offscreen allocation events")
	debugOffscreenFIFO       = flag.Bool("debug_offscreen_fifo", false, "reuse textures as late as possible, rather than as early as possible (MIGHT increase parallelism, costs nothing, uses different textures each frame); mutually exclusive with --debug_offscreen_by_name and --debug_offscreen_unmanaged")
	debugOffscreenSeparate   = flag.Bool("debug_offscreen_separate", false, "disable texture reuse within frame for different purpose (MIGHT increase parallelism, costs VRAM, uses different textures each frame); mutually exclusive with --debug_offscreen_by_name and --debug_offscreen_unmanaged")
	debugOffscreenByName     = flag.Bool("debug_offscreen_by_name", false, "always reuse same texture for same purpose even across frames (MIGHT increase parallelism, costs even more VRAM, very consistent texture use)")
	debugOffscreenUnmanaged  = flag.Bool("debug_offscreen_unmanaged", false, "reallocate textures every frame (probably really SLOW, may even LEAK); mutually exclusive with --debug_offscreen_by_name")
	debugOffscreenAvoidReuse = flag.Bool("debug_offscreen_avoid_reuse", false, "avoid reuse of offscreen textures within a frame for the same purpose (MIGHT increase parallelism, costs VRAM)")
)

func AvoidReuse() bool {
	return *debugOffscreenAvoidReuse
}

func id(img *ebiten.Image) (id int64) {
	// Note: this uses Ebitengine internals. This is expected to break on internal API changes.
	// Only used when *debugOffscreen or *debugOffscreenEvents is set.
	defer func() {
		r := recover()
		if r != nil {
			log.Errorf("could not get image id: operation %d failed: %v", id, r)
		}
	}()
	id = -1
	vp := reflect.ValueOf(img) // *ebiten.Image
	id = -2
	v := vp.Elem() // ebiten.Image
	id = -3
	up := v.FieldByName("image") // *ui.Image
	id = -4
	u := up.Elem() // ui.Image
	id = -5
	mp := u.FieldByName("mipmap") // *mipmap.Image
	id = -6
	m := mp.Elem() // mipmap.Mipmap
	id = -7
	bup := m.FieldByName("orig") // *buffered.Image
	id = -8
	bu := bup.Elem() // buffered.Image
	id = -9
	ap := bu.FieldByName("img") // *atlas.Image
	id = -10
	a := ap.Elem() // atlas.Image
	id = -11
	bp := a.FieldByName("backend") // *backend
	id = -12
	if bp.IsZero() { // No ID assigned yet.
		return
	}
	id = -13
	b := bp.Elem() // backend
	id = -14
	rp := b.FieldByName("restorable") // *restorable.Image
	id = -15
	r := rp.Elem() // restorable.Image
	id = -16
	gp := r.FieldByName("image") // *graphicscommand.Image
	id = -17
	g := gp.Elem() // graphicscommand.Image
	id = -18
	i := g.FieldByName("id")
	id = -19
	return i.Int()
}

type manager interface {
	New(name string, explicit bool) *ebiten.Image
	Dispose(img *ebiten.Image)
	Collect()
	Report()
}

type baseManager struct {
	w, h      int
	names     map[*ebiten.Image]string
	pastNames map[int64][]string
	allocated int
	freed     int
}

func newBaseManager(w, h int) baseManager {
	return baseManager{
		w:         w,
		h:         h,
		names:     map[*ebiten.Image]string{},
		pastNames: map[int64][]string{},
	}
}

func (m *baseManager) recordName(name string, img *ebiten.Image) *ebiten.Image {
	if *debugOffscreenEvents {
		log.Infof("offscreen: alloc %v -> %v", name, id(img))
	}
	if _, found := m.names[img]; found {
		log.Fatalf("to be created texture %v was already allocated by this", img)
	}
	m.names[img] = name
	if *debugOffscreen {
		m.pastNames[id(img)] = append(m.pastNames[id(img)], name)
	}
	m.allocated++
	return img
}

func (m *baseManager) clearName(img *ebiten.Image) string {
	name, found := m.names[img]
	if !found {
		log.Fatalf("to be disposed texture %v was never allocated by this", img)
	}
	if *debugOffscreenEvents {
		log.Infof("offscreen: dispose %v -> %v", name, id(img))
	}
	delete(m.names, img)
	m.freed++
	return name
}

func (m *baseManager) Report() {
	if *debugOffscreen {
		log.Infof("offscreen: %d textures allocated, %d textures freed, %d textures in use", m.allocated, m.freed, len(m.names))
		var ids []int64
		for id := range m.pastNames {
			ids = append(ids, id)
		}
		sort.Slice(ids, func(a, b int) bool { return ids[a] < ids[b] })
		for _, id := range ids {
			log.Infof("offscreen: texture %d was used for: %v", id, m.pastNames[id])
		}
	}
	m.allocated = 0
	m.freed = 0
	m.pastNames = map[int64][]string{}
}

type unManager struct {
	baseManager
}

func newUnManager(w, h int) *unManager {
	return &unManager{
		baseManager: newBaseManager(w, h),
	}
}

func (m *unManager) New(name string, explicit bool) *ebiten.Image {
	return m.recordName(name, ebiten.NewImage(m.w, m.h))
}

func (m *unManager) Dispose(img *ebiten.Image) {
	m.clearName(img)
	img.Dispose()
}

func (m *unManager) Collect() {}

type listManager struct {
	baseManager

	available     []*ebiten.Image
	inUse         []*ebiten.Image
	inUseExplicit []*ebiten.Image

	fifo     bool
	separate bool
}

func newListManager(w, h int) *listManager {
	return &listManager{
		baseManager: newBaseManager(w, h),
		fifo:        *debugOffscreenFIFO,
		separate:    *debugOffscreenSeparate,
	}
}

func (m *listManager) New(name string, explicit bool) *ebiten.Image {
	var img *ebiten.Image
	n := len(m.available)
	if n == 0 {
		// New.
		img = ebiten.NewImage(m.w, m.h)
	} else if m.fifo {
		// Shift.
		img = m.available[0]
		copy(m.available[0:], m.available[1:])
		m.available = m.available[:(n - 1)]
	} else {
		// Pop.
		img = m.available[n-1]
		m.available = m.available[:(n - 1)]
	}
	if explicit {
		m.inUseExplicit = append(m.inUseExplicit, img)
	} else {
		m.inUse = append(m.inUse, img)
	}
	return m.recordName(name, img)
}

func (m *listManager) Dispose(img *ebiten.Image) {
	for i, t := range m.inUse {
		if t == img {
			if m.separate {
				// No Dispose within a frame.
				// Collect does the job instead.
				return
			}
			m.clearName(img)
			m.inUse = append(m.inUse[:i], m.inUse[(i+1):]...)
			m.available = append(m.available, img)
			return
		}
	}
	for i, t := range m.inUseExplicit {
		if t == img {
			m.clearName(img)
			m.inUseExplicit = append(m.inUseExplicit[:i], m.inUseExplicit[(i+1):]...)
			m.available = append(m.available, img)
			return
		}
	}
	name := m.clearName(img)
	log.Fatalf("attempted to dispose of unmanaged texture %v", name)
}

func (m *listManager) Collect() {
	// Behaves as if all were disposed separately.
	// Usually disposing is done in reverse order of creation, so mimic that here too.
	for i := len(m.inUse) - 1; i >= 0; i-- {
		img := m.inUse[i]
		m.clearName(img)
		m.available = append(m.available, img)
	}
	m.inUse = m.inUse[:0]
}

func (m *listManager) Report() {
	m.baseManager.Report()
	if *debugOffscreen {
		log.Infof("offscreen: %d textures available, %d textures in use, %d textures explicitly in use", len(m.available), len(m.inUse), len(m.inUseExplicit))
	}
}

type byNameManager struct {
	baseManager

	byName map[string]*ebiten.Image
}

func newByNameManager(w, h int) *byNameManager {
	return &byNameManager{
		baseManager: newBaseManager(w, h),
		byName:      map[string]*ebiten.Image{},
	}
}

func (m *byNameManager) New(name string, explicit bool) *ebiten.Image {
	img, found := m.byName[name]
	if found && img == nil {
		log.Fatalf("unexpected offscreen name reuse for %v", name)
	}
	m.byName[name] = nil
	if img == nil {
		img = ebiten.NewImage(m.w, m.h)
	}
	return m.recordName(name, img)
}

func (m *byNameManager) Dispose(img *ebiten.Image) {
	name := m.clearName(img)
	other, found := m.byName[name]
	if other != nil {
		log.Fatalf("double dispose of texture %v", name)
	}
	if !found {
		log.Fatalf("attempted to dispose of unmanaged texture %v", name)
	}
	m.byName[name] = img
}

func (m *byNameManager) Collect() {}

func (m *byNameManager) Report() {
	m.baseManager.Report()
	if *debugOffscreen {
		available, inUse := 0, 0
		for _, img := range m.byName {
			if img == nil {
				inUse++
			} else {
				available++
			}
		}
		log.Infof("offscreen: %d textures available, %d textures in use", available, inUse)
	}
}

type size struct {
	w, h int
}

var (
	managers = map[size]manager{}
)

func managerForSize(w, h int) manager {
	if w != 640 || h != 360 {
		log.Fatalf("unexpected size")
	}
	key := size{w: w, h: h}
	m, found := managers[key]
	if !found {
		if *debugOffscreenUnmanaged {
			m = newUnManager(w, h)
		} else if *debugOffscreenByName {
			m = newByNameManager(w, h)
		} else {
			m = newListManager(w, h)
		}
		managers[key] = m
	}
	return m
}

func New(name string, w, h int) *ebiten.Image {
	return managerForSize(w, h).New(name, false)
}

func NewExplicit(name string, w, h int) *ebiten.Image {
	return managerForSize(w, h).New(name, true)
}

func Dispose(img *ebiten.Image) {
	w, h := img.Size()
	managerForSize(w, h).Dispose(img)
}

func Collect() {
	for _, m := range managers {
		m.Report()
		m.Collect()
	}
}

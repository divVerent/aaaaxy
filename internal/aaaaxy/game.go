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
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/audiowrap"
	"github.com/divVerent/aaaaxy/internal/demo"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/menu"
	"github.com/divVerent/aaaaxy/internal/music"
	"github.com/divVerent/aaaaxy/internal/noise"
	"github.com/divVerent/aaaaxy/internal/shader"
	"github.com/divVerent/aaaaxy/internal/timing"
)

var RegularTermination = menu.RegularTermination

var (
	externalCapture = flag.Bool("external_dump", false, "assume an external dump application like apitrace is running; makes game run in lock step with rendering")
	screenFilter    = flag.String("screen_filter", "linear2xcrt", "filter to use for rendering the screen; current possible values are 'simple', 'linear', 'linear2x', 'linear2xcrt' and 'nearest'")
	// TODO(divVerent): Remove this flag when https://github.com/hajimehoshi/ebiten/issues/1772 is resolved.
	screenFilterMaxScale    = flag.Float64("screen_filter_max_scale", 4.0, "maximum scale-up factor for the screen filter")
	screenFilterScanLines   = flag.Float64("screen_filter_scan_lines", 0.1, "strength of the scan line effect in the linear2xcrt filter")
	screenFilterCRTStrength = flag.Float64("screen_filter_crt_strength", 0.5, "strength of CRT deformation in the linear2xcrt filter")
	screenFilterJitter      = flag.Float64("screen_filter_jitter", 0.0, "for any filter other than simple, amount of jitter to add to the filter")
)

type Game struct {
	Menu menu.Controller

	offScreens        chan *ebiten.Image
	linear2xShader    *ebiten.Shader
	linear2xCRTShader *ebiten.Shader
	framesToDump      int
}

var _ ebiten.Game = &Game{}

func (g *Game) Update() error {
	g.framesToDump++

	timing.ReportRegularly()

	defer timing.Group()()
	timing.Section("update")
	defer timing.Group()()

	timing.Section("input")
	input.Update()

	timing.Section("demo_pre")
	if demo.Update() {
		log.Infof("demo playback ended, exiting")
		return RegularTermination
	}

	defer func() {
		timing.Section("demo_post")
		demo.PostUpdate(g.Menu.World.Player.Rect.Origin)
	}()

	timing.Section("menu")
	err := g.Menu.Update()
	if err != nil {
		return err
	}

	timing.Section("world")
	err = g.Menu.UpdateWorld()
	if err != nil {
		return err
	}

	// As the world's Update method may change the sound system info,
	// run this part last to reduce sound latency.

	timing.Section("noise")
	noise.Update()

	timing.Section("audiowrap")
	audiowrap.Update()

	return nil
}

func (g *Game) drawAtGameSizeThenReturnTo(screen *ebiten.Image, to chan *ebiten.Image) {
	sw, sh := screen.Size()
	if sw != engine.GameWidth || sh != engine.GameHeight {
		log.Infof("skipping frame as sizes do not match up: got %vx%v, want %vx%v",
			sw, sh, engine.GameWidth, engine.GameHeight)
		return
	}

	timing.Section("fontcache")
	font.KeepInCache(screen)

	timing.Section("world")
	g.Menu.DrawWorld(screen)

	timing.Section("menu")
	g.Menu.Draw(screen)

	timing.Section("dump")
	dumpFrameThenReturnTo(screen, to, g.framesToDump)
	g.framesToDump = 0
}

func (g *Game) drawOffscreen() *ebiten.Image {
	if g.offScreens == nil {
		n := 1
		if dumping() {
			// When dumping, cycle between two offscreen images so we can dump in the background thread.
			n = 2
		}
		g.offScreens = make(chan *ebiten.Image, n)
		for i := 0; i < n; i++ {
			g.offScreens <- ebiten.NewImage(engine.GameWidth, engine.GameHeight)
		}
	}
	offScreen := <-g.offScreens
	g.drawAtGameSizeThenReturnTo(offScreen, g.offScreens)
	// Note: following code of the draw code may still use the image, but that's OK as long as drawOffscreen() isn't called again.
	return offScreen
}

func (g *Game) setOffscreenGeoM(screen *ebiten.Image, geoM *ebiten.GeoM, w, h int) {
	sw, sh := screen.Size()
	fw := float64(sw) / float64(w)
	fh := float64(sh) / float64(h)
	f := fw
	if fh < fw {
		f = fh
	}
	dx := (float64(sw) - f*float64(w)) * 0.5
	dy := (float64(sh) - f*float64(h)) * 0.5
	geoM.Scale(f, f)
	geoM.Translate(dx, dy)
	geoM.Translate((rand.Float64()-0.5)**screenFilterJitter, (rand.Float64()-0.5)**screenFilterJitter)
}

// First two terms of the Taylor expansion of asin(strength*x)/strength.
func crtK1() float64 {
	return 1.0 / 6.0 * math.Pow(*screenFilterCRTStrength, 2)
}

func crtK2() float64 {
	return 3.0 / 40.0 * math.Pow(*screenFilterCRTStrength, 4)
}

func (g *Game) Draw(screen *ebiten.Image) {
	defer timing.Group()()
	timing.Section("draw")
	defer timing.Group()()

	switch *screenFilter {
	case "simple":
		if dumping() {
			// We're dumping, so we NEED an offscreen.
			// This is actually just like "nearest", except that to ebiten we have a game-sized and not screen-sized screen.
			// So we can use an identity matrix and need not clear the screen.
			options := &ebiten.DrawImageOptions{
				CompositeMode: ebiten.CompositeModeCopy,
				Filter:        ebiten.FilterNearest,
			}
			screen.DrawImage(g.drawOffscreen(), options)
		} else {
			// It's all sync, so just provide a dummy channel to discard it.
			g.drawAtGameSizeThenReturnTo(screen, make(chan *ebiten.Image, 1))
		}
	case "linear":
		screen.Clear()
		options := &ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeCopy,
			Filter:        ebiten.FilterLinear,
		}
		g.setOffscreenGeoM(screen, &options.GeoM, engine.GameWidth, engine.GameHeight)
		screen.DrawImage(g.drawOffscreen(), options)
	case "linear2x":
		if g.linear2xShader == nil {
			var err error
			g.linear2xShader, err = shader.Load("linear2x.kage", nil)
			if err != nil {
				log.Errorf("BROKEN RENDERER, WILL FALLBACK: could not load linear2x shader: %v", err)
				*screenFilter = "simple"
				return
			}
		}
		options := &ebiten.DrawRectShaderOptions{
			CompositeMode: ebiten.CompositeModeCopy,
			Images: [4]*ebiten.Image{
				g.drawOffscreen(),
				nil,
				nil,
				nil,
			},
		}
		g.setOffscreenGeoM(screen, &options.GeoM, engine.GameWidth, engine.GameHeight)
		screen.DrawRectShader(engine.GameWidth, engine.GameHeight, g.linear2xShader, options)
	case "linear2xcrt":
		if g.linear2xCRTShader == nil {
			var err error
			g.linear2xCRTShader, err = shader.Load("linear2xcrt.kage", nil)
			if err != nil {
				log.Errorf("BROKEN RENDERER, WILL FALLBACK: could not load linear2xcrt shader: %v", err)
				*screenFilter = "linear2x"
				return
			}
		}
		options := &ebiten.DrawRectShaderOptions{
			CompositeMode: ebiten.CompositeModeCopy,
			Images: [4]*ebiten.Image{
				g.drawOffscreen(),
				nil,
				nil,
				nil,
			},
			Uniforms: map[string]interface{}{
				"ScanLineEffect": float32(*screenFilterScanLines * 2.0),
				"CRTK1":          float32(crtK1()),
				"CRTK2":          float32(crtK2()),
			},
		}
		g.setOffscreenGeoM(screen, &options.GeoM, engine.GameWidth, engine.GameHeight)
		screen.DrawRectShader(engine.GameWidth, engine.GameHeight, g.linear2xCRTShader, options)
	case "nearest":
		screen.Clear()
		options := &ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeCopy,
			Filter:        ebiten.FilterNearest,
		}
		g.setOffscreenGeoM(screen, &options.GeoM, engine.GameWidth, engine.GameHeight)
		screen.DrawImage(g.drawOffscreen(), options)
	default:
		log.Errorf("WARNING: unknown screen filter type: %q; reverted to simple", *screenFilter)
		*screenFilter = "simple"
	}

	// Once this has run, we can start fading in music.
	music.Enable()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if *screenFilter == "simple" {
		return engine.GameWidth, engine.GameHeight
	}
	d := ebiten.DeviceScaleFactor()
	// TODO: when https://github.com/hajimehoshi/ebiten/issues/1772 is resolved,
	// change this back to int(float64(outsideWidth) * d), int(float64(outsideHeight) * d).
	f := math.Min(
		math.Min(
			float64(outsideWidth)*d/engine.GameWidth,
			float64(outsideHeight)*d/engine.GameHeight),
		*screenFilterMaxScale)
	return int(engine.GameWidth * f), int(engine.GameHeight * f)
}

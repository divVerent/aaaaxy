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

package aaaaaa

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/flag"
	"github.com/divVerent/aaaaaa/internal/font"
	"github.com/divVerent/aaaaaa/internal/input"
	"github.com/divVerent/aaaaaa/internal/menu"
	"github.com/divVerent/aaaaaa/internal/music"
	"github.com/divVerent/aaaaaa/internal/noise"
	"github.com/divVerent/aaaaaa/internal/shader"
	"github.com/divVerent/aaaaaa/internal/timing"
)

var RegularTermination = menu.RegularTermination

var (
	externalCapture    = flag.Bool("external_dump", false, "assume an external dump application like apitrace is running; makes game run in lock step with rendering")
	screenFilter       = flag.String("screen_filter", "linear2x", "filter to use for rendering the screen; current possible values are 'simple', 'linear', 'linear2x' and 'nearest'")
	screenFilterJitter = flag.Float64("screen_filter_jitter", 0.0, "for any filter other than simple, amount of jitter to add to the filter")
)

type Game struct {
	Menu menu.Menu

	offScreen      *ebiten.Image
	linear2xShader *ebiten.Shader
}

var _ ebiten.Game = &Game{}

func (g *Game) Update() error {
	timing.ReportRegularly()

	defer timing.Group()()
	timing.Section("update")
	defer timing.Group()()

	timing.Section("input")
	input.Update()

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
	timing.Section("music")
	music.Update()

	timing.Section("noise")
	noise.Update()

	return nil
}

func (g *Game) drawAtGameSize(screen *ebiten.Image) {
	timing.Section("fontcache")
	font.KeepInCache(screen)

	timing.Section("world")
	g.Menu.DrawWorld(screen)

	timing.Section("menu")
	g.Menu.Draw(screen)

	timing.Section("dump")
	dumpFrame(screen)
}

func (g *Game) drawOffscreen() *ebiten.Image {
	if g.offScreen == nil {
		g.offScreen = ebiten.NewImage(engine.GameWidth, engine.GameHeight)
	}
	g.drawAtGameSize(g.offScreen)
	return g.offScreen
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

func (g *Game) Draw(screen *ebiten.Image) {
	defer timing.Group()()
	timing.Section("draw")
	defer timing.Group()()

	switch *screenFilter {
	case "simple":
		g.drawAtGameSize(screen)
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
				log.Printf("BROKEN RENDERER, WILL FALLBACK: could not load linear2x shader: %v", err)
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
	case "nearest":
		screen.Clear()
		options := &ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeCopy,
			Filter:        ebiten.FilterNearest,
		}
		g.setOffscreenGeoM(screen, &options.GeoM, engine.GameWidth, engine.GameHeight)
		screen.DrawImage(g.drawOffscreen(), options)
	default:
		log.Printf("WARNING: unknown screen filter type: %q; reverted to simple", *screenFilter)
		*screenFilter = "simple"
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if *screenFilter != "simple" {
		return outsideWidth, outsideHeight
	}
	return engine.GameWidth, engine.GameHeight
}

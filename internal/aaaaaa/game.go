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

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/flag"
	"github.com/divVerent/aaaaaa/internal/font"
	"github.com/divVerent/aaaaaa/internal/input"
	"github.com/divVerent/aaaaaa/internal/menu"
	"github.com/divVerent/aaaaaa/internal/music"
	"github.com/divVerent/aaaaaa/internal/noise"
	"github.com/divVerent/aaaaaa/internal/timing"
)

var RegularTermination = menu.RegularTermination

var (
	externalCapture = flag.Bool("external_dump", false, "assume an external dump application like apitrace is running; makes game run in lock step with rendering")
	screenFilter    = flag.String("screen_filter", "default", "filter to use for rendering the screen; current possible values are 'default', 'linear' and 'nearest'")
)

type Game struct {
	Menu menu.Menu

	offScreen *ebiten.Image
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

func (g *Game) setOffscreenGeoM(screen *ebiten.Image, geoM *ebiten.GeoM) {
	sw, sh := screen.Size()
	fw := float64(sw) / float64(engine.GameWidth)
	fh := float64(sh) / float64(engine.GameHeight)
	f := fw
	if fh < fw {
		f = fh
	}
	dx := (float64(sw) - f*engine.GameWidth) * 0.5
	dy := (float64(sh) - f*engine.GameHeight) * 0.5
	geoM.Scale(f, f)
	geoM.Translate(dx, dy)
}

func (g *Game) Draw(screen *ebiten.Image) {
	defer timing.Group()()
	timing.Section("draw")
	defer timing.Group()()

	switch *screenFilter {
	case "default":
		g.drawAtGameSize(screen)
	case "linear":
		screen.Clear()
		options := &ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeCopy,
			Filter:        ebiten.FilterLinear,
		}
		g.setOffscreenGeoM(screen, &options.GeoM)
		screen.DrawImage(g.drawOffscreen(), options)
	case "nearest":
		screen.Clear()
		options := &ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeCopy,
			Filter:        ebiten.FilterNearest,
		}
		g.setOffscreenGeoM(screen, &options.GeoM)
		screen.DrawImage(g.drawOffscreen(), options)
		// TODO: Add some shader based filters.
	default:
		log.Printf("WARNING: unknown screen filter type: %q; reverted to default", *screenFilter)
		*screenFilter = "default"
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if *screenFilter != "default" {
		return outsideWidth, outsideHeight
	}
	return engine.GameWidth, engine.GameHeight
}

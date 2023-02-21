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
	"errors"
	"fmt"
	go_image "image"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/audiowrap"
	"github.com/divVerent/aaaaxy/internal/demo"
	"github.com/divVerent/aaaaxy/internal/dump"
	"github.com/divVerent/aaaaxy/internal/engine"
	"github.com/divVerent/aaaaxy/internal/exitstatus"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/font"
	"github.com/divVerent/aaaaxy/internal/fun"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/menu"
	"github.com/divVerent/aaaaxy/internal/music"
	"github.com/divVerent/aaaaxy/internal/noise"
	"github.com/divVerent/aaaaxy/internal/offscreen"
	"github.com/divVerent/aaaaxy/internal/palette"
	"github.com/divVerent/aaaaxy/internal/shader"
	"github.com/divVerent/aaaaxy/internal/timing"
)

var (
	screenFilter = flag.String("screen_filter", flag.SystemDefault(map[string]string{
		"android/*": "linear2x",
		"js/*":      "linear2x",
		"*/*":       "linear2xcrt",
	}), "filter to use for rendering the screen; current possible values are 'nearest', 'linear', 'linear2x' and 'linear2xcrt'")
	screenFilterScanLines   = flag.Float64("screen_filter_scan_lines", 0.1, "strength of the scan line effect in the linear2xcrt filters")
	screenFilterCRTStrength = flag.Float64("screen_filter_crt_strength", 0.5, "strength of CRT deformation in the linear2xcrt filters")
	screenStretch           = flag.Bool("screen_stretch", false, "stretch screen content instead of letterboxing")
	paletteFlag             = flag.String("palette", flag.SystemDefault(map[string]string{
		"android/*": "none",
		"js/*":      "none",
		"*/*":       "vga",
	}), "render with palette; can be set to '"+strings.Join(palette.Names(), "', '")+"' or 'none'")
	paletteRemapOnly          = flag.Bool("palette_remap_only", false, "only apply the palette's color remapping, do not actually reduce color set")
	paletteRemapColors        = flag.Bool("palette_remap_colors", true, "remap input colors to close palette colors on load (less dither but wrong colors)")
	paletteDitherSize         = flag.Int("palette_dither_size", 4, "dither pattern size (really should be a power of two when using the bayer dither mode)")
	paletteDitherMode         = flag.String("palette_dither_mode", "plastic2", "dither type (none, bayer, bayer2, halftone, halftone2, plastic, plastic2, random or random2)")
	paletteDitherWorldAligned = flag.Bool("palette_dither_world_aligned", true, "align dither pattern to world as opposed to screen")
	debugEnableDrawing        = flag.Bool("debug_enable_drawing", true, "enable drawing the display; set to false for faster demo processing or similar")
	showFPS                   = flag.Bool("show_fps", false, "show fps counter")
	showTime                  = flag.Bool("show_time", false, "show game time")
)

type ditherMode int

const (
	bayerDither ditherMode = iota
	bayer2Dither
	halftoneDither
	halftone2Dither
	plasticDither
	plastic2Dither
	randomDither
	random2Dither
)

type Game struct {
	Menu menu.Controller

	init      initState
	canUpdate bool
	canDraw   bool

	// screenWidth and screenHeight are updated by Layout().
	screenWidth  int
	screenHeight int

	offscreenTokens   chan int
	offscreenReturns  chan *ebiten.Image
	offscreenIndexes  map[*ebiten.Image]int
	linear2xShader    *ebiten.Shader
	linear2xCRTShader *ebiten.Shader

	// Copies of parameters so we know when to update.
	palette           *palette.Palette
	paletteDitherSize int
	paletteDitherMode ditherMode

	paletteLUT       *ebiten.Image  // Updates when palette changes.
	paletteLUTSize   int            // Updates when palette changes.
	paletteLUTPerRow int            // Updates when palette changes.
	paletteLUTWidth  int            // Updates when palette changes.
	paletteBayern    []float32      // Updates when palette or paletteDitherSize change.
	paletteShader    *ebiten.Shader // Updates when paletteDitherSize changes.

	framesToDump int
}

var _ ebiten.Game = &Game{}

func NewGame() *Game {
	return &Game{
		offscreenIndexes: map[*ebiten.Image]int{},
	}
}

func (g *Game) updateFrame() error {
	timing.Section("input")
	input.Update(g.screenWidth, g.screenHeight, engine.GameWidth, engine.GameHeight, crtK1(), crtK2())

	timing.Section("demo_pre")
	if demo.Update() {
		log.Infof("demo playback ended, exiting")
		return exitstatus.RegularTermination
	}

	defer func() {
		timing.Section("demo_post")
		if g.Menu.World.Player != nil {
			demo.PostUpdate(g.Menu.World.Player.Rect.Origin)
		}
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

func (g *Game) Update() error {
	if !g.canUpdate {
		return nil
	}

	if !g.init.done {
		return g.InitStep()
	}
	g.canDraw = true

	g.framesToDump++

	timing.Update()

	defer timing.Group()()
	timing.Section("update")

	defer timing.Group()()

	for frame := 0; frame < *fpsDivisor; frame++ {
		if err := g.updateFrame(); err != nil {
			if errors.Is(err, exitstatus.RegularTermination) {
				log.Infof("exiting normally")
			} else {
				log.Infof("exiting due to: %v", err)
			}
			return err
		}
	}

	return nil
}

func (g *Game) palettePrepare(maybeScreen *ebiten.Image, tmp *ebiten.Image) (*ebiten.Image, func() *ebiten.Image) {
	// This is an extra pass so it can still run at low-res.
	pal := palette.ByName(*paletteFlag)

	if pal == nil {
		// No palette.
		*paletteFlag = "none"
		screen := g.maybeAcquireOffscreen(maybeScreen)
		return screen, func() *ebiten.Image { return screen }
	}

	if *paletteRemapOnly {
		// Color reduction disabled.
		screen := g.maybeAcquireOffscreen(maybeScreen)
		return screen, func() *ebiten.Image { return screen }
	}

	// Shaders depend on Bayer pattern size, and this should usually not change at runtime.
	ditherSize := *paletteDitherSize
	if ditherSize < 2 {
		*paletteDitherSize = 2
		ditherSize = 2
	}

	var ditherMode ditherMode
	switch *paletteDitherMode {
	case "none":
		// No dither is the same as a 1x1 Bayer dither.
		// That way, we can use the same shader.
		ditherMode = bayerDither
		ditherSize = 1
	case "bayer":
		ditherMode = bayerDither
	case "halftone":
		ditherMode = halftoneDither
	case "bayer2":
		ditherMode = bayer2Dither
	case "halftone2":
		ditherMode = halftone2Dither
	case "plastic":
		ditherMode = plasticDither
		ditherSize = 0
	case "plastic2":
		ditherMode = plastic2Dither
		ditherSize = 0
	case "random":
		ditherMode = randomDither
		ditherSize = 0
	case "random2":
		ditherMode = random2Dither
		ditherSize = 0
	default:
		log.Errorf("unknown dither mode %v, switching to bayer", *paletteDitherMode)
		*paletteDitherMode = "bayer"
		ditherMode = bayerDither
	}

	// Need images?
	if g.paletteLUT == nil {
		g.paletteLUT = ebiten.NewImage(engine.GameWidth, engine.GameHeight)
	}

	// Bayer pattern changed?
	if ditherSize != g.paletteDitherSize || g.paletteDitherMode != ditherMode {
		if g.paletteShader != nil {
			g.paletteShader.Dispose()
		}
		g.paletteShader = nil
	}

	// Need a new shader?
	if g.paletteShader == nil {
		var err error
		params := map[string]interface{}{}
		switch ditherMode {
		case bayerDither, halftoneDither:
			params["BayerSize"] = ditherSize
		case bayer2Dither, halftone2Dither:
			params["BayerSize"] = ditherSize
			params["TwoColor"] = true
		case plasticDither:
			params["PlasticDither"] = true
		case plastic2Dither:
			params["PlasticDither"] = true
			params["TwoColor"] = true
		case randomDither:
			params["RandomDither"] = true
		case random2Dither:
			params["RandomDither"] = true
			params["TwoColor"] = true
		}
		g.paletteShader, err = shader.Load("dither.kage.tmpl", params)
		if err != nil {
			log.Errorf("BROKEN RENDERER, WILL FALLBACK: could not load palette shader for dither size %d: %v", *paletteDitherSize, err)
			*paletteFlag = "none"
			screen := g.maybeAcquireOffscreen(maybeScreen)
			return screen, func() *ebiten.Image { return screen }
		}
		g.paletteDitherSize = ditherSize
		g.paletteDitherMode = ditherMode
		g.palette = nil
	}

	// Need a LUT?
	if g.palette != pal {
		var lut go_image.Image
		switch ditherMode {
		case bayerDither, halftoneDither, randomDither, plasticDither:
			lut, g.paletteLUTSize, g.paletteLUTPerRow, g.paletteLUTWidth = pal.ToLUT(g.paletteLUT.Bounds(), 1)
		case bayer2Dither, halftone2Dither, random2Dither, plastic2Dither:
			lut, g.paletteLUTSize, g.paletteLUTPerRow, g.paletteLUTWidth = pal.ToLUT(g.paletteLUT.Bounds(), 2)
		}
		if nrgba, ok := lut.(*go_image.NRGBA); ok {
			g.paletteLUT.SubImage(nrgba.Rect).(*ebiten.Image).ReplacePixels(nrgba.Pix)
		} else {
			log.Fatalf("palette LUT isn't NRGBA, got %T, please fix game data", lut)
		}
		switch ditherMode {
		case bayerDither, bayer2Dither:
			g.paletteBayern = pal.BayerPattern(g.paletteDitherSize)
		case halftoneDither, halftone2Dither:
			g.paletteBayern = pal.HalftonePattern(g.paletteDitherSize)
		case randomDither, random2Dither, plasticDither, plastic2Dither:
			g.paletteBayern = nil
		}
		g.palette = pal
	}

	paletteOffscreen := tmp
	if tmp == nil {
		paletteOffscreen = offscreen.New("PaletteOffscreen", engine.GameWidth, engine.GameHeight)
	}

	return paletteOffscreen, func() *ebiten.Image {
		var scroll m.Delta
		if *paletteDitherWorldAligned {
			scroll = g.Menu.World.ScrollPos().Delta(m.Pos{X: engine.GameWidth / 2, Y: engine.GameHeight / 2})
			if ditherSize > 0 {
				scroll = scroll.Mod(ditherSize)
			}
		}
		options := &ebiten.DrawRectShaderOptions{
			CompositeMode: ebiten.CompositeModeCopy,
			Images: [4]*ebiten.Image{
				paletteOffscreen,
				g.paletteLUT,
				nil,
				nil,
			},
			Uniforms: map[string]interface{}{
				"LUTSize":   float32(g.paletteLUTSize),
				"LUTPerRow": float32(g.paletteLUTPerRow),
				"LUTWidth":  float32(g.paletteLUTWidth),
				"Offset": []float32{
					float32(scroll.DX),
					float32(scroll.DY),
				},
			},
		}
		if ditherSize > 0 {
			options.Uniforms["Bayern"] = g.paletteBayern
		}
		screen := g.maybeAcquireOffscreen(maybeScreen)
		screen.DrawRectShader(engine.GameWidth, engine.GameHeight, g.paletteShader, options)
		if tmp == nil {
			offscreen.Dispose(paletteOffscreen)
		}
		return screen
	}
}

func (g *Game) drawAtGameSizeThenReturnTo(maybeScreen *ebiten.Image, to chan *ebiten.Image, tmp *ebiten.Image) *ebiten.Image {
	drawDest, finishDrawing := g.palettePrepare(maybeScreen, tmp)

	sw, sh := drawDest.Size()
	if sw != engine.GameWidth || sh != engine.GameHeight {
		log.Infof("skipping frame as sizes do not match up: got %vx%v, want %vx%v",
			sw, sh, engine.GameWidth, engine.GameHeight)
		screen := finishDrawing()
		to <- screen
		return screen
	}

	timing.Section("fontcache")
	font.KeepInCache(drawDest)

	if !g.canDraw {
		text, fraction := g.init.Current()
		bg := palette.EGA(palette.Blue, uint8(m.Rint(255*(1-fraction))))
		fg := palette.EGA(palette.LightGrey, 255)
		ol := palette.EGA(palette.Black, 255)
		drawDest.Fill(bg)
		if font.ByName["MenuSmall"].Face != nil && text != "" {
			r := font.ByName["MenuSmall"].BoundString(text)
			y := m.Rint(float64((engine.GameHeight-r.Size.DY))*(1-fraction)) - r.Origin.Y
			font.ByName["MenuSmall"].Draw(drawDest, text, m.Pos{
				X: engine.GameWidth / 2,
				Y: y,
			}, font.Center, fg, ol)
		}
		screen := finishDrawing()
		to <- screen
		return screen
	}

	timing.Section("world")
	g.Menu.DrawWorld(drawDest)

	timing.Section("menu")
	g.Menu.Draw(drawDest)

	timing.Section("input")
	input.Draw(drawDest)

	timing.Section("global_overlays")
	if *showFPS {
		timing.Section("fps")
		font.ByName["DebugSmall"].Draw(drawDest,
			locale.G.Get("%.1f fps, %.1f tps", ebiten.CurrentFPS(), ebiten.CurrentTPS()),
			m.Pos{X: engine.GameWidth - 1, Y: engine.GameHeight - 4}, font.Right,
			palette.EGA(palette.White, 255), palette.EGA(palette.Black, 255))
	}
	if *showTime {
		timing.Section("time")
		font.ByName["DebugSmall"].Draw(drawDest,
			fun.FormatText(&g.Menu.World.PlayerState, "{{GameTime}}"),
			m.Pos{X: 0, Y: engine.GameHeight - 4}, font.Left,
			palette.EGA(palette.White, 255), palette.EGA(palette.Black, 255))
	}

	timing.Section("demo_postdraw")
	demo.PostDraw(drawDest)

	timing.Section("dump")
	screen := finishDrawing()
	dump.ProcessFrameThenReturnTo(screen, to, g.framesToDump)
	g.framesToDump = 0

	// Once this has run, we can start fading in music.
	music.Enable()

	return screen
}

func (g *Game) maybeAcquireOffscreen(screen *ebiten.Image) *ebiten.Image {
	if screen != nil {
		return screen
	}
	i := <-g.offscreenTokens
	offscreen := offscreen.NewExplicit(fmt.Sprintf("Offscreen.%d", i), engine.GameWidth, engine.GameHeight)
	g.offscreenIndexes[offscreen] = i
	return offscreen
}

func (g *Game) drawOffscreen(tmp *ebiten.Image) *ebiten.Image {
	if g.offscreenTokens == nil {
		n := 1
		if dump.Active() {
			// When dumping, cycle between two offscreen images so we can dump in the background thread.
			n = 2
		}
		g.offscreenTokens = make(chan int, n)
		for i := 0; i < n; i++ {
			g.offscreenTokens <- i
		}
		g.offscreenReturns = make(chan *ebiten.Image, n)
	}
	offscreen := g.drawAtGameSizeThenReturnTo(nil, g.offscreenReturns, tmp)
	// Note: following code of the draw code may still use the image, but that's OK as long as drawOffscreen() isn't called again.
	return offscreen
}

// First two terms of the Taylor expansion of asin(strength*x)/strength.
func crtK1() float64 {
	if *screenFilter != "linear2xcrt" {
		return 0
	}
	return 1.0 / 6.0 * math.Pow(*screenFilterCRTStrength, 2)
}

func crtK2() float64 {
	if *screenFilter != "linear2xcrt" {
		return 0
	}
	return 3.0 / 40.0 * math.Pow(*screenFilterCRTStrength, 4)
}

func (g *Game) Draw(screen *ebiten.Image) {
	w, h := screen.Size()
	if w != engine.GameWidth || h != engine.GameHeight {
		// NOTE: This implies *screenStretch.
		screen = screen.SubImage(go_image.Rectangle{
			Min: go_image.Point{},
			Max: go_image.Point{X: engine.GameWidth, Y: engine.GameHeight},
		}).(*ebiten.Image)
	}

	defer timing.Group()()
	timing.Section("draw")
	defer timing.Group()()

DoneDisposing:
	for {
		select {
		case off := <-g.offscreenReturns:
			offscreen.Dispose(off)
			g.offscreenTokens <- g.offscreenIndexes[off]
			delete(g.offscreenIndexes, off)
		default:
			break DoneDisposing
		}
	}
	offscreen.Collect()

	if !*debugEnableDrawing {
		return
	}

	if !dump.Active() {
		// No offscreen needed. Just render.
		g.drawAtGameSizeThenReturnTo(screen, make(chan *ebiten.Image, 1), nil)
		return
	}

	// When dumping, draw to an offscreen, hand it to the dumper and also copy it to the screen.

	var tmp *ebiten.Image
	if !offscreen.AvoidReuse() {
		tmp = screen
	}
	srcImage := g.drawOffscreen(tmp)
	options := &ebiten.DrawImageOptions{
		CompositeMode: ebiten.CompositeModeCopy,
		Filter:        ebiten.FilterNearest,
	}
	screen.DrawImage(srcImage, options)
}

func (g *Game) DrawFinalScreen(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
	w, h := offscreen.Size()
	if w != engine.GameWidth || h != engine.GameHeight {
		// NOTE: This implies *screenStretch.
		offscreen = offscreen.SubImage(go_image.Rectangle{
			Min: go_image.Point{},
			Max: go_image.Point{X: engine.GameWidth, Y: engine.GameHeight},
		}).(*ebiten.Image)
	}

	if *screenStretch {
		// Note that due to the code in Layout(), this changes almost nothing;
		// differences are 1 pixel or less.
		// Doing this override anyway to remove possible small black bars on some displays.

		ssz := screen.Bounds().Size()
		sw, sh := ssz.X, ssz.Y
		fw := float64(sw) / float64(engine.GameWidth)
		fh := float64(sh) / float64(engine.GameHeight)
		geoM.Reset()
		geoM.Scale(fw, fh)
	}

	switch *screenFilter {
	case "nearest":
		// Normal nearest blitting.
		options := &ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeCopy,
			Filter:        ebiten.FilterNearest,
			GeoM:          geoM,
		}
		screen.DrawImage(offscreen, options)
	case "linear":
		// Normal linear scaling.
		options := &ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeCopy,
			Filter:        ebiten.FilterLinear,
			GeoM:          geoM,
		}
		screen.DrawImage(offscreen, options)
	case "linear2x":
		if g.linear2xShader == nil {
			var err error
			g.linear2xShader, err = shader.Load("linear2xcrt.kage.tmpl", map[string]interface{}{
				"CRT": false,
			})
			if err != nil {
				log.Errorf("BROKEN RENDERER, WILL FALLBACK: could not load linear2x shader: %v", err)
				*screenFilter = "linear"
				return
			}
		}
		options := &ebiten.DrawRectShaderOptions{
			CompositeMode: ebiten.CompositeModeCopy,
			Images: [4]*ebiten.Image{
				offscreen,
				nil,
				nil,
				nil,
			},
			GeoM: geoM,
		}
		screen.DrawRectShader(engine.GameWidth, engine.GameHeight, g.linear2xShader, options)
	case "linear2xcrt":
		if g.linear2xCRTShader == nil {
			var err error
			g.linear2xCRTShader, err = shader.Load("linear2xcrt.kage.tmpl", map[string]interface{}{
				"CRT": true,
			})
			if err != nil {
				log.Errorf("BROKEN RENDERER, WILL FALLBACK: could not load linear2xcrt shader: %v", err)
				*screenFilter = "linear2x"
				return
			}
		}
		options := &ebiten.DrawRectShaderOptions{
			CompositeMode: ebiten.CompositeModeCopy,
			Images: [4]*ebiten.Image{
				offscreen,
				nil,
				nil,
				nil,
			},
			Uniforms: map[string]interface{}{
				"ScanLineEffect": float32(*screenFilterScanLines * 2.0),
				"CRTK1":          float32(crtK1()),
				"CRTK2":          float32(crtK2()),
			},
			GeoM: geoM,
		}
		screen.DrawRectShader(engine.GameWidth, engine.GameHeight, g.linear2xCRTShader, options)
	default:
		log.Errorf("unknown screen filter type: %q; reverted to simple", *screenFilter)
		*screenFilter = "linear2x"
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.screenWidth = engine.GameWidth
	g.screenHeight = engine.GameHeight
	if *screenStretch {
		if g.screenWidth*outsideHeight > g.screenHeight*outsideWidth {
			g.screenHeight = g.screenWidth * outsideHeight / outsideWidth
		} else {
			g.screenWidth = g.screenHeight * outsideWidth / outsideHeight
		}
	}
	g.canUpdate = true
	return g.screenWidth, g.screenHeight
}

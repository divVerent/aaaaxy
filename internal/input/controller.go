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

package input

import (
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/m"
)

type ImpulseState struct {
	Held    bool `json:",omitempty"`
	JustHit bool `json:",omitempty"`
}

func (i *ImpulseState) Empty() bool {
	return !i.Held && !i.JustHit
}

func (i *ImpulseState) OrEmpty() ImpulseState {
	if i == nil {
		return ImpulseState{}
	}
	return *i
}

func (i *ImpulseState) UnlessEmpty() *ImpulseState {
	if i.Empty() {
		return nil
	}
	return i
}

type InputMap int

func (i InputMap) ContainsAny(o InputMap) bool {
	return i&o != 0
}

type impulse struct {
	ImpulseState
	Name string

	keys              map[ebiten.Key]InputMap
	padControls       padControls
	mouseControl      bool
	touchRect         *m.Rect
	touchImage        *ebiten.Image
	externallyPressed bool
}

const (
	NoInput InputMap = 0

	// Allocated input bits.
	DOSKeyboardWithEscape    InputMap = 1
	NESKeyboardWithEscape    InputMap = 2
	FPSKeyboardWithEscape    InputMap = 4
	ViKeyboardWithEscape     InputMap = 8
	Gamepad                  InputMap = 16
	DOSKeyboardWithBackspace InputMap = 32
	NESKeyboardWithBackspace InputMap = 64
	FPSKeyboardWithBackspace InputMap = 128
	ViKeyboardWithBackspace  InputMap = 256
	Touchscreen              InputMap = 512

	// Computed helpers values.
	AnyKeyboardWithEscape    = DOSKeyboardWithEscape | NESKeyboardWithEscape | FPSKeyboardWithEscape | ViKeyboardWithEscape
	AnyKeyboardWithBackspace = DOSKeyboardWithBackspace | NESKeyboardWithBackspace | FPSKeyboardWithBackspace | ViKeyboardWithBackspace
	DOSKeyboard              = DOSKeyboardWithEscape | DOSKeyboardWithBackspace
	NESKeyboard              = NESKeyboardWithEscape | NESKeyboardWithBackspace
	FPSKeyboard              = FPSKeyboardWithEscape | FPSKeyboardWithBackspace
	ViKeyboard               = ViKeyboardWithEscape | ViKeyboardWithBackspace
	AnyKeyboard              = AnyKeyboardWithEscape | AnyKeyboardWithBackspace
	AnyInput                 = AnyKeyboard | Gamepad | Touchscreen
)

var (
	Left       = (&impulse{Name: "Left", keys: leftKeys, padControls: leftPad, touchRect: touchRectLeft}).register()
	Right      = (&impulse{Name: "Right", keys: rightKeys, padControls: rightPad, touchRect: touchRectRight}).register()
	Up         = (&impulse{Name: "Up", keys: upKeys, padControls: upPad, touchRect: touchRectUp}).register()
	Down       = (&impulse{Name: "Down", keys: downKeys, padControls: downPad, touchRect: touchRectDown}).register()
	Jump       = (&impulse{Name: "Jump", keys: jumpKeys, padControls: jumpPad, touchRect: touchRectJump}).register()
	Action     = (&impulse{Name: "Action", keys: actionKeys, padControls: actionPad, touchRect: touchRectAction}).register()
	Exit       = (&impulse{Name: "Exit", keys: exitKeys, padControls: exitPad, mouseControl: true, touchRect: touchRectExit}).register()
	Fullscreen = (&impulse{Name: "Fullscreen", keys: fullscreenKeys /* no padControls */}).register()

	impulses = []*impulse{}

	inputMap InputMap

	// Wait for first frame to detect initial gamepad situation.
	firstUpdate = true

	// Current mouse/finger hover pos, if any.
	hoverPos *m.Pos

	// Last mouse/finger click/release pos, if any.
	clickPos *m.Pos
)

func (i *impulse) register() *impulse {
	impulses = append(impulses, i)
	return i
}

func (i *impulse) update() {
	keyboardHolders := i.keyboardPressed()
	gamepadHolders := i.gamepadPressed()
	touchHolders := i.touchPressed()
	mouseHolders := i.mousePressed()
	holders := keyboardHolders | gamepadHolders | touchHolders | mouseHolders
	held := holders != NoInput || i.externallyPressed
	if held && !i.Held {
		i.JustHit = true
		// Whenever a new key is pressed, update the flag whether we're actually
		// _using_ the gamepad. Used for some in-game text messages.
		if holders != NoInput {
			inputMap &= holders
		}
		if inputMap == NoInput {
			inputMap = holders
		}
		// Hide mouse pointer if using another input device in the menu.
		if mouseHolders == NoInput {
			mouseCancel()
		}
	} else {
		i.JustHit = false
	}
	i.Held = held
	i.externallyPressed = false
}

func Init() error {
	gamepadInit()
	return touchInit()
}

func Update(screenWidth, screenHeight, gameWidth, gameHeight int, crtK1, crtK2, borderStretchPower float64) {
	gamepadScan()
	if firstUpdate {
		// At first, assume gamepad whenever one is present.
		switch {
		case len(gamepads) > 0:
			inputMap = Gamepad
		case runtime.GOOS == "android":
			inputMap = Touchscreen
		case runtime.GOOS == "ios":
			inputMap = Touchscreen
		case runtime.GOOS == "js":
			inputMap = Touchscreen
		default:
			inputMap = AnyKeyboard
		}
		firstUpdate = false
	}
	clickPos, hoverPos = nil, nil
	mouseUpdate(screenWidth, screenHeight, gameWidth, gameHeight, crtK1, crtK2, borderStretchPower)
	touchUpdate(screenWidth, screenHeight, gameWidth, gameHeight, crtK1, crtK2, borderStretchPower)
	for _, i := range impulses {
		i.update()
	}
	easterEggUpdate()
}

type Mode int

const (
	PlayingMode Mode = iota
	EndingMode
	MenuMode
	TouchEditMode
)

func SetMode(mode Mode) {
	switch mode {
	case PlayingMode:
		mouseSetWantClicks(false)
		touchSetUsePad(true)
		touchSetShowPad(true)
		touchSetEditor(false)
	case EndingMode:
		mouseSetWantClicks(false)
		touchSetUsePad(true)
		touchSetShowPad(false)
		touchSetEditor(false)
	case MenuMode:
		mouseSetWantClicks(true)
		touchSetUsePad(false)
		touchSetShowPad(false)
		touchSetEditor(false)
	case TouchEditMode:
		mouseSetWantClicks(true)
		touchSetUsePad(false)
		touchSetShowPad(false)
		touchSetEditor(true)
	}
}

func EasterEggJustHit() bool {
	return easterEgg.justHit || snesEasterEgg.justHit
}

func KonamiCodeJustHit() bool {
	return konamiCode.justHit || snesKonamiCode.justHit || kbdKonamiCode.justHit || literalKbdKonamiCode.justHit
}

type ExitButtonID int

const (
	Start ExitButtonID = iota
	Back
	Escape
	Backspace
)

func ExitButton() ExitButtonID {
	if inputMap.ContainsAny(Gamepad) {
		return Start
	}
	if inputMap.ContainsAny(Touchscreen) {
		return Back
	}
	if runtime.GOOS != "js" {
		// On JS, the Esc key is kinda "reserved" for leaving fullsreeen.
		// Thus we never recommend it, even if the user used it before.
		if inputMap.ContainsAny(AnyKeyboardWithEscape) {
			return Escape
		}
	}
	return Backspace
}

type ActionButtonID int

const (
	BX ActionButtonID = iota
	Elsewhere
	B
	CtrlShift
	Z
	ShiftETab
	EnterShift
)

func ActionButton() ActionButtonID {
	if inputMap.ContainsAny(Gamepad) {
		return BX
	}
	if inputMap.ContainsAny(Touchscreen) {
		if Action.touchRect.Size.IsZero() {
			return Elsewhere
		}
		return B
	}
	if inputMap.ContainsAny(DOSKeyboard) {
		return CtrlShift
	}
	if inputMap.ContainsAny(NESKeyboard) {
		return Z
	}
	if inputMap.ContainsAny(FPSKeyboard) {
		return ShiftETab
	}
	if inputMap.ContainsAny(ViKeyboard) {
		return EnterShift
	}
	// Should never hit this.
	return CtrlShift
}

func HoverPos() (m.Pos, bool) {
	if hoverPos == nil {
		return m.Pos{}, false
	}
	return *hoverPos, true
}

func ClickPos() (m.Pos, bool) {
	if clickPos == nil {
		return m.Pos{}, false
	}
	return *clickPos, true
}

func CancelHover() {
	mouseCancel()
}

type MouseStatus int

const (
	NoMouse MouseStatus = iota
	HoveringMouse
	ClickingMouse
)

func Mouse() (m.Pos, MouseStatus) {
	if clickPos != nil {
		return *clickPos, ClickingMouse
	}
	if hoverPos != nil {
		return *hoverPos, HoveringMouse
	}
	return m.Pos{}, NoMouse
}

// Demo code.

type DemoState struct {
	InputMap          InputMap      `json:",omitempty"`
	Left              *ImpulseState `json:",omitempty"`
	Right             *ImpulseState `json:",omitempty"`
	Up                *ImpulseState `json:",omitempty"`
	Down              *ImpulseState `json:",omitempty"`
	Jump              *ImpulseState `json:",omitempty"`
	Action            *ImpulseState `json:",omitempty"`
	Exit              *ImpulseState `json:",omitempty"`
	HoverPos          *m.Pos        `json:",omitempty"`
	ClickPos          *m.Pos        `json:",omitempty"`
	EasterEggJustHit  bool          `json:",omitempty"`
	KonamiCodeJustHit bool          `json:",omitempty"`
}

func LoadFromDemo(state *DemoState) {
	if state == nil {
		state = &DemoState{}
	}
	inputMap = state.InputMap
	Left.ImpulseState = state.Left.OrEmpty()
	Right.ImpulseState = state.Right.OrEmpty()
	Up.ImpulseState = state.Up.OrEmpty()
	Down.ImpulseState = state.Down.OrEmpty()
	Jump.ImpulseState = state.Jump.OrEmpty()
	Action.ImpulseState = state.Action.OrEmpty()
	Exit.ImpulseState = state.Exit.OrEmpty()
	hoverPos = state.HoverPos
	clickPos = state.ClickPos
	easterEgg.justHit = state.EasterEggJustHit
	snesEasterEgg.justHit = state.EasterEggJustHit
	konamiCode.justHit = state.KonamiCodeJustHit
	snesKonamiCode.justHit = state.KonamiCodeJustHit
	kbdKonamiCode.justHit = state.KonamiCodeJustHit
	literalKbdKonamiCode.justHit = state.KonamiCodeJustHit
}

func SaveToDemo() *DemoState {
	return &DemoState{
		InputMap:          inputMap,
		Left:              Left.ImpulseState.UnlessEmpty(),
		Right:             Right.ImpulseState.UnlessEmpty(),
		Up:                Up.ImpulseState.UnlessEmpty(),
		Down:              Down.ImpulseState.UnlessEmpty(),
		Jump:              Jump.ImpulseState.UnlessEmpty(),
		Action:            Action.ImpulseState.UnlessEmpty(),
		Exit:              Exit.ImpulseState.UnlessEmpty(),
		HoverPos:          hoverPos,
		ClickPos:          clickPos,
		EasterEggJustHit:  EasterEggJustHit(),
		KonamiCodeJustHit: KonamiCodeJustHit(),
	}
}

func Draw(screen *ebiten.Image) {
	touchDraw(screen)
}

func DrawEditor(screen *ebiten.Image) {
	touchEditDraw(screen)
}

func ExitPressed() {
	Exit.externallyPressed = true
}

func SetActionButtonAvailable(avail bool) {
	actionButtonAvailable = avail
}

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
	"github.com/divVerent/aaaaxy/internal/sequence"
)

const (
	easterEggA      = 1
	easterEggB      = 2
	easterEggX      = 4
	easterEggY      = 8
	easterEggLeft   = 16
	easterEggRight  = 32
	easterEggUp     = 64
	easterEggDown   = 128
	easterEggJump   = 256
	easterEggAction = 512
)

type easterEggState struct {
	mask          int
	sequence      *sequence.Sequence
	sequenceFrame int  // Frames since last key.
	justHit       bool // If it was just activated this frame.
	prevState     int  // Previous key state.
}

const easterEggSequenceMaxFrames = 60 // At most one sec between key presses.

func (s *easterEggState) update(state int) {
	s.justHit = false

	presses := (state & ^s.prevState) & s.mask
	s.prevState = state

	// Count frames since last key.
	s.sequenceFrame++

	// Too long ago = reset.
	if s.sequenceFrame > easterEggSequenceMaxFrames {
		s.sequence.Reset()
		s.sequenceFrame = 0
		return
	}

	// Nothing pressed = no change.
	if presses == 0 {
		return
	}

	// Reset time since last key.
	s.sequenceFrame = 0

	// Add the byte to the sequence.
	s.sequence.Add(presses)

	// Check if it is hit.
	s.justHit = s.sequence.Match()
}

var (
	easterEgg = easterEggState{
		mask: easterEggA | easterEggB | easterEggX | easterEggY,
		sequence: sequence.New(
			easterEggA,
			easterEggA,
			easterEggA,
			easterEggA,
			easterEggX,
			easterEggY,
		)}
	snesEasterEgg = easterEggState{
		mask: easterEggA | easterEggB | easterEggX | easterEggY,
		sequence: sequence.New(
			easterEggB,
			easterEggB,
			easterEggB,
			easterEggB,
			easterEggY,
			easterEggX,
		)}
	konamiCode = easterEggState{
		mask: easterEggUp | easterEggDown | easterEggLeft | easterEggRight | easterEggJump | easterEggAction,
		sequence: sequence.New(
			easterEggUp,
			easterEggUp,
			easterEggDown,
			easterEggDown,
			easterEggLeft,
			easterEggRight,
			easterEggLeft,
			easterEggRight,
			easterEggAction,
			easterEggJump,
		)}
	snesKonamiCode = easterEggState{
		mask: easterEggUp | easterEggDown | easterEggLeft | easterEggRight | easterEggJump | easterEggAction,
		sequence: sequence.New(
			easterEggUp,
			easterEggUp,
			easterEggDown,
			easterEggDown,
			easterEggLeft,
			easterEggRight,
			easterEggLeft,
			easterEggRight,
			easterEggJump,
			easterEggAction,
		)}
)

func easterEggButtonState() int {
	s := 0
	if Left.Held {
		s |= easterEggLeft
	}
	if Right.Held {
		s |= easterEggRight
	}
	if Up.Held {
		s |= easterEggUp
	}
	if Down.Held {
		s |= easterEggDown
	}
	if Jump.Held {
		s |= easterEggJump
	}
	if Action.Held {
		s |= easterEggAction
	}
	return s
}

func easterEggUpdate() {
	gamepadState := gamepadEasterEggKeyState()
	state := keyboardEasterEggKeyState() | gamepadState | easterEggButtonState()
	easterEgg.update(state)
	snesEasterEgg.update(gamepadState) // Only allow reversing on gamepads as this is literal.
	konamiCode.update(state)
	snesKonamiCode.update(state) // Allow reversing the actions on keyboard too.
}

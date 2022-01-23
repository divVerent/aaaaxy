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

type easterEggKeyState int

const (
	easterEggA      easterEggKeyState = 1
	easterEggB      easterEggKeyState = 2
	easterEggX      easterEggKeyState = 4
	easterEggY      easterEggKeyState = 8
	easterEggLeft   easterEggKeyState = 16
	easterEggRight  easterEggKeyState = 32
	easterEggUp     easterEggKeyState = 64
	easterEggDown   easterEggKeyState = 128
	easterEggJump   easterEggKeyState = 256
	easterEggAction easterEggKeyState = 512
)

type easterEggState struct {
	mask          easterEggKeyState
	sequence      []easterEggKeyState
	sequencePos   int               // Last pos we _have_.
	sequenceFrame int               // Frames since last bump.
	justHit       bool              // If it was just activated this frame.
	prevState     easterEggKeyState // Previous key state.
}

const easterEggSequenceMaxFrames = 180 // 3 seconds should be enough to enter that.

func (s *easterEggState) update(state easterEggKeyState) {
	// TODO: unify state machine with internal/game/target/sequence?
	s.justHit = false

	presses := (state & ^s.prevState) & s.mask
	s.prevState = state

	// Count frames since start of sequence.
	if s.sequencePos == 0 {
		// Reset timer if at start.
		s.sequenceFrame = 0
	} else {
		s.sequenceFrame++
		// Too long ago = reset.
		if s.sequenceFrame > easterEggSequenceMaxFrames {
			s.sequencePos = 0
			return
		}
	}

	// Nothing pressed = no change.
	if presses == 0 {
		return
	}

	// Wrong key = reset.
	if presses != s.sequence[s.sequencePos] {
		s.sequencePos = 0
		return
	}

	// Advance.
	s.sequencePos++

	// End of sequence = good.
	if s.sequencePos == len(s.sequence) {
		s.justHit = true
		s.sequencePos = 0
	}
}

var (
	easterEgg = easterEggState{
		mask: easterEggA | easterEggB | easterEggX | easterEggY,
		sequence: []easterEggKeyState{
			easterEggA,
			easterEggA,
			easterEggA,
			easterEggA,
			easterEggX,
			easterEggY,
		}}
	snesEasterEgg = easterEggState{
		mask: easterEggA | easterEggB | easterEggX | easterEggY,
		sequence: []easterEggKeyState{
			easterEggB,
			easterEggB,
			easterEggB,
			easterEggB,
			easterEggY,
			easterEggX,
		}}
	konamiCode = easterEggState{
		mask: easterEggUp | easterEggDown | easterEggLeft | easterEggRight | easterEggJump | easterEggAction,
		sequence: []easterEggKeyState{
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
		}}
	snesKonamiCode = easterEggState{
		mask: easterEggUp | easterEggDown | easterEggLeft | easterEggRight | easterEggJump | easterEggAction,
		sequence: []easterEggKeyState{
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
		}}
)

func easterEggButtonState() easterEggKeyState {
	var s easterEggKeyState
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
	gamepadState := keyboardEasterEggKeyState() | gamepadEasterEggKeyState() | easterEggButtonState()
	state := keyboardEasterEggKeyState() | gamepadState | easterEggButtonState()
	easterEgg.update(state)
	snesEasterEgg.update(gamepadState) // Only allow reversing on gamepads as this is literal.
	konamiCode.update(state)
	snesKonamiCode.update(state) // Allow reversing the actions on keyboard too.
}

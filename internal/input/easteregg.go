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
	easterEggA easterEggKeyState = 1
	easterEggX easterEggKeyState = 2
	easterEggY easterEggKeyState = 4
)

var (
	easterEggSequence = []easterEggKeyState{
		easterEggA,
		easterEggA,
		easterEggA,
		easterEggA,
		easterEggX,
		easterEggY,
	}
	easterEggSequencePos                     = 0 // Last pos we _have_.
	easterEggSequenceFrame                   = 0 // Frames since last bump.
	easterEggJustHit                         = false
	easterEggPrevState     easterEggKeyState = 0
)

const easterEggSequenceMaxFrames = 180 // 3 seconds should be enough to enter that.

func easterEggUpdate() {
	easterEggJustHit = false
	state := keyboardEasterEggKeyState() | gamepadEasterEggKeyState()
	presses := state & ^easterEggPrevState
	easterEggPrevState = state

	// Count frames since start of sequence.
	if easterEggSequencePos == 0 {
		// Reset timer if at start.
		easterEggSequenceFrame = 0
	} else {
		easterEggSequenceFrame++
		// Too long ago = reset.
		if easterEggSequenceFrame > easterEggSequenceMaxFrames {
			easterEggSequencePos = 0
			return
		}
	}

	// Nothing pressed = no change.
	if presses == 0 {
		return
	}

	// Wrong key = reset.
	if presses != easterEggSequence[easterEggSequencePos] {
		easterEggSequencePos = 0
		return
	}

	// Advance.
	easterEggSequencePos++

	// End of sequence = good.
	if easterEggSequencePos == len(easterEggSequence) {
		easterEggJustHit = true
		easterEggSequencePos = 0
	}
}

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

package splash

import (
	"fmt"
	"time"

	"github.com/divVerent/aaaaxy/internal/flag"
)

var (
	loadingScreen                = flag.Bool("loading_screen", true, "show a detailed loading screen")
	debugLoadingScreenSkipFrames = flag.Int("debug_loading_screen_skip_frames", flag.SystemDefault(map[string]int{
		"android/*": 1, // Android seems to delay by one frame.
		"windows/*": 1, // Direct3D seems to delay by one frame.
		"*/*":       0,
	}), "number of frames to wait on the loading screen for each step; needed if the loading screen behaves erratically due to render delay")
)

type Status int

const (
	EndFrame Status = iota
	Continue
)

// State represents the current splash screen state.
type State struct {
	// knownFractions are the known progress bar fractions for each step as loaded.
	knownFractions map[string]float64

	// started is when the first init step started.
	started time.Time

	// startTimes are the actual start times of each step.
	startTimes map[string]time.Time

	// done is created when each step finishes.
	done map[string]struct{}

	// curStep is the currently executing step.
	curStep string

	// curStepName is the localized name of the currently executing step.
	curStepName string

	// curFraction is the current progress bar fraction.
	curFraction float64

	// skipFrames is the number of frames to skip before the next step.
	skipFrames int
}

// ProvideFractions loads a known fractions map.
// This should have been dumped by a previous run.
func (s *State) ProvideFractions(fractions map[string]float64) {
	s.knownFractions = fractions
}

// RunImmediately runs the given status-ish function as a single step.
// Useful for doing stuff w/o an actual loading screen.
func RunImmediately(errPrefix string, f func(s *State) (Status, error)) (Status, error) {
	// Simpler implementation that never updates the loading screen and does all init in one frame.
	for {
		status, err := f(nil)
		if err != nil {
			return EndFrame, fmt.Errorf("%v: %w", errPrefix, err)
		}
		if status == EndFrame {
			// f did not terminate yet - we need to call it again.
			continue
		}
		return status, nil
	}
}

// Enter enters a splash screen section.
// step must be an unique string identifying what is being loaded.
// f is allowed to call Enter too, but must return false, nil if its own Enter calls returned false.
// f must repeat all Enter calls it does, but will never be called again once it returned true.
func (s *State) Enter(step string, stepName string, errPrefix string, f func(s *State) (Status, error)) (Status, error) {
	if !*loadingScreen || s == nil {
		return RunImmediately(errPrefix, f)
	}

	if s.skipFrames > 0 {
		s.skipFrames--
		return EndFrame, nil
	}

	if s.startTimes == nil {
		s.startTimes = map[string]time.Time{}
	}
	if s.done == nil {
		s.done = map[string]struct{}{}
	}
	if _, have := s.done[step]; have {
		// Already done!
		return Continue, nil
	}
	if _, have := s.startTimes[step]; !have {
		s.startTimes[step] = time.Now()
		s.curStep = step
		s.curStepName = stepName
		frac := s.knownFractions[step]
		if frac > s.curFraction {
			s.curFraction = frac
		}
		// Must force a refresh so the new text actually can show.
		s.skipFrames = *debugLoadingScreenSkipFrames
		return EndFrame, nil
	}
	status, err := f(s)
	if err != nil {
		return EndFrame, fmt.Errorf("%v: %w", errPrefix, err)
	}
	if status == EndFrame {
		// f did not terminate yet - we need to call it again next frame.
		s.skipFrames = *debugLoadingScreenSkipFrames
		return EndFrame, nil
	}
	s.done[step] = struct{}{}
	return status, nil
}

// Single wraps a simple function into a splash screen step.
func Single(f func() error) func(s *State) (Status, error) {
	return func(s *State) (Status, error) {
		return Continue, f()
	}
}

// Current returns the current progress bar content.
func (s *State) Current() (string, float64) {
	return s.curStepName, s.curFraction
}

// ToFractions returns the init fraction map by step. This can be provided via ProvideFractions next time.
func (s *State) ToFractions() map[string]float64 {
	ended := time.Now()
	var started time.Time
	for _, v := range s.startTimes {
		if started.IsZero() || v.Before(started) {
			started = v
		}
	}
	fractions := make(map[string]float64, len(s.startTimes))
	for k, v := range s.startTimes {
		fractions[k] = float64(v.Sub(started)) / float64(ended.Sub(started))
	}
	return fractions
}

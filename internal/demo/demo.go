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

package demo

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
)

var (
	demoRecord = flag.String("demo_record", "", "local file path for demo to record to")
	demoPlay   = flag.String("demo_play", "", "local file path for demo to play back")
)

type frame struct {
	SaveGame *level.SaveGame `json:",omitempty"`
	Input    input.DemoState

	// The following data is not actually played back, but compared at playback time.
	// TODO: On first N regressions (i.e. regression frames where previous frame had no regression yet), dump a screenshot and link to it in stderr?
	SavedGames []uint64 `json:",omitempty"`
	PlayerPos  m.Pos
}

var (
	demoPlayerFile     *os.File
	demoPlayer         *json.Decoder
	demoPlayerFrame    frame
	demoPlayerFrameIdx int
	demoRecorderFrame  frame
	demoRecorderFile   *os.File
	demoRecorder       *json.Encoder
)

func Init() error {
	if *demoPlay != "" {
		var err error
		demoPlayerFile, err = os.Open(*demoPlay)
		if err != nil {
			return err
		}
		demoPlayer = json.NewDecoder(demoPlayerFile)
	}
	if *demoRecord != "" {
		if flag.Cheating() {
			return fmt.Errorf("cannot record a demo while cheating")
		}
		var err error
		demoRecorderFile, err = os.Create(*demoRecord)
		if err != nil {
			return err
		}
		demoRecorder = json.NewEncoder(demoRecorderFile)
		demoRecorder.SetIndent("", "")
	}
	return nil
}

func BeforeExit() {
	if demoPlayer != nil {
		if demoPlayer.More() {
			regression("game ended but demo would still go on")
		}
		err := demoPlayerFile.Close()
		if err != nil {
			log.Fatalf("failed to close played demo from %v: %v", *demoPlay, err)
		}
		regressionBeforeExit()
	}
	if demoRecorder != nil {
		recordFrame()
		err := demoRecorderFile.Close()
		if err != nil {
			log.Fatalf("failed to save demo to %v: %v", *demoRecord, err)
		}
	}
}

func Playing() bool {
	return demoPlayer != nil
}

func Update() bool {
	wantQuit := false
	if demoPlayer != nil {
		wantQuit = playFrame()
	}
	if demoRecorder != nil {
		recordFrame()
	}
	return wantQuit
}

func PostUpdate(playerPos m.Pos) {
	if demoPlayer != nil {
		postPlayFrame(playerPos)
	}
	if demoRecorder != nil {
		postRecordFrame(playerPos)
	}
}

func PostDraw(screen *ebiten.Image) {
	if demoPlayer != nil {
		regressionPostDrawFrame(screen)
	}
}

func playFrame() bool {
	if !demoPlayer.More() {
		regression("demo ended but game didn't quit")
		return true
	}
	s := demoPlayerFrame.SaveGame
	demoPlayerFrame = frame{}
	err := demoPlayer.Decode(&demoPlayerFrame)
	if err != nil {
		regression("could not decode demo frame: %v", err)
	}
	// Restore save game, so loading always succeeds even if we've regressed.
	if demoPlayerFrame.SaveGame == nil {
		demoPlayerFrame.SaveGame = s
	}
	input.LoadFromDemo(&demoPlayerFrame.Input)
	return false
}

func postPlayFrame(playerPos m.Pos) {
	if len(demoPlayerFrame.SavedGames) != 0 {
		regression("saved game: got no saves, want %v", demoPlayerFrame.SavedGames)
	}
	if playerPos != demoPlayerFrame.PlayerPos {
		regression("player pos: got %v, want %v", playerPos, demoPlayerFrame.PlayerPos)
	}
	regressionPostPlayFrame()
	demoPlayerFrameIdx++
}

func recordFrame() {
	demoRecorderFrame = frame{
		Input: *input.SaveToDemo(),
	}
}

func postRecordFrame(playerPos m.Pos) {
	demoRecorderFrame.PlayerPos = playerPos
	err := demoRecorder.Encode(&demoRecorderFrame)
	if err != nil {
		log.Fatalf("could not encode demo frame: %v", err)
	}
}

func InterceptSaveGame(save *level.SaveGame) bool {
	// While playing back, we only save to memory to allow later recalling.
	if demoPlayer != nil {
		// Ensure next load event will be handled right according to this save game.
		// This shoulnd't be needed - InterceptPostLoadGame should have ensured the save game is always updated on every load event.
		// Still there to have better chance of being in sync during playback with regression.
		demoPlayerFrame.SaveGame = save
		if len(demoPlayerFrame.SavedGames) == 0 {
			regression("saved game: got hash %v, want no saves", save.Hash)
		} else {
			if save.Hash != demoPlayerFrame.SavedGames[0] {
				regression("saved game: got hash %v, want %v", save.Hash, demoPlayerFrame.SavedGames[0])
			}
			demoPlayerFrame.SavedGames = demoPlayerFrame.SavedGames[1:]
		}
		return true
	}
	if demoRecorder != nil {
		demoRecorderFrame.SavedGames = append(demoRecorderFrame.SavedGames, save.Hash)
	}
	return false
}

func InterceptPreLoadGame() (*level.SaveGame, bool) {
	// While playing back, we always return the last save game from the demo.
	if demoPlayer != nil {
		return demoPlayerFrame.SaveGame, true
	}
	return nil, false
}

func InterceptPostLoadGame(save *level.SaveGame) {
	// While recording, store the current save game.
	if demoRecorder != nil {
		demoRecorderFrame.SaveGame = save
	}
}

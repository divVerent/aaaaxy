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
	demoPlayerFile       *os.File
	demoPlayer           *json.Decoder
	demoPlayerSave       *level.SaveGame
	demoPlayerSavedGames []uint64
	demoPlayerPos        m.Pos
	demoRecorderFrame    frame
	demoRecorderFile     *os.File
	demoRecorder         *json.Encoder
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
			log.Errorf("REGRESSION: game ended but demo would still go on")
		}
		err := demoPlayerFile.Close()
		if err != nil {
			log.Fatalf("failed to close played demo from %v: %v", *demoPlay, err)
		}
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

func playFrame() bool {
	if !demoPlayer.More() {
		log.Errorf("REGRESSION: demo ended but game didn't quit")
		return true
	}
	var f frame
	err := demoPlayer.Decode(&f)
	if err != nil {
		log.Fatalf("could not decode demo frame: %v", err)
	}
	if f.SaveGame != nil {
		demoPlayerSave = f.SaveGame
	}
	input.LoadFromDemo(&f.Input)
	demoPlayerSavedGames = f.SavedGames
	demoPlayerPos = f.PlayerPos
	return false
}

func postPlayFrame(playerPos m.Pos) {
	if len(demoPlayerSavedGames) != 0 {
		log.Errorf("REGRESSION: saved game: got no saves, want %v", demoPlayerSavedGames)
	}
	if playerPos != demoPlayerPos {
		log.Errorf("REGRESSION: player pos: got %v, want %v", playerPos, demoPlayerPos)
	}
	demoPlayerSavedGames = nil
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
		demoPlayerSave = save
		if len(demoPlayerSavedGames) == 0 {
			log.Errorf("REGRESSION: saved game: got hash %v, want no saves", save.Hash)
		} else {
			if save.Hash != demoPlayerSavedGames[0] {
				log.Errorf("REGRESSION: saved game: got hash %v, want %v", save.Hash, demoPlayerSavedGames[0])
			}
			demoPlayerSavedGames = demoPlayerSavedGames[1:]
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
		return demoPlayerSave, true
	}
	return nil, false
}

func InterceptPostLoadGame(save *level.SaveGame) {
	// While recording, store the current save game.
	if demoRecorder != nil {
		demoRecorderFrame.SaveGame = save
	}
}

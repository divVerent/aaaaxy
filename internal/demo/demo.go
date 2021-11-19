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

	"github.com/google/go-cmp/cmp"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
)

var (
	demoRecord   = flag.String("demo_record", "", "local file path for demo to record to")
	demoPlay     = flag.String("demo_play", "", "local file path for demo to play back")
	demoTimedemo = flag.Bool("demo_timedemo", false, "run demos as fast as possible, only limited by rendering; normally you'd want to pass -vsync=false too when using this")
)

type frame struct {
	SaveGame *level.SaveGame  `json:",omitempty"`
	Input    *input.DemoState `json:",omitempty"`

	// The following data is not actually played back, but compared at playback time.
	SaveGames     []uint64        `json:",omitempty"`
	FinalSaveGame *level.SaveGame `json:",omitempty"`
	PlayerPos     *m.Pos          `json:",omitempty"`
}

var (
	demoPlayerFile            *os.File
	demoPlayer                *json.Decoder
	demoPlayerFrame           frame
	demoPlayerFrameIdx        int
	demoRecorderFrame         frame
	demoRecorderFile          *os.File
	demoRecorderFinalSaveGame *level.SaveGame
	demoRecorder              *json.Encoder
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
	if demoRecorder != nil {
		demoRecorderFrame = frame{
			FinalSaveGame: demoRecorderFinalSaveGame,
		}
		err := demoRecorder.Encode(&demoRecorderFrame)
		if err != nil {
			log.Fatalf("could not encode final demo frame: %v", err)
		}
		err = demoRecorderFile.Close()
		if err != nil {
			log.Fatalf("failed to save demo to %v: %v", *demoRecord, err)
		}
	}
	if demoPlayer != nil {
		if playReadFrame() {
			regression(highPrio, "game ended but demo would still go on")
		}
		err := demoPlayerFile.Close()
		if err != nil {
			log.Fatalf("failed to close played demo from %v: %v", *demoPlay, err)
		}
		regressionBeforeExit()
	}
}

func Playing() bool {
	return demoPlayer != nil
}

func Timedemo() bool {
	return Playing() && *demoTimedemo
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

func playReadFrame() bool {
	s := demoPlayerFrame.SaveGame
	for demoPlayer.More() {
		demoPlayerFrame = frame{}
		err := demoPlayer.Decode(&demoPlayerFrame)
		if err != nil {
			log.Fatalf("could not decode demo frame: %v", err)
		}
		if demoPlayerFrame.FinalSaveGame == nil {
			// Restore save game, so loading always succeeds even if we've regressed.
			if demoPlayerFrame.SaveGame == nil {
				demoPlayerFrame.SaveGame = s
			}
			return true
		}
		diff := cmp.Diff(demoPlayerFrame.FinalSaveGame.State, s.State)
		if diff != "" {
			regression(highPrio, "difference in final save state (-want +got):\n%v", diff)
		}
	}
	return false
}

func playFrame() bool {
	if !playReadFrame() {
		regression(highPrio, "demo ended but game didn't quit")
		return true
	}
	input.LoadFromDemo(demoPlayerFrame.Input)
	return false
}

func postPlayFrame(playerPos m.Pos) {
	if len(demoPlayerFrame.SaveGames) != 0 {
		regression(mediumPrio, "saved game: got no saves, want %v", demoPlayerFrame.SaveGames)
	}
	if demoPlayerFrame.PlayerPos != nil && playerPos != *demoPlayerFrame.PlayerPos {
		d := playerPos.Delta(*demoPlayerFrame.PlayerPos).Norm1()
		dlog := 0
		dpow := 1
		for d >= dpow {
			dlog++
			dpow *= 2
		}
		regression(lowPrio.WithParam(dlog), "player pos: got %v, want %v", playerPos, *demoPlayerFrame.PlayerPos)
	}
	regressionPostPlayFrame()
	demoPlayerFrameIdx++
}

func recordFrame() {
	demoRecorderFrame = frame{
		Input: input.SaveToDemo(),
	}
}

func postRecordFrame(playerPos m.Pos) {
	demoRecorderFrame.PlayerPos = &playerPos
	err := demoRecorder.Encode(&demoRecorderFrame)
	if err != nil {
		log.Fatalf("could not encode demo frame: %v", err)
	}
}

func InterceptSaveGame(save *level.SaveGame) bool {
	// Always record everything.
	if demoRecorder != nil {
		demoRecorderFrame.SaveGames = append(demoRecorderFrame.SaveGames, save.StateHash)
		demoRecorderFinalSaveGame = save
	}
	// While playing back, we only save to memory to allow later recalling.
	if demoPlayer != nil {
		// Ensure next load event will be handled right according to this save game.
		// This shoulnd't be needed - InterceptPostLoadGame should have ensured the save game is always updated on every load event.
		// Still there to have better chance of being in sync during playback with regression.
		demoPlayerFrame.SaveGame = save
		if len(demoPlayerFrame.SaveGames) == 0 {
			regression(mediumPrio, "saved game: got hash %v, want no saves", save.StateHash)
		} else {
			if save.StateHash != demoPlayerFrame.SaveGames[0] {
				regression(mediumPrio, "saved game: got hash %v, want %v", save.StateHash, demoPlayerFrame.SaveGames[0])
			}
			demoPlayerFrame.SaveGames = demoPlayerFrame.SaveGames[1:]
		}
		return true
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

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
)

var (
	demoRecord = flag.String("demo_record", "", "local file path for demo to record to")
	demoPlay   = flag.String("demo_play", "", "local file path for demo to play back")
)

type frame struct {
	SaveGame *level.SaveGame `json:",omitempty"`
	Input    input.DemoState
}

var (
	demoPlayerFile      *os.File
	demoPlayer          *json.Decoder
	demoPlayerSave      *level.SaveGame
	demoRecorderStarted bool
	demoRecorderFrame   frame
	demoRecorderFile    *os.File
	demoRecorder        *json.Encoder
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

func playFrame() bool {
	if !demoPlayer.More() {
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
	return false
}

func recordFrame() {
	if demoRecorderStarted {
		// We are recording demo frames one frame late
		// so that save games loaded during the frame
		// are provided at frame start.
		err := demoRecorder.Encode(&demoRecorderFrame)
		if err != nil {
			log.Fatalf("could not encode demo frame: %v", err)
		}
	}
	demoRecorderFrame = frame{
		Input: *input.SaveToDemo(),
	}
	demoRecorderStarted = true
}

func InterceptSaveGame(save *level.SaveGame) bool {
	// While playing back, we only save to memory to allow later recalling.
	if demoPlayer != nil {
		demoPlayerSave = save
		return true
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

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

package player_state

import (
	"fmt"
	"log"

	"github.com/divVerent/aaaaaa/internal/level"
)

type PlayerState struct {
	Level *level.Level
}

func (s *PlayerState) LastCheckpoint() string {
	return s.Level.Player.PersistentState["last_checkpoint"]
}

func (s *PlayerState) CheckpointsWalked(from, to string) bool {
	return s.Level.Player.PersistentState["checkpoints_walked."+from+"."+to] != ""
}

type SeenState int

const (
	NotSeen SeenState = iota
	SeenNormal
	SeenFlipped
)

func (s *PlayerState) CheckpointSeen(name string) SeenState {
	state := s.Level.Player.PersistentState["checkpoint_seen."+name]
	switch state {
	case "":
		return NotSeen
	case "FlipX":
		return SeenFlipped
	case "Identity":
		return SeenNormal
	default:
		log.Panicf("invalid checkpoint_seen state: %v", state)
	}
	// Unreachable.
	return 0
}

func (s *PlayerState) RecordCheckpoint(name string, flipped bool) bool {
	flip := "Identity"
	if flipped {
		flip = "FlipX"
	}
	updated := false
	if s.Level.Player.PersistentState["checkpoint_seen."+name] != flip {
		s.Level.Player.PersistentState["checkpoint_seen."+name] = flip
		updated = true
	}
	if s.Level.Checkpoints[name].Properties["dead_end"] != "true" {
		if s.Level.Player.PersistentState["last_checkpoint"] != name {
			s.Level.Player.PersistentState["last_checkpoint"] = name
			updated = true
		}
	}
	return updated
}

func (s *PlayerState) RecordCheckpointEdge(name string, flipped bool) bool {
	from := s.Level.Player.PersistentState["last_checkpoint"]
	updated := s.RecordCheckpoint(name, flipped)
	if from != name {
		if s.Level.Player.PersistentState["checkpoints_walked."+from+"."+name] != "true" {
			s.Level.Player.PersistentState["checkpoints_walked."+from+"."+name] = "true"
			updated = true
		}
	}
	return updated
}

func (s *PlayerState) TnihSignsSeen(name string) (seen, total int) {
	seen, total = 0, 0
	for _, sign := range s.Level.TnihSignsByCheckpoint[name] {
		total++
		if sign.PersistentState["seen"] == "true" {
			seen++
		}
	}
	return
}

func (s *PlayerState) Frames() int {
	framesStr := s.Level.Player.PersistentState["frames"]
	var frames int
	if framesStr != "" {
		_, err := fmt.Sscanf(framesStr, "%d", &frames)
		if err != nil {
			log.Panicf("could not parse frames counter: %v", err)
		}
	}
	return frames
}

func (s *PlayerState) AddFrame() {
	s.Level.Player.PersistentState["frames"] = fmt.Sprint(s.Frames() + 1)
}

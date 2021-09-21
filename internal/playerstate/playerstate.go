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

package playerstate

import (
	"fmt"
	"strings"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	cheatFullMapNormal   = flag.Bool("cheat_full_map_normal", false, "Show the full map.")
	cheatFullMapFlipped  = flag.Bool("cheat_full_map_flipped", false, "Show the full map.")
	cheatPlayerAbilities = flag.StringMap("cheat_player_abilities", map[string]string{}, "Override player abilities")
)

type PlayerState struct {
	Level *level.Level
}

func (s *PlayerState) HasAbility(name string) bool {
	switch (*cheatPlayerAbilities)[name] {
	case "true":
		return true
	case "false":
		return false
	}
	key := "can_" + name
	return s.Level.Player.PersistentState[key] == "true"
}

func (s *PlayerState) GiveAbility(name string) bool {
	if (*cheatPlayerAbilities)[name] != "" {
		return false
	}
	key := "can_" + name
	if s.Level.Player.PersistentState[key] == "true" {
		return false
	}
	s.Level.Player.PersistentState[key] = "true"
	return true
}

func (s *PlayerState) LastCheckpoint() string {
	return s.Level.Player.PersistentState["last_checkpoint"]
}

func (s *PlayerState) CheckpointsWalked(from, to string) bool {
	if *cheatFullMapNormal || *cheatFullMapFlipped {
		return true
	}
	// CheckpointsWalked is a symmetric relation.
	return s.Level.Player.PersistentState["checkpoints_walked."+from+"."+to] != "" ||
		s.Level.Player.PersistentState["checkpoints_walked."+to+"."+from] != ""
}

type SeenState int

const (
	NotSeen SeenState = iota
	SeenNormal
	SeenFlipped
)

func (s *PlayerState) CheckpointSeen(name string) SeenState {
	if *cheatFullMapNormal {
		return SeenNormal
	}
	if *cheatFullMapFlipped {
		return SeenFlipped
	}
	state := s.Level.Player.PersistentState["checkpoint_seen."+name]
	switch state {
	case "":
		return NotSeen
	case "FlipX":
		return SeenFlipped
	case "Identity":
		return SeenNormal
	default:
		log.Fatalf("invalid checkpoint_seen state: %v", state)
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
			log.Fatalf("could not parse frames counter: %v", err)
		}
	}
	return frames
}

func (s *PlayerState) AddFrame() {
	s.Level.Player.PersistentState["frames"] = fmt.Sprint(s.Frames() + 1)
}

func (s *PlayerState) Escapes() int {
	escapesStr := s.Level.Player.PersistentState["escapes"]
	var escapes int
	if escapesStr != "" {
		_, err := fmt.Sscanf(escapesStr, "%d", &escapes)
		if err != nil {
			log.Fatalf("could not parse escapes counter: %v", err)
		}
	}
	return escapes
}

func (s *PlayerState) AddEscape() {
	s.Level.Player.PersistentState["escapes"] = fmt.Sprint(s.Escapes() + 1)
}

func (s *PlayerState) Won() bool {
	return s.Level.Player.PersistentState["won"] == "true"
}

func (s *PlayerState) SetWon() {
	s.Level.Player.PersistentState["won"] = "true"
}

type SpeedrunCategories int

const (
	AnyPercentSpeedrun     SpeedrunCategories = 0x01
	HundredPercentSpeedrun SpeedrunCategories = 0x02
	AllSignsSpeedrun       SpeedrunCategories = 0x04
	AllPathsSpeedrun       SpeedrunCategories = 0x08
	AllSecretsSpeedrun     SpeedrunCategories = 0x10
	AllFlippedSpeedrun     SpeedrunCategories = 0x20
	NoEscapeSpeedrun       SpeedrunCategories = 0x40
	allCategoriesSpeedrun  SpeedrunCategories = 0x7F
)

func (c SpeedrunCategories) Strings() (categories string, tryNext string) {
	tryNext = ""
	cats := []string{}
	addCategory := func(cat string, have bool) {
		if have {
			cats = append(cats, cat)
		} else {
			if tryNext == "" {
				tryNext = cat
			}
		}
	}
	if flag.Cheating() {
		addCategory("Cheat%", true)
		addCategory("Without Cheating Of Course", false)
	}
	if c&HundredPercentSpeedrun == 0 {
		addCategory("Any%", c&AnyPercentSpeedrun != 0)
	}
	addCategory("100%", c&HundredPercentSpeedrun != 0)
	addCategory("All Notes", c&AllSignsSpeedrun != 0)
	addCategory("All Paths", c&AllPathsSpeedrun != 0)
	addCategory("All Secrets", c&AllSecretsSpeedrun != 0)
	addCategory("All Flipped", c&AllFlippedSpeedrun != 0)
	noEscape := "No Escape"
	if input.UsingGamepad() {
		noEscape = "No Start"
	}
	addCategory(noEscape, c&NoEscapeSpeedrun != 0)
	l := len(cats)
	switch l {
	case 0:
		categories = "None"
	case 1:
		categories = cats[0]
	default:
		categories = strings.Join(cats[0:l-1], ", ") + " and " + cats[l-1]
	}
	return
}

func (s *PlayerState) SpeedrunCategories() SpeedrunCategories {
	cat := allCategoriesSpeedrun
	if !s.Won() {
		cat &^= AnyPercentSpeedrun
	}
	cat |= HundredPercentSpeedrun | AllFlippedSpeedrun | AllSignsSpeedrun
	for cp, cpSp := range s.Level.Checkpoints {
		if cp == "" {
			// Start is not a real CP.
			continue
		}
		if cpSp.Properties["secret"] == "true" {
			// Secrets are not needed for 100%, all paths or all signs run.
			// However they have their own run category here.
			for _, sign := range s.Level.TnihSignsByCheckpoint[cp] {
				if sign.PersistentState["seen"] != "true" {
					cat &^= AllSecretsSpeedrun
				}
			}
			continue
		}
		switch s.CheckpointSeen(cp) {
		case NotSeen:
			cat &^= HundredPercentSpeedrun
		case SeenNormal:
			// Note: this means AllFlipped is possible without 100%. WAI.
			cat &^= AllFlippedSpeedrun
		}
		for _, next := range s.Level.CheckpointLocations.Locs[cp].NextByDir {
			// Note: this skips dead ends; very much intended.
			if !next.Forward || next.Optional {
				continue
			}
			if !s.CheckpointsWalked(cp, next.Other) {
				cat &^= AllPathsSpeedrun
			}
		}
		// Dead ends not needed for all signs run.
		for _, sign := range s.Level.TnihSignsByCheckpoint[cp] {
			if sign.PersistentState["seen"] != "true" {
				cat &^= AllSignsSpeedrun
			}
		}
	}
	if s.Escapes() != 0 {
		// Note: this is impossible when also having AllSecrets,
		// as secrets typically cannot be left.
		cat &^= NoEscapeSpeedrun
	}
	return cat
}

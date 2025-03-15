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
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/propmap"
)

var (
	cheatFullMapNormal   = flag.Bool("cheat_full_map_normal", false, "show the full map")
	cheatFullMapFlipped  = flag.Bool("cheat_full_map_flipped", false, "show the full map")
	cheatPlayerAbilities = flag.StringMap[bool]("cheat_player_abilities", map[string]bool{}, "override player abilities")
	showSwitchLevel      = flag.Bool("show_switch_level", false, "show the level selector menu")
)

type PlayerState struct {
	Level *level.Level
}

// Init must be called when Level got externally changed, e.g. by loading world or a save state.
func (s *PlayerState) Init() {
	if s.Level.SaveGameVersion != 1 {
		log.Fatalf("please FIXME! On the next SaveGameVersion, please remove the escapes to teleport translation, the JSON compat hack in internal/math/pos.go, and remove this check too")
	}
	// If he savegame has no teleports info, use the escapes counter.
	// Also ensure all new savegames have the teleports counter to not double apply this.
	teleports := propmap.ValueOrP(s.Level.Player.PersistentState, "teleports", -1, nil)
	if teleports < 0 {
		propmap.Set(s.Level.Player.PersistentState, "teleports", s.Escapes())
	}
}

func (s *PlayerState) HasAbility(name string) bool {
	if name == "switch_level" {
		return *showSwitchLevel
	}
	have, found := (*cheatPlayerAbilities)[name]
	if found {
		return have
	}
	have, found = (*cheatPlayerAbilities)["all"]
	if found && s.Level.Abilities[name] {
		return have
	}
	key := "can_" + name
	return propmap.ValueOrP(s.Level.Player.PersistentState, key, false, nil)
}

func (s *PlayerState) GiveAbility(name string) bool {
	if name == "switch_level" {
		if *showSwitchLevel {
			return false
		}
		*showSwitchLevel = true
		return true
	}
	_, found := (*cheatPlayerAbilities)[name]
	if found {
		return false
	}
	key := "can_" + name
	if propmap.ValueOrP(s.Level.Player.PersistentState, key, false, nil) {
		return false
	}
	propmap.Set(s.Level.Player.PersistentState, key, true)
	return true
}

func (s *PlayerState) LastCheckpoint() string {
	return propmap.StringOr(s.Level.Player.PersistentState, "last_checkpoint", "")
}

func (s *PlayerState) CheckpointsWalked(from, to string) bool {
	if *cheatFullMapNormal || *cheatFullMapFlipped {
		return true
	}
	// CheckpointsWalked is a symmetric relation.
	return propmap.StringOr(s.Level.Player.PersistentState, "checkpoints_walked."+from+"."+to, "") != "" ||
		propmap.StringOr(s.Level.Player.PersistentState, "checkpoints_walked."+to+"."+from, "") != ""
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
	state := propmap.StringOr(s.Level.Player.PersistentState, "checkpoint_seen."+name, "")
	switch state {
	case "":
		return NotSeen
	case "FlipX":
		return SeenFlipped
	case "Identity":
		return SeenNormal
	default:
		log.TraceErrorf("invalid checkpoint_seen state: %v", state)
		return NotSeen
	}
}

func (s *PlayerState) RecordCheckpoint(name string, flipped bool) bool {
	flip := "Identity"
	if flipped {
		flip = "FlipX"
	}
	updated := false
	if propmap.StringOr(s.Level.Player.PersistentState, "checkpoint_seen."+name, "") != flip {
		propmap.Set(s.Level.Player.PersistentState, "checkpoint_seen."+name, flip)
		updated = true
	}
	if !propmap.ValueOrP(s.Level.Checkpoints[name].Properties, "dead_end", false, nil) {
		if propmap.StringOr(s.Level.Player.PersistentState, "last_checkpoint", "") != name {
			propmap.Set(s.Level.Player.PersistentState, "last_checkpoint", name)
			updated = true
		}
	}
	return updated
}

func (s *PlayerState) RecordCheckpointEdge(name string, flipped bool) bool {
	from := propmap.StringOr(s.Level.Player.PersistentState, "last_checkpoint", "")
	updated := s.RecordCheckpoint(name, flipped)
	if from != name {
		if !propmap.ValueOrP(s.Level.Player.PersistentState, "checkpoints_walked."+from+"."+name, false, nil) {
			propmap.Set(s.Level.Player.PersistentState, "checkpoints_walked."+from+"."+name, true)
			updated = true
		}
	}
	return updated
}

func (s *PlayerState) TnihSignsSeen(name string) (seen, total int) {
	seen, total = 0, 0
	for _, sign := range s.Level.TnihSignsByCheckpoint[name] {
		total++
		if propmap.ValueOrP(sign.PersistentState, "seen", false, nil) {
			seen++
		}
	}
	return
}

func (s *PlayerState) Frames() int {
	frames, err := propmap.ValueOr(s.Level.Player.PersistentState, "frames", 0)
	if err != nil {
		log.Errorf("could not parse frames counter: %v", err)
		return 60 * 86400 // Takes at least one day.
	}
	return frames
}

func (s *PlayerState) AddFrame() {
	propmap.Set(s.Level.Player.PersistentState, "frames", s.Frames()+1)
}

func (s *PlayerState) Escapes() int {
	escapes, err := propmap.ValueOr(s.Level.Player.PersistentState, "escapes", 0)
	if err != nil {
		log.Errorf("could not parse escapes counter: %v", err)
		return 60 * 86400 // Takes at least one day.
	}
	return escapes
}

func (s *PlayerState) AddEscape() {
	propmap.Set(s.Level.Player.PersistentState, "escapes", s.Escapes()+1)
}

func (s *PlayerState) Teleports() int {
	teleports, err := propmap.ValueOr(s.Level.Player.PersistentState, "teleports", 0)
	if err != nil {
		log.Errorf("could not parse teleports counter: %v", err)
		return 60 * 86400 // Takes at least one day.
	}
	return teleports
}

func (s *PlayerState) AddTeleport() {
	propmap.Set(s.Level.Player.PersistentState, "teleports", s.Teleports()+1)
}

func (s *PlayerState) SetLives(n int) {
	propmap.Set(s.Level.Player.PersistentState, "lives", n)
}

func (s *PlayerState) Won() bool {
	return propmap.ValueOrP(s.Level.Player.PersistentState, "won", false, nil)
}

func (s *PlayerState) SetWon() {
	propmap.Set(s.Level.Player.PersistentState, "won", true)
}

type SpeedrunCategories int

const (
	// Real speedrun categories.
	AnyPercentSpeedrun     SpeedrunCategories = 0x01
	AllCheckpointsSpeedrun SpeedrunCategories = 0x02
	AllSignsSpeedrun       SpeedrunCategories = 0x04
	AllPathsSpeedrun       SpeedrunCategories = 0x08
	AllSecretsSpeedrun     SpeedrunCategories = 0x10
	AllFlippedSpeedrun     SpeedrunCategories = 0x20
	NoEscapeSpeedrun       SpeedrunCategories = 0x40
	NoTeleportsSpeedrun    SpeedrunCategories = 0x80
	NoPushSpeedrun         SpeedrunCategories = 0x100
	// Remapping (reason: one can have all CPs but not Any%, i.e. won the game yet):
	// AnyPercent AllCheckpoints => Result
	// false      false          => 0
	// true       false          => AnyPercent
	// false      true           => AllCheckpoints
	// true       true           => HundredPercent
	hundredPercentSpeedrun SpeedrunCategories = 0x1000
	// The following ones only are used internally when naming.
	withoutCheatsSpeedrun SpeedrunCategories = 0x2000
	cheatingSpeedrun      SpeedrunCategories = 0x4000
	impossibleSpeedrun    SpeedrunCategories = 0x8000
	allCategoriesSpeedrun SpeedrunCategories = 0xFF
)

func (c SpeedrunCategories) Name() string {
	switch c {
	case AnyPercentSpeedrun:
		return locale.GI.Get("Any%")
	case AllCheckpointsSpeedrun:
		return locale.G.Get("All Checkpoints")
	case AllSignsSpeedrun:
		return locale.G.Get("All Notes")
	case AllPathsSpeedrun:
		return locale.G.Get("All Paths")
	case AllSecretsSpeedrun:
		return locale.G.Get("All Secrets")
	case AllFlippedSpeedrun:
		return locale.G.Get("All Flipped")
	case NoTeleportsSpeedrun:
		return locale.G.Get("No Teleports")
	case NoEscapeSpeedrun:
		switch input.ExitButton() {
		default: // case input.Escape:
			return locale.G.Get("No Escape")
		case input.Backspace:
			return locale.G.Get("No Backspace")
		case input.Start:
			return locale.G.Get("No Start")
		case input.Back:
			return locale.G.Get("No Back")
		}
	case NoPushSpeedrun:
		return locale.G.Get("No Coil")
	case hundredPercentSpeedrun:
		return locale.GI.Get("100%")
	case withoutCheatsSpeedrun:
		return locale.G.Get("Without Cheating Of Course")
	case cheatingSpeedrun:
		return locale.GI.Get("Cheat%")
	case impossibleSpeedrun:
		return locale.G.Get("Impossible")
	default:
		return locale.G.Get("???")
	}
}

func (c SpeedrunCategories) ShortName() string {
	switch c {
	case AnyPercentSpeedrun:
		return "%"
	case AllCheckpointsSpeedrun:
		return "C"
	case AllSignsSpeedrun:
		return "N"
	case AllPathsSpeedrun:
		return "P"
	case AllSecretsSpeedrun:
		return "S"
	case AllFlippedSpeedrun:
		return "F"
	case hundredPercentSpeedrun:
		return "&"
	case NoTeleportsSpeedrun:
		return "T"
	case NoEscapeSpeedrun:
		return "E"
	case NoPushSpeedrun:
		return "U"
	case withoutCheatsSpeedrun:
		return "" // Never actually appears other than in tryNext.
	case cheatingSpeedrun:
		return "c"
	case impossibleSpeedrun:
		return "!"
	default:
		return "?"
	}
}

func (c SpeedrunCategories) describeCommon() (categories []SpeedrunCategories, tryNext SpeedrunCategories) {
	addCategory := func(cat, what SpeedrunCategories) {
		if c.ContainAll(what) {
			categories = append(categories, cat)
		} else {
			if tryNext == 0 {
				tryNext = cat
			}
		}
	}
	if is, _ := flag.Cheating(); is {
		addCategory(cheatingSpeedrun, 0)
		addCategory(withoutCheatsSpeedrun, impossibleSpeedrun)
	} else if c.ContainAll(AllCheckpointsSpeedrun) {
		addCategory(hundredPercentSpeedrun, AllCheckpointsSpeedrun /* always true */)
	} else {
		addCategory(AnyPercentSpeedrun, AnyPercentSpeedrun)
		addCategory(AllCheckpointsSpeedrun, AllCheckpointsSpeedrun /* always false */)
	}
	addCategory(AllSignsSpeedrun, AllSignsSpeedrun)
	addCategory(AllPathsSpeedrun, AllPathsSpeedrun)
	addCategory(AllSecretsSpeedrun, AllSecretsSpeedrun)
	addCategory(AllFlippedSpeedrun, AllFlippedSpeedrun)
	addCategory(NoTeleportsSpeedrun, NoTeleportsSpeedrun)
	addCategory(NoEscapeSpeedrun, NoEscapeSpeedrun)
	addCategory(NoPushSpeedrun, NoPushSpeedrun)
	return categories, tryNext
}

func (c SpeedrunCategories) Describe() (categories string, tryNext string) {
	categoryIds, tryNextId := c.describeCommon()
	categoryNames := make([]string, 0, len(categoryIds))
	for _, catId := range categoryIds {
		name := catId.Name()
		if name == "" {
			continue
		}
		categoryNames = append(categoryNames, catId.Name())
	}
	l := len(categoryNames)
	switch l {
	case 0:
		categories = locale.G.Get("None")
	case 1:
		categories = categoryNames[0]
	default:
		categories = strings.Join(categoryNames[0:l-1], locale.G.Get(", ")) + locale.G.Get(" and ") + categoryNames[l-1]
	}
	return categories, tryNextId.Name()
}

func (c SpeedrunCategories) DescribeShort() string {
	categoryIds, _ := c.describeCommon()
	cats := ""
	for _, catId := range categoryIds {
		cats += catId.ShortName()
	}
	return cats
}

func (c SpeedrunCategories) ContainAll(cats SpeedrunCategories) bool {
	return (c & cats) == cats
}

func fakeScore(n, k int) string {
	if k == 0 {
		return ""
	}

	// fakeScore(k) = a k^2 + b k + c
	// Constraints:
	// fakeScore(0) = 0
	// fakeScore(n-1) = 0.997
	// fakeScore(n) = 0.998
	// Formulas below from solving using maxima.
	n64, k64 := int64(n), int64(k)
	denom := 1000 * n64 * (n64 - 1)
	anum := (n64 - 998)
	bnum := -(n64*n64 - 1996*n64 + 998)
	cnum := int64(0)
	num := anum*k64*k64 + bnum*k64 + cnum
	dig3 := (num*1000 + denom/2) / denom
	return fmt.Sprintf(".%03d", dig3)
}

func (s *PlayerState) Score() string {
	score := 0
	for cp := range s.Level.Checkpoints {
		if cp == "" {
			// Start is not a real CP.
			continue
		}
		// Score is just number of TnihSigns seen.
		for _, sign := range s.Level.TnihSignsByCheckpoint[cp] {
			if propmap.ValueOrP(sign.PersistentState, "seen", false, nil) {
				score++
			}
		}
	}
	qCount, qHit := len(s.Level.QuestionBlocks), 0
	for _, q := range s.Level.QuestionBlocks {
		if propmap.ValueOrP(q.PersistentState, "used", false, nil) {
			qHit++
		}
	}
	fake := fakeScore(qCount, qHit)
	return fmt.Sprintf("%d%s", score, fake)
}

func (s *PlayerState) SpeedrunCategories() SpeedrunCategories {
	cat := allCategoriesSpeedrun
	if !s.Won() {
		cat &^= AnyPercentSpeedrun
	}
	cat |= AllCheckpointsSpeedrun | AllFlippedSpeedrun | AllSignsSpeedrun | NoPushSpeedrun
	for cp, cpSp := range s.Level.Checkpoints {
		if cp == "" {
			// Start is not a real CP.
			continue
		}
		if propmap.ValueOrP(cpSp.Properties, "secret", false, nil) {
			// Secrets are not needed for 100%, all paths or all signs run.
			// However they have their own run category here.

			// As all secrets have a sign, this isn't _necessary_, but still,
			// let's be robust.
			if s.CheckpointSeen(cp) == NotSeen {
				cat &^= AllSecretsSpeedrun
			}

			// All signs in secret rooms need to have been seen.
			for _, sign := range s.Level.TnihSignsByCheckpoint[cp] {
				if !propmap.ValueOrP(sign.PersistentState, "seen", false, nil) {
					cat &^= AllSecretsSpeedrun
				}
			}
			continue
		}
		switch s.CheckpointSeen(cp) {
		case NotSeen:
			cat &^= AllCheckpointsSpeedrun
		case SeenNormal:
			// Note: this means AllFlipped is possible without 100%. WAI.
			cat &^= AllFlippedSpeedrun
		}
		if s.CheckpointSeen(cp) != NotSeen {
			for _, next := range s.Level.CheckpointLocations.Locs[cp].NextByDir {
				// Skip non-forward edges.
				if !next.Forward || next.Optional {
					continue
				}
				// Skip secrets.
				nextCpSp := s.Level.Checkpoints[next.Other]
				if propmap.ValueOrP(nextCpSp.Properties, "secret", false, nil) {
					continue
				}
				// Only if the other CP was actually hit.
				// (i.e. Any% All Paths means all paths between the CPs actually
				//  hit, and 100% All Paths means really all paths)
				if s.CheckpointSeen(next.Other) == NotSeen {
					continue
				}
				if !s.CheckpointsWalked(cp, next.Other) {
					// log.Infof("MISSING PATH: %v %v", cp, next.Other)
					cat &^= AllPathsSpeedrun
				}
			}
		}
		for _, sign := range s.Level.TnihSignsByCheckpoint[cp] {
			if !propmap.ValueOrP(sign.PersistentState, "seen", false, nil) {
				cat &^= AllSignsSpeedrun
			}
		}
	}
	if s.Escapes() != 0 {
		// Note: this is impossible when also having AllSecrets,
		// as secrets typically cannot be left.
		cat &^= NoEscapeSpeedrun
	}
	if s.Teleports() != 0 {
		// Note: this can in theory be combined with AllSecrets.
		cat &^= NoTeleportsSpeedrun
	}
	if s.HasAbility("push") {
		// Probably can't be combined with much.
		cat &^= NoPushSpeedrun
	}
	return cat
}

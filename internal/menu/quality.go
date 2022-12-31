// Copyright 2022 Google LLC
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

package menu

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/demo"
	"github.com/divVerent/aaaaxy/internal/dump"
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	autoAdjustQuality = flag.Bool("auto_adjust_quality", true, "automatically adjust graphics quality to keep good fps")
)

type qualitySetting int

const (
	lowestQuality qualitySetting = iota
	lowQuality
	mediumQuality
	highQuality
	maxQuality
	autoQuality
	qualitySettingCount
)

func (s qualitySetting) String() string {
	switch s {
	case autoQuality:
		return locale.G.Get("Auto (%s)", currentActualQuality())
	case maxQuality:
		return locale.G.Get("Max")
	case highQuality:
		return locale.G.Get("High")
	case mediumQuality:
		return locale.G.Get("Medium")
	case lowQuality:
		return locale.G.Get("Low")
	case lowestQuality:
		return locale.G.Get("Lowest")
	}
	return locale.G.Get("???")
}

func currentQuality() qualitySetting {
	if *autoAdjustQuality {
		return autoQuality
	}
	return currentActualQuality()
}

func currentActualQuality() qualitySetting {
	if flag.Get[string]("screen_filter") == "linear2xcrt" {
		return maxQuality
	}
	if flag.Get[bool]("draw_outside") {
		return highQuality
	}
	if flag.Get[bool]("draw_blurs") {
		return mediumQuality
	}
	if flag.Get[bool]("expand_using_vertices_accurately") {
		return lowQuality
	}
	return lowestQuality
}

func (s qualitySetting) apply() error {
	if s == autoQuality {
		*autoAdjustQuality = true
		return maxQuality.applyActual()
	}
	*autoAdjustQuality = false
	return s.applyActual()
}

func (s qualitySetting) applyActual() error {
	switch s {
	case maxQuality:
		flag.Set("draw_blurs", true)
		flag.Set("draw_outside", true)
		flag.Set("expand_using_vertices_accurately", true)
		flag.Set("screen_filter", "linear2xcrt") // <-
	case highQuality:
		flag.Set("draw_blurs", true)
		flag.Set("draw_outside", true) // <-
		flag.Set("expand_using_vertices_accurately", true)
		flag.Set("screen_filter", "simple")
	case mediumQuality:
		flag.Set("draw_blurs", true) // <-
		flag.Set("draw_outside", false)
		flag.Set("expand_using_vertices_accurately", true)
		flag.Set("screen_filter", "simple") // <-
	case lowQuality:
		flag.Set("draw_blurs", false)
		flag.Set("draw_outside", false)
		flag.Set("expand_using_vertices_accurately", true) // <-
		flag.Set("screen_filter", "nearest")
	case lowestQuality:
		flag.Set("draw_blurs", false)
		flag.Set("draw_outside", false)
		flag.Set("expand_using_vertices_accurately", false)
		flag.Set("screen_filter", "nearest")
	}
	return nil
}

const (
	// Must reach about 50fps at least twice every 10 seconds.
	minFPS               = 49
	measureQualityFrames = 600
	minGoodQualityFrames = 100
)

var (
	totalQualityFrames int = 0
	goodQualityFrames  int = 0
)

func performQualityAdjustment() {
	// Don't auto adjust if disabled, dumping, benchmarking or not having focus.
	if !*autoAdjustQuality || dump.Active() || demo.Timedemo() || !ebiten.IsFocused() {
		totalQualityFrames = 0
		goodQualityFrames = 0
		return
	}

	// Check if downgrade is needed.
	totalQualityFrames++
	if ebiten.CurrentFPS() >= minFPS {
		goodQualityFrames++
	}
	if totalQualityFrames < measureQualityFrames {
		return
	}
	ok := goodQualityFrames >= minGoodQualityFrames
	totalQualityFrames = 0
	goodQualityFrames = 0
	if ok {
		return
	}

	// Downgrade quality.
	g := currentActualQuality()
	if g == lowestQuality {
		log.Warningf("couldn't even get good framerate at quality %v - cannot downgrade further", g)
	} else {
		log.Warningf("didn't get good framerate at quality %v - moving to %v", g, g-1)
		g--
		g.applyActual()
	}
}

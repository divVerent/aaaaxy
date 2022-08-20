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
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
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
		return fmt.Sprintf("Auto (%s)", currentActualQuality())
	case maxQuality:
		return "Max"
	case highQuality:
		return "High"
	case mediumQuality:
		return "Medium"
	case lowQuality:
		return "Low"
	case lowestQuality:
		return "Lowest"
	}
	return "???"
}

func currentQuality() qualitySetting {
	if *autoAdjustQuality {
		return autoQuality
	}
	return currentActualQuality()
}

func currentActualQuality() qualitySetting {
	if flag.Get("screen_filter").(string) == "linear2xcrt" {
		return maxQuality
	}
	if flag.Get("draw_outside").(bool) {
		return highQuality
	}
	if flag.Get("draw_blurs").(bool) {
		return mediumQuality
	}
	if flag.Get("expand_using_vertices_accurately").(bool) {
		return lowQuality
	}
	return lowestQuality
}

func (s qualitySetting) apply() error {
	if s == autoQuality {
		*autoAdjustQuality = true
		return nil
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
	if !*autoAdjustQuality {
		return
	}
	if !ebiten.IsFocused() {
		totalQualityFrames = 0
		goodQualityFrames = 0
		return
	}
	totalQualityFrames++
	if ebiten.CurrentFPS() >= minFPS {
		goodQualityFrames++
	}
	if totalQualityFrames < measureQualityFrames {
		return
	}
	if goodQualityFrames < minGoodQualityFrames {
		g := currentActualQuality()
		if g == lowestQuality {
			log.Warningf("couldn't even get good framerate at quality %v - giving up", g)
		} else {
			log.Warningf("didn't get good framerate at quality %v - moving to %v", g, g-1)
			g--
			g.applyActual()
		}
	}
	totalQualityFrames = 0
	goodQualityFrames = 0
}

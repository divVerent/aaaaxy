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

package timing

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	debugProfiling      = flag.Duration("debug_profiling", 0, "enable simple wall-clock profiling to log messages")
	debugFrameProfiling = flag.Bool("debug_frame_profiling", false, "print duration of each frame")
)

const (
	minReportDuration = time.Millisecond
)

type node struct {
	name    string
	started time.Time
}

type entry struct {
	total            time.Duration
	worstFrame       time.Duration
	thisFrame        time.Duration
	touchedThisFrame bool
	count            int
	frames           int
}

func (e *entry) String() string {
	c := e.count
	if c < 1 {
		c = 1
	}
	f := e.frames
	if f < 1 {
		f = 1
	}
	return fmt.Sprintf("%v (calls %d*%v, frames %d*%v, worst frame %v)",
		e.total,
		e.count,
		e.total/time.Duration(c),
		e.frames,
		e.total/time.Duration(f),
		e.worstFrame)
}

var (
	accumulator map[string]*entry
	stack       []node
	nextReport  time.Time
	prevFrame   time.Time
)

func restartProfiling() {
	accumulator, stack = map[string]*entry{}, []node{
		{name: "", started: time.Time{}},
	}
}

func stopProfiling() {
	accumulator, stack = nil, nil
}

func Group() func() {
	if stack != nil {
		sameName := stack[len(stack)-1].name
		stack = append(stack, node{name: sameName, started: time.Now()})
	}
	return endGroup
}

func endGroup() {
	if stack == nil {
		return
	}
	accountTime(time.Now())
	stack = stack[:len(stack)-1]
}

func Section(section string) {
	if stack == nil {
		return
	}
	now := time.Now()
	accountTime(now)
	newName := stack[len(stack)-2].name
	if section != "" {
		newName += "/" + section
	}
	stack[len(stack)-1] = node{name: newName, started: now}
	accountCount()
}

func current() (*node, *entry) {
	n := &stack[len(stack)-1]
	section := n.name
	var e *entry
	e, found := accumulator[section]
	if !found {
		e = &entry{}
		accumulator[section] = e
	}
	return n, e
}

func accountCount() {
	_, e := current()
	e.count++
}

func accountTime(now time.Time) {
	n, e := current()
	e.thisFrame += now.Sub(n.started)
	e.touchedThisFrame = true
}

func Update() {
	now := time.Now()
	if *debugFrameProfiling {
		if !prevFrame.IsZero() {
			delta := now.Sub(prevFrame)
			log.Infof("frame time: %v", delta)
		}
	}
	prevFrame = now
	if *debugProfiling != 0 && stack == nil {
		restartProfiling()
		return
	}
	if *debugProfiling == 0 {
		stopProfiling()
		return
	}
	for _, entry := range accumulator {
		if !entry.touchedThisFrame {
			continue
		}
		entry.total += entry.thisFrame
		if entry.thisFrame > entry.worstFrame {
			entry.worstFrame = entry.thisFrame
		}
		entry.frames++
		entry.thisFrame = 0
		entry.touchedThisFrame = false
	}
	if now.After(nextReport) {
		PrintReport()
		nextReport = now.Add(*debugProfiling)
		restartProfiling()
	}
}

func PrintReport() {
	if *debugProfiling == 0 {
		return
	}
	if !nextReport.IsZero() {
		report := make([]string, 0, len(accumulator))
		for section, entry := range accumulator {
			if entry.total < minReportDuration {
				continue
			}
			report = append(report, fmt.Sprintf("%-48s %v", section, entry))
		}
		sort.Strings(report)
		log.Infof("timing report:\n%v", strings.Join(report, "\n"))
	}
}

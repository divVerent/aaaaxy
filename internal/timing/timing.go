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
	"log"
	"sort"
	"strings"
	"time"
)

const (
	minReportDuration = time.Millisecond
)

type node struct {
	name    string
	started time.Time
}

type entry struct {
	duration time.Duration
	count    int
}

func (e *entry) String() string {
	c := e.count
	if c < 1 {
		c = 1
	}
	return fmt.Sprintf("%v (%d*%v)",
		e.duration,
		e.count,
		e.duration/time.Duration(c))
}

var (
	accumulator = map[string]*entry{}
	stack       []node
	nextReport  time.Time
)

func reset() {
	accumulator, stack = map[string]*entry{}, []node{
		{name: "", started: time.Time{}},
	}
}

func Group() func() {
	sameName := stack[len(stack)-1].name
	stack = append(stack, node{name: sameName, started: time.Now()})
	return endGroup
}

func endGroup() {
	accountTime(time.Now())
	stack = stack[:len(stack)-1]
}

func Section(section string) {
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
	e.duration += now.Sub(n.started)
}

func ReportRegularly() {
	now := time.Now()
	if now.After(nextReport) {
		if !nextReport.IsZero() {
			report := make([]string, 0, len(accumulator))
			for section, entry := range accumulator {
				if entry.duration < minReportDuration {
					continue
				}
				report = append(report, fmt.Sprintf("  %-48s %v", section, entry))
			}
			sort.Strings(report)
			log.Printf("Timing report:\n%v", strings.Join(report, "\n"))
		}
		reset()
		nextReport = now.Add(time.Second)
	}
}

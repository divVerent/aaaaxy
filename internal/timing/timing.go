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
	"log"
	"time"
)

const (
	minReportDuration = time.Millisecond
)

type node struct {
	name    string
	started time.Time
}

var (
	accumulator = map[string]time.Duration{}
	stack       []node
	nextReport  time.Time
)

func reset() {
	accumulator, stack = map[string]time.Duration{}, []node{
		{name: "", started: time.Time{}},
	}
}

func Group() func() {
	stack = append(stack, node{name: stack[len(stack)-1].name, started: time.Now()})
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
}

func accountTime(now time.Time) {
	accumulator[stack[len(stack)-1].name] += now.Sub(stack[len(stack)-1].started)
}

func ReportRegularly() {
	now := time.Now()
	if now.After(nextReport) {
		if !nextReport.IsZero() {
			for section, duration := range accumulator {
				if duration < minReportDuration {
					delete(accumulator, section)
				}
			}
			log.Printf("Timing report: %v", accumulator)
		}
		reset()
		nextReport = now.Add(time.Second)
	}
}

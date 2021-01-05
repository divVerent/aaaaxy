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

// Copyright 2021 Google LLC
//
// Licensed under the Apache Livense, Version 2.0 (the "License");
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

package level

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

type (
	CheckpointLocations struct {
		Locs map[string]*CheckpointLocation
		Rect m.Rect
	}
	CheckpointLocation struct {
		MapPos       m.Pos
		NextByDir    map[m.Delta]CheckpointEdge // Note: two sided.
		NextDeadEnds []CheckpointEdge
	}
	CheckpointEdge struct {
		Other    string
		Forward  bool
		Optional bool
	}
)

type fraction struct {
	num, denom int
}

func (a fraction) Less(b fraction) bool {
	s := 1
	if a.denom*b.denom < 0 {
		s = -1
	}
	return a.num*b.denom*s < b.num*a.denom*s
}

type edge struct {
	a, b           string
	priority       int
	unstraightness fraction
}

func unstraightness(d m.Delta) fraction {
	dx := d.DX
	if dx < 0 {
		dx = -dx
	}
	dy := d.DY
	if dy < 0 {
		dy = -dy
	}
	if dx < dy {
		return fraction{dx, dy}
	} else {
		return fraction{dy, dx}
	}
}

var AllCheckpointDirs = []m.Delta{
	m.North(),
	m.East(),
	m.South(),
	m.West(),
}

func possibleDirs(d m.Delta) (m.Delta, m.Delta) {
	dx := 0
	if d.DX > 0 {
		dx = 1
	}
	if d.DX < 0 {
		dx = -1
	}
	dy := 0
	if d.DY > 0 {
		dy = 1
	}
	if d.DY < 0 {
		dy = -1
	}
	prefersX := d.DX*dx > d.DY*dy // Sorry, have no abs.
	if prefersX {
		if dy == 0 {
			// x-only.
			return m.Delta{DX: dx, DY: 0}, m.Delta{DX: dx, DY: 0}
		} else {
			// x better than y.
			return m.Delta{DX: dx, DY: 0}, m.Delta{DX: 0, DY: dy}
		}
	} else {
		if dx == 0 {
			// y-only.
			return m.Delta{DX: 0, DY: dy}, m.Delta{DX: 0, DY: dy}
		} else {
			// y better than x.
			return m.Delta{DX: 0, DY: dy}, m.Delta{DX: dx, DY: 0}
		}
	}
}

func (l *Level) LoadCheckpointLocations(filename string) (*CheckpointLocations, error) {
	r, err := vfs.Load("generated", filename+".cp.json")
	if err != nil {
		return nil, fmt.Errorf("could not load checkpoint locations for %q: %v", filename, err)
	}
	var g JSONCheckpointGraph
	if err := json.NewDecoder(r).Decode(&g); err != nil {
		return nil, fmt.Errorf("could not decode checkpoint locations for %q: %v", filename, err)
	}
	var loc0 *CheckpointLocations
	var err0 error
	tryAtAngle := func(x, y int) {
		if loc0 != nil {
			return
		}
		loc, err := l.loadCheckpointLocations(filename, g, m.Delta{DX: x, DY: y}, m.Delta{DX: -y, DY: x})
		if err == nil {
			log.Infof("note: loading checkpoint locations required rotation by 1 0 -> %v %v", x, y)
			loc0, err0 = loc, nil
		} else if err0 == nil {
			err0 = err
		}
	}
	// Try known solution.
	// tryAtAngle(32, -9, false)
	// Brute force possible rotations.
	tryAtAngle(1, 0)
	b := 32
	for a := 1; a < b; a++ {
		tryAtAngle(b, a)
		tryAtAngle(b, -a)
	}
	return loc0, err0
}

// loadCheckpointLocations loads the checkpoint locations for the given level, possibly with a matrix transform.
func (l *Level) loadCheckpointLocations(filename string, g JSONCheckpointGraph, right, down m.Delta) (*CheckpointLocations, error) {
	id2name := map[EntityID]string{}
	loc := &CheckpointLocations{
		Locs: map[string]*CheckpointLocation{},
	}
	var minPos, maxPos m.Pos
	for _, o := range g.Objects {
		if o.Name == "" {
			// Not a CP, but the player initial spawn.
			continue
		}
		cp := l.Checkpoints[o.Name]
		if cp == nil {
			return nil, fmt.Errorf("could not find checkpoint referenced by locations for %q in %q", o.Name, filename)
		}
		rawPos, err := o.MapPos()
		if err != nil {
			return nil, fmt.Errorf("could not parse checkpoint location %q for %q in %q: %v", o.Pos, o.Name, filename, err)
		}
		pos := m.Pos{
			X: rawPos.Delta(m.Pos{}).Dot(right),
			Y: rawPos.Delta(m.Pos{}).Dot(down),
		}
		if len(loc.Locs) == 0 || pos.X < minPos.X {
			minPos.X = pos.X
		}
		if len(loc.Locs) == 0 || pos.Y < minPos.Y {
			minPos.Y = pos.Y
		}
		if len(loc.Locs) == 0 || pos.X > maxPos.X {
			maxPos.X = pos.X
		}
		if len(loc.Locs) == 0 || pos.Y > maxPos.Y {
			maxPos.Y = pos.Y
		}
		loc.Locs[o.Name] = &CheckpointLocation{
			MapPos: pos,
		}
	}
	loc.Rect = m.Rect{
		Origin: minPos,
		Size:   maxPos.Delta(minPos),
	}
	for name, cp := range l.Checkpoints {
		if name == "" {
			// Not a real CP, but the player initial spawn.
			continue
		}
		id2name[cp.ID] = name
	}
	edges := []edge{}
	nodeDegrees := make(map[string]int, len(l.Checkpoints))
	for name, cp := range l.Checkpoints {
		if name == "" {
			// Not a real CP, but the player initial spawn.
			continue
		}
		cpLoc := loc.Locs[name]
		if cpLoc == nil {
			return nil, fmt.Errorf("could not find checkpoint location for %q in %q", name, filename)
		}
		cpDeadEnd := cp.Properties["dead_end"] == "true"
		for propname, propval := range cp.Properties {
			if !strings.HasPrefix(propname, "next_") {
				continue
			}
			var nextID EntityID
			if _, err := fmt.Sscanf(propval, "%d", &nextID); err != nil {
				return nil, fmt.Errorf("could not parse next checkpoint ID %q for %q property %q in %q", propval, name, propname, filename)
			}
			other := id2name[nextID]
			if other == "" {
				return nil, fmt.Errorf("next checkpoint ID for %q property %q in %q is not a checkpoint", name, propname, filename)
			}
			otherDeadEnd := l.Checkpoints[other].Properties["dead_end"] == "true"
			otherLoc := loc.Locs[other]
			if otherLoc == nil {
				return nil, fmt.Errorf("next checkpoint %q in %q has no location yet", other, filename)
			}
			otherPos := otherLoc.MapPos
			moveDelta := otherPos.Delta(cpLoc.MapPos)
			unstraight := unstraightness(moveDelta)
			if cpDeadEnd || otherDeadEnd {
				cpLoc.NextDeadEnds = append(cpLoc.NextDeadEnds, CheckpointEdge{
					Other:   other,
					Forward: true,
				})
				otherLoc.NextDeadEnds = append(otherLoc.NextDeadEnds, CheckpointEdge{
					Other:   name,
					Forward: false,
				})
			} else {
				edges = append(edges, edge{
					a:              name,
					b:              other,
					unstraightness: unstraight,
				})
				nodeDegrees[name] += 1
				nodeDegrees[other] += 1
			}
		}
	}

	// alreadyAssigned returns if edge a -> b already has some assignment.
	alreadyAssigned := func(a, b string) bool {
		for _, edge := range loc.Locs[a].NextByDir {
			if edge.Other == b {
				return true
			}
		}
		return false
	}

	// assignEdge creates a required two-sided edge from a to b.
	assignEdge := func(a, b string, dir m.Delta, forward bool) bool {
		la := loc.Locs[a]
		if _, found := la.NextByDir[dir]; found {
			return false
		}
		lb := loc.Locs[b]
		revDir := dir.Mul(-1)
		if _, found := lb.NextByDir[revDir]; found {
			return false
		}
		la.NextByDir[dir] = CheckpointEdge{
			Other:   b,
			Forward: forward,
		}
		lb.NextByDir[revDir] = CheckpointEdge{
			Other:   a,
			Forward: !forward,
		}
		return true
	}

	// assignOptionalEdge assigns the both directions of edge from a to b, but ignores failure.
	assignOptionalEdge := func(a, b string, dir m.Delta) {
		la := loc.Locs[a]
		if _, found := la.NextByDir[dir]; !found {
			la.NextByDir[dir] = CheckpointEdge{
				Other:    b,
				Forward:  true,
				Optional: true,
			}
		}
		lb := loc.Locs[b]
		revDir := dir.Mul(-1)
		if _, found := lb.NextByDir[revDir]; !found {
			lb.NextByDir[revDir] = CheckpointEdge{
				Other:    a,
				Forward:  false,
				Optional: true,
			}
		}
	}

	// Group by quadrant. Insert perfectly straight edges into both neighboring quadrants.
	type quadrant struct {
		cp  string
		dir m.Delta
	}
	type quadEdge struct {
		other   string
		forward bool
	}
	quadMap := make(map[quadrant][]quadEdge)
	for _, edge := range edges {
		la := loc.Locs[edge.a]
		lb := loc.Locs[edge.b]
		delta := lb.MapPos.Delta(la.MapPos)
		maybeAddToQuadrant := func(dir m.Delta) {
			if delta.DX*dir.DX < 0 {
				return
			}
			if delta.DY*dir.DY < 0 {
				return
			}
			key := quadrant{
				cp:  edge.a,
				dir: dir,
			}
			quadMap[key] = append(quadMap[key], quadEdge{other: edge.b, forward: true})
			revDir := dir.Mul(-1)
			key = quadrant{
				cp:  edge.b,
				dir: revDir,
			}
			quadMap[key] = append(quadMap[key], quadEdge{other: edge.a, forward: false})
		}
		maybeAddToQuadrant(m.Delta{DX: 1, DY: 1})
		maybeAddToQuadrant(m.Delta{DX: 1, DY: -1})
		maybeAddToQuadrant(m.Delta{DX: -1, DY: 1})
		maybeAddToQuadrant(m.Delta{DX: -1, DY: -1})
	}

reprioritize:
	// Assign all edges to keyboard mapping.
	// Initialize map.
	for _, loc := range loc.Locs {
		loc.NextByDir = map[m.Delta]CheckpointEdge{}
	}

	var errorStrings []string
	collectError := func(format string, args ...interface{}) {
		errorStrings = append(errorStrings, fmt.Sprintf(format, args...))
	}

	// Every quadrant with three edges: GIVE UP.
	// Every quadrant with two edges: assign right away.
	for quad, others := range quadMap {
		if len(others) >= 3 {
			collectError("three checkpoint edges are in the same quadrant: %v -> %v", quad, others)
			continue
		}
		if len(others) < 2 {
			// Assign later.
			continue
		}
		// Precisely two others. The assignment is unique and well defined.
		la := loc.Locs[quad.cp]
		lb0 := loc.Locs[others[0].other]
		lb1 := loc.Locs[others[1].other]
		delta0 := lb0.MapPos.Delta(la.MapPos)
		delta1 := lb1.MapPos.Delta(la.MapPos)
		// Assign the straighter one to its preferred dir, and the less straight one to the remaining dir.
		if unstraightness(delta0).Less(unstraightness(delta1)) {
			bestDir, _ := possibleDirs(delta0)
			if !assignEdge(quad.cp, others[0].other, bestDir, others[0].forward) {
				collectError("could not fulfill forced first assignment in a quadrant: %v -> %v (%v -> %v)", quad, others, la, lb0)
			}
			if !assignEdge(quad.cp, others[1].other, quad.dir.Sub(bestDir), others[1].forward) {
				collectError("could not fulfill forced second assignment in a quadrant: %v -> %v (%v -> %v)", quad, others, la, lb1)
			}
		} else {
			bestDir, _ := possibleDirs(delta1)
			if !assignEdge(quad.cp, others[1].other, bestDir, others[0].forward) {
				collectError("could not fulfill forced first assignment in a quadrant: %v -> %v (%v -> %v)", quad, others, la, lb1)
			}
			if !assignEdge(quad.cp, others[0].other, quad.dir.Sub(bestDir), others[1].forward) {
				collectError("could not fulfill forced second assignment in a quadrant: %v -> %v (%v -> %v)", quad, others, la, lb0)
			}
		}
	}

	// Sort edges by unstraightness.
	sort.SliceStable(edges, func(a, b int) bool {
		dp := edges[a].priority - edges[b].priority
		if dp != 0 {
			// Highest priority first.
			return dp > 0
		}
		ua := edges[a].unstraightness
		ub := edges[b].unstraightness
		if ua.Less(ub) {
			// Straightest edges first.
			return true
		}
		if ub.Less(ua) {
			return false
		}
		// Tie breaker.
		na := fmt.Sprintf("%v -> %v", edges[a].a, edges[a].b)
		nb := fmt.Sprintf("%v -> %v", edges[b].a, edges[b].b)
		return na < nb
	})

	// Assign anything remaining in this preference order.
	for i := range edges {
		edge := &edges[i]
		if alreadyAssigned(edge.a, edge.b) {
			continue
		}
		// Try assigning the edge to its preferred dir, and if impossible, to its other dir.
		la := loc.Locs[edge.a]
		lb := loc.Locs[edge.b]
		delta := lb.MapPos.Delta(la.MapPos)
		bestDir, otherDir := possibleDirs(delta)

		if !assignEdge(edge.a, edge.b, bestDir, true) {
			if !assignEdge(edge.a, edge.b, otherDir, true) {
				if edge.priority < 1 {
					edge.priority++
					goto reprioritize
				}
				collectError("could not assign edge %v: no remaining assignments", edge)
			}
		}
	}

	// Finally fill up the keyboard directions.
	for i := len(edges) - 1; i >= 0; i-- {
		edge := edges[i]
		// Try unidirectionally assigning the remaining directions.
		la := loc.Locs[edge.a]
		lb := loc.Locs[edge.b]
		delta := lb.MapPos.Delta(la.MapPos)
		bestDir, otherDir := possibleDirs(delta)
		assignOptionalEdge(edge.a, edge.b, bestDir)
		assignOptionalEdge(edge.a, edge.b, otherDir)
	}

	if len(errorStrings) != 0 {
		return nil, errors.New(strings.Join(errorStrings, "; "))
	}

	return loc, nil
}

type JSONCheckpointGraph struct {
	Objects []JSONCheckpointObject
}

type JSONCheckpointObject struct {
	Name string
	Pos  string
}

func (o *JSONCheckpointObject) MapPos() (m.Pos, error) {
	var x, y float64
	if _, err := fmt.Sscanf(o.Pos, "%f,%f", &x, &y); err != nil {
		return m.Pos{}, err
	}
	// Note: reverse Y coordinate between graphviz and ebiten.
	return m.Pos{X: int(x), Y: -int(y)}, nil
}

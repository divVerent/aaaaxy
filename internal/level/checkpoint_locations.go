// Copyright 2021 Google LLC
//
// Licensed under the Apache License, SaveGameVersion 2.0 (the "License");
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
	"fmt"
	"github.com/divVerent/aaaaxy/internal/log"
	"math"
	"sort"
	"strings"

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

type edge struct {
	a, b           string
	priority       int
	unstraightness float64
}

func unstraightness(d m.Delta) float64 {
	dx := math.Abs(float64(d.DX))
	dy := math.Abs(float64(d.DY))
	return math.Min(dx, dy) / math.Max(dx, dy)
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
	loc, err := l.loadCheckpointLocations(filename, g, m.Delta{DX: 1, DY: 0}, m.Delta{DX: 0, DY: 1})
	if err == nil {
		return loc, nil
	}
	err0 := err
	for d := 16; d > 0; d-- {
		// Try some rotation.
		loc, err = l.loadCheckpointLocations(filename, g, m.Delta{DX: d, DY: 1}, m.Delta{DX: -1, DY: d})
		if err == nil {
			log.Infof("Note: loading checkpoint locations required rotation by %d 1", d)
			return loc, nil
		}
		// Try the opposite.
		loc, err = l.loadCheckpointLocations(filename, g, m.Delta{DX: d, DY: -1}, m.Delta{DX: 1, DY: d})
		if err == nil {
			log.Infof("Note: loading checkpoint locations required rotation by %d -1", d)
			return loc, nil
		}
	}
	return nil, err0
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
	// Assign all edges to keyboard mapping.
	// Note: there MIGHT be a shorter algorithm for all this, not sure.
	// Those three separate steps look suspicious.
	// This one is sure correct though, as whenever we choose the unpreferred direction,
	// we MUST chose it or we'd fail (so order of assigning the unpreferred ones does not matter).
	// However, if we assign the unpreferred one once we have to,
	// this helps choosing the unpreferred one in further assignments, so it is necessary.
	// Now translate to NextByDir. Successively map the "most straight direction" to the closest remaining available direction.
again:
	for _, loc := range loc.Locs {
		loc.NextByDir = map[m.Delta]CheckpointEdge{}
	}
	sort.Slice(edges, func(a, b int) bool {
		dp := edges[a].priority - edges[b].priority
		if dp != 0 {
			// Largest priority first.
			return dp > 0
		}
		da := nodeDegrees[edges[a].a] + nodeDegrees[edges[a].b]
		db := nodeDegrees[edges[b].a] + nodeDegrees[edges[b].b]
		dd := db - da
		if dd != 0 {
			// Largest degrees first.
			return dd > 0
		}
		du := edges[a].unstraightness - edges[b].unstraightness
		if du != 0 {
			// Straightest edges first.
			return du < 0
		}
		na := fmt.Sprintf("%v -> %v", edges[a].a, edges[a].b)
		nb := fmt.Sprintf("%v -> %v", edges[b].a, edges[b].b)
		return na < nb
	})
nextEdge:
	for i := range edges {
		edge := &edges[i]
		a := loc.Locs[edge.a]
		b := loc.Locs[edge.b]
		delta := b.MapPos.Delta(a.MapPos)
		bestDir, otherDir := possibleDirs(delta)
		for _, dir := range []m.Delta{bestDir, otherDir} {
			if _, found := a.NextByDir[dir]; found {
				continue
			}
			if _, found := b.NextByDir[dir.Mul(-1)]; found {
				continue
			}
			a.NextByDir[dir] = CheckpointEdge{
				Other:   edge.b,
				Forward: true,
			}
			b.NextByDir[dir.Mul(-1)] = CheckpointEdge{
				Other:   edge.a,
				Forward: false,
			}
			continue nextEdge
		}
		if edge.priority < 1 {
			log.Debugf("Prioritizing edge %v...", edge)
			edge.priority += 1
			goto again
		}
		return nil, fmt.Errorf("could not map edge %v to keyboard direction in %q", edge, filename)
	}
	// Now add the preferred direction unidirectionally whereever not there yet.
	for _, edge := range edges {
		a := loc.Locs[edge.a]
		b := loc.Locs[edge.b]
		delta := b.MapPos.Delta(a.MapPos)
		dir, _ := possibleDirs(delta)
		if _, found := a.NextByDir[dir]; !found {
			a.NextByDir[dir] = CheckpointEdge{
				Other:    edge.b,
				Forward:  true,
				Optional: true,
			}
		}
		if _, found := b.NextByDir[dir.Mul(-1)]; !found {
			b.NextByDir[dir.Mul(-1)] = CheckpointEdge{
				Other:    edge.a,
				Forward:  false,
				Optional: true,
			}
		}
	}
	// Now add the unpreferred direction undirectionally whereever not there yet.
	for i := len(edges) - 1; i >= 0; i-- {
		edge := edges[i]
		a := loc.Locs[edge.a]
		b := loc.Locs[edge.b]
		delta := b.MapPos.Delta(a.MapPos)
		_, dir := possibleDirs(delta)
		if _, found := a.NextByDir[dir]; !found {
			a.NextByDir[dir] = CheckpointEdge{
				Other:    edge.b,
				Forward:  true,
				Optional: true,
			}
		}
		if _, found := b.NextByDir[dir.Mul(-1)]; !found {
			b.NextByDir[dir.Mul(-1)] = CheckpointEdge{
				Other:    edge.a,
				Forward:  false,
				Optional: true,
			}
		}
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

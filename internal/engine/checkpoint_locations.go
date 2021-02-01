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

package engine

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"

	m "github.com/divVerent/aaaaaa/internal/math"
	"github.com/divVerent/aaaaaa/internal/vfs"
)

type (
	CheckpointLocations struct {
		Locs   map[string]*CheckpointLocation
		MinPos m.Pos
		MaxPos m.Pos
	}
	CheckpointLocation struct {
		MapPos    m.Pos
		NextByDir map[m.Delta]CheckpointEdge // Note: two sided.
	}
	CheckpointEdge struct {
		Other   string
		Forward bool
	}
)

type edge struct {
	a, b           string
	unstraightness float64
}

func unstraightness(d m.Delta) float64 {
	dx := math.Abs(float64(d.DX))
	dy := math.Abs(float64(d.DY))
	return math.Min(dx, dy) / math.Max(dx, dy)
}

var allDirs = []m.Delta{
	m.North(),
	m.East(),
	m.South(),
	m.West(),
}

// LoadCheckpointLocations loads the checkpoint locations for the given level.
func (l *Level) LoadCheckpointLocations(filename string) (*CheckpointLocations, error) {
	r, err := vfs.Load("maps", filename+".cp.json")
	if err != nil {
		return nil, fmt.Errorf("could not load checkpoint locations for %q: %v", filename, err)
	}
	var g JSONCheckpointGraph
	if err := json.NewDecoder(r).Decode(&g); err != nil {
		return nil, fmt.Errorf("could not decode checkpoint locations for %q: %v", filename, err)
	}
	id2name := map[EntityID]string{}
	loc := &CheckpointLocations{
		Locs: map[string]*CheckpointLocation{},
	}
	for _, o := range g.Objects {
		if o.Name == "" {
			// Not a CP, but the player initial spawn.
			continue
		}
		cp := l.Checkpoints[o.Name]
		if cp == nil {
			return nil, fmt.Errorf("could not find checkpoint referenced by locations for %q in %q", o.Name, filename)
		}
		pos, err := o.MapPos()
		if err != nil {
			return nil, fmt.Errorf("could not parse checkpoint location %q for %q in %q: %v", o.Pos, o.Name, filename)
		}
		if len(loc.Locs) == 0 || pos.X < loc.MinPos.X {
			loc.MinPos.X = pos.X
		}
		if len(loc.Locs) == 0 || pos.Y < loc.MinPos.Y {
			loc.MinPos.Y = pos.Y
		}
		if len(loc.Locs) == 0 || pos.X > loc.MaxPos.X {
			loc.MaxPos.X = pos.X
		}
		if len(loc.Locs) == 0 || pos.Y > loc.MaxPos.Y {
			loc.MaxPos.Y = pos.Y
		}
		loc.Locs[o.Name] = &CheckpointLocation{
			MapPos:    pos,
			NextByDir: map[m.Delta]CheckpointEdge{},
		}
	}
	for name, cp := range l.Checkpoints {
		if name == "" {
			// Not a real CP, but the player initial spawn.
			continue
		}
		id2name[cp.ID] = name
	}
	edges := []edge{}
	for name, cp := range l.Checkpoints {
		if name == "" {
			// Not a real CP, but the player initial spawn.
			continue
		}
		l := loc.Locs[name]
		if l == nil {
			return nil, fmt.Errorf("could not find checkpoint location for %q in %q", name, filename)
		}
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
			otherLoc := loc.Locs[other]
			if otherLoc == nil {
				return nil, fmt.Errorf("next checkpoint %q in %q has no location yet", other, filename)
			}
			otherPos := otherLoc.MapPos
			moveDelta := otherPos.Delta(l.MapPos)
			edges = append(edges, edge{
				a:              name,
				b:              other,
				unstraightness: unstraightness(moveDelta),
			})
		}
	}
	// Now translate to NextByDir. Successively map the "most straight direction" to the closest remaining available direction.
	sort.Slice(edges, func(a, b int) bool {
		return edges[a].unstraightness < edges[b].unstraightness
	})
	for _, edge := range edges {
		bestDir := m.Delta{}
		bestScore := 0
		a := loc.Locs[edge.a]
		b := loc.Locs[edge.b]
		delta := b.MapPos.Delta(a.MapPos)
		for _, dir := range allDirs {
			if _, found := a.NextByDir[dir]; found {
				continue
			}
			if _, found := b.NextByDir[dir.Mul(-1)]; found {
				continue
			}
			score := dir.Dot(delta)
			if score <= bestScore {
				continue
			}
			bestDir, bestScore = dir, score
		}
		if (bestDir == m.Delta{}) {
			return nil, fmt.Errorf("could not map edge %v to keyboard direction in %q", edge, filename)
		}
		a.NextByDir[bestDir] = CheckpointEdge{
			Other:   edge.b,
			Forward: true,
		}
		b.NextByDir[bestDir.Mul(-1)] = CheckpointEdge{
			Other:   edge.a,
			Forward: false,
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

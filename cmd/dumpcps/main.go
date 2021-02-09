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

package main

import (
	"fmt"
	"log"
	"math"

	"github.com/divVerent/aaaaaa/internal/engine"
	"github.com/divVerent/aaaaaa/internal/flag"
	m "github.com/divVerent/aaaaaa/internal/math"
)

var (
	initialCP = flag.String("initial_cp", "leap_of_faith", "Name of the initial checkpoint to base the layout on.")
)

type (
	Edge struct {
		WantDelta m.Delta
		From, To  *Vertex
	}
	Vertex struct {
		Name     string
		HasPos   bool
		WantPos  m.Pos
		OutEdges []*Edge
		InEdges  []*Edge
	}
)

func CalcPos(v *Vertex) {
	// Already done?
	if v.HasPos {
		return
	}
	// Nothing to do?
	if len(v.InEdges) == 0 {
		return
	}
	// First do all in-edges.
	var d m.Delta
	for _, in := range v.InEdges {
		CalcPos(in.From)
		d = d.Add(in.From.WantPos.Add(in.WantDelta).Delta(m.Pos{}))
	}
	d = d.Div(len(v.InEdges))
	v.WantPos = m.Pos{}.Add(d)
}

func main() {
	flag.Parse(flag.NoConfig)
	level, err := engine.LoadLevel("level")
	if err != nil {
		log.Panicf("Could not load level: %v", err)
	}
	// Gather a checkpoint ID to name map.
	cpMap := map[engine.EntityID]*engine.Spawnable{}
	for name, sp := range level.Checkpoints {
		if name == "" {
			// Ignore initial player spawn.
			continue
		}
		cpMap[sp.ID] = sp
	}
	// Generate all edges and vertices.
	vertices := map[engine.EntityID]*Vertex{}
	for id, sp := range cpMap {
		vertices[id] = &Vertex{Name: sp.Properties["name"]}
	}
	for id, sp := range cpMap {
		v := vertices[id]
		for _, conn := range []struct {
			name string
			dir  m.Delta
		}{
			{"next_left", m.West()},
			{"next_right", m.East()},
			{"next_up", m.North()},
			{"next_down", m.South()},
		} {
			next := sp.Properties[conn.name]
			if next == "" {
				continue
			}
			var nextID engine.EntityID
			if _, err := fmt.Sscanf(next, "%d", &nextID); err != nil {
				log.Panicf("Could not parse next CP %q -> %q: %v", sp.Properties["name"], next, err)
			}
			nextVert := vertices[nextID]
			if nextVert == nil {
				log.Panicf("Checkpoint %q doesn't point at a checkpoint but entity %d", sp.Properties["name"], nextID)
			}
			distance := 20
			if sp.Properties["dead_end"] == "true" {
				distance = 10
			}
			edge := &Edge{
				WantDelta: conn.dir.Mul(distance),
				From:      v,
				To:        nextVert,
			}
			v.OutEdges = append(v.OutEdges, edge)
			nextVert.InEdges = append(nextVert.InEdges, edge)
		}
	}
	// Build a .dot input file from all CPs.
	fmt.Print(`
		digraph G {
			layout = "neato";
			size = "256,256";
			overlap = false;
			splines = false;
		`)
	// Emit all nodes.
	for _, v := range vertices {
		CalcPos(v)
		fmt.Printf(`
				%s [width=2.0, height=2.0, fixedsize=true, shape=box, label="%s", pos="%d,%d"];
			`, v.Name, v.Name, v.WantPos.X, -v.WantPos.Y)
	}
	// Emit all edges.
	for _, v := range vertices {
		for _, e := range v.OutEdges {
			l := math.Sqrt(float64(e.To.WantPos.Delta(v.WantPos).Length2()))
			w := 1.0 / l
			port := func(d m.Delta) string {
				if d.DX < 0 {
					return "w"
				}
				if d.DX > 0 {
					return "e"
				}
				if d.DY < 0 {
					return "n"
				}
				if d.DY > 0 {
					return "s"
				}
				return "x"
			}
			headport := port(e.WantDelta.Mul(-1))
			tailport := port(e.WantDelta)
			fmt.Printf(`
					%s -> %s [len=%f, w=%f, headport=%s, tailport=%s];
				`, v.Name, e.To.Name, l, w, headport, tailport)
		}
	}
	fmt.Print(`
		}
		`)
}

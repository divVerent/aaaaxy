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
	"sort"
	"strings"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/level"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/propmap"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

type (
	Edge struct {
		WantDelta m.Delta
		From, To  *Vertex
	}
	Vertex struct {
		Name       string
		HasPos     bool
		CalcingPos bool
		MapPos     m.Pos
		WantPos    m.Pos
		OutEdges   []*Edge
		InEdges    []*Edge
	}
)

func CalcPos(v *Vertex) {
	// TODO: This algorithm is BAD. Need a better way that also kinda handles cycles.
	// Or maybe give up on hinting and find working options for dot that don't need it?
	// Already done?
	if v.HasPos {
		return
	}
	v.CalcingPos = true
	// First do all in-edges.
	var d m.Delta
	n := 0
	for _, in := range v.InEdges {
		if !in.From.CalcingPos {
			CalcPos(in.From)
			d = d.Add(in.From.WantPos.Add(in.WantDelta).Delta(m.Pos{}))
			n++
		}
	}
	// Nothing to do?
	if n > 0 {
		d = d.Div(n)
		v.WantPos = m.Pos{}.Add(d)
	}
	v.CalcingPos = false
}

func main() {
	log.Debugf("initializing VFS...")
	err := vfs.Init()
	if err != nil {
		log.Fatalf("could not initialize VFS: %v", err)
	}
	log.Debugf("parsing flags...")
	flag.Parse(flag.NoConfig)
	log.Debugf("loading level...")
	lvl, err := level.NewLoader("level").SkipCheckpointLocations(true).Load()
	if err != nil {
		log.Fatalf("could not load level: %v", err)
	}
	log.Debugf("generating checkpoint ID to name map...")
	cpMap := map[level.EntityID]*level.Spawnable{}
	for name, sp := range lvl.Checkpoints {
		if name == "" {
			// Ignore initial player spawn.
			continue
		}
		cpMap[sp.ID] = sp
	}
	log.Debugf("listing vertices...")
	vertices := map[level.EntityID]*Vertex{}
	var parseErr error
	for id, sp := range cpMap {
		vertices[id] = &Vertex{
			Name:   propmap.StringOr(sp.Properties, "name", ""),
			MapPos: sp.LevelPos.Mul(level.TileSize).Add(sp.RectInTile.Center().Delta(m.Pos{})),
		}
	}
	log.Debugf("listing entity IDs...")
	entityIDs := make([]level.EntityID, 0, len(cpMap))
	for id := range cpMap {
		entityIDs = append(entityIDs, id)
	}
	log.Debugf("sorting entity IDs...")
	sort.SliceStable(entityIDs, func(a, b int) bool {
		return entityIDs[a] < entityIDs[b]
	})
	log.Debugf("computing edges...")
	for _, id := range entityIDs {
		sp := cpMap[id]
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
			next := propmap.ValueOrP(sp.Properties, conn.name, -1, &parseErr)
			if next == -1 {
				continue
			}
			nextID := level.EntityID(next)
			nextSp := cpMap[nextID]
			nextVert := vertices[nextID]
			if nextVert == nil {
				log.Fatalf("checkpoint %q doesn't point at a checkpoint but entity %d", next, nextID)
			}
			distance := 10
			if propmap.ValueOrP(nextSp.Properties, "dead_end", false, &parseErr) {
				distance = 15
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
	if parseErr != nil {
		log.Fatalf("failed to parse: %v", parseErr)
	}
	log.Debugf("calculating positions...")
	for _, id := range entityIDs {
		v := vertices[id]
		CalcPos(v)
	}
	log.Debugf("writing header...")
	_, err = fmt.Print(`
		digraph G {
			layout = "neato";
			start = 4;  // Consistent random seed. Decided by fair dice roll.
			overlap = false;
			splines = false;
			maxiter = 131072;
			epsilon = 0.000001;
			// mode = KK;
			// model = circuit;
			// model = subset;
		`)
	if err != nil {
		log.Fatalf("failed to write to output: %v", err)
	}
	log.Debugf("writing vertices...")
	for _, id := range entityIDs {
		v := vertices[id]
		nameReadable := strings.ReplaceAll(v.Name, "_", "_\\n")
		_, err := fmt.Printf(`
				%s [width=2.0, height=2.0, fixedsize=true, shape=box, label="%s", pos="%d,%d"];
			`, v.Name, nameReadable, v.MapPos.X, -v.MapPos.Y)
		if err != nil {
			log.Fatalf("failed to write to output: %v", err)
		}
	}
	log.Debugf("writing edges...")
	for _, id := range entityIDs {
		v := vertices[id]
		for _, e := range v.OutEdges {
			_, err := fmt.Printf(`
					%s -> %s [len=%f];
				`, v.Name, e.To.Name, e.WantDelta.Length())
			if err != nil {
				log.Fatalf("failed to write to output: %v", err)
			}
		}
	}
	log.Debugf("writing footer...")
	_, err = fmt.Print(`
		}
		`)
	if err != nil {
		log.Fatalf("failed to write to output: %v", err)
	}
	log.Debugf("done.")
}

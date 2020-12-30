package aaaaaa

import (
	m "github.com/divVerent/aaaaaa/internal/math"
)

func appendLineToTraces(traces map[m.Delta]struct{}, start, end m.Delta) {
	delta := end.Sub(start)
	length := delta.Norm1()
	traces[start] = struct{}{}
	for i := MinEntitySize; i < length; i += MinEntitySize {
		pos := start.Add(delta.Mul(i).Div(length))
		traces[pos] = struct{}{}
	}
	traces[end] = struct{}{}
}

// TraceBox moves a size-sized box from from to to and yields info about where it hits solid etc.
func TraceBox(w *World, from m.Pos, size m.Delta, to m.Pos, o TraceOptions) TraceResult {
	// TODO make a real implementation.
	traces := map[m.Delta]struct{}{}
	delta := to.Delta(from)
	if delta.DX < 0 {
		appendLineToTraces(traces, m.Delta{DX: 0, DY: 0}, m.Delta{DX: 0, DY: size.DY - 1})
	} else {
		appendLineToTraces(traces, m.Delta{DX: size.DX - 1, DY: 0}, m.Delta{DX: size.DX - 1, DY: size.DY - 1})
	}
	if delta.DY < 0 {
		appendLineToTraces(traces, m.Delta{DX: 0, DY: 0}, m.Delta{DX: size.DX - 1, DY: 0})
	} else {
		appendLineToTraces(traces, m.Delta{DX: 0, DY: size.DY - 1}, m.Delta{DX: size.DX - 1, DY: size.DY - 1})
	}
	var result TraceResult
	var shortest int
	haveTrace := false
	for delta := range traces {
		trace := TraceLine(w, from.Add(delta), to.Add(delta), o)
		adjustedEnd := trace.EndPos.Sub(delta)
		length := adjustedEnd.Delta(from).Norm1()
		if !haveTrace || length < shortest {
			shortest = length
			haveTrace = true
			result.EndPos = adjustedEnd
			result.HitTilePos = trace.HitTilePos
			result.HitTile = trace.HitTile
			result.HitEntity = trace.HitEntity
			result.HitFogOfWar = trace.HitFogOfWar
		}
	}
	return result
}

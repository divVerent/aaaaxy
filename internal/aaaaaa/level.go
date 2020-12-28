package aaaaaa

import (
	"fmt"
	"io"

	"github.com/fardog/tmx"
)

// Level is a parsed form of a loaded level.
type Level struct {
	Tiles map[Pos]*LevelTile
}

// LevelTile is a single tile in the level.
type LevelTile struct {
	Tile     Tile
	Warpzone *Warpzone
}

// Warpzone represents a warp tile. Whenever anything enters this tile, it gets
// moved to "to" and the direction transformed by "transform". For the game to
// work, every warpzone must be paired with an exact opposite elsewhere. This
// is ensured at load time.
type Warpzone struct {
	ToTile    Pos
	Transform Orientation
}

type Spawnable struct {
	// Entity ID. Used to decide what needs spawning. Unique within a level.
	ID EntityID
	// Entity type. Used to spawn it on demand.
	EntityType  string
	LevelPos    Pos
	PosInTile   Delta
	Size        Delta
	Orientation Orientation
}

func LoadLevel(r io.Reader) (*Level, error) {
	t, err := tmx.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("invalid map: %v", err)
	}
	if t.Orientation != "orthogonal" {
		return nil, fmt.Errorf("unsupported map: got orientation %q, want orthogonal", t.Orientation)
	}
	if t.TileWidth != TileSize || t.TileHeight != TileSize {
		return nil, fmt.Errorf("unsupported map: got tile size %dx%d, want %dx%d", t.TileWidth, t.TileHeight, TileSize, TileSize)
	}
	if len(t.TileSets) != 1 {
		return nil, fmt.Errorf("unsupported map: got %d embedded tilesets, want 1", len(t.TileSets))
	}
	if len(t.Layers) != 1 {
		return nil, fmt.Errorf("unsupported map: got %d layers, want 1", len(t.Layers))
	}
	if len(t.ImageLayers) != 0 {
		return nil, fmt.Errorf("unsupported map: got %d image layers, want 0", len(t.ImageLayers))
	}
	tds, err := t.Layers[0].TileDefs(t.TileSets)
	if err != nil {
		return nil, fmt.Errorf("invalid map layer: %v", err)
	}
	Level := Level{}
	for i, td := range tds {
		pos := Pos{Y: i / t.Layers[0].Width, X: i % t.Layers[0].Width}
		orientation := Identity()
		if td.HorizontallyFlipped {
			orientation = FlipX().Concat(orientation)
		}
		if td.VerticallyFlipped {
			orientation = FlipY().Concat(orientation)
		}
		if td.DiagonallyFlipped {
			orientation = FlipD().Concat(orientation)
		}
		solid, err := td.Tile.Properties.Bool("solid")
		if err != nil {
			return nil, fmt.Errorf("invalid map: could not parse solid: %v", err)
		}
		Level.Tiles[pos] = &LevelTile{
			Tile: Tile{
				Solid:       solid,
				LevelPos:    pos,
				Image:       nil, // CachePic(td.Tile.Image),
				Orientation: orientation,
			},
		}
	}
	type RawWarpzone struct {
		StartTile, EndTile Pos
		Orientation        Orientation
	}
	warpzones := map[string][]RawWarpzone{}
	for _, og := range t.ObjectGroups {
		for _, o := range og.Objects {
			objProps := o.Properties
			if o.GlobalID != 0 {
				tile := t.TileSets[0].TileWithID(o.GlobalID.TileID(&t.TileSets[0]))
				objProps = append(append(tmx.Properties{}, objProps...), tile.Properties...)
			}
			objType := o.Type
			if objType == "" {
				objTypeProp := objProps.WithName("type")
				if objTypeProp != nil {
					objType = objTypeProp.Value
				}
			}
			startTile := Pos{X: int(o.X) / TileSize, Y: int(o.Y) / TileSize}
			endTile := Pos{X: int(o.X+o.Width-1) / TileSize, Y: int(o.Y+o.Height-1) / TileSize}
			orientation := Identity()
			orientationProp := objProps.WithName("orientation")
			if orientationProp != nil {
				orientation, err = ParseOrientation(orientationProp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid orientation: %v", err)
				}
			}
			if objType == "warpzone" {
				// Warpzones must be paired by name.
				warpzones[o.Name] = append(warpzones[o.Name], RawWarpzone{
					StartTile:   startTile,
					EndTile:     endTile,
					Orientation: orientation,
				})
			}
			delta := Delta{DX: int(o.X) % TileSize, DY: int(o.Y) % TileSize}
			ent := Spawnable{
				ID:          EntityID(o.ObjectID),
				EntityType:  objType,
				LevelPos:    startTile,
				PosInTile:   delta,
				Size:        Delta{DX: int(o.Width), DY: int(o.Height)},
				Orientation: orientation,
			}
			for y := startTile.Y; y <= endTile.Y; y++ {
				for x := startTile.X; x <= endTile.X; x++ {
					pos := Pos{X: x, Y: y}
					Level.Tiles[pos].Tile.Spawnables = append(Level.Tiles[pos].Tile.Spawnables, &ent)
				}
			}
		}
	}
	for warpname, warppair := range warpzones {
		if len(warppair) != 2 {
			return nil, fmt.Errorf("unpaired warpzone %q: got %d, want 2", warpname, len(warppair))
		}
		for a := 0; a < 2; a++ {
			from := warppair[a]
			to := warppair[1-a]
			// Warp orientation: right = direction to walk the warp, down = orientation (for mirroring).
			// Transform is identity transform iff the warps are reverse in right and identical in down.
			// T = to * flipx * from^-1
			// T' = from * flipx * to^-1
			// T T' = id
			transform := to.Orientation.Concat(FlipX()).Concat(from.Orientation.Inverse())
			fromCenter2 := from.StartTile.Add(from.EndTile.Delta(Pos{}))
			toCenter2 := to.StartTile.Add(to.EndTile.Delta(Pos{}))
			for fromy := from.StartTile.Y; fromy <= from.EndTile.Y; fromy++ {
				for fromx := from.StartTile.X; fromx <= from.EndTile.X; fromx++ {
					fromPos := Pos{X: fromx, Y: fromy}
					fromPos2 := fromPos.Add(fromPos.Delta(Pos{}))
					toPos2 := toCenter2.Add(transform.Apply(fromPos2.Delta(fromCenter2)))
					toPos := toPos2.Scale(1, 2).Add(to.Orientation.Apply(East()))
					Level.Tiles[fromPos].Warpzone = &Warpzone{
						ToTile:    toPos,
						Transform: transform,
					}
				}
			}
		}
	}
	return &Level, nil
}

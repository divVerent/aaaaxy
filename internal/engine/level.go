package engine

import (
	"fmt"
	"log"
	"strings"

	"github.com/fardog/tmx"
	"github.com/hajimehoshi/ebiten/v2"

	m "github.com/divVerent/aaaaaa/internal/math"
	"github.com/divVerent/aaaaaa/internal/vfs"
)

// Level is a parsed form of a loaded level.
type Level struct {
	Tiles  map[m.Pos]*LevelTile
	Player *Spawnable
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
	ToTile    m.Pos
	Transform m.Orientation
}

func LoadLevel(filename string) (*Level, error) {
	r, err := vfs.Load("maps", filename+".tmx")
	if err != nil {
		return nil, fmt.Errorf("could not open map: %v", err)
	}
	defer r.Close()
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
	level := Level{
		Tiles: map[m.Pos]*LevelTile{},
	}
	imgCache := map[string]*ebiten.Image{}
	cachePic := func(path string) (*ebiten.Image, error) {
		img := imgCache[path]
		if img == nil {
			var err error
			img, err = LoadImage("tiles", path)
			if err != nil {
				return nil, err
			}
			imgCache[path] = img
		}
		return img, nil
	}
	for i, td := range tds {
		if td.Nil {
			continue
		}
		pos := m.Pos{X: i % t.Layers[0].Width, Y: i / t.Layers[0].Width}
		orientation := m.Identity()
		if td.HorizontallyFlipped {
			orientation = m.FlipX().Concat(orientation)
		}
		if td.VerticallyFlipped {
			orientation = m.FlipY().Concat(orientation)
		}
		if td.DiagonallyFlipped {
			orientation = m.FlipD().Concat(orientation)
		}
		solid, err := td.Tile.Properties.Bool("solid")
		if err != nil {
			return nil, fmt.Errorf("invalid map: could not parse solid: %v", err)
		}
		opaque, err := td.Tile.Properties.Bool("opaque")
		if err != nil {
			return nil, fmt.Errorf("invalid map: could not parse opaque: %v", err)
		}
		img, err := cachePic(td.Tile.Image.Source)
		if err != nil {
			return nil, fmt.Errorf("invalid image: %v", err)
		}
		imgByOrientation := map[m.Orientation]*ebiten.Image{}
		for _, prop := range td.Tile.Properties {
			if oStr := strings.TrimPrefix(prop.Name, "img."); oStr != prop.Name {
				o, err := m.ParseOrientation(oStr)
				if err != nil {
					return nil, fmt.Errorf("invalid map: could not parse orientation tile: %v", err)
				}
				imgByOrientation[o], err = cachePic(prop.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid image: %v", err)
				}
			}
		}
		level.Tiles[pos] = &LevelTile{
			Tile: Tile{
				Solid:              solid,
				Opaque:             opaque,
				LevelPos:           pos,
				Image:              img,
				ImageByOrientation: imgByOrientation,
				Orientation:        orientation,
			},
		}
	}
	type RawWarpzone struct {
		StartTile, EndTile m.Pos
		Orientation        m.Orientation
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
			// TODO actually support object orientation.
			entRect := m.Rect{
				Origin: m.Pos{
					X: int(o.X),
					Y: int(o.Y),
				},
				Size: m.Delta{
					DX: int(o.Width),
					DY: int(o.Height),
				},
			}
			if o.GlobalID != 0 {
				// Tile entities are given by their bottom left coordinate in tiled.
				entRect.Origin.Y -= entRect.Size.DY
			}
			startTile := entRect.Origin.Div(TileSize)
			endTile := entRect.OppositeCorner().Div(TileSize)
			orientation := m.Identity()
			orientationProp := objProps.WithName("orientation")
			if orientationProp != nil {
				orientation, err = m.ParseOrientation(orientationProp.Value)
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
				continue
			}
			properties := map[string]string{}
			for _, prop := range objProps {
				properties[prop.Name] = prop.Value
			}
			ent := Spawnable{
				ID:         EntityID(o.ObjectID),
				EntityType: objType,
				LevelPos:   startTile,
				RectInTile: m.Rect{
					Origin: entRect.Origin.Sub(
						startTile.Mul(TileSize).Delta(m.Pos{})),
					Size: entRect.Size,
				},
				Orientation: orientation,
				Properties:  properties,
			}
			for y := startTile.Y; y <= endTile.Y; y++ {
				for x := startTile.X; x <= endTile.X; x++ {
					pos := m.Pos{X: x, Y: y}
					levelTile := level.Tiles[pos]
					if levelTile == nil {
						log.Panicf("invalid entity location: outside map bounds: %v in %v", pos, ent)
					}
					levelTile.Tile.Spawnables = append(levelTile.Tile.Spawnables, &ent)
				}
			}
			if objType == "player" {
				level.Player = &ent
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
			transform := to.Orientation.Concat(m.FlipX()).Concat(from.Orientation.Inverse())
			fromCenter2 := from.StartTile.Add(from.EndTile.Delta(m.Pos{}))
			toCenter2 := to.StartTile.Add(to.EndTile.Delta(m.Pos{}))
			for fromy := from.StartTile.Y; fromy <= from.EndTile.Y; fromy++ {
				for fromx := from.StartTile.X; fromx <= from.EndTile.X; fromx++ {
					fromPos := m.Pos{X: fromx, Y: fromy}
					fromPos2 := fromPos.Add(fromPos.Delta(m.Pos{}))
					toPos2 := toCenter2.Add(transform.Apply(fromPos2.Delta(fromCenter2)))
					toPos := toPos2.Div(2).Add(to.Orientation.Apply(m.West()))
					levelTile := level.Tiles[fromPos]
					if levelTile == nil {
						log.Panicf("invalid warpzone location: outside map bounds: %v in %v", fromPos, warppair)
					}
					toTile := level.Tiles[toPos]
					if toTile == nil {
						log.Panicf("invalid warpzone destination location: outside map bounds: %v in %v", toPos, warppair)
					}
					levelTile.Warpzone = &Warpzone{
						ToTile:    toPos,
						Transform: transform,
					}
				}
			}
		}
	}
	return &level, nil
}

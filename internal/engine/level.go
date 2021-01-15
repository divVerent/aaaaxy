package engine

import (
	"fmt"
	"log"
	"strings"

	"github.com/fardog/tmx"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/divVerent/aaaaaa/internal/image"
	m "github.com/divVerent/aaaaaa/internal/math"
	"github.com/divVerent/aaaaaa/internal/vfs"
)

// Level is a parsed form of a loaded level.
type Level struct {
	Tiles       map[m.Pos]*LevelTile
	Player      *Spawnable
	Checkpoints map[string]*Spawnable
}

// LevelTile is a single tile in the level.
type LevelTile struct {
	Tile     Tile
	WarpZone *WarpZone
}

// WarpZone represents a warp tile. Whenever anything enters this tile, it gets
// moved to "to" and the direction transformed by "transform". For the game to
// work, every warpZone must be paired with an exact opposite elsewhere. This
// is ensured at load time.
type WarpZone struct {
	ToTile    m.Pos
	Transform m.Orientation
}

// SaveGame is the data structure we save game state with.
// It contains all needed (in addition to loading the level) to reset to the last visited checkpoint.
type SaveGame struct {
	Data      map[EntityID]PersistentState
	DataHash  uint64
	LevelHash uint64
}

func (l *Level) Hash() uint64 {
	// TODO actually hash.
	return 42
}

func (save *SaveGame) Hash() uint64 {
	// TODO actually hash.
	return 42
}

func (l *Level) SaveGame() SaveGame {
	save := SaveGame{
		Data:      map[EntityID]PersistentState{},
		LevelHash: l.Hash(),
	}
	saveOne := func(s *Spawnable) {
		if len(s.PersistentState) > 0 {
			save.Data[s.ID] = s.PersistentState
		}
	}
	for _, tile := range l.Tiles {
		for _, s := range tile.Tile.Spawnables {
			saveOne(s)
		}
	}
	saveOne(l.Player)
	save.DataHash = save.Hash()
	return save
}

func (l *Level) LoadGame(save SaveGame) error {
	if save.DataHash != save.Hash() {
		return fmt.Errorf("someone tampered with the save game")
	}
	if save.LevelHash != l.Hash() {
		return fmt.Errorf("save game does not match level: got %v, want %v", save.LevelHash, l.Hash())
	}
	loadOne := func(s *Spawnable) {
		// Do not reallocate the map! Works better with already loaded entities.
		for key := range s.PersistentState {
			delete(s.PersistentState, key)
		}
		for key, value := range save.Data[s.ID] {
			s.PersistentState[key] = value
		}
	}
	for _, tile := range l.Tiles {
		for _, s := range tile.Tile.Spawnables {
			loadOne(s)
		}
	}
	loadOne(l.Player)
	return nil
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
	for i, ts := range t.TileSets {
		if ts.Source != "" {
			r, err := vfs.Load("tiles", ts.Source)
			if err != nil {
				return nil, fmt.Errorf("could not open tileset: %v", err)
			}
			defer r.Close()
			decoded, err := tmx.DecodeTileset(r)
			if err != nil {
				return nil, fmt.Errorf("could not decode tileset: %v", err)
			}
			decoded.FirstGlobalID = ts.FirstGlobalID
			t.TileSets[i] = *decoded
		}
	}
	tds, err := t.Layers[0].TileDefs(t.TileSets)
	if err != nil {
		return nil, fmt.Errorf("invalid map layer: %v", err)
	}
	level := Level{
		Tiles:       map[m.Pos]*LevelTile{},
		Checkpoints: map[string]*Spawnable{},
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
		properties := map[string]string{}
		for _, prop := range td.Tile.Properties {
			properties[prop.Name] = prop.Value
		}
		solid := properties["solid"] != "false"
		opaque := properties["opaque"] != "false"
		img, err := image.Load("tiles", td.Tile.Image.Source)
		if err != nil {
			return nil, fmt.Errorf("invalid image: %v", err)
		}
		imgByOrientation := map[m.Orientation]*ebiten.Image{}
		for propName, propValue := range properties {
			if oStr := strings.TrimPrefix(propName, "img."); oStr != propName {
				o, err := m.ParseOrientation(oStr)
				if err != nil {
					return nil, fmt.Errorf("invalid map: could not parse orientation tile: %v", err)
				}
				if o == m.Identity() && propValue != td.Tile.Image.Source {
					return nil, fmt.Errorf("invalid tileset: unrotated image isn't same as img: got %q, want %q", propValue, td.Tile.Image.Source)
				}
				imgByOrientation[o], err = image.Load("tiles", propValue)
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
	type RawWarpZone struct {
		StartTile, EndTile m.Pos
		Orientation        m.Orientation
	}
	warpZones := map[string][]RawWarpZone{}
	for _, og := range t.ObjectGroups {
		for _, o := range og.Objects {
			properties := map[string]string{}
			if o.Name != "" {
				properties["name"] = o.Name
			}
			for _, prop := range o.Properties {
				properties[prop.Name] = prop.Value
			}
			var tile *tmx.Tile
			if o.GlobalID != 0 {
				tile = t.TileSets[0].TileWithID(o.GlobalID.TileID(&t.TileSets[0]))
				properties["image"] = tile.Image.Source
				for _, prop := range tile.Properties {
					properties[prop.Name] = prop.Value
				}
			}
			objType := o.Type
			if objType == "" {
				objType = properties["type"]
			}
			if objType == "" && tile != nil {
				// This entity renders a tile.
				objType = "Sprite"
				properties["image_dir"] = "tiles"
				properties["image"] = tile.Image.Source
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
			startTile := entRect.Origin.Div(TileSize)
			endTile := entRect.OppositeCorner().Div(TileSize)
			orientation := m.Identity()
			if orientationProp := properties["orientation"]; orientationProp != "" {
				orientation, err = m.ParseOrientation(orientationProp)
				if err != nil {
					return nil, fmt.Errorf("invalid orientation: %v", err)
				}
			}
			if objType == "WarpZone" {
				// WarpZones must be paired by name.
				warpZones[properties["name"]] = append(warpZones[properties["name"]], RawWarpZone{
					StartTile:   startTile,
					EndTile:     endTile,
					Orientation: orientation,
				})
				continue
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
				Orientation:     orientation,
				Properties:      properties,
				PersistentState: PersistentState{},
			}
			if objType == "Player" {
				level.Player = &ent
				level.Checkpoints[""] = &ent
				// Do not link to tiles.
				continue
			}
			if objType == "Checkpoint" {
				level.Checkpoints[properties["name"]] = &ent
				// These do get linked.
			}
			for y := startTile.Y; y <= endTile.Y; y++ {
				for x := startTile.X; x <= endTile.X; x++ {
					pos := m.Pos{X: x, Y: y}
					levelTile := level.Tiles[pos]
					if levelTile == nil {
						log.Panicf("Invalid entity location: outside map bounds: %v in %v", pos, ent)
					}
					levelTile.Tile.Spawnables = append(levelTile.Tile.Spawnables, &ent)
				}
			}
		}
	}
	for warpname, warppair := range warpZones {
		if len(warppair) != 2 {
			return nil, fmt.Errorf("unpaired WarpZone %q: got %d, want 2", warpname, len(warppair))
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
						log.Panicf("Invalid WarpZone location: outside map bounds: %v in %v", fromPos, warppair)
					}
					toTile := level.Tiles[toPos]
					if toTile == nil {
						log.Panicf("Invalid WarpZone destination location: outside map bounds: %v in %v", toPos, warppair)
					}
					levelTile.WarpZone = &WarpZone{
						ToTile:    toPos,
						Transform: transform,
					}
				}
			}
		}
	}
	return &level, nil
}

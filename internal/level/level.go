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
	"fmt"
	"strconv"
	"strings"

	"github.com/fardog/tmx"
	"github.com/mitchellh/hashstructure/v2"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/version"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

var (
	debugCheckTnihSigns = flag.Bool("debug_check_tnih_signs", false, "if set, we verify that all checkpoints have a TnihSign")
)

// Level is a parsed form of a loaded level.
type Level struct {
	Player                *Spawnable
	Checkpoints           map[string]*Spawnable
	TnihSignsByCheckpoint map[string][]*Spawnable
	CheckpointLocations   *CheckpointLocations
	SaveGameVersion       int
	CreditsMusic          string
	Hash                  uint64 `hash:"-"`

	tiles []LevelTile
	width int
}

// Tile returns the tile at the given position.
func (l *Level) Tile(pos m.Pos) *LevelTile {
	i := pos.X + pos.Y*l.width
	t := &l.tiles[i]
	if !t.Valid {
		return nil
	}
	return t
}

// setTile sets the tile at the given position.
func (l *Level) setTile(pos m.Pos, t *LevelTile) {
	i := pos.X + pos.Y*l.width
	l.tiles[i] = *t
}

// ForEachTile iterates over all tiles in the level.
func (l *Level) ForEachTile(f func(pos m.Pos, t *LevelTile)) {
	for i := range l.tiles {
		f(m.Pos{X: i % l.width, Y: i / l.width}, &l.tiles[i])
	}
}

// LevelTile is a single tile in the level.
type LevelTile struct {
	Tile      Tile
	WarpZones []*WarpZone
	Valid     bool
}

// WarpZone represents a warp tile. Whenever anything enters this tile, it gets
// moved to "to" and the direction transformed by "transform". For the game to
// work, every warpZone must be paired with an exact opposite elsewhere. This
// is ensured at load time. Warpzones can be temporarily toggled by name; this
// state is lost on checkpoint restore.
type WarpZone struct {
	Name       string
	Invert     bool
	Switchable bool
	PrevTile   m.Pos
	ToTile     m.Pos
	Transform  m.Orientation
}

// SaveGameDataV1 is a not-yet-hashed SaveGame.
type SaveGameDataV1 struct {
	State        map[EntityID]PersistentState `hash:"-"`
	GameVersion  string
	LevelVersion int
	LevelHash    uint64
}

// SaveGame is the data structure we save game state with.
// It contains all needed (in addition to loading the level) to reset to the last visited checkpoint.
// Separate hashes govern the info parts and the state itself so demo regression testing can work across version changes.
type SaveGame struct {
	SaveGameDataV1
	InfoHash  uint64
	StateHash uint64

	// Legacy hash for v0 save games.
	Hash uint64 `json:",omitempty"`
}

// SaveGameData is a not-yet-hashed SaveGame.
type SaveGameData struct {
	State        map[EntityID]PersistentState
	LevelVersion int
	LevelHash    uint64
}

// SaveGame returns the current state as a SaveGame.
func (l *Level) SaveGame() (*SaveGame, error) {
	if l.SaveGameVersion != 1 {
		log.Fatalf("FIXME! On the next SaveGameVersion, please remove the SaveGameData v0 support, make all uint64 hashes `json:\",string\"`, and remove this check too!")
	}
	save := &SaveGame{
		SaveGameDataV1: SaveGameDataV1{
			State:        map[EntityID]PersistentState{},
			GameVersion:  version.Revision(),
			LevelVersion: l.SaveGameVersion,
			LevelHash:    l.Hash,
		},
	}
	saveOne := func(sp *Spawnable) {
		if len(sp.PersistentState) > 0 {
			save.State[sp.ID] = sp.PersistentState
		}
	}
	l.ForEachTile(func(_ m.Pos, tile *LevelTile) {
		for _, sp := range tile.Tile.Spawnables {
			saveOne(sp)
		}
	})
	saveOne(l.Player)
	var err error
	save.StateHash, err = hashstructure.Hash(save.State, hashstructure.FormatV2, nil)
	if err != nil {
		return nil, err
	}
	save.InfoHash, err = hashstructure.Hash(save.SaveGameDataV1, hashstructure.FormatV2, nil)
	if err != nil {
		return nil, err
	}
	return save, nil
}

// Clone clones the given Level to a new Level struct sharing no persistent state data.
func (l *Level) Clone() *Level {
	// First make a shallow copy.
	out := new(Level)
	*out = *l
	// Now "deepen" all that's needed. This means all Spawnable objects.
	// As they heavily alias each other, keep a map so the aliasing stays intact.
	clones := map[*Spawnable]*Spawnable{}
	clone := func(sp *Spawnable) *Spawnable {
		cloned := clones[sp]
		if cloned == nil {
			cloned = sp.Clone()
			clones[sp] = cloned
		}
		return cloned
	}
	out.Player = clone(out.Player)
	out.Checkpoints = make(map[string]*Spawnable, len(l.Checkpoints))
	for cp, cpSp := range l.Checkpoints {
		out.Checkpoints[cp] = clone(cpSp)
	}
	out.TnihSignsByCheckpoint = make(map[string][]*Spawnable, len(l.TnihSignsByCheckpoint))
	for cp, signs := range l.TnihSignsByCheckpoint {
		outSigns := make([]*Spawnable, len(signs))
		out.TnihSignsByCheckpoint[cp] = outSigns
		for i, sign := range signs {
			outSigns[i] = clone(sign)
		}
	}
	out.tiles = make([]LevelTile, len(l.tiles))
	for i := range l.tiles {
		tile := &l.tiles[i]
		outTile := &out.tiles[i]
		*outTile = *tile
		outTile.Tile.Spawnables = make([]*Spawnable, len(tile.Tile.Spawnables))
		for i, sp := range tile.Tile.Spawnables {
			outTile.Tile.Spawnables[i] = clone(sp)
		}
	}
	return out
}

// LoadGame loads the given SaveGame into the map.
// Note that when this returns an error, the SaveGame might have been partially loaded and the world may need to be reset.
func (l *Level) LoadGame(save *SaveGame) error {
	if save.Hash != 0 && save.InfoHash == 0 && save.StateHash == 0 {
		saveV0 := &SaveGameData{
			State:        save.State,
			LevelVersion: save.LevelVersion,
			LevelHash:    save.LevelHash,
		}
		saveHash, err := hashstructure.Hash(saveV0, hashstructure.FormatV2, nil)
		if err != nil {
			return err
		}
		if saveHash != save.Hash {
			return fmt.Errorf("someone tampered with the save game: got %v, want %v", saveHash, save.Hash)
		}
	} else {
		infoHash, err := hashstructure.Hash(save.SaveGameDataV1, hashstructure.FormatV2, nil)
		if err != nil {
			return err
		}
		if infoHash != save.InfoHash {
			return fmt.Errorf("someone tampered with the save game info")
		}
		stateHash, err := hashstructure.Hash(save.State, hashstructure.FormatV2, nil)
		if err != nil {
			return err
		}
		if stateHash != save.StateHash {
			return fmt.Errorf("someone tampered with the save game state")
		}
	}
	if save.GameVersion != version.Revision() {
		log.Warningf("save game does not match game version: got %v, want %v", save.GameVersion, version.Revision())
	}
	if save.LevelVersion != l.SaveGameVersion {
		return fmt.Errorf("save game does not match level version: got %v, want %v", save.LevelVersion, l.SaveGameVersion)
	}
	if save.LevelHash != l.Hash {
		log.Warningf("save game does not match level hash: got %v, want %v; trying to load anyway", save.LevelHash, l.Hash)
	}
	loadOne := func(sp *Spawnable) {
		// Do not reallocate the map! Works better with already loaded entities.
		for key := range sp.PersistentState {
			delete(sp.PersistentState, key)
		}
		for key, value := range save.State[sp.ID] {
			sp.PersistentState[key] = value
		}
	}
	l.ForEachTile(func(_ m.Pos, tile *LevelTile) {
		for _, sp := range tile.Tile.Spawnables {
			loadOne(sp)
		}
	})
	loadOne(l.Player)
	return nil
}

func (l *Level) applyTileMod(startTile, endTile m.Pos, mods map[string]string) {
	var add, remove Contents
	switch mods["object_solid"] {
	case "true":
		add |= ObjectSolidContents
	case "false":
		remove |= ObjectSolidContents
	}
	switch mods["opaque"] {
	case "true":
		add |= OpaqueContents
	case "false":
		remove |= OpaqueContents
	}
	switch mods["player_solid"] {
	case "true":
		add |= PlayerSolidContents
	case "false":
		remove |= PlayerSolidContents
	}
	switch mods["solid"] {
	case "true":
		add |= SolidContents
	case "false":
		remove |= SolidContents
	}
	for y := startTile.Y; y <= endTile.Y; y++ {
		for x := startTile.X; x <= endTile.X; x++ {
			t := l.Tile(m.Pos{X: x, Y: y})
			if t == nil {
				continue
			}
			t.Tile.Contents |= add
			t.Tile.Contents &= ^remove
		}
	}
}

func FetchTileset(ts *tmx.TileSet) error {
	if ts.Source != "" {
		r, err := vfs.Load("", ts.Source)
		if err != nil {
			return fmt.Errorf("could not open tileset: %v", err)
		}
		defer r.Close()
		decoded, err := tmx.DecodeTileset(r)
		if err != nil {
			return fmt.Errorf("could not decode tileset: %v", err)
		}
		decoded.FirstGlobalID = ts.FirstGlobalID
		decoded.Source = ts.Source
		*ts = *decoded
	}
	if ts.TileWidth != TileSize || ts.TileHeight != TileSize {
		return fmt.Errorf("unsupported tileset: got tile size %dx%d, want %dx%d", ts.TileWidth, ts.TileHeight, TileSize, TileSize)
	}
	// ts.Spacing, ts.Margin, ts.TileCount, ts.Columns doesn't matter (we only support multi image tilesets).
	if ts.ObjectAlignment != "topleft" {
		return fmt.Errorf("unsupported tileset: got objectalignment %q, want topleft", ts.ObjectAlignment)
	}
	// ts.Properties doesn't matter.
	if ts.TileOffset.X != 0 || ts.TileOffset.Y != 0 {
		return fmt.Errorf("unsupported tileset: got a tile offset")
	}
	if ts.Image.Source != "" {
		return fmt.Errorf("unsupported tileset: got single image, want image collection")
	}
	// ts.TerrainTypes doesn't matter (editor only).
	// ts.Tiles used later.
	return nil
}

func Load(filename string) (*Level, error) {
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
	// t.RenderOrder doesn't matter.
	// t.Width, t.Height used later.
	if t.TileWidth != TileSize || t.TileHeight != TileSize {
		return nil, fmt.Errorf("unsupported map: got tile size %dx%d, want %dx%d", t.TileWidth, t.TileHeight, TileSize, TileSize)
	}
	// t.HexSideLength doesn't matter.
	// t.StaggerAxis doesn't matter.
	// t.StaggerIndex doesn't matter.
	// t.BackgroundColor doesn't matter.
	// t.NextObjectID doesn't matter.
	// t.TileSets used later.
	// t.Properties used later.
	if len(t.Layers) != 1 {
		return nil, fmt.Errorf("unsupported map: got %d layers, want 1", len(t.Layers))
	}
	// t.ObjectGroups used later.
	if len(t.ImageLayers) != 0 {
		return nil, fmt.Errorf("unsupported map: got %d image layers, want 0", len(t.ImageLayers))
	}
	for i := range t.TileSets {
		err := FetchTileset(&t.TileSets[i])
		if err != nil {
			return nil, fmt.Errorf("unsupported map: failed to decode tileset %d: %v", i, err)
		}
	}
	layer := &t.Layers[0]
	if layer.X != 0 || layer.Y != 0 {
		return nil, fmt.Errorf("unsupported map: layer has been shifted")
	}
	// layer.Width, layer.Height used later.
	// layer.Opacity, layer.Visible not used (we allow it though as it may help in the editor).
	if layer.OffsetX != 0 || layer.OffsetY != 0 {
		return nil, fmt.Errorf("unsupported map: layer has an offset")
	}
	// layer.Properties not used.
	// layer.RawData not used.
	tds, err := layer.TileDefs(t.TileSets)
	if err != nil {
		return nil, fmt.Errorf("invalid map layer: %v", err)
	}
	saveGameVersion, err := t.Properties.Int("save_game_version")
	if err != nil {
		return nil, fmt.Errorf("unsupported map: could not read save_game_version: %v", err)
	}
	var creditsMusic string
	if prop := t.Properties.WithName("credits_music"); prop != nil {
		creditsMusic = prop.Value
	}
	level := Level{
		Checkpoints:           map[string]*Spawnable{},
		TnihSignsByCheckpoint: map[string][]*Spawnable{},
		SaveGameVersion:       int(saveGameVersion),
		CreditsMusic:          creditsMusic,
		tiles:                 make([]LevelTile, layer.Width*layer.Height),
		width:                 layer.Width,
	}
	var tnihSigns []*Spawnable
	checkpoints := map[EntityID]*Spawnable{}
	for i, td := range tds {
		if td.Nil {
			continue
		}
		if td.Tile == nil {
			return nil, fmt.Errorf("invalid tiledef: %v [%s]", td, td.TileSet.Source)
		}
		// td.Tile.Probability not used (editor only).
		// td.Tile.Properties used later.
		// td.Tile.Image used later.
		if len(td.Tile.Animation) != 0 {
			return nil, fmt.Errorf("unsupported tileset: got an animation")
		}
		if len(td.Tile.ObjectGroup.Objects) != 0 {
			return nil, fmt.Errorf("unsupported tileset: got objects in a tile")
		}
		// td.Tile.RawTerrainType not used (editor only).
		pos := m.Pos{X: i % layer.Width, Y: i / layer.Width}
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
		for i := range td.Tile.Properties {
			prop := &td.Tile.Properties[i]
			properties[prop.Name] = prop.Value
		}
		var contents Contents
		if properties["solid"] != "false" {
			contents |= SolidContents
		}
		if properties["opaque"] != "false" {
			contents |= OpaqueContents
		}
		imgSrc := td.Tile.Image.Source
		imgSrcByOrientation, err := ParseImageSrcByOrientation(imgSrc, properties)
		if err != nil {
			return nil, fmt.Errorf("invalid map: %v", err)
		}
		level.setTile(pos, &LevelTile{
			Tile: Tile{
				Contents:              contents,
				LevelPos:              pos,
				ImageSrc:              imgSrc,
				imageSrcByOrientation: imgSrcByOrientation,
				Orientation:           orientation,
			},
			Valid: true,
		})
	}
	type RawWarpZone struct {
		StartTile, EndTile m.Pos
		Orientation        m.Orientation
		Invert             bool
		Switchable         bool
	}
	warpZones := map[string][]*RawWarpZone{}
	for i := range t.ObjectGroups {
		og := &t.ObjectGroups[i]
		// og.Name, og.Color not used (editor only).
		if og.X != 0 || og.Y != 0 {
			return nil, fmt.Errorf("unsupported map: object group has been shifted")
		}
		// og.Width, og.Height not used.
		// og.Opacity, og.Visible not used (we allow it though as it may help in the editor).
		if og.OffsetX != 0 || og.OffsetY != 0 {
			return nil, fmt.Errorf("unsupported map: object group has an offset")
		}
		// og.DrawOrder not used (we use our own z index).
		// og.Properties not used.
		for j := range og.Objects {
			o := &og.Objects[j]
			// o.ObjectID used later.
			properties := map[string]string{}
			if o.Name != "" {
				properties["name"] = o.Name
			}
			// o.X, o.Y, o.Width, o.Height used later.
			if o.Rotation != 0 {
				return nil, fmt.Errorf("unsupported map: object %v has a rotation (maybe implement this?)", o.ObjectID)
			}
			if o.GlobalID != 0 {
				var tile *tmx.Tile
				for k := range t.TileSets {
					ts := &t.TileSets[k]
					tile = ts.TileWithID(o.GlobalID.TileID(ts))
					if tile != nil {
						break
					}
				}
				if tile == nil {
					return nil, fmt.Errorf("unsupported map: object %v references nonexisting tile %d", o.ObjectID, o.GlobalID)
				}
				if tile.Type == "" {
					properties["type"] = "Sprite"
				} else {
					properties["type"] = tile.Type
				}
				properties["image_dir"] = "tiles"
				properties["image"] = tile.Image.Source
				for k := range tile.Properties {
					prop := &tile.Properties[k]
					properties[prop.Name] = prop.Value
				}
			}
			// o.Visible not used (we allow it though as it may help in the editor).
			if o.Polygons != nil {
				return nil, fmt.Errorf("unsupported map: object %v has polygons", o.ObjectID)
			}
			if o.Polylines != nil {
				return nil, fmt.Errorf("unsupported map: object %v has polylines", o.ObjectID)
			}
			if o.Image.Source != "" {
				properties["type"] = "Sprite"
				properties["image_dir"] = "sprites"
				properties["image"] = o.Image.Source
			}
			if o.Type != "" {
				properties["type"] = o.Type
			}
			for k := range o.Properties {
				prop := &o.Properties[k]
				properties[prop.Name] = prop.Value
			}
			// o.RawExtra not used.
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
			if properties["type"] == "WarpZone" {
				// WarpZones must be paired by name.
				name := properties["name"]
				invert := properties["invert"] == "true"         // Default false.
				switchable := properties["switchable"] == "true" // Default false.
				warpZones[name] = append(warpZones[name], &RawWarpZone{
					StartTile:   startTile,
					EndTile:     endTile,
					Orientation: orientation,
					Switchable:  switchable,
					Invert:      invert,
				})
				continue
			}
			ent := Spawnable{
				ID:         EntityID(o.ObjectID),
				EntityType: properties["type"],
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
			if properties["type"] == "_TileMod" {
				level.applyTileMod(startTile, endTile, properties)
				// Do not link to tiles.
				continue
			}
			if properties["type"] == "Player" {
				level.Player = &ent
				level.Checkpoints[""] = &ent
				// Do not link to tiles.
				continue
			}
			if properties["type"] == "Checkpoint" {
				level.Checkpoints[properties["name"]] = &ent
				checkpoints[ent.ID] = &ent
				// These do get linked.
			}
			if properties["type"] == "TnihSign" {
				tnihSigns = append(tnihSigns, &ent)
				// These do get linked.
			}
			for y := startTile.Y; y <= endTile.Y; y++ {
				for x := startTile.X; x <= endTile.X; x++ {
					pos := m.Pos{X: x, Y: y}
					levelTile := level.Tile(pos)
					if levelTile == nil {
						return nil, fmt.Errorf("invalid entity location: outside map bounds: %v in %v", pos, ent)
					}
					levelTile.Tile.Spawnables = append(levelTile.Tile.Spawnables, &ent)
				}
			}
		}
	}
	for _, sign := range tnihSigns {
		id, err := strconv.Atoi(sign.Properties["reached_from"])
		if err != nil {
			return nil, fmt.Errorf("invalid TnihSign: reached_from not set: %v: %v", sign, err)
		}
		cp := checkpoints[EntityID(id)]
		if cp == nil {
			return nil, fmt.Errorf("invalid TnihSign: checkpoint ID %v not found", id)
		}
		name := cp.Properties["name"]
		level.TnihSignsByCheckpoint[name] = append(level.TnihSignsByCheckpoint[name], sign)
	}
	for name, cpSp := range level.Checkpoints {
		if name == "" {
			// This isn't a real CP.
			continue
		}
		if *debugCheckTnihSigns {
			got := len(level.TnihSignsByCheckpoint[name]) != 0
			want := cpSp.Properties["tnih_sign_expected"] != "false" // default true
			if !got && want {
				log.Fatalf("note: checkpoint %v has no TnihSign - intended?", name)
			}
			if got && !want {
				log.Fatalf("note: checkpoint %v unexpectedly has TnihSign - intended?", name)
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
					prevPos := fromPos.Add(from.Orientation.Apply(m.West()))
					fromPos2 := fromPos.Add(fromPos.Delta(m.Pos{}))
					toPos2 := toCenter2.Add(transform.Apply(fromPos2.Delta(fromCenter2)))
					toPos := toPos2.Div(2).Add(to.Orientation.Apply(m.West()))
					levelTile := level.Tile(fromPos)
					if levelTile == nil {
						return nil, fmt.Errorf("invalid WarpZone location: outside map bounds: %v in %v", fromPos, warppair)
					}
					toTile := level.Tile(toPos)
					if toTile == nil {
						return nil, fmt.Errorf("invalid WarpZone destination location: outside map bounds: %v in %v", toPos, warppair)
					}
					levelTile.WarpZones = append(levelTile.WarpZones, &WarpZone{
						Name:       warpname,
						Invert:     from.Invert,
						Switchable: from.Switchable,
						PrevTile:   prevPos,
						ToTile:     toPos,
						Transform:  transform,
					})
				}
			}
		}
	}
	level.CheckpointLocations, err = level.LoadCheckpointLocations(filename)
	if err != nil {
		log.Errorf("could not load checkpoint locations: %v", err)
	}
	level.Hash, err = hashstructure.Hash(&level, hashstructure.FormatV2, nil)
	if err != nil {
		return nil, fmt.Errorf("could not hash level: %v", err)
	}
	log.Errorf("level hash: %v", level.Hash)
	return &level, nil
}

// VerifyHash returns an error if the level hash changed.
func (l *Level) VerifyHash() error {
	hash, err := hashstructure.Hash(l, hashstructure.FormatV2, nil)
	if err != nil {
		return fmt.Errorf("could not hash level: %v", err)
	}
	if hash != l.Hash {
		return fmt.Errorf("level hash mismatch: got %v, want %v", hash, l.Hash)
	}
	return nil
}

// ParseImageSrcByOrientation parses the imgSrcByOrientation map.
func ParseImageSrcByOrientation(defaultSrc string, properties map[string]string) (map[m.Orientation]string, error) {
	imgSrcByOrientation := map[m.Orientation]string{}
	for propName, propValue := range properties {
		if oStr := strings.TrimPrefix(propName, "img."); oStr != propName {
			o, err := m.ParseOrientation(oStr)
			if err != nil {
				return nil, fmt.Errorf("could not parse orientation tile: %v", err)
			}
			if o == m.Identity() && propValue != defaultSrc {
				return nil, fmt.Errorf("unrotated image isn't same as img: got %q, want %q", propValue, defaultSrc)
			}
			imgSrcByOrientation[o] = propValue
		}
	}
	return imgSrcByOrientation, nil
}

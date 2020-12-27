package aaaaaa

import (
	"fmt"
	"io"

	"github.com/fardog/tmx"
)

// level is a parsed form of a loaded level.
type level struct {
	tiles map[levelPos]*levelTile
}

// levelPos is a position in the level.
type levelPos struct {
	c, r int
}

// levelTile is a single tile in the level.
type levelTile struct {
	tile       Tile
	spawnables []*spawnable
	warp       *warpzone
}

// warpzone represents a warp tile. Whenever anything enters this tile, it gets
// moved to "to" and the direction transformed by "transform". For the game to
// work, every warpzone must be paired with an exact opposite elsewhere. This
// is ensured at load time.
type warpzone struct {
	to        levelPos
	transform Orientation
}

type spawnable struct {
	// Entity ID. Used to decide what needs spawning. Unique within a level.
	id EntityID
}

func LoadLevel(r io.Reader) (*level, error) {
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
	level := level{}
	for i, td := range tds {
		pos := levelPos{r: i / t.Layers[0].Width, c: i % t.Layers[0].Width}
		orientation := Orientation{Right: Delta{1, 0}, Down: Delta{0, 1}}
		if td.HorizontallyFlipped {
			orientation.Right.DX = -orientation.Right.DX
			orientation.Down.DX = -orientation.Down.DX
		}
		if td.VerticallyFlipped {
			orientation.Right.DY = -orientation.Right.DY
			orientation.Down.DY = -orientation.Down.DY
		}
		if td.DiagonallyFlipped {
			orientation.Right.DX, orientation.Right.DY = orientation.Right.DY, orientation.Right.DX
			orientation.Down.DX, orientation.Down.DY = orientation.Down.DY, orientation.Down.DX
		}
		solid, err := td.Tile.Properties.Bool("solid")
		if err != nil {
			return nil, fmt.Errorf("invalid map: could not parse solid: %v", err)
		}
		level.tiles[pos] = &levelTile{
			tile: Tile{
				Solid:       solid,
				levelPos:    pos,
				image:       nil, // CachePic(td.Tile.Image),
				orientation: orientation,
			},
		}
	}
	for _, og := range t.ObjectGroups {
		for _, o := range og.Objects {
			o = o
			// TODO: decode objects
		}
	}
	return &level, nil
}

package sound

import (
	"fmt"
	"io/ioutil"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"

	"github.com/divVerent/aaaaaa/internal/vfs"
)

// Sound represents a sound effect.
type Sound struct {
	sound   []byte
	players []*audio.Player
}

// Sounds are preloaded as byte streams.
var cache = map[string]*Sound{}

// Load loads a sound effect.
// Multiple Load calls to the same sound effect return the same cached instance.
func Load(name string) (*Sound, error) {
	if sound, found := cache[name]; found {
		return sound, nil
	}
	data, err := vfs.Load("sounds", name)
	if err != nil {
		return nil, fmt.Errorf("could not load: %v", err)
	}
	defer data.Close()
	stream, err := vorbis.Decode(audio.CurrentContext(), data)
	if err != nil {
		return nil, fmt.Errorf("could not start decoding: %v", err)
	}
	decoded, err := ioutil.ReadAll(stream)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %v", err)
	}
	sound := &Sound{sound: decoded}
	cache[name] = sound
	return sound, nil
}

// Play plays the given sound effect.
func (s *Sound) Play() {
	audio.NewPlayerFromBytes(audio.CurrentContext(), s.sound).Play()
}

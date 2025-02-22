// Copyright 2022 Google LLC
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

package propmap

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/mitchellh/hashstructure/v2"

	"github.com/divVerent/aaaaxy/internal/log"
)

// Map is a helper to easier deal with maps.
type Map struct {
	m map[string]string
}

// New creates a new map.
func New() Map {
	return Map{
		m: map[string]string{},
	}
}

// ForEach iterates over a map.
func ForEach(pm Map, f func(k, v string) error) error {
	for k, v := range pm.m {
		err := f(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// Empty tells if the map is empty.
func Empty(pm Map) bool {
	return len(pm.m) == 0
}

// Delete deletes a properties map entry.
func Delete(pm Map, key string) {
	delete(pm.m, key)
}

// Set sets a properties map entry.
func Set[V any](pm Map, key string, value V) {
	pm.m[key], _ = printValue(value)
}

// SetDefault sets a properties map entry if not already set.
func SetDefault[V any](pm Map, key string, value V) {
	_, found := pm.m[key]
	if found {
		return
	}
	pm.m[key], _ = printValue(value)
}

// Has returns whether the given key exists.
func Has(pm Map, key string) bool {
	_, found := pm.m[key]
	return found
}

// Value returns the requested value, or fails if not found.
func Value[V any](pm Map, key string, def V) (V, error) {
	debugLogDefault(pm, key, def, false)
	str, found := pm.m[key]
	if !found {
		return def, fmt.Errorf("key %q is missing", key)
	}
	var ret V
	ret, err := parseValue[V](str)
	if err != nil {
		return def, fmt.Errorf("failed to parse key %q value %q: %v", key, str, err)
	}
	return ret, nil
}

// StringOr returns the string value, or the given value if not found.
//
// This is separate as this special case cannot ever fail.
func StringOr(pm Map, key string, def string) string {
	debugLogDefault(pm, key, def, true)
	str, found := pm.m[key]
	if !found {
		return def
	}
	return str
}

// ValueOr returns the requested value, or the given value if not found.
func ValueOr[V any](pm Map, key string, def V) (V, error) {
	debugLogDefault(pm, key, def, true)
	str, found := pm.m[key]
	if !found {
		return def, nil
	}
	var ret V
	ret, err := parseValue[V](str)
	if err != nil {
		return def, fmt.Errorf("failed to parse %q value %q: %v", key, str, err)
	}
	return ret, nil
}

func errToP(err error, errp *error) {
	if err != nil {
		if errp == nil {
			log.Fatalf("unhandled parse error: %v", err)
		} else if *errp == nil {
			*errp = err
		}
	}
}

// ValueP returns the requested value, or fails if not found.
func ValueP[V any](pm Map, key string, def V, errp *error) V {
	v, err := Value(pm, key, def)
	errToP(err, errp)
	return v
}

// ValueOrP returns the requested value, or the given value if not found.
func ValueOrP[V any](pm Map, key string, def V, errp *error) V {
	v, err := ValueOr(pm, key, def)
	errToP(err, errp)
	return v
}

// Interface functions so Maps behave just like map[string]string for cmp, encoding/json and hashstructure.
// This is needed for savegame compatibility.

func (pm Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(pm.m)
}

func (pm *Map) UnmarshalJSON(data []byte) error {
	pm.m = nil
	return json.Unmarshal(data, &pm.m)
}

func (pm Map) Hash() (uint64, error) {
	return hashstructure.Hash(pm.m, hashstructure.FormatV2, nil)
}

func (pm Map) Equal(other Map) bool {
	return cmp.Equal(pm.m, other.m)
}

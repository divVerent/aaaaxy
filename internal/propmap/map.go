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
	"fmt"
	"image/color"
	"time"

	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/palette"
)

// Map is a helper to easier deal with maps.
type Map map[string]string

func parseValue[V any](str string) (V, error) {
	var ret V
	var err error
	switch retP := (interface{})(&ret).(type) {
	case *string:
		*retP = str
	case *color.NRGBA:
		*retP, err = palette.Parse(str, "entity field")
	case *time.Duration:
		*retP, err = time.ParseDuration(str)
	default:
		_, err = fmt.Sscan(str, retP)
	}
	return ret, err
}

// Set sets a properties map entry.
func Set(pm Map, key string, value interface{}) {
	pm[key] = fmt.Sprint(value)
}

// SetDefault sets a properties map entry if not already set.
func SetDefault(pm Map, key string, value interface{}) {
	_, found := pm[key]
	if found {
		return
	}
	pm[key] = fmt.Sprint(value)
}

// Value returns the requested value, or fails if not found.
func Value[V any](pm Map, key string, def V) (V, error) {
	debugLogDefault(pm, key, def, false)
	str, found := pm[key]
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
	str, found := pm[key]
	if !found {
		return def
	}
	return str
}

// ValueOr returns the requested value, or the given value if not found.
func ValueOr[V any](pm Map, key string, def V) (V, error) {
	debugLogDefault(pm, key, def, true)
	str, found := pm[key]
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

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
	"io"
	"time"

	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/palette"
)

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
	case *bool, *int, *float64, *m.Delta, *m.Orientation, *m.Orientations, *TriState:
		_, err = fmt.Sscan(str, retP)
	default:
		log.Fatalf("missing support for type %T", ret)
	}
	return ret, err
}

func printValue[V any](v V) (ret, tmxType string) {
	switch vT := (interface{})(v).(type) {
	case string:
		return vT, "string"
	case color.NRGBA:
		return fmt.Sprintf("#%02x%02x%02x%02x", vT.A, vT.R, vT.G, vT.B), "color"
	case time.Duration:
		return fmt.Sprint(vT), "string"
	case bool:
		return fmt.Sprint(vT), "bool"
	case int:
		return fmt.Sprint(vT), "int"
	case float64, m.Delta, m.Orientation, m.Orientations, TriState:
		return fmt.Sprint(vT), "string"
	default:
		log.Fatalf("missing support for type %T", v)
		return "", ""
	}
}

type TriState struct {
	Active bool
	Value  bool
}

func (t TriState) String() string {
	if !t.Active {
		return ""
	}
	if t.Value {
		return "true"
	}
	return "false"
}

func (t *TriState) Scan(state fmt.ScanState, verb rune) error {
	token, err := state.Token(true, nil)
	if err == io.EOF {
		t.Active = false
		return nil
	}
	t.Active = true
	switch string(token) {
	case "true":
		t.Value = true
	case "false":
		t.Value = false
	default:
		return fmt.Errorf("unexpected token: %s", token)
	}
	return nil
}

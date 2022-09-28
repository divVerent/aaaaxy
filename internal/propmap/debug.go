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

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
)

var (
	debugLogEntityDefaults = flag.Bool("debug_log_entity_defaults", false, "log default values of entities (can be turned into objecttypes.xml)")
)

// DebugSetType enables logging for the given map and annotates logs with the given name.
func DebugSetType(pm Map, name string) {
	if !*debugLogEntityDefaults {
		return
	}
	pm["__type__"] = name
}

var linesPrinted = map[string]bool{}

func debugLogDefault[V any](pm Map, key string, def V, useDefault bool) {
	if !*debugLogEntityDefaults {
		return
	}
	switch key {
	// These keys don't go into .tmx <properties>, they always use special properties.
	case "name", "type":
		return
	}
	name, found := pm["__type__"]
	if !found {
		return
	}
	var typeName, defStr string
	defStr = fmt.Sprint(def)
	switch defT := (interface{})(def).(type) {
	case bool:
		typeName = "bool"
	case int:
		typeName = "int"
	case float64, string:
		typeName = "string"
	case color.NRGBA:
		typeName = "color"
		defStr = fmt.Sprintf("#%02x%02x%02x%02x", defT.A, defT.R, defT.G, defT.B)
	case m.Delta, m.Orientation, m.Orientations, time.Duration:
		typeName = "string"
	default:
		log.Fatalf("missing support for type %T", def)
	}
	if typeName != "string" || defStr != "" {
		defStr = fmt.Sprintf(" default=\"%s\"", defStr)
	}
	if !useDefault {
		defStr = ""
	}

	// To not make the log too large, print each such line only once.
	line := fmt.Sprintf("default for %q: <property name=\"%v\" type=\"%v\"%v/>",
		name, key, typeName, defStr)
	_, printed := linesPrinted[line]
	if printed {
		return
	}
	linesPrinted[line] = true
	log.Infof("%s", line)
}

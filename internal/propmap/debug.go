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

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	debugLogEntityDefaults = flag.Bool("debug_log_entity_defaults", false, "log default values of entities (can be turned into objecttypes.xml)")
)

// DebugSetType enables logging for the given map and annotates logs with the given name.
func DebugSetType(pm Map, name string) {
	if !*debugLogEntityDefaults {
		return
	}
	pm.m["__type__"] = name
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
	name, found := pm.m["__type__"]
	if !found {
		return
	}
	defStr, typeName := printValue(def)
	if typeName != "string" || defStr != "" {
		defStr = fmt.Sprintf(" default=\"%s\"", defStr)
	}
	if !useDefault {
		defStr = ""
	}

	// To not make the log too large, print each such line only once.
	line := fmt.Sprintf("entity %q property %q type %q default %q",
		name, key, typeName, defStr)
	_, printed := linesPrinted[line]
	if printed {
		return
	}
	linesPrinted[line] = true
	log.Infof("%s", line)
}

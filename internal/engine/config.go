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

package engine

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

// LoadConfig loads the current configuration.
func LoadConfig() (*flag.Config, error) {
	data, err := vfs.ReadState(vfs.Config, "config.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil // Not loading anything due to there being no config to load is OK.
		}
		return nil, err
	}
	config := &flag.Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// SaveConfig writes the current configuration to a file.
func SaveConfig() error {
	config := flag.Marshal()
	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}
	return vfs.WriteState(vfs.Config, "config.json", data)
}

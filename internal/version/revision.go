// Copyright 2021 Google LLC
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

package version

import (
	"io/ioutil"
	"strings"

	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

var revision string = "unknown"

func Init() error {
	fh, err := vfs.Load("generated", "version.txt")
	if err != nil {
		log.Errorf("cannot open out my version: %v", err)
		return err
	}
	defer fh.Close()
	revStr, err := ioutil.ReadAll(fh)
	if err != nil {
		log.Errorf("cannot read out my version: %v", err)
		return err
	}
	revision = strings.TrimSpace(string(revStr))
	log.Infof("AAAAXY %v", revision)
	return nil
}

func Revision() string {
	return revision
}

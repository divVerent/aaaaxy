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

package atexit

import (
	"fmt"
	"os"
)

var (
	toDelete []string
)

func Delete(path string) {
	toDelete = append(toDelete, path)
}

func Finish() {
	for _, path := range toDelete {
		err := os.RemoveAll(path)
		if err != nil && !os.IsNotExist(err) {
			// Can't use log here due to cyclic dependency otherwise.
			fmt.Fprintf(os.Stderr, "atexit: could not delete %v: %v", path, err)
		}
	}
}

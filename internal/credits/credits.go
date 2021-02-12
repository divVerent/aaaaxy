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

package credits

import (
	"bufio"
	"log"

	"github.com/divVerent/aaaaaa/internal/vfs"
)

var (
	Lines []string
)

func init() {
	credits, err := vfs.ReadDir("credits")
	if err != nil {
		log.Panicf("Could not list credits: %v", err)
	}
	log.Printf("Credits files: %v", credits)
	var lines []string
	for _, file := range credits {
		if lines != nil {
			lines = append(lines, "")
		}
		rd, err := vfs.Load("credits", file)
		if err != nil {
			log.Panicf("Could not load credits: %v", err)
		}
		defer rd.Close()
		scanner := bufio.NewScanner(rd)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err = scanner.Err(); err != nil {
			log.Panicf("Could not scan credits from %q: %v", file, err)
		}
	}
	log.Print(lines)
	Lines = lines
}

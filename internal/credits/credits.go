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
	"fmt"
	"regexp"

	"github.com/divVerent/aaaaxy/internal/vfs"
)

var (
	Lines    []string
	Licenses []string

	wordWrapRE = regexp.MustCompile(`(?:^\s*|\b)\S.{1,80}(?:\b|$)|^$`)
)

func concatenateLines(dir string) ([]string, error) {
	files, err := vfs.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("could not list files in %v: %w", dir, err)
	}
	var lines []string
	for _, file := range files {
		rd, err := vfs.Load(dir, file)
		if err != nil {
			return nil, fmt.Errorf("could not load file %v in %v: %w", file, dir, err)
		}
		defer rd.Close()
		scanner := bufio.NewScanner(rd)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			lines = append(lines, wordWrapRE.FindAllString(scanner.Text(), -1)...)
		}
		if err = scanner.Err(); err != nil {
			return nil, fmt.Errorf("could not scan file %v in %v: %w", file, dir, err)
		}
		// Make sure items are separated by empty lines.
		if len(lines) > 0 && lines[len(lines)-1] != "" {
			lines = append(lines, "")
		}
	}
	return lines, nil
}

func loadCredits() ([]string, error) {
	return concatenateLines("credits")
}

func loadLicenses() ([]string, error) {
	return concatenateLines("licenses")
}

func Precache() error {
	var err error
	Lines, err = loadCredits()
	if err != nil {
		return err
	}
	Licenses, err = loadLicenses()
	if err != nil {
		return err
	}
	return nil
}

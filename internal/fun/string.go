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

package fun

import (
	"bytes"
	"log"
	"text/template"
	"time"

	m "github.com/divVerent/aaaaxy/internal/math"
)

var formatTextMap = map[string]interface{}{
	"BigCity": func() string {
		// We are guessing a nearby large city with traffic problems by the user's time zone.
		// Rather inaccurate.
		// Sourced by searching for "<city> Road Rage" on Google and maximizing result count.
		// If no road rage found, any larger city will do.
		_, offset := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local).Zone()
		switch m.Div(offset, 3600) {
		case -12:
			return "Baker Island"
		case -11:
			return "Alofi"
		case -10:
			return "Honolulu"
		case -9:
			return "Anchorage"
		case -8:
			return "Vancouver"
		case -7:
			return "Edmonton"
		case -6:
			return "Chicago"
		case -5:
			return "New York"
		case -4:
			return "Halifax"
		case -3:
			return "Nuuk"
		case -2:
			return "Fernando de Noronha"
		case -1:
			return "Ponta Delgada"
		case 0:
			return "London"
		case 1:
			return "Berlin"
		case 2:
			return "Tel Aviv"
		case 3:
			return "Moscow"
		case 4:
			return "Tehran"
		case 5:
			return "Turkistan"
		case 6:
			return "Bishkek"
		case 7:
			return "Hanoi"
		case 8:
			return "Beijing"
		case 9:
			return "Tokyo"
		case 10:
			return "Sydney"
		case 11:
			return "Port Moresby"
		case 12:
			return "Auckland"
		case 13:
			return "Nuku'alofa"
		case 14:
			return "Kiritimati"
		default:
			// Boston is a great default for a bad place to drive.
			return "Boston"
		}
	},
}

// FormatText replaces placeholders in the given text.
func FormatText(s string) string {
	/*
		// Fast path if the template is trivial.
		if !strings.Contains(s, "{") {
			return s
		}
	*/
	tmpl := template.New("")
	tmpl.Funcs(formatTextMap)
	_, err := tmpl.Parse(s)
	if err != nil {
		log.Printf("Failed to parse text template: %v", s)
		return s
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	return buf.String()
}

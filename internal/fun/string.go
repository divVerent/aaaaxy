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
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"

	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/player_state"
)

// TryFormatText replaces placeholders in the given text.
func TryFormatText(ps *player_state.PlayerState, s string) (string, error) {
	/*
		// Fast path if the template is trivial.
		if !strings.Contains(s, "{") {
			return s
		}
	*/
	tmpl := template.New("")
	tmpl.Funcs(map[string]interface{}{
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
		"GameTime": func() (string, error) {
			if ps == nil {
				return "", fmt.Errorf("Cannot use {{GameTime}} in static elements.")
			}
			frames := ps.Frames()
			ss, ms := frames/60, (frames%60)*1000/60
			mm, ss := ss/60, ss%60
			hh, mm := mm/60, mm%60
			return fmt.Sprintf("%d:%02d:%02d.%03d", hh, mm, ss, ms), nil
		},
		"Abilities": func() (string, error) {
			if ps == nil {
				return "", fmt.Errorf("Cannot use {{Abilities}} in static elements.")
			}
			abilities := make([]string, 0, 4)
			if ps.HasAbility("carry") {
				abilities = append(abilities, "The Gloves")
			}
			if ps.HasAbility("control") {
				abilities = append(abilities, "The Remote")
			}
			if ps.HasAbility("push") {
				abilities = append(abilities, "The Coil")
			}
			if ps.HasAbility("stand") {
				abilities = append(abilities, "The Cleats")
			}
			switch len(abilities) {
			case 0:
				return "nothing", nil
			case 1:
				return abilities[0], nil
			default:
				return strings.Join(abilities[:len(abilities)-1], ", ") + " and " + abilities[len(abilities)-1], nil
			case 4:
				return "everything", nil
			}
		},
	})
	_, err := tmpl.Parse(s)
	if err != nil {
		return s, fmt.Errorf("Failed to parse text template: %v", s)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	return buf.String(), err
}

// FormatText replaces placeholders in the given text.
func FormatText(ps *player_state.PlayerState, s string) string {
	result, err := TryFormatText(ps, s)
	if err != nil {
		log.Printf("Failed to execute text template: %v", s)
		return s
	}
	return result
}

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
	"errors"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/divVerent/aaaaxy/internal/input"
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/log"
	m "github.com/divVerent/aaaaxy/internal/math"
	"github.com/divVerent/aaaaxy/internal/playerstate"
)

var (
	timeZoneHours int = 0xDEAD
)

// SetTimeZoneHours overrides the current time zone.
//
// Should be the time zone offset of Jan 1 this year, rounded DOWN (as in, floor) to whole hours.
//
// Used on Android to get the time zone from Java code.
func SetTimeZoneHours(h int) {
	timeZoneHours = h
}

func init() {
	_, offset := time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local).Zone()
	timeZoneHours = m.Div(offset, 3600)
}

// TryFormatText replaces placeholders in the given text.
func TryFormatText(ps *playerstate.PlayerState, s string) (string, error) {
	// Fast path if the template is trivial.
	if !strings.Contains(s, "{{") {
		return s, nil
	}
	tmpl := template.New("")
	tmpl.Funcs(map[string]interface{}{
		"Lang": func() string {
			return string(locale.Active)
		},
		"BR": func() string {
			return "\n"
		},
		"Year": func() string {
			return time.Now().Format("2006")
		},
		"BigCity": func() string {
			// We are guessing a nearby large city with traffic problems by the user's time zone.
			// Rather inaccurate.
			// Sourced by searching for "<city> Road Rage" on Google and maximizing result count.
			// If no road rage found, any larger city will do.
			switch timeZoneHours {
			case -12:
				return locale.G.Get("Baker Island")
			case -11:
				return locale.G.Get("Alofi")
			case -10:
				return locale.G.Get("Honolulu")
			case -9:
				return locale.G.Get("Anchorage")
			case -8:
				return locale.G.Get("San Francisco")
			case -7:
				return locale.G.Get("Edmonton")
			case -6:
				return locale.G.Get("Chicago")
			case -5:
				return locale.G.Get("New York")
			case -4:
				return locale.G.Get("Halifax")
			case -3:
				return locale.G.Get("São Paulo")
			case -2:
				return locale.G.Get("Fernando de Noronha")
			case -1:
				return locale.G.Get("Ponta Delgada")
			case 0:
				return locale.G.Get("London")
			case 1:
				return locale.G.Get("Berlin")
			case 2:
				return locale.G.Get("Tel Aviv")
			case 3:
				return locale.G.Get("Istanbul")
			case 4:
				return locale.G.Get("Tehran")
			case 5:
				return locale.G.Get("Turkistan")
			case 6:
				return locale.G.Get("Bishkek")
			case 7:
				return locale.G.Get("Hanoi")
			case 8:
				return locale.G.Get("Beijing")
			case 9:
				return locale.G.Get("Tokyo")
			case 10:
				return locale.G.Get("Sydney")
			case 11:
				return locale.G.Get("Port Moresby")
			case 12:
				return locale.G.Get("Auckland")
			case 13:
				return locale.G.Get("Nuku'alofa")
			case 14:
				return locale.G.Get("Kiritimati")
			default:
				// Boston is a great default for a bad place to drive.
				return locale.G.Get("Boston")
			}
		},
		"GameTime": func() (string, error) {
			if ps == nil {
				return "", errors.New("cannot use {{GameTime}} in static elements")
			}
			frames := ps.Frames()
			ss, ms := frames/60, (frames%60)*1000/60
			mm, ss := ss/60, ss%60
			hh, mm := mm/60, mm%60
			return locale.G.Get("%d:%02d:%02d.%03d", hh, mm, ss, ms), nil
		},
		"Score": func() (string, error) {
			if ps == nil {
				return "", errors.New("cannot use {{Score}} in static elements")
			}
			return fmt.Sprintf("%d", ps.Score()), nil
		},
		"SpeedrunCategoriesShort": func() (string, error) {
			if ps == nil {
				return "", errors.New("cannot use {{SpeedrunsShort}} in static elements")
			}
			return ps.SpeedrunCategories().DescribeShort(), nil
		},
		"Abilities": func() (string, error) {
			if ps == nil {
				return "", errors.New("cannot use {{Abilities}} in static elements")
			}
			abilities := make([]string, 0, 4)
			if ps.HasAbility("carry") {
				abilities = append(abilities, locale.G.Get("The Gloves"))
			}
			if ps.HasAbility("control") {
				abilities = append(abilities, locale.G.Get("The Remote"))
			}
			if ps.HasAbility("push") {
				abilities = append(abilities, locale.G.Get("The Coil"))
			}
			if ps.HasAbility("stand") {
				abilities = append(abilities, locale.G.Get("The Cleats"))
			}
			switch len(abilities) {
			case 0:
				return locale.G.Get("nothing"), nil
			case 1:
				return abilities[0], nil
			default:
				return strings.Join(abilities[:len(abilities)-1], locale.G.Get(", ")) + locale.G.Get(" and ") + abilities[len(abilities)-1], nil
			case 4:
				return locale.G.Get("everything"), nil
			}
		},
		"ExitButton": func() (string, error) {
			if ps == nil {
				return "", errors.New("cannot use {{ExitButton}} in static elements")
			}
			switch input.ExitButton() {
			case input.Start:
				return locale.G.Get("Start"), nil
			case input.Back:
				return locale.G.Get("Back"), nil
			case input.Escape:
				return locale.G.Get("Escape"), nil
			default: // case input.Backspace:
				return locale.G.Get("Backspace"), nil
			}
		},
		"ActionButton": func() (string, error) {
			if ps == nil {
				return "", errors.New("cannot use {{ActionButton}} in static elements")
			}
			switch input.ActionButton() {
			case input.BX:
				return locale.G.Get("B/X"), nil
			case input.Elsewhere:
				return locale.G.Get("elsewhere"), nil
			case input.B:
				return locale.G.Get("B"), nil
			case input.CtrlShift:
				return locale.G.Get("Ctrl/Shift"), nil
			case input.Z:
				return locale.G.Get("Z"), nil
			case input.ShiftETab:
				return locale.G.Get("Shift/E/Tab"), nil
			case input.EnterShift:
				return locale.G.Get("Enter/Shift"), nil
			default: // case input.Shift:
				return locale.G.Get("Shift"), nil
			}
		},
		"SpeedrunCategories": func() (string, error) {
			if ps == nil {
				return "", errors.New("cannot use {{SpeedrunCategories}} in static elements")
			}
			categories, _ := ps.SpeedrunCategories().Describe()
			return categories, nil
		},
		"SpeedrunTryNext": func() (string, error) {
			if ps == nil {
				return "", errors.New("cannot use {{SpeedrunTryNext}} in static elements")
			}
			_, tryNext := ps.SpeedrunCategories().Describe()
			return tryNext, nil
		},
	})
	_, err := tmpl.Parse(s)
	if err != nil {
		return s, fmt.Errorf("failed to parse text template: %v", s)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	return buf.String(), err
}

// FormatText replaces placeholders in the given text.
func FormatText(ps *playerstate.PlayerState, s string) string {
	result, err := TryFormatText(ps, s)
	if err != nil {
		log.Warningf("failed to execute text template: %v: %v", s, err)
		return s
	}
	return result
}

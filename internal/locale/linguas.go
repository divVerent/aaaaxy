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

package locale

import (
	"sort"
	"strings"
)

// Lingua identifies a language.
type Lingua string

// Linguas is the set of current languages that are translated "enough".
var Linguas = map[Lingua]struct{}{}

// Name returns the human readable name of a language.
func (l Lingua) Name() string {
	switch l {
	case "":
		return "English"
	case "ar":
		return "العربية"
	case "ar-EG":
		return "العربية (مصر)"
	case "be":
		return "Беларуская"
	case "be@tarask":
		return "Biełaruskaja"
	case "de":
		return "Deutsch"
	case "de-CH":
		return "Deutsch (Schweiz)"
	case "he":
		return "עִבְרִית"
	case "la":
		return "Latina"
	case "pt":
		return "Português"
	case "pt-BR":
		return "Português (Brasil)"
	case "uk":
		return "Українська"
	case "zh-Hans":
		return "简体中文"
	default:
		return string(l)
	}
}

func (l Lingua) SortKey() string {
	switch l {
	// Sort both Belarusian variants together.
	case "be@tarask":
		return "Беларуская (Łacinka)"
	default:
		return l.Name()
	}
}

func (l Lingua) Font() string {
	switch l {
	case "ar", "ar-EG":
		return "bitmapfont"
	case "he", "zh-Hans":
		return "unifont"
	default:
		return "gofont"
	}
}

func (l Lingua) AuditHeight() bool {
	switch l {
	// These languages use unifont which is a tad bit too high.
	// Accept anyway, as some minor vertical overflow is not noticeable ingame.
	case "zh-Hans":
		return false
	default:
		return true
	}
}

// Directory returns the directory containing the language.
func (l Lingua) Directory() string {
	switch l {
	// Handle groups.
	case "de-CH":
		return "de"
	case "pt-BR":
		return "pt"
	// This one doesn't get an underscore.
	case "zh-Hans":
		return "zh-Hans"
	// Otherwise Transifex uses underscores, but x/text/language uses dashes.
	default:
		return strings.ReplaceAll(string(l), "-", "_")
	}
}

// GroupMembers returns all additional members of a language group.
// Groups use the same file, but have {{if eq Lang ...}} template commands for minor differences.
func (l Lingua) GroupMembers() []Lingua {
	switch l {
	case "de":
		return []Lingua{"de-CH"}
	case "pt":
		return []Lingua{"pt-BR"}
	default:
		return nil
	}
}

func (l Lingua) Canonical() Lingua {
	// Handle aliases.
	// Aliases are different names for the same language.
	switch l {
	case "iw": // Deprecated common alias for Hebrew. Java still uses it, and thus Android likely, too.
		return "he"
	case "zh-CHS", "zh-CN", "zh-SG": // Language specific Chinese aliases.
		return "zh-Hans"
	default:
		return l
	}
}

// Shape performs glyph shaping on a given string.
func (l Lingua) Shape(s string) string {
	switch l {
	case "ar", "ar-EG", "he":
		return l.shapeArabic(s)
	default:
		return s
	}
}

// UseEbitenText returns whether using ebiten/text for font drawing is safe.
func (l Lingua) UseEbitenText() bool {
	return true
}

// LinguasSorted returns the languages sorted by humanly expected ordering.
func LinguasSorted() []Lingua {
	ret := []Lingua{""}
	for l := range Linguas {
		ret = append(ret, l)
	}
	sort.Slice(ret, func(a, b int) bool {
		return ret[a].SortKey() < ret[b].SortKey()
	})
	return ret
}

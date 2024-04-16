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

	"github.com/divVerent/aaaaxy/internal/log"
)

// Lingua identifies a language.
type Lingua string

// Builtin is a reserved language name for the included language.
const Builtin Lingua = ""

// UserProvided is a reserved language name for language packs.
const UserProvided Lingua = "."

// Linguas is the set of current languages that are translated "enough".
var Linguas = map[Lingua]struct{}{}

// name returns the human readable name of a language.
func (l Lingua) name() string {
	switch l {
	case Builtin:
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
	case "ja":
		return "日本語"
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
	case "zh-Hant":
		return "繁體中文"
	case UserProvided:
		return "user provided"
	default:
		return string(l)
	}
}

// Name returns the human readable name of a language for the current font.
func (l Lingua) Name() string {
	if l != Active {
		switch l {
		case "ar":
			return "Arabic"
		case "ar-EG":
			return "Arabic (Egypt)"
		case "he":
			return "Hebrew"
		case "ja":
			return "Japanese"
		case "zh-Hans":
			return "Chinese (Simplified)"
		case "zh-Hant":
			return "Chinese (Traditional)"
		}
	}
	return l.name()
}

func (l Lingua) SortKey() string {
	switch l {
	case "be@tarask":
		// Sort both Belarusian variants together.
		return "Беларуская (Łacinka)"
	case UserProvided:
		// User provided comes at the far left.
		return ""
	default:
		return l.name()
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
	case "zh-Hant":
		return "zh-Hant"
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
	case "be@latin", "be-Latn": // Actually more correct, but using be@tarask as it is the closest match on Transifex.
		return "be@tarask"
	case "zh-CHS", "zh-CN", "zh-SG": // Language specific Chinese aliases.
		return "zh-Hans"
	case "zh-CHT", "zh-HK", "zh-MO", "zh-TW": // Language specific Chinese aliases.
		return "zh-Hant"
	default:
		return l
	}
}

// LinguasSorted returns the languages sorted by humanly expected ordering.
func LinguasSorted() []Lingua {
	ret := []Lingua{Builtin}
	for l := range Linguas {
		ret = append(ret, l)
	}
	sort.Slice(ret, func(a, b int) bool {
		return ret[a].SortKey() < ret[b].SortKey()
	})
	return ret
}

// ActiveFont returns the font this locale uses.
//
// This is only accessible for the active locale so it can later be defined by the language file itself.
func ActiveFont() string {
	po := G // Workaround for xgotext otherwise not finding the call.
	setting := po.Get("_locale_info:font")
	if setting == "_locale_info:font" {
		return "gofont"
	}
	switch setting {
	case "_locale_info:font", "default":
		return "gofont"
	case "bitmapfont", "gofont", "unifont":
		return setting
	default:
		log.Fatalf("Invalid value of _locale_info:font: got %q, want default, bitmapfont, gofont or unifont", setting)
		return "gofont"
	}
}

type VerticalTextPreference int

const (
	NeverPreferVerticalText VerticalTextPreference = iota
	DefaultPreferVerticalText
	AlwaysPreferVerticalText
)

// ActivePrefersVerticalText returns whether this locale prefers vertical text.
//
// This is only accessible for the active locale so it can later be defined by the language file itself.
func ActivePrefersVerticalText() VerticalTextPreference {
	po := G // Workaround for xgotext otherwise not finding the call.
	setting := po.Get("_locale_info:prefers_vertical_text")
	switch setting {
	case "_locale_info:prefers_vertical_text", "default":
		return DefaultPreferVerticalText
	case "true":
		return AlwaysPreferVerticalText
	case "false":
		return NeverPreferVerticalText
	default:
		log.Fatalf("Invalid value of _locale_info:prefers_vertical_text: got %q, want default, true or false", setting)
		return DefaultPreferVerticalText
	}
}

// ActiveFitsHeight returns whether height auditing will be performed.
//
// This is only accessible for the active locale so it can later be defined by the language file itself.
func ActiveFitsHeight() bool {
	po := G // Workaround for xgotext otherwise not finding the call.
	setting := po.Get("_locale_info:fits_height")
	switch setting {
	case "_locale_info:fits_height":
		return true
	case "true":
		return true
	case "false":
		return false
	default:
		log.Fatalf("Invalid value of _locale_info:fits_height: got %q, want true or false", setting)
		return true
	}
}

// ActiveUsesArabicShaping returns whether Arabic shaping will be performed.
//
// This is only accessible for the active locale so it can later be defined by the language file itself.
func ActiveUsesArabicShaping() bool {
	po := G // Workaround for xgotext otherwise not finding the call.
	setting := po.Get("_locale_info:uses_arabic_shaping")
	switch setting {
	case "_locale_info:uses_arabic_shaping":
		return false
	case "true":
		return true
	case "false":
		return false
	default:
		log.Fatalf("Invalid value of _locale_info:uses_arabic_shaping: got %q, want true or false", setting)
		return false
	}
}

// ActiveShape performs glyph shaping on a given string.
//
// This is only accessible for the active locale so it can later be defined by the language file itself.
func ActiveShape(s string) string {
	switch {
	case ActiveUsesArabicShaping():
		return Active.shapeArabic(s)
	default:
		return s
	}
}

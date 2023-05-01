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
	"github.com/jeandeaual/go-locale"
	"golang.org/x/text/language"

	"github.com/divVerent/aaaaxy/internal/log"
)

// Current returns the preferred system locale, intersected with available locales.
var Current Lingua

// InitCurrent identifies the system locale. Requires Linguas to have been set.
func InitCurrent() {
	Current = ""
	locs, err := locale.GetLocales()
	if err != nil {
		log.Errorf("could not detect current locales: %v", err)
		return
	}
	for _, loc := range locs {
		lang, err := language.Parse(loc)
		if err != nil {
			continue
		}
		for lang != language.Und {
			lingua := Lingua(lang.String()).Canonical()
			if lingua == "en" {
				log.Infof("detected language %s (not translating)", Lingua("").Name())
				return
			}
			if _, found := Linguas[lingua]; found {
				log.Infof("detected language %s", lingua.Name())
				Current = lingua
				return
			}
			lang = lang.Parent()
		}
	}
	log.Infof("detected no supported language (not translating)")
}

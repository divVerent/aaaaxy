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

package menu

import (
	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/locale"
)

type languageSetting struct {
	linguas []locale.Lingua
	current int
}

func (l *languageSetting) init() {
	l.linguas = append([]locale.Lingua{"auto"}, locale.LinguasSorted()...)
	for i, lingua := range l.linguas {
		if flag.Get[string]("language") == string(lingua) {
			l.current = i
		}
	}
}

func (l *languageSetting) name() string {
	lingua := l.linguas[l.current]
	if lingua == "auto" {
		return locale.G.Get("Auto (%s)", locale.Current.Name())
	}
	return lingua.Name()
}

func (l *languageSetting) apply(m *Controller) error {
	lingua := l.linguas[l.current]
	if string(lingua) == flag.Get[string]("language") {
		return nil
	}
	flag.Set("language", string(lingua))
	return m.ReinitLevelNextFrame()
}

func (l *languageSetting) toggle(m *Controller, delta int) error {
	switch delta {
	case 0:
		l.current++
		if l.current >= len(l.linguas) {
			l.current = 0
		}
	case -1:
		if l.current > 0 {
			l.current--
		}
	case +1:
		l.current++
		if l.current >= len(l.linguas) {
			l.current--
		}
	}
	l.apply(m)
	return nil
}

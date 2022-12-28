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

package initlocale

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/locale"
	"github.com/divVerent/aaaaxy/internal/log"
	"github.com/divVerent/aaaaxy/internal/vfs"
)

var (
	language = flag.String("language", "auto", "language to translate the game into; if set to 'auto', it will be detected using the system locale; set to '' to not translate")
)

func initLinguas() {
	data, err := vfs.Load("locales", "LINGUAS")
	if err != nil {
		log.Errorf("could not open LINGUAS file: %v", err)
		return
	}
	defer data.Close()
	buf, err := io.ReadAll(data)
	for _, line := range bytes.Split(buf, []byte{'\n'}) {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		locale.Linguas[locale.Lingua(line)] = struct{}{}
	}
}

func initLocaleDomain(lang locale.Lingua, l locale.Type, domain string) {
	if lang == "" {
		return
	}
	data, err := vfs.Load(fmt.Sprintf("locales/%s", strings.ReplaceAll(string(lang), "-", "_")), fmt.Sprintf("%s.po", domain))
	if err != nil {
		log.Errorf("could not open %s translation for language %s: %v", domain, lang.Name(), err)
		return
	}
	defer data.Close()
	buf, err := io.ReadAll(data)
	if err != nil {
		log.Errorf("could not read %s translation for language %s: %v", domain, lang.Name(), err)
		return
	}
	l.Parse(buf)
	log.Infof("%s translated to language %s", domain, lang.Name())
}

func Init() error {
	initLinguas()
	locale.InitCurrent()
	lang := locale.Lingua(*language)
	if lang == "auto" {
		lang = locale.Current
	}
	initLocaleDomain(lang, locale.G, "game")
	initLocaleDomain(lang, locale.L, "level")
	locale.Active = lang
	return nil
}

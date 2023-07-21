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
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	debugCheckTranslations = flag.Bool("debug_check_translations", false, "fail startup if a translation contains a format string mismatch or a too big text box")
)

var (
	formatRE = regexp.MustCompile(`({{[^}]*}})|%(?:\[(\d+)\])?([-+# 0-9.]*[a-zA-Z%])`)
	badRE    = regexp.MustCompile(` {{BR}}|{{BR}} |^ | $|^\n|\n$`)
)

func formats(s string) map[string]int {
	index := 1
	out := map[string]int{}
	for _, match := range formatRE.FindAllStringSubmatch(s, -1) {
		if match[1] != "" {
			if match[1] == "{{BR}}" {
				// Newlines are allowed to vary.
				continue
			}
			if strings.HasPrefix(match[1], "{{if ") || strings.HasPrefix(match[1], "{{else") || match[1] == "{{end}}" {
				// Conditionals are allowed to vary.
				continue
			}
			out[match[0]]++
		} else if match[2] != "" {
			// Has an explicit index. No change needed.
			out[match[0]]++
			i, err := strconv.Atoi(match[2])
			if err != nil {
				log.Fatalf("failed to parse format string %q: %v", s, err)
			}
			index = i + 1
		} else {
			out[fmt.Sprintf("%%[%d]%s", index, match[3])]++
			index++
		}
	}
	return out
}

func auditPo(po Type) error {
	for k, vs := range po.GetDomain().GetTranslations() {
		if k == "" {
			// Not a real string, just a header.
			continue
		}
		kbads := map[string]struct{}{}
		for _, kbad := range badRE.FindAllString(k, -1) {
			kbads[kbad] = struct{}{}
		}
		kf := formats(k)
		for _, v := range vs.Trs {
			vf := formats(v)
			if !reflect.DeepEqual(kf, vf) {
				err := fmt.Errorf("translation format string mismatch: %q (%v) -> %q (%v)", k, kf, v, vf)
				if *debugCheckTranslations {
					return err
				} else {
					log.Errorf("%v", err)
				}
			}
			for _, vbad := range badRE.FindAllString(v, -1) {
				if _, found := kbads[vbad]; found {
					// Same as original - probably OK then.
					continue
				}
				err := fmt.Errorf("translation contains bad substring: %q -> %q (%q), matched by regexp %v", k, v, vbad, badRE)
				if *debugCheckTranslations {
					return err
				} else {
					log.Errorf("%v", err)
				}
			}
		}
	}
	return nil
}

func Audit() error {
	err := auditPo(G)
	if err != nil {
		return err
	}
	return auditPo(L)
}

func Errorf(format string, args ...interface{}) {
	if *debugCheckTranslations {
		log.Fatalf(format, args...)
	} else {
		log.Errorf(format, args...)
	}
}

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

	"github.com/divVerent/aaaaxy/internal/flag"
	"github.com/divVerent/aaaaxy/internal/log"
)

var (
	debugCheckTranslations = flag.Bool("debug_check_translations", false, "fail startup if a translation contains a format string mismatch")
)

var (
	formatRE = regexp.MustCompile(`({{[^}]*}})|%(?:(\d+)\$)?([^a-z]*[a-z])`)
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
			out[match[0]]++
		} else if match[2] != "" {
			// Go doesn't support %1$s syntax yet, sadly.
			out["UNSUPPORTED:"+match[0]]++
		} else {
			out[fmt.Sprintf("%%%d$%s", index, match[3])]++
		}
	}
	return out
}

func auditPo(po Type) error {
	for k, vs := range po.GetDomain().GetTranslations() {
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

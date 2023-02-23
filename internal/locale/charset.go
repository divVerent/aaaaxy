// Copyright 2023 Google LLC
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

	"github.com/divVerent/aaaaxy/internal/log"
)

func CharSet(base string, baseWeight, maxCount int) []rune {
	weights := map[rune]int{}
	for _, r := range base {
		weights[r] = baseWeight
	}
	for _, po := range []Type{G, L} {
		for k, vs := range po.GetDomain().GetTranslations() {
			if k == "" {
				// Not a real string, just a header.
				continue
			}
			kbads := map[string]struct{}{}
			for _, kbad := range badRE.FindAllString(k, -1) {
				kbads[kbad] = struct{}{}
			}
			for _, v := range vs.Trs {
				for _, r := range formatRE.ReplaceAllString(v, "") {
					if r < ' ' {
						continue
					}
					weights[r]++
				}
			}
		}
	}
	totalWeight := 0
	var out []rune
	for r, w := range weights {
		out = append(out, r)
		totalWeight += w
	}
	sort.Slice(out, func(i, j int) bool {
		// Prefer those with higher weight.
		if d := weights[out[i]] - weights[out[j]]; d != 0 {
			return d > 0
		}
		// At equal weight, prefer those first in ASCII. They're punctuation and digits.
		return out[i] < out[j]
	})
	if len(out) > maxCount {
		out = out[:maxCount]
	}
	chosenWeight := 0
	for _, r := range out {
		chosenWeight += weights[r]
	}
	log.Infof("font pinning coverage: %v%%", float64(chosenWeight)*100.0/float64(totalWeight))
	sort.Slice(out, func(i, j int) bool {
		return out[i] < out[j]
	})
	return out
}

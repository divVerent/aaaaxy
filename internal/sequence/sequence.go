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

package sequence

type Sequence struct {
	want  []int
	got   []int
	shift int
}

func New(want ...int) *Sequence {
	return &Sequence{
		want:  want,
		got:   make([]int, len(want)),
		shift: 0,
	}
}

func (s *Sequence) Add(what int) {
	s.got[s.shift] = what
	s.shift = (s.shift + 1) % len(s.got)
}

func (s *Sequence) Reset() {
	for i := range s.got {
		s.got[i] = 0
	}
}

func (s *Sequence) Match() bool {
	for i, w := range s.want {
		if s.got[(s.shift+i)%len(s.got)] != w {
			return false
		}
	}
	return true
}

#!/bin/sh
# Copyright 2021 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Plays some music using the shepard tone samples.

melody='
	9,0.5
	12,0.75
	14,0.25
	16,0.5
	9,0.5
	12,0.75
	14,0.25
	16,0.5
	9,0.5
	12,0.75
	14,0.25
	16,0.5
	17,0.5
	16,0.5
	14,0.5
	12,0.5
	9,0.5
	12,0.75
	14,0.25
	16,0.5
	9,0.5
	12,0.75
	14,0.25
	16,0.5
	9,0.5
	12,0.75
	14,0.25
	16,0.5
	14,0.5
	12,0.5
	11,0.5
	9,2
	17,0.5
	16,0.5
	14,0.5
	12,0.5
	11,2
	19,0.5
	17,0.5
	16,0.5
	14,0.5
	12,0.5
	14,0.5
	16,0.5
	17,0.5
	16,0.5
'
offset=5

transpose=0
for i in $(seq 1 12); do
	for note in $melody; do
		duration=${note#*,}
		duration=$(echo "$duration * 0.5" | bc -l)
		note=${note%%,*}
		note=$(((note + transpose) % 12))
		set -- "$@" --{ -endpos "$duration" `printf ../shepard_%02d.ogg $note` --}
	done
	transpose=$((transpose + offset))
done
mpv "$@"

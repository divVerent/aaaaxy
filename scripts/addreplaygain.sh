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

find assets third_party -path '*/music/*.ogg' | while read -r file; do
	silence=$(sox "$file" -t dat - | awk 'BEGIN { n = 0; } /^\s*[0-9]/ { if ($2 != 0 || $3 != 0) { nextfile; } ++n; } END { print(n); }')
	gain=$(loudgain --maxtpl=0 --pregain=-5 --output-new "$file" | tail -n 1 | cut -f 9)
	gain=$(echo "e(l(10) * ${gain% dB} / 20)" | bc -l)
	echo "$file -> $gain, $silence"
	data=$(cat "$file.json" || echo "{}")
	echo "$data" | jq ". + {play_start: $silence, replay_gain: $gain}" > "$file.json"
done

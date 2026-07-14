#!/bin/sh
# Copyright 2026 Google LLC
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

# Generates a soundfont to play with Shepard Tone.


# -17.5dB chosen so Rush E barely doesn't clip.

cat <<EOF
<global>
volume=-17.5
ampeg_hold=0.001
ampeg_release=0.5
EOF

for n in $(seq 0 127); do
  i=$((n % 12))
  i0=$(printf %02d $i)
  cat <<EOF
<region>
trigger=attack
lokey=$n
hikey=$n
lovel=0
hivel=127
sample=$(realpath "../shepard_$i0.ogg")
pitch_keycenter=$n
EOF
done

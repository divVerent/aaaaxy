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

# Generates "Shepard tone" audio files.

center=0
basevol=0.5
falloff=0.98
attack=0.01
decay=2.95
len=3
reverbdecay=1.5
for basenote in $(seq 0 11); do
	set --
	vols=
	ch=0
	for i in $(seq -32 32); do
		note=$((basenote + i * 12))
		#if [ $note -lt -48 ] || [ $note -gt 39 ]; then
		#	continue
		#fi
		delta=$((note - center))
		delta=${delta#-}
		vol=$(echo "$basevol * $falloff ^ $delta" | bc -l)
		ch=$((ch+1))
		vols=$vols${ch}v$vol,
		set -- "$@" "$len" sine "%$note"
	done
	name=$(printf ../shepard_%02d $basenote)
	sox -n -r 44.1k -e signed -b 16 -c 2 \
	  "$name.wav" \
	  synth "$@" \
	  remix -m "${vols%,}" \
	  fade l "$attack" "$len" "$decay" \
	  reverb 95 50 100 100 \
	  fade q 0 "$len" "$reverbdecay"
	oggenc -q3 -o "$name.ogg" "$name.wav"
done

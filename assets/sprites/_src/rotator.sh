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

steps=30
size=64
minwidth=4

images="can_carry can_control can_push can_stand can_switch_level"

for i in $(seq 0 $((steps/2))); do
	width=$(echo "c(8 * a(1) * $i / $steps) * $size" | bc -l)
	width=${width%.*}
	case "$width" in
		-*)
			width=${width#-}
			flip=-flop
			;;
		*)
			flip=
			;;
	esac
	case "$width" in
		''|0)
			width=1
			;;
	esac
	if [ $width -lt $minwidth ]; then
		width=$minwidth
	fi
	echo "$i -> $width $flip"
	for img in $images; do
		convert "$img.png" \
			-filter Point \
			-resize "${width}x${size}!" \
			$flip \
			-gravity center \
			-background none \
			-extent "${size}x${size}" \
			\( \
				+clone \
				+level 0,0 \
				-channel A \
				-morphology Dilate:1 Square \
			\) \
			+swap \
			-composite \
			../"${img}_default_${i}.png"
	done
done

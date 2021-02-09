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


colors='
	8#555555
	9#5555ff
	a#55ff55
	b#55ffff
	c#ff5555
	d#ff55ff
	e#ffff55
	f#ffffff
'

for kc1 in $colors; do
	k1=${kc1%%#*}
	c1=${kc1#?}
	for kc2 in $colors; do
		k2=${kc2%%#*}
		c2=${kc2#?}
		convert bgtransition.png \
			-fill "$c1" -opaque "#00fe00" \
			-fill "$c2" -opaque "#ff00fe" \
			../bg_"$k1$k2"_v.png
		convert bgtransition.png \
			-fill "$c1" -opaque "#00fe00" \
			-fill "$c2" -opaque "#ff00fe" \
			-rotate 270 \
			../bg_"$k1$k2"_h.png
		if [ x"$k1" = x"8" ]; then
			for img in l m nl nr r; do
				convert train_"$img".png \
					-fill "$c1" -opaque "#00fe00" \
					-fill "$c2" -opaque "#ff00fe" \
					../train_"$k1$k2"_"$img".png
			done
		fi
	done
	convert -size 16x16 xc:"$c1" ../bg_"$k1".png
	if [ x"$k1" = x"8" ]; then
		convert track.png \
			-fill "$c1" -opaque "#00fe00" \
			../track_"$k1"_v.png
	fi
done

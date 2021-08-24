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

convert -size 128x128 gradient:black-white -colorspace sRGB -depth 8 -alpha copy ../gradient_top_bottom.png
convert -size 128x128 gradient:black-white -rotate 270 -colorspace sRGB -depth 8 -alpha copy ../gradient_left_right.png
convert -size 128x128 radial-gradient:white-black -colorspace sRGB -depth 8 -alpha copy ../gradient_outside_inside.png

for i in $(seq 0 7); do
	convert \
		\( -size 128x128 xc:black -colorspace gray +noise Random -gaussian-blur 2x2 -level 40%,60% \) \
		\( ../gradient_outside_inside.png -alpha extract \) \
		-compose CopyOpacity -composite \
		../magic_idle_$i.png
done

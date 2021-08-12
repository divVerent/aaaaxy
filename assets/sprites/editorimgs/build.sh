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

convert ../arrow32d.png -rotate 180 arrow32_nw.png
convert ../arrow32d.png -rotate 270 arrow32_ne.png
convert ../arrow32d.png -rotate 90 arrow32_sw.png
convert ../arrow32s.png -rotate 180 arrow32_w.png
convert ../arrow32s.png -rotate 270 arrow32_n.png
convert ../arrow32s.png -rotate 90 arrow32_s.png
convert ../bullet8d_idle_0.png -rotate 180 bullet8_nw.png
convert ../bullet8d_idle_0.png -rotate 270 bullet8_ne.png
convert ../bullet8d_idle_0.png -rotate 90 bullet8_sw.png
convert ../bullet8s_idle_0.png -rotate 180 bullet8_w.png
convert ../bullet8s_idle_0.png -rotate 270 bullet8_n.png
convert ../bullet8s_idle_0.png -rotate 90 bullet8_s.png
convert ../car_idle_0.png -geometry 16x32 car.png
convert ../forcefield.png -crop 16x16+0+0 +repage -rotate 90 forcefield_v.png
convert ../forcefield.png -crop 16x16+0+0 +repage forcefield_h.png
convert ../gradient_left_right.png -alpha off -geometry 32x32 gradient_left_right.png
convert ../gradient_outside_inside.png -alpha off -geometry 32x32 gradient_outside_inside.png
convert ../gradient_top_bottom.png -alpha off -geometry 32x32 gradient_top_bottom.png
convert ../movingdoor.png -geometry 16x32 movingdoor.png
convert ../oneway_idle_0.png -rotate 180 oneway_w.png
convert ../oneway_idle_0.png -rotate 270 oneway_n.png
convert ../oneway_idle_0.png -rotate 90 oneway_s.png
convert ../questionblock.png -colorspace gray kaizoblock.png
convert ../spike.png -rotate 180 spike_w.png
convert ../spike.png -rotate 270 spike_n.png
convert ../spike.png -rotate 90 spike_s.png
convert ../v.png -rotate 180 v_w.png
convert ../v.png -rotate 270 v_n.png
convert ../v.png -rotate 90 v_s.png

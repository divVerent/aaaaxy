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

convert -size 128x128 gradient:transparent-white -colorspace sRGB -depth 8 ../gradient_top_bottom.png
convert -size 128x128 gradient:transparent-white -rotate 270 -colorspace sRGB -depth 8 ../gradient_left_right.png
convert -size 128x128 radial-gradient:white-transparent -colorspace sRGB -depth 8 ../gradient_outside_inside.png

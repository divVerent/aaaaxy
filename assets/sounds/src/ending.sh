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

# Generates the ending sound effect.

sox ../questionblock.ogg -c 2 youre_winner.wav \
	pad 0 5 \
	gain -3 \
	reverb 75 50 100 100 \
	fade q 0 5 0.5
oggenc -q3 -o ../youre_winner.ogg youre_winner.wav

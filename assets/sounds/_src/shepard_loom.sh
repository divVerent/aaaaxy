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

# Generates "Loom" sound effects.

# Better for shepard tone:
# Major: C G7 G = C E G B D F G F D B G E C
# Minor: d Bb d = d f a bb d bb a f d

sox \
  --combine merge \
  ../shepard_{03,07,10,01,05,08,10,08,05,01,10,07,03,07,10,01,05,08,10,08,05,01,10,07,03}.wav \
  ../loom.wav \
  delay $(seq 1 0.2 5.8 | sed -e 's,.*,& &,g') \
  remix - - \
  gain -n -1

sox \
  --combine merge \
  -c 2 \
  ../shepard_{05,08,00,01,05,01,00,08,05,08,00,01,05,01,00,08,05}.wav \
  ../loom_minor.wav \
  delay $(seq 1 0.3 5.8 | sed -e 's,.*,& &,g') \
  remix - - \
  gain -n -1

oggenc -q3 -o ../loom.ogg ../loom.wav
oggenc -q3 -o ../loom_minor.ogg ../loom_minor.wav

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

sox \
  --combine merge \
  ../shepard_{00,04,07,00,04,07,00,07,04,00,07,04,00,04,07,00,04,07,00,07,04,00,07,04,00}.wav \
  ../loom.wav \
  delay $(seq 1 0.2 5.8) \
  remix - \
  remix - - \
  reverb 95 50 100 100 \
  fade q 0 7.8 1.5 \
  gain -n -1

sox \
  --combine merge \
  ../shepard_{09,00,04,09,04,00,09,00,04,09,04,00,09}.wav \
  ../loom_minor.wav \
  delay $(seq 1 0.4 5.8) \
  remix - \
  remix - - \
  reverb 95 50 100 100 \
  fade q 0 7.8 1.5 \
  gain -n -1

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

# Renders a MIDI to a WAV using Shepard Tone.

set -ex

midi=$1; shift
wav=$1; shift

converted=$(mktemp --suffix=.mid)
trap 'rm -f "$converted"' EXIT

config=$(mktemp --suffix=.sfz)
trap 'rm -f "$converted" "$config"' EXIT

midicopy -nodrums "$midi" "$converted"
sh shepard.sfz.sh > "$config"

sfizz_render --sfz "$config" --polyphony 4096 --midi "$converted" --wav "$wav"

#!/bin/sh
# Copyright 2022 Google LLC
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

set -ex

d0=$PWD
out=$PWD/aaaaxy.dat

rm -f "$out"

# Reverse order so just in case, assets override everything else.

cd "$d0"/licenses
7za a -tzip -mx=9 "$out" [!_]*/[!_]*

for d in "$d0"/third_party/*/assets; do
	cd "$d"
	7za a -tzip -mx=9 "$out" [!_]*/[!_]*
done

cd "$d0"/assets
7za a -tzip -mx=9 "$out" [!_]*/[!_]*

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

set -ex

rm -rf \
	packaging/AAAAXY.iconset \
	packaging/AAAAXY.app/Contents/Resources \
	packaging/AAAAXY.app/Contents/Info.plist

mkdir -p packaging/AAAAXY.app/Contents/Resources
mkdir -p packaging/AAAAXY.iconset

for res in 16 32 128 256 512; do
	convert assets/sprites/riser_small_up_0.png \
		-filter Point -geometry ${res}x${res} \
		-define png:bit-depth=8 \
		-define png:color-type=6 \
		-define png:format=png32 \
		packaging/AAAAXY.iconset/icon_${res}x${res}.png
done
png2icns packaging/AAAAXY.app/Contents/Resources/icon.icns packaging/AAAAXY.iconset/*.png ||\
iconutil --convert icns --output packaging/AAAAXY.app/Contents/Resources/icon.icns packaging/AAAAXY.iconset
scripts/Info.plist.sh $(scripts/version.sh macos) packaging/AAAAXY.app
ln -f aaaaxy.dat packaging/AAAAXY.app/Contents/Resources/aaaaxy.dat

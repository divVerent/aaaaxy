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

export ANDROID_HOME=$HOME/Android/Sdk

export CGO_CPPFLAGS=-DNDEBUG
export CGO_CFLAGS='-g0 -O3'
export CGO_CXXFLAGS='-g0 -O3'
export CGO_LDFLAGS='-g0 -s'
mkdir -p AndroidStudioProjects/AAAAXY/app/libs
GOOS=darwin go generate \
	-tags zip
go run github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile bind \
	-target android \
	-javapkg io.github.divverent.aaaaxy \
	-o AndroidStudioProjects/AAAAXY/app/libs/aaaaxy.aar \
	-androidapi 21 \
	-tags zip \
	-gcflags=all=-dwarf=false \
	-ldflags=all=-s \
	-ldflags=all=-w \
	-a \
	-trimpath \
	github.com/divVerent/aaaaxy/AndroidStudioProjects/AAAAXY/app/src/main/go/aaaaxy

mkdir -p AndroidStudioProjects/AAAAXY/app/src/main/assets
ln -f aaaaxy.dat AndroidStudioProjects/AAAAXY/app/src/main/assets/aaaaxy.dat

cd AndroidStudioProjects/AAAAXY
./gradlew assembleDebug
./gradlew assembleRelease
./gradlew bundleRelease

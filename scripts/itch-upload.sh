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

# This script uploads a new release to Itch.

version=$1; shift

set -ex

../butler/butler push aaaaxy-linux-amd64-"$version".zip divVerent/aaaaxy:linux-amd64
../butler/butler push AAAAXY-x86_64.AppImage divVerent/aaaaxy:linux-amd64-appimage
../butler/butler push aaaaxy-windows-amd64-"$version".zip divVerent/aaaaxy:windows-amd64
../butler/butler push aaaaxy-windows-386-"$version".zip divVerent/aaaaxy:windows-386
../butler/butler push aaaaxy-darwin-"$version".zip divVerent/aaaaxy:mac
../butler/butler push aaaaxy-js-wasm-"$version".zip divVerent/aaaaxy:js-wasm

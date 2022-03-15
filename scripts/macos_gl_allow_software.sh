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

# Preload hack to allow GLFW with software rendering on macOS.
# Allows running the binary inside a VM for testing on CI.

set -ex
cd scripts
clang -framework Foundation -dynamiclib -o macos_gl_allow_software.dylib macos_gl_allow_software.m

# Normally we could just use DYLD_INSERT_LIBRARIES...
# However it seems like SIP is active on GitHub Actions,
# so instead, let's edit the binary.
git clone https://github.com/Tyilo/insert_dylib || true
cd insert_dylib/insert_dylib
clang main.c -o insert_dylib
cd ../../..

for binary in "$@"; do
	scripts/insert_dylib/insert_dylib/insert_dylib --all-yes --inplace scripts/macos_gl_allow_software.dylib "$binary"
done

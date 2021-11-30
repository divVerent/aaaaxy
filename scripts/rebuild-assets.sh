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

# This script rebuilds all regeneratable assets.
#
# Ideally these assets should move into assets/generated and thus no longer be
# in the git repo; will be some work to support this in Tiled.
#
# Also currently this isn't reproducible - the files differ in some metadata.
# That needs fixing, too.

set -ex

for src in "$(pwd)"/assets/*/_src/*.sh; do
	cd "${src%/*}"
	sh "$src"
done

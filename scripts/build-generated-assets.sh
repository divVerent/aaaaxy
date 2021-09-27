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

: ${GO:=go}

# Run go natively.
export GOOS=
export GOARCH=

mkdir -p assets/generated
${GO} run github.com/divVerent/aaaaxy/cmd/dumpcps assets/maps/level.tmx > assets/generated/level.cp.dot
neato -Tjson assets/generated/level.cp.dot > assets/generated/level.cp.json
scripts/image-load-order.sh assets/tiles assets/sprites third_party/grafxkid_classic_hero_and_baddies_pack/assets/sprites > assets/generated/image_load_order.txt
scripts/version.sh semver > assets/generated/version.txt

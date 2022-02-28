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

if [ x"$AAAAXY_GENERATE_ASSETS" = x'true' ]; then
	if ! [ "assets/generated/level.cp.json" -nt "assets/maps/level.tmx" ]; then
		trap 'rm -f assets/generated/level.cp.json' EXIT
		# Using |cat> instead of > because snapcraft for some reason doesn't allow using a regular > shell redirection with "go run".
		${GO} run github.com/divVerent/aaaaxy/cmd/dumpcps |cat> assets/generated/level.cp.dot
		grep -c . assets/generated/level.cp.dot
		neato -Tjson assets/generated/level.cp.dot > assets/generated/level.cp.json
		grep -c . assets/generated/level.cp.json
		trap - EXIT
	fi
	diff -bu -I'.*"width".*' assets/_saved/level.cp.json assets/generated/level.cp.json

	scripts/image-load-order.sh assets/generated/image_load_order.txt assets/tiles assets/sprites third_party/grafxkid_classic_hero_and_baddies_pack/assets/sprites
	diff -u assets/_saved/image_load_order.txt assets/generated/image_load_order.txt
else
	cp assets/_saved/* assets/generated/
fi

scripts/version.sh semver > assets/generated/version.txt

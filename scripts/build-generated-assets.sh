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
cp assets/_saved/* assets/generated/

if [ x"$AAAAXY_GENERATE_ASSETS" = x'true' ]; then
	if [ x"$AAAAXY_GENERATE_CHECKPOINT_LOCATIONS" = x'true' ]; then
		for lfile in assets/maps/*.tmx; do
			lname=${lfile%.tmx}
			lname=${lname##*/}
			if [ x"$AAAAXY_FORCE_GENERATE_ASSETS" = x'true' ] || ! [ "assets/generated/$lname.cp.json" -nt "assets/maps/$lname.tmx" ]; then
				trap 'rm -f "assets/generated/$lname.cp.json"' EXIT
				# Using |cat> instead of > because snapcraft for some reason doesn't allow using a regular > shell redirection with "go run".
				${GO} run ${GO_FLAGS} github.com/divVerent/aaaaxy/cmd/dumpcps -level="$lname" |cat> "assets/generated/$lname.cp.dot"
				grep -c . "assets/generated/$lname.cp.dot"
				neato -Tjson assets/generated/$lname.cp.dot > assets/generated/$lname.cp.json
				grep -c . "assets/generated/$lname.cp.json"
				trap - EXIT
			fi
			if [ x"$AAAAXY_DIFF_ASSETS" != x'false' ]; then
				diff -bu -I'.*"width".*' assets/_saved/level.cp.json assets/generated/level.cp.json
			fi
		done
	fi

	if [ x"${AAAAXY_FORCE_GENERATE_ASSETS}" = x'true' ] || [ x"${AAAAXY_DIFF_ASSETS}" != x'false' ]; then
		rm -f assets/generated/lut_*.png
	fi
	${GO} run ${GO_FLAGS} github.com/divVerent/aaaaxy/cmd/dumpluts --palette_max_cycles=inf

	sh scripts/image-load-order.sh assets/generated/image_load_order.txt assets/tiles assets/sprites third_party/grafxkid_classic_hero_and_baddies_pack/assets/sprites
	if [ x"$AAAAXY_DIFF_ASSETS" != x'false' ]; then
		diff -u assets/_saved/image_load_order.txt assets/generated/image_load_order.txt
	fi

	if [ x"$AAAAXY_DIFF_ASSETS" != x'false' ]; then
		for f in assets/_saved/lut_*.png; do
			g=assets/generated/"${f##*/}"
			result=$(convert "$f" "$g" \
				-channel RGBA \
				-metric RMSE -format '%[distortion]' -compare \
				INFO:)
			if [ x"$result" != x'0' ]; then
				echo >&2 "$f and $g differ."
				exit 1
			fi
		done
	fi
fi

sh scripts/version.sh semver > assets/generated/version.txt

# Prepare compressed font.
gzip -9 < ./third_party/gnu_unifont/assets/fonts/_unifont-15.1.04.bdf > assets/generated/unifont.bdf.gz

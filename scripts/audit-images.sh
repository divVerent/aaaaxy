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

find ./assets -name \*.png | sort | while read -r file; do
	# Exceptions.
	case "$file" in
		# Editing only.
		*/_editorimgs/*) continue ;;
		*/_saved/*) continue ;;
		*/_src/*) continue ;;
		*/generated/*) continue ;;
		./assets/sprites/warpzone_*.png) continue ;;
		# Intentionally violating.
		./assets/sprites/clock_*.png) continue ;;
		./assets/sprites/gradient_*.png) continue ;;
		./assets/sprites/magic_*.png) continue ;;
		# Screenshots etc.
		./docs/*) continue ;;
		# SDL.
		./third_party/SDL_GameControllerDB/*) continue ;;
	esac
	set -- \
		"$file" -depth 8 +dither \
		-write MPR:orig \
		-channel RGB -remap scripts/cga-palette.pnm +channel \
		MPR:orig -alpha set -compose copy-opacity -composite \
		-channel A -threshold 25% +channel
	f=$(
		convert \
			\( "$file" -depth 8 +dither \) \
			\( "$@" \) \
			-channel RGBA \
			-metric RMSE -format '%[distortion]\n' -compare \
			INFO:
	)
	if [ "$f" !=  0 ]; then
		echo "convert "$@" "$file"  # $f"
	fi
done

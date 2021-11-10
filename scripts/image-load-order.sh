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

set -e

root=$PWD
out=$1; shift

need=false
for dir in "$@"; do
	if ! [ "$out" -nt "$dir" ]; then
		need=true
		break
	fi
	for img in "$dir"/*.png; do
		if ! [ "$out" -nt "$img" ]; then
			need=true
			break
		fi
	done
	if $need; then
		break
	fi
done
if ! $need; then
	exit 0
fi

# Load largest images first to optimize the BSP-based atlas ebiten generates.
trap 'rm -f "$out"' EXIT
for dir in "$@"; do
	for img in "$dir"/*.png; do
		vfsimg=${img#third_party/*/}
		vfsimg=${vfsimg#assets/}
		set -- $(identify -format '%[width] %[height]' "$img")
		echo "$(($1 * $2)) $1 $2 $vfsimg"
	done
done | sort -r -n -k 1,3 -s | cut -d ' ' -f 4 > "$out"
trap - EXIT

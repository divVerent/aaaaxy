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

destdir="$1"; shift

root=$PWD

tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

logged() {
	printf >&2 '+ %s\n' "$*"
	"$@"
}

for sourcedir in assets third_party/*/assets; do
	logged cd "$root/$sourcedir"
	find . -type f | while read -r file; do
		mkdir -p "$tmpdir/${file%/*}"
		logged ln -snf "$root/$sourcedir/$file" "$tmpdir/$file"
	done
	# Also copy over our license and copyright files to make
	# the copyright situation really clear to anyone extracting the
	# VFS data.
	directory=${sourcedir%/*}
	directory=${directory##*/}
	{
		echo "Applying to the following files:"
		logged find . -type f -print
		if [ -f ../COPYRIGHT.md ]; then
			echo
			logged cat ../COPYRIGHT.md
		fi
		echo
		echo "License file: $directory.LICENSE"
	} | logged dd status=none of="$tmpdir/$directory.COPYRIGHT"
	logged ln -snf "$root/$sourcedir/../LICENSE" "$tmpdir/$directory.LICENSE"
done
logged go run github.com/rakyll/statik -m -f -src "$tmpdir/" -dest "$root/$destdir"

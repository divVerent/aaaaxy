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

: ${GO:=go}

destdir="$1"; shift

root=$PWD

tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

logged() {
	printf >&2 '+ %s\n' "$*"
	"$@"
}

for sourcedir in assets third_party/*/assets licenses; do
	case "$sourcedir" in
		licenses)
			prefix=licenses/
			;;
		*)
			prefix=
			;;
	esac
	logged cd "$root/$sourcedir"
	find . -name src -prune -or -name editorimgs -prune -or -type f -print | while read -r file; do
		mkdir -p "$tmpdir/$prefix${file%/*}"
		logged ln -snf "$root/$sourcedir/$file" "$tmpdir/$prefix$file"
	done
done
logged $GO run github.com/rakyll/statik -m -f -src "$tmpdir/" -dest "$root/$destdir"

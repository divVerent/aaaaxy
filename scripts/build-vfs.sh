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

root=$PWD
destdir="$root"/internal/vfs/_embedroot

rm -rf "$destdir"
for sourcedir in assets third_party/*/assets licenses; do
	case "$sourcedir" in
		licenses)
			prefix=licenses/
			;;
		*)
			prefix=
			;;
	esac
	cd "$root/$sourcedir"
	find . -name _src -prune -or -name _editorimgs -prune -or -type f -print | while read -r file; do
		mkdir -p "$destdir/$prefix${file%/*}"
		cp "$root/$sourcedir/$file" "$destdir/$prefix$file"
	done
done

echo "Checkpoints file embedded:"
cat "$destdir"/generated/level.cp.json
echo "."

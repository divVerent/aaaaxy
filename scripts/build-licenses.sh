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

out=$1

rm -rf licenses/*.txt

# Note: ignoring errors here, as some golang.org packages
# do not have a discoverable license file. As they're all under Go's license,
# that is fine.
$GO run ${GO_FLAGS} github.com/google/go-licenses save github.com/divVerent/aaaaxy --force --save_path=licenses/software-licenses || true

# This will fail if go-licenses wrote no output.
for d in licenses/software-licenses/*/; do
	[ -d "$d" ]
done

# Translate to a single directory.
find licenses/software-licenses -type f | while read -r path; do
	file=${path##*/}
	path=${path%/*}
	name=${path#licenses/software-licenses/}
	case "$name" in
		github.com/divVerent/aaaaxy)
			cleanname="aaaaxy"
			;;
		*)
			cleanname=software-$(echo -n "$name" | tr -c '0-9A-Za-z-' '_')
			;;
	esac
	echo "$name:" > "licenses/$cleanname-COPYRIGHT.txt"
	mv "$path/$file" "licenses/$cleanname-$file.txt"
done

rm -rf licenses/software-licenses

# Add our own third party stuff.
rm -rf licenses/asset-licenses
find third_party -name COPYRIGHT.md | while read -r path; do
	path=${path%/*}
	name=${path##*/}
	cleanname=asset-$(echo -n "$name" | tr -c '0-9A-Za-z-' '_')
	{
		echo "$name:"
		echo
		cat "$path/COPYRIGHT.md"
	} > "licenses/$cleanname-COPYRIGHT.txt"
	cp "$path/LICENSE" "licenses/$cleanname-LICENSE.txt"
done

# TODO: change structure as follows:
#   licenses/item.txt         - contains just the name of the item, plus extra info
#   licenses/item-LICENSE.txt - contains the license
# That would be something we can easily display in a license dialog.

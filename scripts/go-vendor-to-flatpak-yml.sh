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

flatpakdir=$1

TAB="	"
LF="
"

rm -rf vendor
go mod vendor

d=$(mktemp -d)

yml="$flatpakdir/io.github.divverent.aaaaxy.yml"

sed -i -e '1,/# --- GO MODULES START HERE.* ---/!d' "$yml"
cp vendor/modules.txt "$flatpakdir/modules.txt"
exec 3>>"$yml"

d0=$PWD
while read -r command pkg ver _; do
	[ x"$command" = x'#' ] || continue
	suffix=
	case "$pkg" in
		*/v?)
			pkg=${pkg%/v?}
			suffix=/${pkg##*/v?}
			;;
	esac
	case "$pkg" in
		github.com/*)
			# Cut off subdirectory paths.
			# go-gl seems to use that.
			pkg=$(echo "$pkg" | cut -d / -f 1-3)
			;;
	esac
	url=https://$pkg
	case "$pkg" in
		golang.org/x/*)
			# These modules don't use their real git URL.
			# Must be some special handling in "go get".
			url=https://go.googlesource.com/${pkg##*/}
			;;
		*)
			url=https://$pkg
			;;
	esac
	rm -rf "$d/git"
	git clone "$url" "$d/git"
	cd "$d/git"
	case "$ver" in
		*-*-*)
			tag=
			commit=$(git rev-parse "${ver##*-}")
			version="commit: $commit"
			;;
		*)
			commit=$(git rev-parse "$ver")
			version="tag: $ver$LF        commit: $commit"
			;;
	esac
	cat >&3 <<EOF
      - type: git
        url: $url
        $version
        dest: vendor/$pkg$suffix
EOF
	rm -rf "$d/git"
	cd "$d0"
done < "$d0/vendor/modules.txt"

exec 3>&-

rm -rf "$d" vendor

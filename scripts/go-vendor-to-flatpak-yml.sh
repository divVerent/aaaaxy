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
while read -r command pkg ver _ replacementpkg replacementver; do
	[ x"$command" = x'#' ] || continue
	[ x"$ver" != x'=>' ] || continue
	dir=$pkg
	if [ -n "$replacementpkg" ]; then
		pkg=$replacementpkg
	fi
	if [ -n "$replacementver" ]; then
		ver=$replacementver
	fi
	suffix=
	case "$pkg" in
		*/v?)
			suffix=/${pkg##*/}
			pkg=${pkg%/v?}
			dir=${dir%/v?}
			;;
	esac
	case "$pkg" in
		github.com/*)
			# Cut off subdirectory paths.
			# go-gl seems to use that.
			pkg=$(echo "$pkg" | cut -d / -f 1-3)
			dir=$(echo "$dir" | cut -d / -f 1-3)
			;;
		golang.org/x/exp/shiny)
			# For some reason this has to be fetched from x/exp.
			pkg=$(echo "$pkg" | cut -d / -f 1-3)
			dir=$(echo "$dir" | cut -d / -f 1-3)
			;;
	esac
	url=https://$pkg
	case "$pkg" in
		golang.org/x/*)
			# These modules don't use their real git URL.
			# Must be some special handling in "go get".
			url=https://go.googlesource.com/${pkg##*/}
			;;
		go.opencensus.io)
			url=https://github.com/census-instrumentation/opencensus-go
			;;
		k8s.io/klog)
			url=https://github.com/kubernetes/klog
			;;
		*)
			url=https://$pkg
			;;
	esac
	rm -rf "$d/git"
	git clone "$url" "$d/git"
	cd "$d/git"
	case "$pkg":"$ver" in
		github.com/hajimehoshi/oto:v2.3.0-alpha.4)
			# Retracted version.
			# See https://github.com/hajimehoshi/oto/issues/177.
			ver=v2.2.0-alpha.4
			;;
	esac
	if git rev-parse "$ver" >/dev/null 2>&1; then
		commit=$(git rev-parse "$ver")
		version="tag: $ver$LF        commit: $commit"
	else
		tag=
		commit=$(git rev-parse "${ver##*-}")
		version="commit: $commit"
	fi
	cat >&3 <<EOF
      - type: git
        url: $url
        $version
        dest: vendor/$dir$suffix
EOF
	rm -rf "$d/git"
	cd "$d0"
done < "$d0/vendor/modules.txt"

exec 3>&-

rm -rf "$d" vendor

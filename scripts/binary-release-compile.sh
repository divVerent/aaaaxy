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

GOOS=$($GO env GOOS)
GOARCH=$($GO env GOARCH)
GOEXE=$($GO env GOEXE)
zip="$PWD/aaaaxy-$GOOS-$GOARCH-$(scripts/version.sh gittag).zip"

exec 3>&1
exec >&2

case "$GOOS" in
	darwin)
		appdir=packaging/
		app=AAAAXY.app
		prefix=packaging/AAAAXY.app/Contents/MacOS/
		;;
	js)
		appdir=.
		app="aaaaxy-$GOOS-$GOARCH$GOEXE aaaaxy.html wasm_exec.js"
		prefix=
		;;
	*)
		appdir=.
		app=aaaaxy-$GOOS-$GOARCH$GOEXE
		prefix=
		;;
esac

make clean
make BUILDTYPE=release PREFIX="$prefix"

rm -f "$zip"
7za a -tzip -mx=9 "$zip" \
	README.md LICENSE CONTRIBUTING.md \
	licenses
(
	cd "$appdir"
	7za a -tzip -mx=9 "$zip" \
		$app
)

case "$GOOS" in
	linux)
		arch=$GOARCH
		case "$arch" in
			amd64)
				arch=x86_64
				;;
			386)
				arch=x86
				;;
		esac
		rm -rf packaging/AAAAXY.AppDir
		linuxdeploy-$(uname -m).AppImage \
			--appdir=packaging/AAAAXY.AppDir \
			-e "$app" \
			-d packaging/"$app".desktop \
			-i packaging/"$app".png
		mkdir -p packaging/AAAAXY.AppDir/usr/share/metainfo
		id=io.github.divverent.aaaaxy_$($GO env GOARCH)
		cp packaging/"$id".metainfo.xml packaging/AAAAXY.AppDir/usr/share/metainfo/
		appimagetool-$(uname -m).AppImage \
			-u "gh-releases-zsync|divVerent|aaaaxy|latest|AAAAXY-$arch.AppImage.zsync" \
			packaging/AAAAXY.AppDir \
			"AAAAXY-$arch.AppImage"
		;;
esac

make clean

echo >&3 "$zip"

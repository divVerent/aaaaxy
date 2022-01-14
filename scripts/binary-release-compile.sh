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
: ${LIPO:=lipo}

GOOS=$($GO env GOOS)
GOEXE=$($GO env GOEXE)

case "$GOOS" in
	js)
		# HACK: Itch and Apache want a .wasm file extension, but GOEXE doesn't actually have that.
		GOEXE=.wasm
		;;
esac

case "$#" in
	1)
		GOARCH_SUFFIX=-$1
		;;
	*)
		GOARCH_SUFFIX=
		;;
esac

zip="$PWD/aaaaxy-$GOOS$GOARCH_SUFFIX-$(scripts/version.sh gittag).zip"

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
		app="aaaaxy-$GOOS$GOARCH_SUFFIX$GOEXE index.html wasm_exec.js"
		prefix=
		;;
	*)
		appdir=.
		app=aaaaxy-$GOOS$GOARCH_SUFFIX$GOEXE
		prefix=
		;;
esac

make clean

case "$prefix" in
	*/*)
		mkdir -p "${prefix%/*}"
		;;
esac

if [ -n "$GOARCH_SUFFIX" ]; then
	eval "export CGO_ENV=\$CGO_ENV_$1"
	binary=${prefix}aaaaxy-$GOOS$GOARCH_SUFFIX$GOEXE
	GOARCH=$(GOARCH=$1 $GO env GOARCH) make BUILDTYPE=release BINARY="$binary"
	unset CGO_ENV
else
	lipofiles=
	for arch in "$@"; do
		eval "export CGO_ENV=\$CGO_ENV_$arch"
		binary=${prefix}aaaaxy-$GOOS-$arch$GOEXE
		GOARCH=$(GOARCH=$arch $GO env GOARCH) make BUILDTYPE=release BINARY="$binary"
		unset CGO_ENV
		lipofiles="$lipofiles $binary"
	done
	binary=${prefix}aaaaxy$GOARCH_SUFFIX$GOEXE
	$LIPO -create $lipofiles -output "$binary"
	rm -f $lipofiles
fi

case "$GOOS" in
	js)
		# Pack in a form itch.io can use.
		cp aaaaxy.html index.html
		cp "$(GOOS=js GOARCH=wasm go env GOROOT)"/misc/wasm/wasm_exec.js .
		;;
esac

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
		arch=${GOARCH_SUFFIX#-}
		case "$arch" in
			amd64)
				arch=x86_64
				;;
			386)
				arch=x86
				;;
		esac
		scripts/build-appimage-resources.sh
		rm -rf packaging/AAAAXY.AppDir
		linuxdeploy-$(uname -m).AppImage \
			--appdir=packaging/AAAAXY.AppDir \
			-e "$app" \
			-d packaging/"$app".desktop \
			-i packaging/"$app".png
		mkdir -p packaging/AAAAXY.AppDir/usr/share/metainfo
		id=io.github.divverent.aaaaxy_${GOARCH_SUFFIX#-}
		cp packaging/"$id".metainfo.xml packaging/AAAAXY.AppDir/usr/share/metainfo/
		appimagetool-$(uname -m).AppImage \
			-u "gh-releases-zsync|divVerent|aaaaxy|latest|AAAAXY-$arch.AppImage.zsync" \
			packaging/AAAAXY.AppDir \
			"AAAAXY-$arch.AppImage"
		;;
esac

make clean

echo >&3 "$zip"

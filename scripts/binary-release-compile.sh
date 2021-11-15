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
		app="aaaaxy-$GOOS$GOARCH_SUFFIX$GOEXE aaaaxy.html wasm_exec.js"
		prefix=
		;;
	*)
		appdir=.
		app=aaaaxy-$GOOS$GOARCH_SUFFIX$GOEXE
		prefix=
		;;
esac

make clean

if [ -n "$GOARCH_SUFFIX" ]; then
	eval "export CGO_ENV=\$CGO_ENV_$1"
	GOARCH=$(GOARCH=$1 $GO env GOARCH) make BUILDTYPE=release PREFIX="$prefix"
	unset CGO_ENV
else
	lipofiles=
	for arch in "$@"; do
		eval "export CGO_ENV=\$CGO_ENV_$arch"
		GOARCH=$(GOARCH=$arch $GO env GOARCH) make BUILDTYPE=release PREFIX="$prefix"
		unset CGO_ENV
		lipofiles="$lipofiles ${prefix}aaaaxy-$GOOS-$arch$GOEXE"
	done
	$LIPO -create $lipofiles -output "${prefix}aaaaxy-$GOOS"
	rm -f $lipofiles
fi

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
		id=io.github.divverent.aaaaxy_${GOARCH_SUFFIX#-}
		cp packaging/"$id".metainfo.xml packaging/AAAAXY.AppDir/usr/share/metainfo/
		appimagetool-$(uname -m).AppImage \
			-u "gh-releases-zsync|divVerent|aaaaxy|latest|AAAAXY-$arch.AppImage.zsync" \
			packaging/AAAAXY.AppDir \
			"AAAAXY$GOARCH_SUFFIX.AppImage"
		;;
esac

make clean

echo >&3 "$zip"

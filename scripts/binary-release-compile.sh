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
zip="aaaaxy-$GOOS-$GOARCH-$(scripts/version.sh gittag).zip"

exec 3>&1
exec >&2

make clean
make BUILDTYPE=release

# Then pack it all together.
case "$GOOS" in
	darwin) app=AAAAXY.app ;;
	js) app="aaaaxy-$GOOS-$GOARCH$GOEXE aaaaxy.html wasm_exec.js" ;;
	*) app=aaaaxy-$GOOS-$GOARCH$GOEXE ;;
esac

rm -f "$zip"
7za a -tzip -mx=9 "$zip" \
	$app \
	README.md LICENSE CONTRIBUTING.md \
	licenses

make clean

echo >&3 "$zip"

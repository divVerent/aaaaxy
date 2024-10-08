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

version=$1
app=$2

binary=
for x in "$app/Contents/MacOS"/*; do
	[ -f "$x" ] || continue
	[ -x "$x" ] || continue
	[ -z "$binary" ] || {
		echo >&2 'More than one executable found!'
		exit 1
	}
	binary=${x##*/}
done

cat > "$app/Contents/Info.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>$binary</string>
	<key>CFBundleIconFile</key>
	<string>icon.icns</string>
	<key>CFBundleName</key>
	<string>AAAAXY</string>
	<key>CFBundleIdentifier</key>
	<string>io.github.divverent.aaaaxy</string>
	<key>CFBundleVersion</key>
	<string>$version</string>
	<key>LSMinimumSystemVersion</key>
	<string>10.13.0</string>
	<key>LSApplicationCategoryType</key>
	<string>public.app-category.puzzle-games</string>
	<key>NSHighResolutionCapable</key>
	<true/>
</dict>
</plist>
EOF

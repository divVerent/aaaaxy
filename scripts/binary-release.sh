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

prev=$(git describe --always --long --match 'v*.*' --exclude 'v[0-9].[0-9]' --exclude 'v[0-9].[0-9].0-alpha' --exclude 'v[0-9].[0-9].0-beta' --exclude 'v[0-9].[0-9].0-rc')
# We want to exclude v*.* and v*.*.0-(alpha/beta).
prev=${prev%-*-g*}

new=$(sh scripts/version.sh gittag)

echo "Releaseing: $prev -> $new."

GOOS=linux GOARCH=amd64 scripts/binary-release-compile.sh
GOOS=windows GOARCH=amd64 scripts/binary-release-compile.sh
GOOS=linux GOARCH=386 scripts/binary-release-compile.sh
GOOS=darwin GOARCH=amd64 CGO_ENV="PATH=$HOME/src/osxcross-sdk/bin:$PATH CGO_ENABLED=1 CC=o64-clang CXX=o64-clang++ MACOSX_DEPLOYMENT_TARGET=10.12" scripts/binary-release-compile.sh

git tag -a "$new" -m "$(
	echo "Release $new"
	echo
	echo "Changes since $prev:"
	git log --format='%w(72,2,4)- %s' "$prev"..
)" -e

echo "Now run:"
echo "  git push origin tag $new"
echo "Then create the release on GitHub with the following message:"
git show -s "$new"
echo "In the release, upload aaaaxy-*-$new.zip"

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

prev=$(git describe --always --long)
prev=${prev%-*-g*}

new=$(sh scripts/version.sh gittag)

echo "Releaseing: $prev -> $new."

make allrelease ZIPFILE="aaaaxy-$new.zip"

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
echo "In the release, upload aaaaxy-$new.zip"

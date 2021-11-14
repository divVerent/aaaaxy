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

echo "Releasing: $prev -> $new."

GOOS=linux scripts/binary-release-compile.sh amd64
GOOS=windows scripts/binary-release-compile.sh amd64
GOOS=windows scripts/binary-release-compile.sh 386
GOOS=darwin CGO_ENV_amd64="PATH=$HOME/src/osxcross-sdk/bin:$PATH CGO_ENABLED=1 CC=o64-clang CXX=o64-clang++ MACOSX_DEPLOYMENT_TARGET=10.12" CGO_ENV_arm64="PATH=$HOME/src/osxcross-sdk/bin:$PATH CGO_ENABLED=1 CC=oa64e-clang CXX=oa64e-clang++ MACOSX_DEPLOYMENT_TARGET=10.12" LIPO="$HOME/src/osxcross-sdk/bin/lipo" scripts/binary-release-compile.sh amd64 arm64

cat <<EOF >.commitmsg
Release $new

Changes since $prev:
$(git log --format='%w(72,2,4)- %s' "$prev"..)
EOF
vi .commitmsg

VERSION=$new perl -0777 -pi -e '
	use strict;
	use warnings;
	my $version = $ENV{VERSION};
	/(?<=<!-- BEGIN DOWNLOAD LINKS TEMPLATE\n)(.*)(?=\nEND DOWNLOAD LINKS TEMPLATE -->)/s
		or die "Template not found.";
	my $template = $1;
	$template =~ s/VERSION/$version/g;
	s/(?<=<!-- BEGIN DOWNLOAD LINKS -->\n)(.*)(?=\n<!-- END DOWNLOAD LINKS -->)/$template/gs;
' docs/index.md

git commit -a -m "$(cat .commitmsg)"
git tag -a "$new" -m "$(cat .commitmsg)"

echo "Now run:"
echo "  git push origin main"
echo "  git push origin tag $new"
echo "Then create the release on GitHub with the following message:"
git show -s "$new"
echo "In the release, upload aaaaxy-*-$new.zip and AAAAXY-*.AppImage*."

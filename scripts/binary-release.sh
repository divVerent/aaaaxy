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
# Note: sync the MACOSX_DEPLOYMENT_TARGET with current Go requirements and Info.plist.sh.
GOOS=darwin CGO_ENV_amd64="PATH=$HOME/src/osxcross-sdk/bin:$PATH CGO_ENABLED=1 CC=o64-clang CXX=o64-clang++ MACOSX_DEPLOYMENT_TARGET=10.13" CGO_ENV_arm64="PATH=$HOME/src/osxcross-sdk/bin:$PATH CGO_ENABLED=1 CC=oa64-clang CXX=oa64-clang++ MACOSX_DEPLOYMENT_TARGET=10.13" LIPO="$HOME/src/osxcross-sdk/bin/lipo" scripts/binary-release-compile.sh amd64 arm64
GOOS=js scripts/binary-release-compile.sh wasm

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

set +x

cat <<EOF
Now run:
  git push origin tag $new
Then create the release on GitHub with the following message:
EOF
git show -s $new
cat <<EOF
In the release, upload aaaaxy-*-$new.zip (except for wasm) and
AAAAXY-*.AppImage*.
Once the release is published, finally run:
  git push origin main
If all this is done, consider also updating the snap:
  rm -f *.snap
  snap run snapcraft clean && snap run snapcraft && snap run snapcraft upload *.snap
And the FlatPak:
  go-vendor-to-flatpak-yml.sh ../io.github.divverent.aaaaxy
  cd ../io.github.divverent.aaaaxy
  vi io.github.divverent.aaaaxy.yml
  ... update commit and version number ...
  git commit -a
  git push origin HEAD:beta
  ... watch for build progress on https://flathub.org/builds/#/ ...
  ... test the build ...
  git push origin HEAD
Finally, update the URLs and wasm zip on itch.io.
EOF

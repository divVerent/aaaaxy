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

cat <<EOF >.commitmsg
Release $new

Changes since $prev:
$(git log --format='%w(72,2,4)- %s' "$prev"..)
EOF
vi .commitmsg

# Update gamecontroller mappings.
git submodule update --remote

# Update metainfo with current date and version already, and replace the text by a placeholder.
VERSION=$new DATE=$(date +%Y-%m-%d) MSG=$(cat .commitmsg) perl -0777 -pi -e '
	use strict;
	use warnings;
	my $version = $ENV{VERSION};
	my $date = $ENV{DATE};
	my $msg = $ENV{MSG};
	$msg =~ s/^Release .*//gm;
	$msg =~ s/^Changes since .*//gm;
	$msg =~ s/^  - /<\/li><li>/gm;
	$msg =~ s/^    //gm;
	$msg =~ s/^\n*<\/li>/<ul>/s;
	$msg =~ s/\n*$/<\/li><\/ul>/s;
	$msg =~ s/\n*<\/li>/<\/li>/g;
	$msg =~ s/\n/ /g;
	s/releases\/[^\/<]*<\/url>/releases\/$version<\/url>/g;
	s/<release version="[^"]*" date="[0-9-]*">/<release version="$version" date="$date">/g;
	s/<description>.*<\/description>/<description>$msg<\/description>/g;
' io.github.divverent.aaaaxy.metainfo.xml

GOOS=linux scripts/binary-release-compile.sh amd64
GOOS=windows scripts/binary-release-compile.sh amd64
GOOS=windows scripts/binary-release-compile.sh 386
# Note: sync the MACOSX_DEPLOYMENT_TARGET with current Go requirements and Info.plist.sh.
GOOS=darwin CGO_ENV_amd64="PATH=$HOME/src/osxcross-sdk/bin:$PATH CGO_ENABLED=1 CC=o64-clang CXX=o64-clang++ MACOSX_DEPLOYMENT_TARGET=10.13" CGO_ENV_arm64="PATH=$HOME/src/osxcross-sdk/bin:$PATH CGO_ENABLED=1 CC=oa64-clang CXX=oa64-clang++ MACOSX_DEPLOYMENT_TARGET=10.13" LIPO="$HOME/src/osxcross-sdk/bin/lipo" scripts/binary-release-compile.sh amd64 arm64
GOOS=js scripts/binary-release-compile.sh wasm

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
cat .commitmsg
cat <<EOF

In the release, upload aaaaxy-*-$new.zip (except for wasm) and
AAAAXY-*.AppImage*.
Once the release is published, finally run:
  git push origin main
If all this is done, consider also updating the snap:
  rm -f *.snap
  snap run snapcraft clean && snap run snapcraft && snap run snapcraft upload *.snap
Then go to https://snapcraft.io/aaaaxy/releases and activate the new build as
beta, later as stable.
And the FlatPak:
  scripts/go-vendor-to-flatpak-yml.sh ../io.github.divverent.aaaaxy
  cd ../io.github.divverent.aaaaxy
  vi io.github.divverent.aaaaxy.yml
  ... update commit and version number ...
  git commit -a
  git push origin HEAD:beta
  ... https://flathub.org/builds/#/apps/io.github.divverent.aaaaxy~2Fbeta ...
  ... watch, test, publish the build ...
  git push origin HEAD
  ... https://flathub.org/builds/#/apps/io.github.divverent.aaaaxy ...
  ... watch, test, publish the build ...
Finally, update the URLs and wasm zip on https://itch.io/game/edit/1199736.
EOF

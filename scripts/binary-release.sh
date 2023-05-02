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

: ${ADVZIP:=advzip -4}

if git ls-files -dmo | grep -q .; then
	echo >&2 'Working directory is not clean. Please commit or clean first.'
	exit 1
fi

if [ -z "$prev" ]; then
	prev=$(git describe --always --long --match 'v*.*' --exclude 'v[0-9].[0-9]' --exclude 'v[0-9].[0-9].0-alpha' --exclude 'v[0-9].[0-9].0-beta' --exclude 'v[0-9].[0-9].0-rc')
	# We want to exclude v*.* and v*.*.0-(alpha/beta).
	prev=${prev%-*-g*}
fi

new=$(sh scripts/version.sh gittag)

cat <<EOF >.commitmsg
Release $new

Changes since $prev:
$(git log --format='%w(72,2,4)- %s' "$prev"..)
EOF
vi .commitmsg

# Update gamecontroller mappings.
git submodule update --remote

# Include exact versions of submodules so that the source tarball on github
# contains the exact submodule version info.
git submodule > .gitmoduleversions

# Also store the current semver in the checkout. Used for compiling from
# source tarballs.
sh scripts/version.sh semver > .lastreleaseversion
sh scripts/version.sh android > .lastreleaseversioncode

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

# Provide changelog for Android.
mkdir -p fastlane/metadata/android/en-US/changelogs
log=fastlane/metadata/android/en-US/changelogs/"$(sh scripts/version.sh android)".txt
tail -n +4 .commitmsg > "$log"
git add "$log"

# Provide version for iOS.
sed -i -e "
	s,CURRENT_PROJECT_VERSION = .*;,CURRENT_PROJECT_VERSION = 1;,g;
	s,MARKETING_VERSION = .*;,MARKETING_VERSION = $(sh scripts/version.sh ios);,g;
" XcodeProjects/iOS/aaaaxy.xcodeproj/project.pbxproj

# Also pack the SDL game controller DB at the exact version used for the
# release. Used for compiling from source tarballs.
zip -r sdl-gamecontrollerdb-for-aaaaxy-$new.zip third_party/SDL_GameControllerDB/assets/input/*
$ADVZIP -z sdl-gamecontrollerdb-for-aaaaxy-$new.zip

# Also pack the files that do NOT get embedded into a mapping pack.
(
	cd assets/
	zip -r ../mappingsupport-for-aaaaxy-$new.zip ../LICENSE objecttypes.xml _* */_*
	$ADVZIP -z ../mappingsupport-for-aaaaxy-$new.zip
)

GOOS=linux sh scripts/binary-release-compile.sh amd64
GOOS=windows sh scripts/binary-release-compile.sh amd64
GOOS=windows GO386=sse2 sh scripts/binary-release-compile.sh 386
# Note: sync the MACOSX_DEPLOYMENT_TARGET with current Go requirements and Info.plist.sh.
GOOS=darwin CGO_ENV_amd64="PATH=$HOME/src/osxcross/target/bin:$PATH CGO_ENABLED=1 CC=o64-clang CXX=o64-clang++ MACOSX_DEPLOYMENT_TARGET=10.13" CGO_ENV_arm64="PATH=$HOME/src/osxcross/target/bin:$PATH CGO_ENABLED=1 CC=oa64-clang CXX=oa64-clang++ MACOSX_DEPLOYMENT_TARGET=10.13" LIPO="$HOME/src/osxcross/target/bin/lipo" sh scripts/binary-release-compile.sh amd64 arm64
GOOS=js sh scripts/binary-release-compile.sh wasm
(
	cd AndroidStudioProjects/AAAAXY/
	export ANDROID_HOME=$HOME/Android/Sdk
	./gradlew assembleRelease bundleRelease
)
cp AndroidStudioProjects/AAAAXY/app/build/outputs/apk/release/app-release.apk aaaaxy.apk

git commit -a -m "$(cat .commitmsg)"
git tag -a "$new" -m "$(cat .commitmsg)"
newrev=$(git rev-parse HEAD)
git push -f origin HEAD:binary-release-test

set +x

cat <<EOF
Please wait for automated tests on
https://github.com/divVerent/aaaaxy/actions

If these all pass, proceed by running

  sh scripts/publish-release.sh $new $newrev
EOF

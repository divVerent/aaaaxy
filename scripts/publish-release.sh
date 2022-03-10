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

new=$1; shift
newrev=$1; shift
dir=$PWD

if [ x"$(git rev-parse "$new"^0)" != x"$newrev" ]; then
	echo >&2 'Usage: $0 new-version new-git-revision'
	exit 1
fi

if [ x"$(git rev-parse "$new"^0)" != x"$(git rev-parse HEAD)" ]; then
	echo >&2 'Must be on the release tag.'
	exit 1
fi

if [ x"$(git symbolic-ref HEAD)" != x'refs/heads/main' ]; then
	echo >&2 'Must be on the main branch.'
	exit 1
fi

if ! [ -f .commitmsg ]; then
	echo >&2 'scripts/binary-release.sh has to be run first.'
	exit 1
fi

# Upload the binaries to GitHub.
hub release create \
	-a aaaaxy-linux-amd64-"$new".zip \
	-a AAAAXY-x86_64.AppImage \
	-a AAAAXY-x86_64.AppImage.zsync \
	-a aaaaxy-windows-amd64-"$new".zip \
	-a aaaaxy-windows-386-"$new".zip \
	-a aaaaxy-darwin-"$new".zip \
	-a sdl-gamecontrollerdb-for-aaaaxy-"$new".zip \
	-m "$(cat .commitmsg)" \
	"$new"

# Mark the release done.
git push origin main

# Publish it to gh-pages.
git worktree add /tmp/gh-pages gh-pages
(
	cd /tmp/gh-pages
	git reset --hard '@{u}'
	VERSION=$new perl -0777 -pi -e '
		use strict;
		use warnings;
		my $version = $ENV{VERSION};
		/(?<=<!-- BEGIN DOWNLOAD LINKS TEMPLATE\n)(.*)(?=\nEND DOWNLOAD LINKS TEMPLATE -->)/s
			or die "Template not found.";
		my $template = $1;
		$template =~ s/VERSION/$version/g;
		s/(?<=<!-- BEGIN DOWNLOAD LINKS -->\n)(.*)(?=\n<!-- END DOWNLOAD LINKS -->)/$template/gs;
	' index.md
	git commit -a -m "$(cat "$dir"/.commitmsg)"
	git push origin HEAD
)
git worktree remove /tmp/gh-pages

# Snap. Got kicked off by this git push.
xdg-open https://snapcraft.io/aaaaxy/builds

# Flatpak - first push a new build.
scripts/go-vendor-to-flatpak-yml.sh ../io.github.divverent.aaaaxy
(
	cd ../io.github.divverent.aaaaxy
	sed -i -e "/--- TAG GOES HERE ---/,+1 s/: .*/: $new/" io.github.divverent.aaaaxy.yml
	sed -i -e "/--- REV GOES HERE ---/,+1 s/: .*/: $newrev/" io.github.divverent.aaaaxy.yml
	git commit -a -m "Release $new."
	git push origin HEAD:beta
	git push origin HEAD
)

# Then let the user test and publish it.
xdg-open 'https://flathub.org/builds/#/apps/io.github.divverent.aaaaxy~2Fbeta'
xdg-open 'https://flathub.org/builds/#/apps/io.github.divverent.aaaaxy'

# Arch Linux.
(
	cd ../aur-aaaaxy
	sed -i -e "s/^pkgver=.*/pkgver=${new#v}/; s/^pkgrel=.*/pkgrel=1/;" PKGBUILD
	doas /root/archlinux/archlinux-testing-build-aaaaxy.sh
	git commit -a -m "Release $new."
	git push
)

# Itch.
scripts/itch-upload.sh "$new"
xdg-open https://itch.io/dashboard/game/1199736/devlog

# Gamejolt.
xdg-open https://gamejolt.com/dashboard/games/682854/packages
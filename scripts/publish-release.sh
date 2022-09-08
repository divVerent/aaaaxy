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
	echo >&2 'sh scripts/binary-release.sh has to be run first.'
	exit 1
fi

# First send the new tag to GitHub.
git push origin tag "$new"

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

# Alpine Linux comes first as it has some automated testing.
(
	cd ~/src/aports/community/aaaaxy
	git checkout master
	git fetch origin
	git reset --hard origin/master
	sed -i -e "s/^pkgver=.*/pkgver=${new#v}/; s/^pkgrel=.*/pkgrel=0/;" APKBUILD
	podman run --pull=always --rm --mount=type=bind,source=$PWD,target=/aaaaxy docker.io/library/alpine:edge /bin/sh -c '
		set -e
		apk add alpine-sdk sudo
		abuild-keygen -i -a -n
		cd /aaaaxy
		abuild -F checksum
		abuild -F -r
	'
	git commit -a -m "testing/aaaaxy: upgrade to $new"
	git push -f divVerent HEAD:aaaaxy
	# TODO is there a more direct URL to create a MR right away?
	xdg-open 'https://gitlab.alpinelinux.org/divVerent/aports/-/merge_requests/new?merge_request%5Bsource_branch%5D=aaaaxy'
)

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
sh scripts/go-vendor-to-flatpak-yml.sh ../io.github.divverent.aaaaxy
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
	podman run --pull=always --rm --mount=type=bind,source=$PWD,target=/aaaaxy docker.io/library/archlinux:latest /bin/sh -c '
		set -e
		pacman --noconfirm -Syu base-devel namcap pacman-contrib sudo
		useradd -m builder
		cp -Rv /aaaaxy /aaaaxy.rw
		chown -R builder:builder /aaaaxy.rw
		echo "builder ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers
		su builder -c "
			set -e
			cd /aaaaxy.rw
			updpkgsums
			makepkg -f --noconfirm --syncdeps
			makepkg --printsrcinfo > .SRCINFO
			namcap PKGBUILD
			namcap *.pkg.*
		"
		cat /aaaaxy.rw/PKGBUILD > /aaaaxy/PKGBUILD
		cat /aaaaxy.rw/.SRCINFO > /aaaaxy/.SRCINFO
	'
	git commit -a -m "Release $new."
	git push
)

# Itch.
sh scripts/itch-upload.sh "$new"
xdg-open https://itch.io/dashboard/game/1199736/devlog

# Pling.
xdg-open https://www.pling.com/p/1758731/edit

# Google Play.
xdg-open https://play.google.com/console/u/0/developers/7023357085064533456/app/4972666448425908274/tracks/open-testing

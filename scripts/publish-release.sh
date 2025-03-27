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

open_url() {
	if [ -n "$DISPLAY" ]; then
		xdg-open "$@"
	else
		echo "Open URL: $*"
		echo "Press ENTER when done..."
		read -r ok
	fi
}

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
	-a aaaaxy.apk \
	-m "$(cat .commitmsg)" \
	"$new"

# Alpine Linux comes first as it has some automated testing.
(
	cd ~/src/aports/community/aaaaxy
	git checkout master
	git fetch origin
	git reset --hard origin/master
	sed -i -e "s/^pkgver=.*/pkgver=${new#v}/; s/^pkgrel=.*/pkgrel=0/;" APKBUILD
	podman run --network=slirp4netns:enable_ipv6=false --pull=newer --rm --mount=type=bind,source=$PWD,target=/aaaaxy docker.io/library/alpine:edge /bin/sh -c '
		set -e
		apk add alpine-sdk sudo
		abuild-keygen -i -a -n
		cd /aaaaxy
		abuild -F checksum
		abuild -F -r
	'
	podman run --network=slirp4netns:enable_ipv6=false --pull=newer --rm --mount=type=bind,source=$PWD,target=/aaaaxy docker.io/library/alpine:3.21 /bin/sh -c '
		set -e
		apk add alpine-sdk sudo
		abuild-keygen -i -a -n
		cd /aaaaxy
		abuild -F -r
	'
	git commit -a -m "community/aaaaxy: upgrade to ${new#v}"
	git push -f divVerent HEAD:aaaaxy
	# TODO is there a more direct URL to create a MR right away?
	open_url 'https://gitlab.alpinelinux.org/divVerent/aports/-/merge_requests/new?merge_request%5Bsource_branch%5D=aaaaxy'
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
open_url https://snapcraft.io/aaaaxy/builds

# Flatpak - first push a new build.
(
	cd ../io.github.divverent.aaaaxy
	git checkout master
	git fetch
	git reset --hard '@{u}'
	git merge --no-edit beta
)
sh scripts/go-vendor-to-flatpak-yml.sh ../io.github.divverent.aaaaxy
(
	cd ../io.github.divverent.aaaaxy
	sed -i -e "/--- TAG GOES HERE ---/,+1 s/: .*/: $new/" io.github.divverent.aaaaxy.yml
	sed -i -e "/--- REV GOES HERE ---/,+1 s/: .*/: $newrev/" io.github.divverent.aaaaxy.yml
	git commit -a -m "Release $new."
	git push -f origin HEAD:aaaaxy
	open_url 'https://github.com/flathub/io.github.divverent.aaaaxy/compare/master...aaaaxy'
	open_url 'https://github.com/flathub/io.github.divverent.aaaaxy/compare/beta...aaaaxy'
	open_url 'https://flathub.org/builds/#/apps/io.github.divverent.aaaaxy~2Fbeta'
	open_url 'https://flathub.org/builds/#/apps/io.github.divverent.aaaaxy'
)

# Arch Linux.
(
	cd ../aur-aaaaxy
	git checkout master
	git fetch origin
	git reset --hard origin/master
	sed -i -e "s/^pkgver=.*/pkgver=${new#v}/; s/^pkgrel=.*/pkgrel=1/;" PKGBUILD
	podman run --network=slirp4netns:enable_ipv6=false --pull=always --rm --mount=type=bind,source=$PWD,target=/aaaaxy docker.io/library/archlinux:latest /bin/sh -c '
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
open_url https://itch.io/dashboard/game/1199736/devlog

# Pling.
open_url https://www.pling.com/p/1758731/edit

# Google Play.
open_url https://play.google.com/console/u/0/developers/7023357085064533456/app/4972666448425908274/tracks/production

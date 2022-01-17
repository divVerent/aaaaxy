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

# A script I run on my machine to test the Arch Linux PKGBUILD before
# committing. Not expected to work anywhere else.

mkdir /root/archlinux/root.x86_64/aur-aaaaxy
umount /root/archlinux/root.x86_64
mount --bind /root/archlinux/root.x86_64 /root/archlinux/root.x86_64
umount /root/archlinux/root.x86_64/aur-aaaaxy
mount --bind /home/rpolzer/src/aur-aaaaxy /root/archlinux/root.x86_64/aur-aaaaxy
/root/archlinux/root.x86_64/bin/arch-chroot /root/archlinux/root.x86_64  \
	pacman -Syu alsa-lib hicolor-icon-theme libglvnd libx11 git go graphviz imagemagick libxcursor libxinerama libxi libxrandr
/root/archlinux/root.x86_64/bin/arch-chroot /root/archlinux/root.x86_64 \
	su builder -c '
		set -ex
		rm -rf /tmp/aur-aaaaxy;
		cp -rv /aur-aaaaxy /tmp/aur-aaaaxy;
		cd /tmp/aur-aaaaxy;
		makepkg;
		namcap PKGBUILD;
		namcap *.pkg.*;
	'

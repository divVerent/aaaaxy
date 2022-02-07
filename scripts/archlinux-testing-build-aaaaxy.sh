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

set -ex

umount /root/archlinux/root.x86_64 || true
mount --bind /root/archlinux/root.x86_64 /root/archlinux/root.x86_64 || true
rsync -vaSHPAX /home/rpolzer/src/aur-aaaaxy/. /root/archlinux/root.x86_64/aur-aaaaxy
/root/archlinux/root.x86_64/bin/arch-chroot /root/archlinux/root.x86_64 pacman -Syu pacman-contrib
/root/archlinux/root.x86_64/bin/arch-chroot /root/archlinux/root.x86_64 chown -R builder /aur-aaaaxy
/root/archlinux/root.x86_64/bin/arch-chroot /root/archlinux/root.x86_64 \
	su builder -c '
		set -ex;
		cd /aur-aaaaxy;
		updpkgsums;
		makepkg --syncdeps;
		makepkg --printsrcinfo > .SRCINFO;
		namcap PKGBUILD;
		namcap *.pkg.*;
	'
cat /root/archlinux/root.x86_64/aur-aaaaxy/PKGBUILD > /home/rpolzer/src/aur-aaaaxy/PKGBUILD
cat /root/archlinux/root.x86_64/aur-aaaaxy/.SRCINFO > /home/rpolzer/src/aur-aaaaxy/.SRCINFO

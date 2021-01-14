#!/bin/sh

set -e

root=$PWD
destdir="$root/assets-vfs"

for sourcedir in assets third_party/*/assets; do
	cd "$root/$sourcedir"
	find . -type f | while read -r file; do
		dir=${file%/*}
		if [ x"$dir" != x"$prevdir" ]; then
			mkdir -vp "$destdir/$dir"
			prevdir=$dir
		fi
		ln -vsnf "$root/$sourcedir/$file" "$destdir/$file"
	done
done

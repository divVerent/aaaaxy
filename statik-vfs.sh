#!/bin/sh

set -ex

destdir="$1"; shift

root=$PWD

tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

if $NO_VFS; then
	touch "$tmpdir/use-local-assets.stamp"
else
	for sourcedir in assets third_party/*/assets; do
		cd "$root/$sourcedir"
		find . -type f | while read -r file; do
			mkdir -p "$tmpdir/${file%/*}"
			ln -snf "$root/$sourcedir/$file" "$tmpdir/$file"
		done
	done
fi
statik -m -f -src "$tmpdir/" -dest "$root/$destdir"

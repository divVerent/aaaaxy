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
	cd "$root"
	for license in third_party/*/LICENSE; do
		directory=${license%/*}
		directory=${directory#third_party/}
		cd "$root/third_party/$directory"
		[ -d assets ] || continue
		{
			echo "Applying to the following files:"
			find assets -type f -print
			echo
			cat LICENSE
		} > "$tmpdir/$directory.LICENSE"
	done
fi
statik -m -f -src "$tmpdir/" -dest "$root/$destdir"

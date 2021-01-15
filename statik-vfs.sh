#!/bin/sh

set -e

destdir="$1"; shift

root=$PWD

tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

logged() {
	printf >&2 '+ %s\n' "$*"
	"$@"
}

if $NO_VFS; then
	logged touch "$tmpdir/use-local-assets.stamp"
else
	for sourcedir in assets third_party/*/assets; do
		logged cd "$root/$sourcedir"
		find . -type f | while read -r file; do
			mkdir -p "$tmpdir/${file%/*}"
			logged ln -snf "$root/$sourcedir/$file" "$tmpdir/$file"
		done
	done
	cd "$root"
	for license in third_party/*/LICENSE; do
		directory=${license%/*}
		directory=${directory#third_party/}
		logged cd "$root/third_party/$directory"
		[ -d assets ] || continue
		{
			echo "Applying to the following files:"
			logged find assets -type f -print
			echo
			logged cat LICENSE
		} > "$tmpdir/$directory.LICENSE"
	done
fi
logged statik -m -f -src "$tmpdir/" -dest "$root/$destdir"

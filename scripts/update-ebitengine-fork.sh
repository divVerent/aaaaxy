#!/bin/sh

set -ex

version=$1
cherrypicks='
	99ffe09b63e0d906cc1f502c24f4d2325e6cc09d
'

cd ../ebiten
git fetch
git fetch divVerent
git checkout refs/tags/"$version"
cps=
for cp in $cherrypicks; do
	cp=$(git rev-parse "$cp")
	rev0=$(git rev-parse HEAD)
	git cherry-pick --keep-redundant-commits --allow-empty "$cp" || true
	rev=$(git rev-parse HEAD)
	if [ x"$rev" != x"$rev0" ]; then
		cps="$cps"-and-"$cp"
	fi
done
cps=${cps##-and-}
sed -i -e '/^\/\/ update-ebitengine-fork\.sh changes:$/,$d' ../aaaaxy/go.mod
if [ -n "$cps" ]; then
	tag="$version"-with-"$cps"
	git tag -f -a -m'update-ebitengine-fork.sh' "$tag"
	git push -f divVerent tag "$tag"
	cat >> ../aaaaxy/go.mod <<EOF

// update-ebitengine-fork.sh changes:
replace github.com/hajimehoshi/ebiten/v2 => github.com/divVerent/ebiten/v2 $tag
EOF
fi
make -C ../aaaaxy mod-update

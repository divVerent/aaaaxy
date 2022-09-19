#!/bin/sh

set -ex

version=$1
cherrypicks='
	47de8027b9bef42a1f0bfb52eb2962126bad8092
'

cd ../ebiten
git fetch
git checkout refs/tags/"$version"
cps=
for cp in $cherrypicks; do
	rev0=$(git rev-parse HEAD)
	git cherry-pick --allow-empty "$cp"
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
	git push -f github tag "$tag"
	cat >> ../aaaaxy/go.mod <<EOF

// update-ebitengine-fork.sh changes:
replace github.com/hajimehoshi/ebiten/v2 => github.com/divVerent/ebiten/v2 $tag
EOF
fi
make -C ../aaaaxy mod-tidy

#!/bin/sh

set -ex

version=$1
cherrypicks='
	b278e5521fc37e3dd9bd1d813aee04d855a36811
	809ad991c278ef43ed04270847536ad1924d8d57
	06c141475c0d85d37221c6055345053f54fb1d6c
'

cd ../ebiten
git fetch
git checkout refs/tags/"$version"
cps=
for cp in $cherrypicks; do
	rev0=$(git rev-parse HEAD)
	git cherry-pick --allow-empty "$cp" || true
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

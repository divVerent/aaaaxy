#!/bin/sh

set -ex

version=$1
cherrypicks='
	574925cf7a72deaf73be4c481348a7a44f7b7e19
	cc247962703eba99eae732876496375191f16cbe
	b96aea70f1559d8a35856bf6a7c814ae168dcce4
	c2038434680951c8ca879613d80665973846fd7d
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

#!/bin/sh

set -ex

package=$1
out=$2

rm -rf "$out"

# Note: ignoring errors here, as some golang.org packages
# do not have a discoverable license file. As they're all under Go's license,
# that is fine.
go get github.com/google/go-licenses
go run github.com/google/go-licenses save github.com/divVerent/aaaaxy/cmd/aaaaxy --save_path="$out" || true

# Add our own third party stuff.
find third_party -name LICENSE -o -name COPYRIGHT.md | while read -r path; do
  mkdir -p "$out/${path%/*}"
  cp "$path" "$out/$path"
done

# List all we got.
find "$out" -type f

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

set -ex

: ${GO:=go}

# Run go natively.
export GOOS=
export GOARCH=

out=$1

rm -rf licenses/software-licenses

# Note: ignoring errors here, as some golang.org packages
# do not have a discoverable license file. As they're all under Go's license,
# that is fine.
$GO run github.com/google/go-licenses save github.com/divVerent/aaaaxy --force --save_path=licenses/software-licenses || true

# This will fail if go-licenses wrote no output.
for d in licenses/software-licenses/*/; do
	[ -d "$d" ]
done

# Add our own third party stuff.
rm -rf licenses/asset-licenses
find third_party -name LICENSE -o -name COPYRIGHT.md | while read -r path; do
  name=${path%/*}
  name=${name##*/}
  mkdir -p "licenses/asset-licenses/$name/"
  cp "$path" "licenses/asset-licenses/$name/"
done

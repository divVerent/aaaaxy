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
target_arch=$(${GO} env GOARCH)
export GOOS=
export GOARCH=

convert \
	-filter Point \
	\( assets/sprites/riser_small_up_0.png -geometry 16x16 \) \
	\( assets/sprites/riser_small_up_0.png -geometry 32x32 \) \
	\( assets/sprites/riser_small_up_0.png -geometry 48x48 \) \
	\( assets/sprites/riser_small_up_0.png -geometry 64x64 \) \
	\( assets/sprites/riser_small_up_0.png -geometry 256x256 \) \
	packaging/aaaaxy.ico

scripts/aaaaxy.exe.manifest.sh $(scripts/version.sh windows) > packaging/aaaaxy.manifest
${GO} run github.com/akavel/rsrc \
	-arch "${target_arch}" \
	-ico packaging/aaaaxy.ico \
	-manifest packaging/aaaaxy.manifest \
	-o aaaaxy.syso

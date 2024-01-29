#!/bin/sh
# Copyright 2024 Google LLC
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

# To be run from the gradle project only.
#
# Works around https://github.com/hajimehoshi/ebiten/issues/2899 (build failure
# due to go-licenses depending on golang.org/x/exp and gomobile not allowing
# that for some reason).

set -ex

git diff --exit-code ../../../internal/builddeps/builddeps.go ../../../go.mod ../../../go.sum

atexit() {
	git checkout ../../../internal/builddeps/builddeps.go ../../../go.mod ../../../go.sum
}
trap atexit EXIT

rm -f ../../../internal/builddeps/builddeps.go
go mod tidy

go run github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile "$@"

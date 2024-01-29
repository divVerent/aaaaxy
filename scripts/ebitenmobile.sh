#!/bin/sh

# To be run from the gradle project only.

set -ex

git diff --exit-code ../../../internal/builddeps/builddeps.go ../../../go.mod ../../../go.sum

atexit() {
	git checkout ../../../internal/builddeps/builddeps.go ../../../go.mod ../../../go.sum
}
trap atexit EXIT

rm -f ../../../internal/builddeps/builddeps.go
go mod tidy

go run github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile "$@"

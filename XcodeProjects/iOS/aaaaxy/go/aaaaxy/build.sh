#!/bin/sh

set -ex

export CGO_ENABLED=1
export CGO_CPPFLAGS=-DNDEBUG
export CGO_CFLAGS='-g0 -O3'
export CGO_CXXFLAGS='-g0 -O3'
export CGO_LDFLAGS='-g0 -O3'
export GOOS=ios
#export AAAAXY_BUILD_USE_VERSION_FILE=true

go generate -tags zip github.com/divVerent/aaaaxy

# TODO combine with ebitenmobile.sh.

git diff --exit-code ../../../../../internal/builddeps/builddeps.go ../../../../../go.mod ../../../../../go.sum

atexit() {
	git checkout ../../../../../internal/builddeps/builddeps.go ../../../../../go.mod ../../../../../go.sum
}
trap atexit EXIT

rm -rf ../../../../../internal/builddeps/builddeps.go
go mod tidy

go run github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile bind \
	-target ios \
	-o aaaaxy.xcframework \
	-iosversion 12.0 \
	-tags zip \
	-gcflags=all=-dwarf=false \
	-ldflags=all=-s \
	-ldflags=all=-w \
	-a \
	-trimpath \
	github.com/divVerent/aaaaxy/XcodeProjects/iOS/aaaaxy/go/aaaaxy

cp ../../../../../aaaaxy.dat ../..

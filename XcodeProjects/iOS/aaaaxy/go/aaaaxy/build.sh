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

../../../../../scripts/ebitenmobile.sh bind \
	-target ios \
	-o aaaaxy.xcframework \
	-iosversion 13.0 \
	-tags zip \
	-gcflags=all=-dwarf=false \
	-ldflags=all='-s -w -buildid=' \
	-a \
	-trimpath \
	-buildvcs=false \
	github.com/divVerent/aaaaxy/XcodeProjects/iOS/aaaaxy/go/aaaaxy

cp ../../../../../aaaaxy.dat ../..

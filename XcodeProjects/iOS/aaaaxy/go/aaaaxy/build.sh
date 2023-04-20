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

version=$(
	 cd ../../../../.. && scripts/version.sh ios
)
sed -i -e "
	s,CURRENT_PROJECT_VERSION = .*;,CURRENT_PROJECT_VERSION = 1;,g;
	s,MARKETING_VERSION = .*;,MARKETING_VERSION = $version;,g;
" ../../../aaaaxy.xcodeproj/project.pbxproj

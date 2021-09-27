# System properties.
GO ?= go
EXE = $(shell $(GO) env GOEXE)
SUFFIX = -$(shell $(GO) env GOOS)-$(shell $(GO) env GOARCH)$(EXE)

# Internal variables.
DUMPCPS = github.com/divVerent/aaaaxy/cmd/dumpcps
VERSION = github.com/divVerent/aaaaxy/internal/version
# TODO glfw is gccgo-built, which still seems to include private paths. Fix.
UPXFLAGS = -9
SOURCES = $(shell git ls-files \*.go)
GENERATED_ASSETS = assets/generated
EMBEDROOT = internal/vfs/_embedroot
EXTRAFILES = README.md LICENSE CONTRIBUTING.md
LICENSES_THIRD_PARTY = licenses
ZIP = 7za -tzip -mx=9 a
CP = cp --reflink=auto
RESOURCE_FILES =

GENERATED_STUFF = $(GENERATED_ASSETS) $(LICENSES_THIRD_PARTY) $(RESOURCE_FILES)

# Output file name when building a release.
ZIPFILE = aaaaxy-$(shell $(GO) env GOOS)-$(shell $(GO) env GOARCH)-$(shell sh scripts/version.sh gittag).zip

# Provide a way to build binaries that are faster at image/video dumping.
# This however makes them slower for normal use, so we're not releasing those.
FASTER_VIDEO_DUMPING = false
ifeq ($(FASTER_VIDEO_DUMPING),true)
BUILDTAGS =
else
BUILDTAGS = ebitensinglethread
endif

# Release/debug flags.
BUILDTYPE = debug
ifeq ($(BUILDTYPE),release)
ifeq ($(shell $(GO) env GOARCH),wasm)
GOFLAGS ?= -tags "embed,$(BUILDTAGS)" -ldflags=all="-s -w" -gcflags=all="-dwarf=false" -trimpath
else
GOFLAGS ?= -tags "embed,$(BUILDTAGS)" -ldflags=all="-s -w" -gcflags=all="-dwarf=false" -trimpath -buildmode=pie
endif
CPPFLAGS ?= -DNDEBUG
CFLAGS ?= -g0 -O3
CXXFLAGS ?= -g0 -O3
LDFLAGS ?= -g0 -s
INFIX =
GENERATED_STUFF += $(EMBEDROOT)
else
ifeq ($(BUILDTYPE),extradebug)
GOFLAGS ?= -tags "ebitendebug,$(BUILDTAGS)"
INFIX = -extradebug
else
GOFLAGS ?= -tags "$(BUILDTAGS)"
INFIX = -debug
endif
endif
BINARY = $(PREFIX)aaaaxy$(INFIX)$(SUFFIX)

# Windows only: include icon.
ifeq ($(shell $(GO) env GOOS),windows)
RESOURCE_FILES = resources.ico resources.manifest resources.syso
endif

# OS X only: app bundle.
APPBUNDLE = AAAAXY.app
APPICONS = $(APPBUNDLE)/Contents/Resources/icon.icns

# OS X releases only: app bundle.
ifeq ($(shell $(GO) env GOOS),darwin)
ifeq ($(BUILDTYPE),release)
PREFIX = $(APPBUNDLE)/Contents/MacOS/
endif
endif

# Include version.
GOFLAGS += -ldflags="-X $(VERSION).revision=$(shell scripts/version.sh semver)"

# cgo support.
CGO_CPPFLAGS ?= $(CPPFLAGS)
CGO_CFLAGS ?= $(CFLAGS)
CGO_CXXFLAGS ?= $(CXXFLAGS)
CGO_LDFLAGS ?= $(LDFLAGS)
CGO_ENV = \
	CGO_CPPFLAGS="$(CGO_CPPFLAGS)" \
	CGO_CFLAGS="$(CGO_CFLAGS)" \
	CGO_CXXFLAGS="$(CGO_CXXFLAGS)" \
	CGO_LDFLAGS="$(CGO_LDFLAGS)"

# OS X cross compile support.
ifeq ($(shell $(GO) env GOOS),darwin)
ifeq ($(shell uname),Linux)
CGO_ENV += PATH=$(HOME)/src/osxcross-sdk/bin:$(PATH)
CGO_ENV += CGO_ENABLED=1
CGO_ENV += CC=o64-clang
CGO_ENV += CXX=o64-clang++
CGO_ENV += MACOSX_DEPLOYMENT_TARGET=10.12
endif
endif

.PHONY: all
all: bin

.PHONY: bin
bin: $(BINARY)

.PHONY: clean
clean:
	$(RM) -r $(BINARY) $(GENERATED_STUFF)

.PHONY: vet
vet:
	$(GO) vet `find ./cmd ./internal -name \*.go -print | sed -e 's,/[^/]*$$,,' | sort -u`

# The actual build process follows.

# Packing the data files.
.PHONY: generate
generate:
	$(GO) generate $(GOFLAGS)

# Binary building.
$(BINARY): generate $(SOURCES)
	$(CGO_ENV) \
	$(GO) build -o $(BINARY) $(GOFLAGS)

# Binary release building.

.PHONY: webprepare
webprepare:
	cp $(shell $(GO) env GOROOT)/misc/wasm/wasm_exec.js .

.PHONY: releaseclean
releaseclean:
	$(RM) $(ZIPFILE)

.PHONY: binrelease
binrelease: releaseclean $(BINARY) $(EXTRAFILES) $(LICENSES_THIRD_PARTY)
	$(ZIP) $(ZIPFILE) $(BINARY) $(EXTRAFILES) $(LICENSES_THIRD_PARTY)

.PHONY: osxbinrelease
osxbinrelease: releaseclean $(BINARY) $(EXTRAFILES) $(LICENSES_THIRD_PARTY)
	$(ZIP) $(ZIPFILE) $(APPBUNDLE) $(EXTRAFILES) $(LICENSES_THIRD_PARTY)

.PHONY: webbinrelease
webbinrelease: releaseclean webprepare $(BINARY) $(EXTRAFILES) $(LICENSES_THIRD_PARTY)
	$(ZIP) $(ZIPFILE) $(BINARY) $(EXTRAFILES) $(LICENSES_THIRD_PARTY) aaaaxy$(INFIX).html wasm_exec.js

.PHONY: allrelease
allrelease:
	GO="$(GO)" GOOS=linux GOARCH=amd64 $(MAKE) binrelease BUILDTYPE=release
	GO="$(GO)" GOOS=windows GOARCH=amd64 $(MAKE) binrelease BUILDTYPE=release
	GO="$(GO)" GOOS=windows GOARCH=386 $(MAKE) binrelease BUILDTYPE=release
	GO="$(GO)" GOOS=darwin GOARCH=amd64 $(MAKE) osxbinrelease BUILDTYPE=release

.PHONY: webdebug
webdebug: webprepare
	GO="$(GO)" GOOS=js GOARCH=wasm $(MAKE) EXE=.wasm debug

.PHONY: webrelease
webrelease: webprepare
	GO="$(GO)" GOOS=js GOARCH=wasm $(MAKE) EXE=.wasm release

# Debugging.
assets/generated/%.cp.pdf: assets/generated/%.cp.dot
	neato -Tpdf $< > $@

# Helper targets.
.PHONY: run
run: $(BINARY)
	EBITEN_INTERNAL_IMAGES_KEY=i EBITEN_SCREENSHOT_KEY=p ./$(BINARY) $(ARGS)

# Prepare git hooks.
.PHONY: setup-git
setup-git:
	git config filter.git-clean-tmx.clean "$$PWD"/scripts/git-clean-tmx.sh
	git config filter.git-clean-md.clean "$$PWD"/scripts/git-clean-md.sh


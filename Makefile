# System properties.
GO ?= go
EXE = $(shell $(GO) env GOEXE)
SUFFIX = -$(shell $(GO) env GOOS)-$(shell $(GO) env GOARCH)$(EXE)

# Internal variables.
VERSION = github.com/divVerent/aaaaxy/internal/version
SOURCES = $(shell git ls-files \*.go)
GENERATED_ASSETS = assets/generated
EMBEDROOT = internal/vfs/_embedroot
EXTRAFILES = README.md LICENSE CONTRIBUTING.md
GENERATED_STUFF = assets/generated licenses

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

# OS X only: app bundle.
APPBUNDLE = AAAAXY.app

# OS X releases only: app bundle.
ifeq ($(shell $(GO) env GOOS),darwin)
ifeq ($(BUILDTYPE),release)
PREFIX = $(APPBUNDLE)/Contents/MacOS/
endif
endif

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
	$(GO) vet $(SOURCES)

# The actual build process follows.

# Packing the data files.
.PHONY: generate
generate:
	$(GO) generate $(GOFLAGS)

# Binary building.
$(BINARY): generate $(SOURCES)
	$(CGO_ENV) $(GO) build -o $(BINARY) $(GOFLAGS)

# Helper targets.
.PHONY: run
run: $(BINARY)
	EBITEN_INTERNAL_IMAGES_KEY=i EBITEN_SCREENSHOT_KEY=p ./$(BINARY) $(ARGS)

# Prepare git hooks.
.PHONY: setup-git
setup-git:
	git config filter.git-clean-tmx.clean "$$PWD"/scripts/git-clean-tmx.sh
	git config filter.git-clean-md.clean "$$PWD"/scripts/git-clean-md.sh


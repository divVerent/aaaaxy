# Settings.
BUILDTYPE = debug
FASTER_VIDEO_DUMPING = false

# System properties.
GO ?= go
EXE = $(shell $(GO) env GOEXE)
SUFFIX = -$(shell $(GO) env GOOS)-$(shell $(GO) env GOARCH)$(EXE)
BINARY = $(PREFIX)aaaaxy$(INFIX)$(SUFFIX)

ifeq ($(BUILDTYPE),release)
INFIX =
BUILDTAGS = embed
else
INFIX = -debug
BUILDTAGS =
endif

# Provide a way to build binaries that are faster at image/video dumping.
# This however makes them slower for normal use, so we're not releasing those.
FASTER_VIDEO_DUMPING = false
ifeq ($(FASTER_VIDEO_DUMPING),false)
BUILDTAGS += ebitensinglethread
endif

# Internal variables.
SOURCES = $(shell git ls-files \*.go)
GENERATED_STUFF = aaaaxy.ico aaaaxy.manifest aaaaxy.syso assets/generated/ internal/vfs/_embedroot/ licenses/

# Configure Go.
GOFLAGS += -tags "$(shell echo $(BUILDTAGS) | tr ' ' ,)"

# Release/debug flags.
BUILDTYPE = debug
ifeq ($(BUILDTYPE),release)
GOFLAGS += -a -ldflags=all="-s -w" -gcflags=all="-dwarf=false" -trimpath
ifneq ($(shell $(GO) env GOARCH),wasm)
GOFLAGS += -buildmode=pie
endif
CPPFLAGS ?= -DNDEBUG
CFLAGS ?= -g0 -O3
CXXFLAGS ?= -g0 -O3
LDFLAGS ?= -g0 -s
endif

# cgo support.
CGO_CPPFLAGS ?= $(CPPFLAGS)
CGO_CFLAGS ?= $(CFLAGS)
CGO_CXXFLAGS ?= $(CXXFLAGS)
CGO_LDFLAGS ?= $(LDFLAGS)
CGO_ENV ?= \
	CGO_CPPFLAGS="$(CGO_CPPFLAGS)" \
	CGO_CFLAGS="$(CGO_CFLAGS)" \
	CGO_CXXFLAGS="$(CGO_CXXFLAGS)" \
	CGO_LDFLAGS="$(CGO_LDFLAGS)"

.PHONY: all
all: bin

.PHONY: bin
bin: $(BINARY)

.PHONY: clean
clean:
	$(RM) -r $(BINARY) $(GENERATED_STUFF)

.PHONY: vet
vet:
	$(GO) vet ./...

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


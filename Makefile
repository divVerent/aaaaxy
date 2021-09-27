# System properties.
GO ?= go
EXE = $(shell $(GO) env GOEXE)
SUFFIX = -$(shell $(GO) env GOOS)-$(shell $(GO) env GOARCH)$(EXE)

# Internal variables.
PACKAGE = github.com/divVerent/aaaaxy
DUMPCPS = github.com/divVerent/aaaaxy/cmd/dumpcps
VERSION = github.com/divVerent/aaaaxy/internal/version
# TODO glfw is gccgo-built, which still seems to include private paths. Fix.
UPXFLAGS = -9
SOURCES = $(shell git ls-files \*.go)
GENERATED_ASSETS = assets/generated/level.cp.json assets/generated/image_load_order.txt
EMBEDROOT = internal/vfs/_embedroot
EXTRAFILES = README.md LICENSE CONTRIBUTING.md
LICENSES_THIRD_PARTY = licenses
ZIP = 7za -tzip -mx=9 a
CP = cp --reflink=auto
RESOURCE_FILES =

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
BINARY_ASSETS = $(EMBEDROOT)
else
ifeq ($(BUILDTYPE),extradebug)
GOFLAGS ?= -tags "ebitendebug,$(BUILDTAGS)"
INFIX = -extradebug
else
GOFLAGS ?= -tags "$(BUILDTAGS)"
INFIX = -debug
endif
BINARY_ASSETS = $(GENERATED_ASSETS)
endif
BINARY = $(PREFIX)aaaaxy$(INFIX)$(SUFFIX)

# Windows only: include icon.
ifeq ($(shell $(GO) env GOOS),windows)
RESOURCE_FILES = cmd/aaaaxy/resources.ico cmd/aaaaxy/resources.manifest cmd/aaaaxy/resources.syso
BINARY_ASSETS += $(RESOURCE_FILES)
endif

# OS X only: app bundle.
APPBUNDLE = AAAAXY.app
APPICONS = $(APPBUNDLE)/Contents/Resources/icon.icns

# OS X releases only: app bundle.
ifeq ($(shell $(GO) env GOOS),darwin)
ifeq ($(BUILDTYPE),release)
BINARY_ASSETS += $(APPICONS)
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
	$(RM) -r $(BINARY) $(GENERATED_ASSETS) $(LICENSES_THIRD_PARTY) $(RESOURCE_FILES)

.PHONY: vet
vet:
	$(GO) vet `find ./cmd ./internal -name \*.go -print | sed -e 's,/[^/]*$$,,' | sort -u`

# The actual build process follows.

# Packing the data files.
assets/generated/image_load_order.txt: assets/tiles assets/sprites $(wildcard third_party/*/assets/sprites)
	mkdir -p assets/generated
	scripts/image-load-order.sh $^ > $@

assets/generated/%.cp.dot: assets/maps/%.tmx cmd/dumpcps/main.go
	mkdir -p assets/generated
	GOOS= GOARCH= $(GO) run $(DUMPCPS) $< > $@

assets/generated/%.cp.json: assets/generated/%.cp.dot
	neato -Tjson $< > $@

.PHONY: $(LICENSES_THIRD_PARTY)
$(LICENSES_THIRD_PARTY):
	GO="$(GO)" GOOS= GOARCH= scripts/collect-licenses.sh $(PACKAGE) $(LICENSES_THIRD_PARTY)

$(EMBEDROOT): $(GENERATED_ASSETS) $(LICENSES_THIRD_PARTY)
	$(RM) -r $(EMBEDROOT)
	CP="$(CP)" scripts/build-vfs.sh $(EMBEDROOT)

# Windows-only stuff.
cmd/aaaaxy/resources.manifest: scripts/aaaaxy.exe.manifest.sh
	$< $(shell scripts/version.sh windows) > $@

%.syso: %.ico %.manifest
	GOOS= GOARCH= $(GO) get -d github.com/akavel/rsrc
	GOOS= GOARCH= $(GO) install github.com/akavel/rsrc
	GOOS= GOARCH= $(GO) run github.com/akavel/rsrc \
		-arch $(shell $(GO) env GOARCH) \
		-ico $*.ico \
		-manifest $*.manifest \
		-o $@

cmd/aaaaxy/resources.ico: assets/sprites/riser_small_up_0.png
	convert \
		-filter Point \
		\( $< -geometry 16x16 \) \
		\( $< -geometry 32x32 \) \
		\( $< -geometry 48x48 \) \
		\( $< -geometry 64x64 \) \
		\( $< -geometry 256x256 \) \
		$@

$(APPICONS): assets/sprites/riser_small_up_0.png
	for res in 16 32 128 256 512 1024; do \
		convert $< \
			-filter Point -geometry $${res}x$${res} \
			-define png:bit-depth=8 \
			-define png:color-type=6 \
			-define png:format=png32 \
			icns$${res}.png; \
	done
	png2icns $@ icns*.png
	$(RM) icns*.png

# Binary building.
$(BINARY): $(BINARY_ASSETS) $(SOURCES)
	$(CGO_ENV) \
	$(GO) build -o $(BINARY) $(GOFLAGS) $(PACKAGE)

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


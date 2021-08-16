# System properties.
GO ?= go
EXE = $(shell $(GO) env GOEXE)
SUFFIX = -$(shell $(GO) env GOOS)-$(shell $(GO) env GOARCH)$(EXE)

# Internal variables.
PACKAGE = github.com/divVerent/aaaaxy/cmd/aaaaxy
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
GOFLAGS ?= -tags embed,$(BUILDTAGS) -ldflags=all="-s -w" -gcflags=all="-dwarf=false" -trimpath
else
GOFLAGS ?= -tags embed,$(BUILDTAGS) -ldflags=all="-s -w" -gcflags=all="-B -dwarf=false" -trimpath -buildmode=pie
endif
CPPFLAGS ?= -DNDEBUG
CFLAGS ?= -g0 -O3
CXXFLAGS ?= -g0 -O3
ifeq ($(shell $(GO) env GOOS),windows)
LDFLAGS ?= -g0 -s -static-libgcc
else
LDFLAGS ?= -g0 -s
endif
INFIX =
BINARY_ASSETS = $(EMBEDROOT)
else
ifeq ($(BUILDTYPE),extradebug)
GOFLAGS ?= -tags ebitendebug,$(BUILDTAGS)
INFIX = -extradebug
else
GOFLAGS ?= -tags $(BUILDTAGS)
INFIX = -debug
endif
BINARY_ASSETS = $(GENERATED_ASSETS)
endif
BINARY = aaaaxy$(INFIX)$(SUFFIX)

# Windows only: include icon.
ifeq ($(shell $(GO) env GOOS),windows)
RESOURCE_FILES = cmd/aaaaxy/resources.ico cmd/aaaaxy/resources.manifest cmd/aaaaxy/resources.syso
BINARY_ASSETS += $(RESOURCE_FILES)
endif

# Include version.
GOFLAGS += -ldflags="-X $(VERSION).revision=$(shell scripts/version.sh semver)"

# cgo support.
CGO_CPPFLAGS ?= $(CPPFLAGS)
CGO_CFLAGS ?= $(CFLAGS)
CGO_CXXFLAGS ?= $(CXXFLAGS)
CGO_LDFLAGS ?= $(LDFLAGS)

.PHONY: bin
bin: $(BINARY)

.PHONY: all
all: debug release

.PHONY: debug
debug:
	$(MAKE) BUILDTYPE=debug bin

.PHONY: release
release:
	$(MAKE) BUILDTYPE=release bin

.PHONY: clean
clean:
	$(RM) -r $(BINARY) $(GENERATED_ASSETS) $(LICENSES_THIRD_PARTY) $(RESOURCE_FILES)

.PHONY: vet
vet:
	$(GO) vet `find ./cmd ./internal -name \*.go -print | sed -e 's,/[^/]*$$,,' | sort -u`

.PHONY: $(EMBEDROOT)
$(EMBEDROOT): $(GENERATED_ASSETS) $(LICENSES_THIRD_PARTY)
	$(RM) -r $(EMBEDROOT)
	CP="$(CP)" scripts/build-vfs.sh $(EMBEDROOT)

cmd/aaaaxy/resources.manifest: scripts/aaaaxy.exe.manifest.sh
	$< $(shell scripts/version.sh windows) > $@

%.syso: %.ico %.manifest
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

$(BINARY): $(BINARY_ASSETS) $(SOURCES)
	CGO_CPPFLAGS="$(CGO_CPPFLAGS)" \
	CGO_CFLAGS="$(CGO_CFLAGS)" \
	CGO_CXXFLAGS="$(CGO_CXXFLAGS)" \
	CGO_LDFLAGS="$(CGO_LDFLAGS)" \
	$(GO) build -o $(BINARY) $(GOFLAGS) $(PACKAGE)

assets/generated/image_load_order.txt: assets/tiles assets/sprites $(wildcard third_party/*/assets/sprites)
	mkdir -p assets/generated
	scripts/image_load_order.sh $^ > $@

assets/generated/%.cp.json: assets/generated/%.cp.dot
	neato -Tjson $< > $@

assets/generated/%.cp.pdf: assets/generated/%.cp.dot
	neato -Tpdf $< > $@

assets/generated/%.cp.dot: assets/maps/%.tmx cmd/dumpcps/main.go
	mkdir -p assets/generated
	GOOS= GOARCH= $(GO) run $(DUMPCPS) $< > $@

.PHONY: $(LICENSES_THIRD_PARTY)
$(LICENSES_THIRD_PARTY):
	GO="$(GO)" GOOS= GOARCH= scripts/collect-licenses.sh $(PACKAGE) $(LICENSES_THIRD_PARTY)

# Building of release zip files starts here.
ZIPFILE = aaaaxy.zip

.PHONY: addextras
addextras: $(EXTRAFILES)
	$(ZIP) $(ZIPFILE) $(EXTRAFILES)

.PHONY: addlicenses
addlicenses: $(LICENSES_THIRD_PARTY)
	$(ZIP) $(ZIPFILE) $(LICENSES_THIRD_PARTY)

.PHONY: addrelease
addrelease: $(BINARY)
	$(ZIP) $(ZIPFILE) $(BINARY)
	$(MAKE) clean

.PHONY: webprepare
webprepare:
	cp $(shell $(GO) env GOROOT)/misc/wasm/wasm_exec.js .

.PHONY: addwebstuff
addwebstuff: webprepare
	$(ZIP) $(ZIPFILE) aaaaxy$(INFIX).html wasm_exec.js

.PHONY: allrelease
allrelease: allreleaseclean
	$(RM) $(ZIPFILE)
	$(MAKE) addextras
	$(MAKE) addlicenses
	GO="$(GO)" GOOS=linux GOARCH=amd64 $(MAKE) BUILDTYPE=release addrelease
	GO="$(GO)" GOOS=windows GOARCH=386 $(MAKE) BUILDTYPE=release addrelease
	GO="$(GO)" GOOS=windows GOARCH=amd64 $(MAKE) BUILDTYPE=release addrelease
	# Disabled because build is WAY too slow to be playable.
	# $(MAKE) BUILDTYPE=release addwebstuff
	# GO="$(GO)" GOOS=js GOARCH=wasm $(MAKE) EXE=.wasm BUILDTYPE=release addrelease

.PHONY: webdebug
webdebug: webprepare
	GO="$(GO)" GOOS=js GOARCH=wasm $(MAKE) EXE=.wasm debug

.PHONY: webrelease
webrelease: webprepare
	GO="$(GO)" GOOS=js GOARCH=wasm $(MAKE) EXE=.wasm release

.PHONY: allreleaseclean
allreleaseclean:
	GO="$(GO)" GOOS=linux GOARCH=amd64 $(MAKE) clean
	GO="$(GO)" GOOS=windows GOARCH=amd64 $(MAKE) clean
	$(RM) $(ZIPFILE)

# Helper targets.
.PHONY: run
run: $(BINARY)
	EBITEN_INTERNAL_IMAGES_KEY=i EBITEN_SCREENSHOT_KEY=p ./$(BINARY) $(ARGS)

.PHONY: setup-git
setup-git:
	git config filter.git-clean-tmx.clean "$$PWD"/scripts/git-clean-tmx.sh
	git config filter.git-clean-md.clean "$$PWD"/scripts/git-clean-md.sh


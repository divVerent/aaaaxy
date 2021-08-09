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
STATIK_ASSETS_ROOT = internal/assets
STATIK_ASSETS = $(STATIK_ASSETS_ROOT)/statik/statik.go
EXTRAFILES = README.md LICENSE CONTRIBUTING.md
LICENSES_THIRD_PARTY = licenses
ZIP = 7za -tzip -mx=9 a

# Release/debug flags.
BUILDTYPE = debug
ifeq ($(BUILDTYPE),release)
ifeq ($(GOARCH),wasm)
GOFLAGS ?= -tags embed,ebitensinglethread -ldflags=all="-s -w" -gcflags=all="-dwarf=false" -trimpath
else
GOFLAGS ?= -tags embed,ebitensinglethread -ldflags=all="-s -w" -gcflags=all="-B -dwarf=false" -trimpath -buildmode=pie
endif
CPPFLAGS ?= -DNDEBUG
CFLAGS ?= -g0 -O3
CXXFLAGS ?= -g0 -O3
LDFLAGS ?= -g0 -s
INFIX =
BINARY_ASSETS = $(EMBEDROOT)
else
ifeq ($(BUILDTYPE),statikrelease)
ifeq ($(GOARCH),wasm)
GOFLAGS ?= -tags statik,ebitensinglethread -ldflags=all="-s -w" -gcflags=all="-dwarf=false" -trimpath
else
GOFLAGS ?= -tags statik,ebitensinglethread -ldflags=all="-s -w" -gcflags=all="-B -dwarf=false" -trimpath -buildmode=pie
endif
CPPFLAGS ?= -DNDEBUG
CFLAGS ?= -g0 -O3
CXXFLAGS ?= -g0 -O3
LDFLAGS ?= -g0 -s
INFIX = statik
BINARY_ASSETS = $(STATIK_ASSETS)
else
ifeq ($(BUILDTYPE),extradebug)
GOFLAGS ?= -tags ebitensinglethread,ebitendebug
INFIX = -extradebug
else
GOFLAGS ?= -tags ebitensinglethread
INFIX = -debug
endif
BINARY_ASSETS = $(GENERATED_ASSETS)
endif
endif
BINARY = aaaaxy$(INFIX)$(SUFFIX)

# Include version.
GOFLAGS += -ldflags="-X $(VERSION).revision=$(shell git describe --always --dirty --first-parent)"

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
	$(RM) -r $(BINARY) $(STATIK_ASSETS) $(GENERATED_ASSETS) $(LICENSES_THIRD_PARTY)

.PHONY: vet
vet:
	$(GO) vet `find ./cmd ./internal -name \*.go -print | sed -e 's,/[^/]*$$,,' | sort -u`

.PHONY: $(EMBEDROOT)
$(EMBEDROOT): $(GENERATED_ASSETS $(LICENSES_THIRD_PARTY)
	$(RM) -r $(EMBEDROOT)
	scripts/build-vfs.sh $(EMBEDROOT) cp

.PHONY: $(STATIK_ASSETS)
$(STATIK_ASSETS): $(GENERATED_ASSETS) $(LICENSES_THIRD_PARTY)
	GO=$(GO) GOOS= GOARCH= scripts/statik-vfs.sh $(STATIK_ASSETS_ROOT)

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
	GO=$(GO) GOOS= GOARCH= scripts/collect-licenses.sh $(PACKAGE) $(LICENSES_THIRD_PARTY)

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
	GO=$(GO) GOOS=linux GOARCH=amd64 $(MAKE) BUILDTYPE=release addrelease
	# Disabled due to Windows Defender FP:
	# GOOS=windows GOARCH=386 $(MAKE) release
	GO=$(GO) GOOS=windows GOARCH=amd64 $(MAKE) BUILDTYPE=release addrelease
	# Disabled because build is WAY too slow to be playable.
	# $(MAKE) BUILDTYPE=release addwebstuff
	# GO=$(GO) GOOS=js GOARCH=wasm $(MAKE) EXE=.wasm BUILDTYPE=release addrelease

.PHONY: webdebug
webdebug: webprepare
	GO=$(GO) GOOS=js GOARCH=wasm $(MAKE) EXE=.wasm debug

.PHONY: webrelease
webrelease: webprepare
	GO=$(GO) GOOS=js GOARCH=wasm $(MAKE) EXE=.wasm release

.PHONY: allreleaseclean
allreleaseclean:
	GO=$(GO) GOOS=linux GOARCH=amd64 $(MAKE) clean
	GO=$(GO) GOOS=windows GOARCH=amd64 $(MAKE) clean
	$(RM) $(ZIPFILE)

# Helper targets.
.PHONY: run
run: $(BINARY)
	EBITEN_INTERNAL_IMAGES_KEY=i EBITEN_SCREENSHOT_KEY=p ./$(BINARY) $(ARGS)

.PHONY: setup-git
setup-git:
	git config filter.git-clean-tmx.clean "$$PWD"/scripts/git-clean-tmx.sh


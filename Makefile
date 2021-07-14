# System properties.
EXE = $(shell go env GOEXE)
SUFFIX = -$(shell go env GOOS)-$(shell go env GOARCH)$(EXE)

# Internal variables.
PACKAGE = github.com/divVerent/aaaaxy/cmd/aaaaxy
DUMPCPS = github.com/divVerent/aaaaxy/cmd/dumpcps
VERSION = github.com/divVerent/aaaaxy/internal/version
# TODO glfw is gccgo-built, which still seems to include private paths. Fix.
UPXFLAGS = -9
SOURCES = $(shell git ls-files \*.go)
GENERATED_ASSETS = assets/maps/level.cp.json assets/generated/image_load_order.txt
STATIK_ASSETS_ROOT = internal/assets
STATIK_ASSETS = $(STATIK_ASSETS_ROOT)/statik/statik.go
EXTRAFILES = README.md LICENSE CONTRIBUTING.md
LICENSES_THIRD_PARTY = licenses
ZIP = 7za -tzip -mx=9 a

# Release/debug flags.
BUILDTYPE = debug
ifeq ($(BUILDTYPE),release)
ifeq ($(GOARCH),wasm)
GOFLAGS ?= -tags statik,ebitensinglethread -ldflags=all="-s -w" -gcflags=all="-dwarf=false" -trimpath
else
GOFLAGS ?= -tags statik,ebitensinglethread -ldflags=all="-s -w" -gcflags=all="-B -dwarf=false" -trimpath -buildmode=pie
endif
CPPFLAGS ?= -DNDEBUG
CFLAGS ?= -g0 -O3
CXXFLAGS ?= -g0 -O3
LDFLAGS ?= -g0 -s
INFIX =
BINARY_ASSETS = $(STATIK_ASSETS)
else
GOFLAGS ?= -tags ebitensinglethread
INFIX = -debug
BINARY_ASSETS = $(GENERATED_ASSETS)
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
	go vet `find ./cmd ./internal -name \*.go -print | sed -e 's,/[^/]*$$,,' | sort -u`

.PHONY: $(STATIK_ASSETS)
$(STATIK_ASSETS): $(GENERATED_ASSETS) $(LICENSES_THIRD_PARTY)
	GOOS= GOARCH= scripts/statik-vfs.sh $(STATIK_ASSETS_ROOT)

$(BINARY): $(BINARY_ASSETS) $(SOURCES)
	CGO_CPPFLAGS="$(CGO_CPPFLAGS)" \
	CGO_CFLAGS="$(CGO_CFLAGS)" \
	CGO_CXXFLAGS="$(CGO_CXXFLAGS)" \
	CGO_LDFLAGS="$(CGO_LDFLAGS)" \
	go build -o $(BINARY) $(GOFLAGS) $(PACKAGE)

# Note: we can't detect whether this needs to be regenerated, so for now let's always do it.
.PHONY: assets/generated/image_load_order.txt
assets/generated/image_load_order.txt:
	mkdir -p assets/generated
	scripts/image_load_order.sh > $@

%.cp.json: %.cp.dot
	neato -Tjson $< > $@

%.cp.pdf: %.cp.dot
	neato -Tpdf $< > $@

%.cp.dot: %.tmx cmd/dumpcps/main.go
	GOOS= GOARCH= go run $(DUMPCPS) $< > $@

.PHONY: $(LICENSES_THIRD_PARTY)
$(LICENSES_THIRD_PARTY):
	GOOS= GOARCH= scripts/collect-licenses.sh $(PACKAGE) $(LICENSES_THIRD_PARTY)

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
	cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js .

.PHONY: addwebstuff
addwebstuff: webprepare
	$(ZIP) $(ZIPFILE) aaaaxy$(INFIX).html wasm_exec.js

.PHONY: allrelease
allrelease: allreleaseclean
	$(RM) $(ZIPFILE)
	$(MAKE) addextras
	$(MAKE) addlicenses
	GOOS=linux GOARCH=amd64 $(MAKE) BUILDTYPE=release addrelease
	# Disabled due to Windows Defender FP:
	# GOOS=windows GOARCH=386 $(MAKE) release
	GOOS=windows GOARCH=amd64 $(MAKE) BUILDTYPE=release addrelease
	# Disabled because build is WAY too slow to be playable.
	# $(MAKE) BUILDTYPE=release addwebstuff
	# GOOS=js GOARCH=wasm $(MAKE) EXE=.wasm BUILDTYPE=release addrelease

.PHONY: webdebug
webdebug: webprepare
	GOOS=js GOARCH=wasm $(MAKE) EXE=.wasm debug

.PHONY: webrelease
webrelease: webprepare
	GOOS=js GOARCH=wasm $(MAKE) EXE=.wasm release

.PHONY: allreleaseclean
allreleaseclean:
	GOOS=linux GOARCH=amd64 $(MAKE) clean
	GOOS=windows GOARCH=amd64 $(MAKE) clean
	$(RM) $(ZIPFILE)

# Helper targets.
.PHONY: run
run: $(BINARY)
	./$(BINARY) $(ARGS)

.PHONY: setup-git
setup-git:
	git config filter.git-clean-tmx.clean "$$PWD"/scripts/git-clean-tmx.sh


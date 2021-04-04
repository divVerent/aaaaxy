# System properties.
EXE = $(shell go env GOEXE)
SUFFIX = -$(shell go env GOOS)-$(shell go env GOARCH)$(EXE)

# Internal variables.
PACKAGE = github.com/divVerent/aaaaaa/cmd/aaaaaa
DUMPCPS = github.com/divVerent/aaaaaa/cmd/dumpcps
# TODO glfw is gccgo-built, which still seems to include private paths. Fix.
UPXFLAGS = -9
SOURCES = $(shell git ls-files \*.go)
GENERATED_ASSETS = assets/maps/level.cp.json
STATIK_ASSETS_ROOT = internal/assets
STATIK_ASSETS = $(STATIK_ASSETS_ROOT)/statik/statik.go
EXTRAFILES = README.md LICENSE CONTRIBUTING.md
LICENSES_THIRD_PARTY = licenses
ZIP = 7za -tzip -mx=9 a

# Release/debug flags.
ifeq ($(BUILDTYPE),release)
GOFLAGS ?= -tags statik -ldflags=all="-s -w -linkmode=external" -gcflags=all="-B -dwarf=false" -trimpath -buildmode=pie
CFLAGS ?= -g0
LDFLAGS ?= -g0 -s
BINARY = aaaaaa$(SUFFIX)
BINARY_ASSETS = $(STATIK_ASSETS)
else
BINARY = aaaaaa-debug$(SUFFIX)
BINARY_ASSETS = $(GENERATED_ASSETS)
endif

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
$(STATIK_ASSETS): $(GENERATED_ASSETS)
	GOOS= GOARCH= ./statik-vfs.sh $(STATIK_ASSETS_ROOT)

$(BINARY): $(BINARY_ASSETS) $(SOURCES)
	CGO_CPPFLAGS="$(CGO_CPPFLAGS)" \
	CGO_CFLAGS="$(CGO_CFLAGS)" \
	CGO_CXXFLAGS="$(CGO_CXXFLAGS)" \
	CGO_LDFLAGS="$(CGO_LDFLAGS)" \
	go build -o $(BINARY) $(GOFLAGS) $(PACKAGE)

%.cp.json: %.cp.dot
	neato -Tjson $< > $@

%.cp.pdf: %.cp.dot
	neato -Tpdf $< > $@

%.cp.dot: %.tmx cmd/dumpcps/main.go
	GOOS= GOARCH= go run $(DUMPCPS) $< > $@

.PHONY: $(LICENSES_THIRD_PARTY)
$(LICENSES_THIRD_PARTY):
	GOOS= GOARCH= ./collect-licenses.sh $(PACKAGE) $(LICENSES_THIRD_PARTY)

# Building of release zip files starts here.
ZIPFILE = aaaaaa.zip

.PHONY: addextras
addextras: $(EXTRAFILES)
	$(ZIP) $(ZIPFILE) $(EXTRAFILES)

.PHONY: addlicenses
addlicenses: $(LICENSES_THIRD_PARTY)
	$(ZIP) $(ZIPFILE) $(LICENSES_THIRD_PARTY)

.PHONY: addrelease
addrelease: $(RELEASE)
	$(ZIP) $(ZIPFILE) $(RELEASE)
	$(MAKE) clean

.PHONY: allrelease
allrelease: allreleaseclean
	$(RM) $(ZIPFILE)
	$(MAKE) addextras
	$(MAKE) addlicenses
	GOOS=linux GOARCH=amd64 $(MAKE) addrelease
	# Disabled due to Windows Defender FP:
	# GOOS=windows GOARCH=386 $(MAKE) release
	GOOS=windows GOARCH=amd64 $(MAKE) addrelease

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
	git config filter.git-clean-tmx.clean "$$PWD"/git-clean-tmx.sh


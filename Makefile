# System properties.
EXE = $(shell go env GOEXE)
SUFFIX = -$(shell go env GOOS)-$(shell go env GOARCH)$(EXE)

# Internal variables.
PACKAGE = github.com/divVerent/aaaaaa/cmd/aaaaaa
DUMPCPS = github.com/divVerent/aaaaaa/cmd/dumpcps
ASSETS = internal/assets
DEBUG = aaaaaa-debug$(SUFFIX)
DEBUG_GOFLAGS =
RELEASE = aaaaaa$(SUFFIX)
RELEASE_GOFLAGS = -ldflags="-s -w" -gcflags="-B -dwarf=false"
UPXFLAGS = -9
SOURCES = $(shell find . -name \*.go)
GENERATED_ASSETS = assets/maps/level.cp.json

.PHONY: default
default: debug

.PHONY: all
all: debug release

.PHONY: debug
debug: $(DEBUG)

.PHONY: release
release: $(RELEASE)

.PHONY: clean
clean:
	$(RM) -r $(DEBUG) $(RELEASE) $(ASSETS) $(GENERATED_ASSETS)

.PHONY: vet
vet:
	go vet `find ./cmd ./internal -name \*.go -print | sed -e 's,/[^/]*$$,,' | sort -u`

.PHONY: $(ASSETS)
$(ASSETS): $(GENERATED_ASSETS)
	GOOS= GOARCH= ./statik-vfs.sh $(ASSETS)

$(DEBUG): $(GENERATED_ASSETS) $(SOURCES)
	go build -o $(DEBUG) $(DEBUG_GOFLAGS) $(PACKAGE)

$(RELEASE): $(ASSETS) $(SOURCES)
	go build -tags statik -o $(RELEASE) $(RELEASE_GOFLAGS) $(PACKAGE)

%.cp.json: %.cp.dot
	neato -Tjson $< > $@

%.cp.pdf: %.cp.dot
	neato -Tpdf $< > $@

%.cp.dot: %.tmx cmd/dumpcps/main.go
	GOOS= GOARCH= go run $(DUMPCPS) $< > $@

# Building of release zip files starts here.
ZIPFILE = aaaaaa.zip

.PHONY: addrelease
addrelease: $(RELEASE)
	zip -9r $(ZIPFILE) $(RELEASE)
	$(MAKE) clean

.PHONY: allrelease
allrelease: allreleaseclean
	$(RM) $(ZIPFILE)
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
run: $(DEBUG)
	$(DEBUG)

.PHONY: setup-git
setup-git:
	git config filter.git-clean-tmx.clean "$$PWD"/git-clean-tmx.sh


# User configrable.
NO_VFS = false

# System properties.
EXE = $(shell go env GOEXE)
SUFFIX = -$(shell go env GOOS)-$(shell go env GOARCH)$(EXE)

# Internal variables.
PACKAGE = github.com/divVerent/aaaaaa/cmd/aaaaaa
ASSETS = internal/assets
DEBUG = aaaaaa-debug$(SUFFIX)
DEBUG_GOFLAGS =
RELEASE = aaaaaa$(SUFFIX)
RELEASE_GOFLAGS = -ldflags="-s -w" -gcflags="-B -dwarf=false"
UPXFLAGS = -9

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
	$(RM) -r $(DEBUG) $(RELEASE) $(ASSETS)

.PHONY: $(ASSETS)
$(ASSETS):
	NO_VFS=$(NO_VFS) ./statik-vfs.sh $(ASSETS)

$(DEBUG): $(ASSETS)
	go build -o $(DEBUG) $(DEBUG_GOFLAGS) $(PACKAGE)

$(RELEASE): $(ASSETS)
	go build -o $(RELEASE) $(RELEASE_GOFLAGS) $(PACKAGE)

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

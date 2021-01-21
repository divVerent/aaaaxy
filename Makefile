# User configrable.
NO_VFS = false

# System properties.
EXE = $(shell go env GOEXE)

# Internal variables.
PACKAGE = github.com/divVerent/aaaaaa/cmd/aaaaaa
ASSETS = internal/assets
DEBUG = aaaaaa-debug$(EXE)
DEBUG_GOFLAGS =
ZIPFILE = aaaaaa.zip
RELEASE = aaaaaa$(EXE)
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

.PHONY: allrelease
allrelease:
	$(RM) $(ZIPFILE)
	GOOS=linux GOARCH=amd64 $(MAKE) release
	GOOS=windows GOARCH=386 $(MAKE) release
	zip -9r $(ZIPFILE) aaaaaa.exe aaaaaa

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
	upx $(UPXFLAGS) $(RELEASE)

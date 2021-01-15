# User configrable.
NO_VFS = false

# System properties.
EXE = $(shell go env GOEXE)

# Internal variables.
PACKAGE = github.com/divVerent/aaaaaa/cmd/aaaaaa
ASSETS = internal/assets
DEBUG = aaaaaa-debug$(EXE)
DEBUG_GOFLAGS =
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

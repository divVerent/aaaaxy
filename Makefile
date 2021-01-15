# User configrable.
NO_VFS = false

# System properties.
EXE = $(go env GOEXE)

# Internal variables.
PACKAGE = github.com/divVerent/aaaaaa/cmd/aaaaaa
ASSETS_SUBPACKAGE = internal/assets
DEBUG = aaaaaa-debug$(EXE)
DEBUG_FLAGS =
RELEASE = aaaaaa$(EXE)
RELEASE_FLAGS = -ldflags="-s -w" -gcflags="-B -dwarf=false"

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
	$(RM) -r $(DEBUG) $(RELEASE) $(ASSETS_SUBPACKAGE)

.PHONY: $(ASSETS_SUBPACKAGE)
$(ASSETS_SUBPACKAGE):
	NO_VFS=$(NO_VFS) ./statik-vfs.sh internal/assets:

$(DEBUG): generate
	go build -o $(DEBUG) $(DEBUG_FLAGS) $(PACKAGE)

$(RELEASE): generate
	go build -o $(RELEASE) $(RELEASE_FLAGS) $(PACKAGE)
	upx --best --ultra-brute $(RELEASE)

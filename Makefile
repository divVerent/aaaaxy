# Settings.
BUILDTYPE = debug

# System properties.
GO ?= go
BINARY = aaaaxy$(shell $(GO) env GOEXE)

ifeq ($(BUILDTYPE),release)
# Smaller binary.
ISRELEASE = true
BUILDTAGS = embed zip
else ifeq ($(BUILDTYPE),ziprelease)
# Separate assets file.
ISRELEASE = true
BUILDTAGS = zip
else ifeq ($(BUILDTYPE),embedrelease)
# Larger binary.
ISRELEASE = true
BUILDTAGS = embed
else
ISRELEASE = false
BUILDTAGS =
endif

ifeq ($(BUILDTYPE),extradebug)
BUILDTAGS += ebitendebug
endif

# Internal variables.
SOURCES = $(shell git ls-files \*.go)
GENERATED_STUFF = aaaaxy.ico aaaaxy.manifest aaaaxy.syso assets/generated/ internal/vfs/_embedroot/ licenses/asset-licenses/ licenses/software-licenses/

# Configure Go.
GO_FLAGS += -tags=$(shell echo $(BUILDTAGS) | tr ' ' ,)

# Configure the Go compiler.
GO_GCFLAGS = -dwarf=false
GO_FLAGS += $(patsubst %,-gcflags=all=%,$(GO_GCFLAGS))

# Configure the Go linker.
GO_LDFLAGS =
ifeq ($(shell $(GO) env GOOS),windows)
GO_LDFLAGS += -H=windowsgui
endif
GO_FLAGS += $(patsubst %,-ldflags=all=%,$(GO_LDFLAGS))

ifeq ($(ISRELEASE),true)
GO_LDFLAGS += -s -w
GO_FLAGS += -a -trimpath
ifneq ($(shell $(GO) env GOARCH),wasm)
GO_FLAGS += -buildmode=pie
endif
CPPFLAGS ?= -DNDEBUG
CFLAGS ?= -g0 -O3
CXXFLAGS ?= -g0 -O3
LDFLAGS ?= -g0 -s
endif

# cgo support.
CGO_CPPFLAGS ?= $(CPPFLAGS)
CGO_CFLAGS ?= $(CFLAGS)
CGO_CXXFLAGS ?= $(CXXFLAGS)
CGO_LDFLAGS ?= $(LDFLAGS)
CGO_ENV ?= \
	CGO_CPPFLAGS="$(CGO_CPPFLAGS)" \
	CGO_CFLAGS="$(CGO_CFLAGS)" \
	CGO_CXXFLAGS="$(CGO_CXXFLAGS)" \
	CGO_LDFLAGS="$(CGO_LDFLAGS)"

.PHONY: all
all: bin

.PHONY: bin
bin: $(BINARY)

.PHONY: clean
clean:
	$(RM) -r $(BINARY) $(GENERATED_STUFF)

.PHONY: vet
vet:
	$(GO) vet ./...
	$(GO) run honnef.co/go/tools/cmd/staticcheck@latest ./...
	# TODO make it bail out when something is found.
	gofmt -d -s $(SOURCES)
	gofmt -d -r 'fmt.Sprintf(s) -> s' $(SOURCES)
	gofmt -d -r 'fmt.Errorf(s) -> errors.New(s)' $(SOURCES)

.PHONY: mod-tidy
mod-tidy:
	$(GO) mod tidy -compat=1.21

.PHONY: mod-update
mod-update:
	$(GO) get -u
	$(MAKE) mod-tidy

.PHONY: assets-update
assets-update:
	AAAAXY_GENERATE_ASSETS=true AAAAXY_FORCE_GENERATE_ASSETS=true AAAAXY_DIFF_ASSETS=false sh -x scripts/build-generated-assets.sh
	cp assets/generated/* assets/_saved/

.PHONY: assets-update-all
assets-update-all:
	AAAAXY_GENERATE_ASSETS=true AAAAXY_GENERATE_CHECKPOINT_LOCATIONS=true AAAAXY_FORCE_GENERATE_ASSETS=true AAAAXY_DIFF_ASSETS=false sh -x scripts/build-generated-assets.sh
	cp assets/generated/* assets/_saved/

# These are not in assets-update as graphics are required for these.
.PHONY: loading-fractions-update
loading-fractions-update: $(BINARY)
	./$(BINARY) -load_config=false -dump_loading_fractions=assets/splash/loading_fractions.json -debug_just_init -debug_enable_drawing=false

pprof-update: $(BINARY)
	./$(BINARY) -load_config=false -debug_cpuprofile=default.pgo -palette=vga -draw_blurs=true -draw_outside=true -expand_using_vertices_accurately=true -screen_filter=linear2xcrt -vsync=false -demo_timedemo -demo_play=assets/demos/_anypercent.dem

# The actual build process follows.

# Packing the data files.
.PHONY: generate
generate:
	$(GO) generate $(GO_FLAGS)

# Binary building.
$(BINARY): generate $(SOURCES)
	$(CGO_ENV) $(GO) build -o $(BINARY) $(GO_FLAGS)

# Helper targets.
.PHONY: run
run: $(BINARY)
	EBITENGINE_INTERNAL_IMAGES_KEY=i EBITENGINE_SCREENSHOT_KEY=p ./$(BINARY) $(ARGS)

# Prepare git hooks.
.PHONY: setup-git
setup-git:
	git config filter.git-clean-tmx.clean "sh \"$$PWD/scripts/git-clean-tmx.sh\""
	git config filter.git-clean-md.clean "sh \"$$PWD/scripts/git-clean-md.sh\""

.PHONY: dangerous-clean-git
dangerous-clean-git:
	git reset --hard
	git clean -xdf
	git submodule update

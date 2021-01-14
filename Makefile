NO_VFS = false

all: aaaaaa-debug

generate:
	rm -rf assets-vfs/
	mkdir assets-vfs/
	set -x; \
	if $(NO_VFS); then \
		touch assets-vfs/use-local-assets.stamp; \
	else \
		./assets-vfs.sh; \
	fi
	statik -m -f -src assets-vfs/ -dest internal/assets/

aaaaaa-debug: generate
	go build -o aaaaaa-debug github.com/divVerent/aaaaaa/cmd/aaaaaa

aaaaaa: generate
	go build -o aaaaaa -ldflags='-s -w' github.com/divVerent/aaaaaa/cmd/aaaaaa
	upx --best --brute ./aaaaaa

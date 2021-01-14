NO_VFS = false

all: aaaaaa-debug

generate:
	NO_VFS=$(NO_VFS) ./statik-vfs.sh internal/assets/

aaaaaa-debug: generate
	go build -o aaaaaa-debug github.com/divVerent/aaaaaa/cmd/aaaaaa

aaaaaa: generate
	go build -o aaaaaa -ldflags='-s -w' github.com/divVerent/aaaaaa/cmd/aaaaaa
	upx --best --brute ./aaaaaa

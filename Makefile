all: aaaaaa-debug

generate:
	go generate ./...
	statik -m -f -src assets/ -dest internal/assets/

aaaaaa-debug: generate
	go build -o aaaaaa-debug github.com/divVerent/aaaaaa/cmd/aaaaaa

aaaaaa: generate
	go build -o aaaaaa -ldflags='-s -w' github.com/divVerent/aaaaaa/cmd/aaaaaa
	upx --best --brute ./aaaaaa

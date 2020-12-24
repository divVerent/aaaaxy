all: aaaaaa

generate:
	go generate ./...
	statik -m -f -src assets/ -dest internal/assets/

aaaaaa: generate
	go build github.com/divVerent/cmd/aaaaaa

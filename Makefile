SHELL = /bin/sh

run: build
	@./bin/myshell

build:
	@go build -o bin/myshell ./cmd/myshell/.

test:
	@go test -v -timeout 30s ./...

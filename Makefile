.DEFAULT_GOAL := test

fmt:
	golangci-lint run --disable-all --no-config -Egofmt --fix
	golangci-lint run --disable-all --no-config -Egofumpt --fix

lint: fmt
	golangci-lint run

build: lint
	go build .

install: build
	go install .

test:
	go test -race -shuffle on ./...

testv:
	go test -race -shuffle on -v ./...

clean:
	go clean -i -r -cache

.PHONY: fmt lint build install test testv clean

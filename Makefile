.DEFAULT_GOAL := test

fmt:
	golangci-lint run --disable-all --no-config -Egofmt --fix
	golangci-lint run --disable-all --no-config -Egofumpt --fix

lint: fmt
	golangci-lint run

test:
	go test -shuffle on github.com/telemachus/gitmirror/cli

testv:
	go test -shuffle on -v github.com/telemachus/gitmirror/cli

testr:
	go test -race -shuffle on github.com/telemachus/gitmirror/cli

build: lint testr
	go build .

install: build
	go install .

clean:
	go clean -i -r -cache

.PHONY: fmt lint build install test testv testr clean

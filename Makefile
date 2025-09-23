.DEFAULT_GOAL := test

PREFIX := $(HOME)/local/gitmirror

fmt:
	golangci-lint run --disable-all --no-config -Egofmt --fix
	golangci-lint run --disable-all --no-config -Egofumpt --fix

staticcheck: fmt
	staticcheck ./...

revive: fmt
	revive -config revive.toml ./...

golangci: fmt
	golangci-lint run

lint: fmt staticcheck revive golangci

test:
	go test -shuffle on github.com/telemachus/gitmirror/internal/cli
	go test -shuffle on github.com/telemachus/gitmirror/internal/git

testv:
	go test -shuffle on -v github.com/telemachus/gitmirror/internal/cli
	go test -shuffle on -v github.com/telemachus/gitmirror/internal/git

testr:
	go test -race -shuffle on github.com/telemachus/gitmirror/internal/cli
	go test -race -shuffle on github.com/telemachus/gitmirror/internal/git

build: lint testr
	go build .

install: build
	go install .

clean:
	rm -f gitmirror
	go clean -i -r -cache

.PHONY: fmt statickcheck revive golangci lint build install test testv testr clean

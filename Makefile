.DEFAULT_GOAL := test

PREFIX := $(HOME)/local/gitmirror

fmt:
	golangci-lint run --disable-all --no-config -Egofmt --fix
	golangci-lint run --disable-all --no-config -Egofumpt --fix

lint: fmt
	staticcheck ./...
	revive -config revive.toml ./...
	golangci-lint run

golangci: fmt
	golangci-lint run

staticcheck: fmt
	staticcheck ./...

revive: fmt
	revive -config revive.toml -exclude internal/optionparser ./...

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
	go build ./cmd/gitmirror

install: build
	go install ./cmd/gitmirror

clean:
	rm -f gitmirror
	go clean -i -r -cache

.PHONY: fmt lint build install test testv testr clean

.DEFAULT_GOAL := build

fmt:
	go fmt ./...

errcheck: fmt
	errcheck ./...

staticcheck: errcheck
	staticcheck ./...

vet: staticcheck
	go vet ./...

lint: vet
	golint ./...

build: lint
	go build .

install: build
	go install .

test:
	go test -race -shuffle on ./...

testv:
	go test -race -shuffle on -v ./...

clean:
	$(RM) gitmirror

.PHONY: fmt errcheck staticcheck vet build install test testv clean

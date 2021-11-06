.DEFAULT_GOAL := build

fmt:
	go fmt ./...

errcheck: fmt
	errcheck ./...

staticcheck: errcheck
	staticcheck ./...

vet: staticcheck
	go vet ./...

build: vet
	go build .

install: build
	go install .

test:
	go test -shuffle on ./...

testv:
	go test -shuffle on -v ./...

clean:
	$(RM) git-backup

.PHONY: fmt errcheck staticcheck vet build install test testv clean

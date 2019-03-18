.PHONY: clean check test build image dependencies fmt

SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')

IMAGE_NAME := containous/whoamitcp

default: clean check test build

clean:
	rm -rf dist/ builds/ cover.out

build: clean
	GO111MODULE=on go build -v .

image:
	docker build -t $(IMAGE_NAME) .

dependencies:
	GO111MODULE=on go mod download

test: clean
	GO111MODULE=on go test -v -cover ./...

check:
	GO111MODULE=on golangci-lint run

fmt:
	gofmt -s -l -w $(SRCS)

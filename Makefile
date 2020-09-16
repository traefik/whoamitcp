.PHONY: clean check test build image fmt

SRCS = $(shell git ls-files '*.go' | grep -v '^vendor/')

IMAGE_NAME := traefik/whoamitcp

default: clean check test build

clean:
	rm -rf dist/ builds/ cover.out

build: clean
	go build -v .

image:
	docker build -t $(IMAGE_NAME) .

test: clean
	go test -v -cover ./...

check:
	golangci-lint run

fmt:
	gofmt -s -l -w $(SRCS)

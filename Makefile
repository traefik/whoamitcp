.PHONY: default check test build image

IMAGE_NAME := traefik/whoamitcp

default: check test build

build:
	CGO_ENABLED=0 go build -v --trimpath .

test:
	go test -v -cover ./...

check:
	golangci-lint run

image:
	docker build -t $(IMAGE_NAME) .

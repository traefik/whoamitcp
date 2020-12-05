.PHONY: clean check test build image publish-images

TAG_NAME := $(shell git tag -l --contains HEAD)

IMAGE_NAME := traefik/whoamitcp

default: clean check test build

clean:
	rm -rf cover.out

build: clean
	CGO_ENABLED=0 go build -v --trimpath .

test: clean
	go test -v -cover ./...

check:
	golangci-lint run

image:
	docker build -t $(IMAGE_NAME) .

publish-images:
	seihon publish -v "$(TAG_NAME)" -v "latest" --image-name $(IMAGE_NAME) --dry-run=false

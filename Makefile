
GOOS?=darwin
GOARCH?=amd64

COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GITHUB_REF?=test

test:
	go test --race ./...

build:
	GOOS=${GOOS} GOARCH=${GOARCH} \
	GO111MODULE=on CGO_ENABLED=0 go build \
		-ldflags "-s -w -X main.Release=${GITHUB_REF} \
		-X main.Commit=${COMMIT} \
		-X main.BuildTime=${BUILD_TIME}" \
		-o bin/tr ./cmd/tr


GOOS?=darwin
GOARCH?=amd64

GOPROXY?=https://gocenter.io
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
RELEASE?=v0.0.0

test:
	go test --race ./...

build:
	GOOS=${GOOS} GOARCH=${GOARCH} \
	GO111MODULE=on GOPROXY=${GOPROXY} \
	CGO_ENABLED=0 go build \
		-ldflags "-s -w -X main.Release=${RELEASE} \
		-X main.Commit=${COMMIT} \
		-X main.BuildTime=${BUILD_TIME}" \
		-o bin/${GOOS}-${GOARCH}/tr ./cmd/tr

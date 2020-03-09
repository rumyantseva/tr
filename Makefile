GOOS?=darwin
GOARCH?=amd64

GOPROXY?=https://gocenter.io
COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GITHUB_REF?=refs/tags/v0.0.0
RELEASE?=$(subst refs/tags/,,${GITHUB_REF})

test:
	go test --race ./...

build:
	GOOS=${GOOS} GOARCH=${GOARCH} \
	GO111MODULE=on GOPROXY=${GOPROXY} \
	CGO_ENABLED=0 go build \
		-ldflags "-s -w -X main.release=${RELEASE} \
		-X main.commit=${COMMIT} \
		-X main.buildTime=${BUILD_TIME}" \
		-o bin/${GOOS}-${GOARCH}/tr ./cmd/tr

artifact: build
	mkdir -p ./bin/artifacts
	tar -zcvf ./bin/artifacts/tr-${GOOS}-${GOARCH}.tar.gz \
		--directory=./bin/${GOOS}-${GOARCH} tr

VERSION = 0.0.6
COMMIT ?= $(shell git rev-parse --short=8 HEAD)

unexport LDFLAGS
LDFLAGS=-ldflags "-s -X main.version=${VERSION} -X main.commit=${COMMIT}"

all: help windows darwin linux

help:
	@echo "Usage: make [all|windows|darwin|linux]"
windows:
	cd cmd/gachifinder && GOOS=windows GOARCH=amd64 GO111MODULE=on go build -o windows/gachifinder.exe ${LDFLAGS}
darwin:
	cd cmd/gachifinder && GOOS=darwin GOARCH=amd64 GO111MODULE=on go build -o osx/gachifinder ${LDFLAGS}
linux:
	cd cmd/gachifinder && GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o linux/gachifinder ${LDFLAGS}

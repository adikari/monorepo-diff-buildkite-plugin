NAME=monorepo-diff-buildkite-plugin
VERSION := $(shell git describe --tags --abbrev=0)
COMMIT=$(shell git rev-parse --short=7 HEAD)
TIMESTAMP:=$(shell date -u '+%Y-%m-%dT%I:%M:%SZ')

LDFLAGS += -X main.BuildTime=${TIMESTAMP}
LDFLAGS += -X main.BuildSHA=${COMMIT}
LDFLAGS += -X main.Version=${VERSION}
HAS_DOCKER=$(shell command -v docker;)

.PHONY: all
all: quality test

.PHONY: test
test:
	go test -race -coverprofile=coverage.out

.PHONY: quality
quality:
	go vet
	go fmt
	go mod tidy
ifneq (${HAS_DOCKER},)
	docker-compose run --rm lint
endif

.PHONY: clean
clean:
	rm -f coverage.out
	rm -rf ${NAME}*

.PHONY: build
build-%: clean
	GOOS=$* GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '${LDFLAGS}' -o ${PWD}/${NAME}-$*

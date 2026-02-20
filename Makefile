# dcctl — build, Docker, release (macOS, Linux, Windows)
# Version from VERSION file; override: make build VERSION=v0.2.0

BINARY    := dcctl
DIST      := dist
VERSION   ?= $(shell v=$$(cat VERSION 2>/dev/null | tr -d '\n\r'); [ -n "$$v" ] && echo "v$$v" || echo "v0.1.0")
COMMIT    := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILDDATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS   := -ldflags "-s -w \
	-X dcctl/cmd/dcctl/root.buildVersion=$(VERSION) \
	-X dcctl/cmd/dcctl/root.buildCommit=$(COMMIT) \
	-X dcctl/cmd/dcctl/root.buildDate=$(BUILDDATE)"

.PHONY: build install test clean docker-build release snapshot build-all

# Build for current platform
build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/dcctl

# Install to $GOBIN (default $HOME/go/bin). Custom: GOBIN=/usr/local/bin make install
install:
	go install $(LDFLAGS) ./cmd/dcctl

test:
	go test ./...

# Docker image with version/commit/builddate from VERSION and git (run from repo root)
docker-build:
	docker build --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) --build-arg BUILDDATE=$(BUILDDATE) -t dcctl .

# Release: only from main. Merge develop → main, update VERSION, then: git tag v0.1.0 && make release
# Requires goreleaser: brew install goreleaser
release:
	@branch=$$(git branch --show-current 2>/dev/null); \
	if [ "$$branch" != "main" ]; then \
	  echo "Error: release only from main (current: $$branch). Merge and checkout main first."; \
	  exit 1; \
	fi; \
	goreleaser release --clean

# Snapshot build (no tag required), outputs to dist/
snapshot:
	goreleaser release --snapshot --clean

# Cross-compile: all platforms (output in dist/)
build-all: build-linux build-darwin build-windows

build-linux:
	@mkdir -p $(DIST)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST)/$(BINARY)-linux-amd64 ./cmd/dcctl
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(DIST)/$(BINARY)-linux-arm64 ./cmd/dcctl

build-darwin:
	@mkdir -p $(DIST)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST)/$(BINARY)-darwin-amd64 ./cmd/dcctl
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST)/$(BINARY)-darwin-arm64 ./cmd/dcctl

build-windows:
	@mkdir -p $(DIST)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST)/$(BINARY)-windows-amd64.exe ./cmd/dcctl
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o $(DIST)/$(BINARY)-windows-arm64.exe ./cmd/dcctl

clean:
	rm -f $(BINARY)
	rm -rf $(DIST)

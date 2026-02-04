.PHONY: build build-all test test-coverage lint lint-fix clean deps

# Module and version info
MODULE := github.com/ravi-technologies/sunday-cli
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# API URL must be provided at build time
ifndef API_URL
  $(error API_URL is required. Usage: make build API_URL=https://api.sunday.app)
endif

# ldflags for build-time variable injection
LDFLAGS := -ldflags "\
	-X '$(MODULE)/internal/version.Version=$(VERSION)' \
	-X '$(MODULE)/internal/version.Commit=$(COMMIT)' \
	-X '$(MODULE)/internal/version.BuildDate=$(BUILD_DATE)' \
	-X '$(MODULE)/internal/version.APIBaseURL=$(API_URL)'"

# ----------------
#    Build
# ----------------

build:
	go build $(LDFLAGS) -o bin/sunday ./cmd/sunday

build-all:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/sunday-darwin-amd64 ./cmd/sunday
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/sunday-darwin-arm64 ./cmd/sunday
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/sunday-linux-amd64 ./cmd/sunday
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/sunday-linux-arm64 ./cmd/sunday
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/sunday-windows-amd64.exe ./cmd/sunday

# ----------------
#    Development
# ----------------

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

clean:
	rm -rf bin/

# ----------------
#    Dependencies
# ----------------

deps:
	go mod download
	go mod tidy

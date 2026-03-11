.PHONY: all build test clean install release lint fmt vet coverage docker help

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X github.com/wesbragagt/gps/internal/version.Version=$(VERSION) \
	-X github.com/wesbragagt/gps/internal/version.GitCommit=$(GIT_COMMIT) \
	-X github.com/wesbragagt/gps/internal/version.BuildDate=$(BUILD_DATE) \
	-s -w"

BINARY := gps
CMD_PATH := ./cmd/gps

# Default target
all: test build

## build: Build the binary for current platform
build:
	@echo "Building $(BINARY)..."
	@mkdir -p bin
	go build $(LDFLAGS) -o bin/$(BINARY) $(CMD_PATH)

## test: Run all tests with race detection
test:
	go test -v -race -coverprofile=coverage.out ./...

## coverage: Generate HTML coverage report
coverage:
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## clean: Remove build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

## install: Install binary to GOPATH/bin
install:
	go install $(LDFLAGS) $(CMD_PATH)

## lint: Run golangci-lint
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not found, please install it" && exit 1)
	golangci-lint run ./...

## fmt: Format Go source files
fmt:
	go fmt ./...

## vet: Run go vet
vet:
	go vet ./...

## check: Run all checks (fmt, vet, test)
check: fmt vet test

## release: Build binaries for all platforms
release: clean
	@echo "Building release binaries..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o bin/$(BINARY)-linux-amd64 $(CMD_PATH)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o bin/$(BINARY)-linux-arm64 $(CMD_PATH)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o bin/$(BINARY)-darwin-amd64 $(CMD_PATH)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o bin/$(BINARY)-darwin-arm64 $(CMD_PATH)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o bin/$(BINARY)-windows-amd64.exe $(CMD_PATH)
	GOOS=windows GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o bin/$(BINARY)-windows-arm64.exe $(CMD_PATH)
	@echo "Release binaries built in bin/"

## archives: Create release archives
archives: release
	@echo "Creating archives..."
	cd bin && \
	tar -czf $(BINARY)-linux-amd64.tar.gz $(BINARY)-linux-amd64 && \
	tar -czf $(BINARY)-linux-arm64.tar.gz $(BINARY)-linux-arm64 && \
	tar -czf $(BINARY)-darwin-amd64.tar.gz $(BINARY)-darwin-amd64 && \
	tar -czf $(BINARY)-darwin-arm64.tar.gz $(BINARY)-darwin-arm64 && \
	zip -q $(BINARY)-windows-amd64.zip $(BINARY)-windows-amd64.exe && \
	zip -q $(BINARY)-windows-arm64.zip $(BINARY)-windows-arm64.exe
	@echo "Archives created in bin/"

## docker: Build Docker image
docker:
	docker build -t $(BINARY):$(VERSION) -t $(BINARY):latest .

## version: Display version info
version:
	@echo "Version: $(VERSION)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Build Date: $(BUILD_DATE)"

## help: Show this help
help:
	@echo "Available targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'

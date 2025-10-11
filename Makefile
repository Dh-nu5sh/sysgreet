.PHONY: fmt lint test test-verbose test-coverage bench build clean help

fmt:
	gofmt -w $(shell find . -name '*.go' -not -path './vendor/*')

lint:
	go vet ./...
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install with: brew install golangci-lint" && exit 1)
	golangci-lint run

test:
	CGO_ENABLED=0 go test ./...

test-verbose:
	CGO_ENABLED=0 go test -v ./...

test-coverage:
	CGO_ENABLED=0 go test -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

bench:
	CGO_ENABLED=0 go test -bench . ./test/benchmarks

build:
	CGO_ENABLED=0 go build -o sysgreet ./cmd/sysgreet

clean:
	rm -f sysgreet coverage.out

help:
	@echo "Available targets:"
	@echo "  fmt            - Format all Go files"
	@echo "  lint           - Run go vet and golangci-lint"
	@echo "  test           - Run all tests"
	@echo "  test-verbose   - Run all tests with verbose output"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  bench          - Run performance benchmarks"
	@echo "  build          - Build the sysgreet binary"
	@echo "  clean          - Remove build artifacts"
	@echo "  help           - Show this help message"

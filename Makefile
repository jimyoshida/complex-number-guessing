.PHONY: build run test coverage clean help

BINARY_NAME=complex-number-guessing
MAIN_PATH=./main.go

help:
	@echo "Available targets:"
	@echo "  build    - Build the Go binary"
	@echo "  run      - Build and run the application"
	@echo "  test     - Run tests"
	@echo "  coverage - Generate test coverage report (HTML)"
	@echo "  clean    - Remove build artifacts"
	@echo "  help     - Show this help message"

build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

run: build
	./bin/$(BINARY_NAME)

test:
	go test -v ./...

coverage:
	mkdir -p build
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o build/coverage.html
	@echo "✓ Coverage report generated: build/coverage.html"

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

.PHONY: build run test clean help

BINARY_NAME=complex-number-guessing
MAIN_PATH=./main.go

help:
	@echo "Available targets:"
	@echo "  build   - Build the Go binary"
	@echo "  run     - Build and run the application"
	@echo "  test    - Run tests"
	@echo "  clean   - Remove build artifacts"
	@echo "  help    - Show this help message"

build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

run: build
	./bin/$(BINARY_NAME)

test:
	go test -v ./...

clean:
	rm -rf bin/
	go clean

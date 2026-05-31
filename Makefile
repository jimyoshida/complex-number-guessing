.PHONY: build run test coverage web-build web-serve web-clean clean help

BINARY_NAME=complex-number-guessing
WASM_OUT=web/game.wasm
WASM_EXEC_SRC=$(shell go env GOROOT)/lib/wasm/wasm_exec.js
WASM_EXEC_DST=web/wasm_exec.js

help:
	@echo "Available targets:"
	@echo "  build      - Build the CLI binary"
	@echo "  run        - Build and run the CLI"
	@echo "  test       - Run tests"
	@echo "  coverage   - Generate test coverage report (HTML)"
	@echo "  web-build  - Compile Go to WASM and copy wasm_exec.js into web/"
	@echo "  web-serve  - Build WASM then serve web/ on http://localhost:8080"
	@echo "  web-clean  - Remove WASM build artifacts from web/"
	@echo "  clean      - Remove all build artifacts"
	@echo "  help       - Show this help message"

build:
	go build -o bin/$(BINARY_NAME) .

run: build
	./bin/$(BINARY_NAME)

test:
	go test -v ./...

coverage:
	mkdir -p build
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o build/coverage.html
	@echo "✓ Coverage report generated: build/coverage.html"

web-build:
	GOOS=js GOARCH=wasm go build -o $(WASM_OUT) .
	cp $(WASM_EXEC_SRC) $(WASM_EXEC_DST)
	@echo "✓ WASM built: $(WASM_OUT)"

web-serve: web-build
	@echo "Serving on http://localhost:8080"
	cd web && python3 -m http.server 8080

web-clean:
	rm -f $(WASM_OUT) $(WASM_EXEC_DST)

clean: web-clean
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -rf build/
	go clean

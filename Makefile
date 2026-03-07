# Village Simulation Makefile

.PHONY: build test run clean

build:
	@echo "Building village simulation..."
	go build -o bin/village ./cmd/village

test:
	@echo "Running tests..."
	go test ./...

run: build
	@echo "Starting village simulation..."
	./bin/village

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/

# Development targets
dev:
	go run cmd/village/main.go

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

vet:
	go vet ./...

fmt:
	go fmt ./...
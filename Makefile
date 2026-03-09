# Village Simulation Makefile

.PHONY: build test run clean dev coverage vet fmt lint lint-fix build-backend build-frontend deploy-dry-run deploy check-go-version

# Default target
all: build

build: build-backend build-frontend
	@echo "Build completed"

build-backend:
	@echo "Building backend..."
	go build -o bin/village ./cmd/village

build-frontend:
	@echo "Building frontend..."
	cd frontend && npm ci --omit=dev
	cd frontend && npm run build -- --configuration production

test:
	@echo "Running tests..."
	go test ./...

run: build
	@echo "Starting village simulation..."
	./bin/village

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf frontend/dist/

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

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

# Docker targets
docker-build:
	docker build --target backend -t village-backend .
	docker build --target frontend -t village-frontend .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Deployment targets
deploy-dry-run:
	@echo "=== Deployment Dry Run ==="
	@echo "1. Build Docker images:"
	@echo "   - Backend: village-backend"
	@echo "   - Frontend: village-frontend"
	@echo "2. Start services with docker-compose:"
	@echo "   - postgres:5432"
	@echo "   - backend:8080"
	@echo "   - frontend:80"
	@echo "3. Apply database migrations (if any)"
	@echo "4. Health checks"
	@echo "Dry run complete - no changes made."

deploy: docker-build docker-run
	@echo "Deployment started. Use 'make docker-logs' to monitor."

# Helper targets
health:
	curl -f http://localhost:8080/health || echo "Backend health check failed"
	curl -f http://localhost:80/health || echo "Frontend health check failed"

migrate:
	@echo "Running database migrations (if any)..."
	# TODO: Add migration command when available

check-go-version:
	@go version | grep -q "go1.25" || (echo "Error: Go version mismatch. Expected 1.25.x" && exit 1)
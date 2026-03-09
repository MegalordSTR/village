# Village Simulation - Resource Economy System

A simulation game featuring a dynamic resource economy with production chains, seasonal variations, and complex agent behaviors.

## Features

- **Resource Economy**: Multiple resource types with quality tiers, spoilage, and base values
- **Production Chains**: Recipes that transform raw materials into processed and advanced goods
- **Seasonal Variations**: Agricultural yields and resource availability change with seasons
- **Agent System**: Villagers with skills, needs, and economic decision-making
- **Real-time Simulation**: Turn-based simulation engine with configurable time steps

## Architecture

- **Backend**: Go service providing simulation engine and REST API
- **Frontend**: Angular web interface for visualization and interaction
- **Database**: PostgreSQL for persistent state (optional)
- **Containerized**: Docker and Docker Compose for easy deployment

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.25+ (for local development)
- Node.js 20+ and npm (for frontend development)

### Using Docker (Recommended)

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd village
   ```

2. Build and start the application:
   ```bash
   docker-compose up --build
   ```

3. Access the application:
   - Frontend: http://localhost:80
   - Backend API: http://localhost:8080
   - PostgreSQL: localhost:5432 (username: village, password: village)

### Local Development

#### Backend

```bash
# Build and run
make build
make run

# Run tests
make test

# Code quality
make vet
make fmt
make lint

## Linting

This project uses [golangci-lint](https://golangci-lint.run/) for static analysis. To install:

```bash
# Install golangci-lint (requires Go 1.25+)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.60.0
```

Run linting:

```bash
make lint          # Run all linters
make lint-fix      # Auto-fix fixable issues
```

Configuration is in `.golangci.yml`. The CI pipeline runs linting on every push and pull request.
```

#### Frontend

```bash
cd frontend

# Install dependencies
npm ci

# Start development server
npm start

# Build for production
npm run build
```

## Deployment

### Single-command Deployment

```bash
make deploy
```

This will:
1. Build Docker images for backend and frontend
2. Start all services using Docker Compose
3. Perform health checks

### Dry Run

To see the deployment plan without executing:

```bash
make deploy-dry-run
```

### Production Deployment

For production deployments, use the included deployment script:

```bash
./scripts/deploy.sh production
```

The script supports environment-specific configurations and health checks.

### CI/CD Pipeline

The project includes a GitHub Actions workflow (`.github/workflows/deploy.yml`) that:
- Runs tests on every push and pull request
- Builds Docker images locally (no external registry dependency)
- Deploys to production on pushes to the main branch either:
  - **Via SSH** to a remote server (if SSH credentials are configured)
  - **Locally on the CI runner** (for validation)

See [deployment documentation](docs/deployment.md) for details.

## Configuration

### Environment Variables

- `DATABASE_URL`: PostgreSQL connection string (default: `postgres://village:village@postgres:5432/village?sslmode=disable`)
- `PORT`: Backend API port (default: `8080`)

Set these in `docker-compose.yml` or via `.env` file.

### Docker Configuration

- `Dockerfile`: Multi-stage build for backend and frontend
- `docker-compose.yml`: Defines backend, frontend, and PostgreSQL services
- `nginx.conf`: Nginx configuration for frontend with API proxying

## Health Checks

Services include health check endpoints:

- Backend: `http://localhost:8080/health`
- Frontend: `http://localhost:80/health`

Use `make health` to check all services.

## Database Migrations

When database schema changes are required:

```bash
make migrate
```

*Note: Migration system is under development.*

## Monitoring and Logging

- **Docker Compose Logs**: `docker-compose logs -f`
- **Application Logs**: Check container stdout
- **Health Monitoring**: Built-in health endpoints

## Project Structure

```
.
├── cmd/village/          # Backend entry point
├── internal/             # Go packages (economy, simulation)
├── frontend/             # Angular application
├── docs/                 # Documentation
├── scripts/              # Deployment and utility scripts
├── Dockerfile            # Multi-stage Docker build
├── docker-compose.yml    # Service orchestration
├── Makefile              # Development commands
└── .github/workflows/    # CI/CD pipelines
```

## Contributing

1. Create a feature branch from `main`
2. Implement changes with tests
3. Ensure all tests pass: `make test`
4. Update documentation as needed
5. Submit a pull request

## License

[License details to be added]
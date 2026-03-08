# Deployment Guide

## Overview

Village Economy System uses a simplified deployment process that eliminates dependency on Docker container registries. Instead of pushing and pulling images from a remote registry, Docker images are built locally (or on the CI runner) and immediately deployed using Docker Compose. This reduces complexity and removes external dependencies.

## Deployment Methods

### 1. Local Deployment (Development)

For local development and testing, use Docker Compose directly:

```bash
# Build images and start all services
docker-compose up --build

# Or use the Makefile target
make deploy
```

This will:
- Build backend and frontend Docker images using the multi‑stage `Dockerfile`
- Start PostgreSQL, backend, and frontend containers
- Expose the services on localhost:
  - Frontend: http://localhost:80
  - Backend API: http://localhost:8080
  - PostgreSQL: localhost:5432

### 2. Production Deployment (Manual)

For production deployments, use the included deployment script:

```bash
./scripts/deploy.sh production
```

The script performs the following steps:
1. Validates the presence of a production environment file (`.env.production`)
2. Builds Docker images with `--no‑cache` to ensure a clean build
3. Stops any existing services (`docker‑compose down`)
4. Starts the new containers (`docker‑compose up -d`)
5. Waits for services to become healthy and runs health checks

### 3. CI/CD Pipeline (GitHub Actions)

The project includes a GitHub Actions workflow (`.github/workflows/deploy.yml`) that automates testing and deployment.

**Trigger:** The workflow runs on every push to the `main` or `master` branch.

**Jobs:**
1. **Test** – runs Go tests, builds the backend and frontend to verify integrity.
2. **Deploy** – (only on push to main/master) builds images and deploys either:
   - **Via SSH** – if the repository secrets `SSH_PRIVATE_KEY`, `DEPLOY_HOST`, `DEPLOY_USER` (and optionally `DEPLOY_PORT`) are set, the workflow uses `appleboy/ssh‑action` to connect to a remote server, pull the latest code, and run `docker‑compose up`.
   - **Locally on the runner** – if no SSH key is configured, the job builds and starts the containers on the CI runner (useful for validating the deployment process without a remote server).

**Simplification:** The pipeline no longer depends on GitHub Container Registry (or any external registry). Images are built directly on the runner or on the remote deployment server.

## Environment Configuration

### Environment Variables

- `DATABASE_URL` – PostgreSQL connection string (default: `postgres://village:village@postgres:5432/village?sslmode=disable`)
- `PORT` – backend API port (default: `8080`)

You can set these variables in:
- `docker‑compose.yml` (for local development)
- `.env.production` file (for production)
- GitHub Actions secrets (for CI/CD)

### Production Environment File

Create a `.env.production` file in the project root with production‑specific values:

```bash
DATABASE_URL=postgres://user:password@host:5432/dbname
PORT=8080
```

The file is optional; if missing, default values from `docker‑compose.yml` will be used.

## Health Checks

Each service includes a health‑check endpoint:

- **Backend:** `http://localhost:8080/health`
- **Frontend:** `http://localhost:80/health`

You can verify all services with:

```bash
make health
```

## Security Considerations

- **Environment Variables**: Never commit `.env.production` to version control. Use GitHub Actions secrets for CI/CD and keep the file secure on production servers.
- **PostgreSQL Credentials**: The default credentials (`village:village`) in `docker‑compose.yml` are for local development only. In production, always set `POSTGRES_USER`, `POSTGRES_PASSWORD`, and `POSTGRES_DB` via `.env.production` or equivalent.
- **SSH Deployment**: When using SSH deployment, store the private key as a GitHub Actions secret. The workflow adds the remote host’s fingerprint via `ssh‑keyscan`, but the first connection remains vulnerable to MITM if the fingerprint is not verified independently. For higher security, pre‑populate `known_hosts` with the verified fingerprint.
- **Network Exposure**: By default, the frontend is accessible on port 80 and the backend on port 8080. In production, consider placing a reverse proxy (nginx, Traefik) in front of the services, enabling TLS/HTTPS, and restricting access with a firewall.
- **Health Checks**: The health endpoints (`/health`) are exposed internally within the Docker network. They should not be publicly accessible.
- **Docker Security**: The backend container runs as a non‑root user (`appuser`). The PostgreSQL container runs as the default `postgres` user. Keep images updated to avoid known vulnerabilities.

## Verification Commands

Before deploying, you can run the following verification steps:

```bash
# Validate Docker Compose configuration
docker-compose config

# Dry‑run of the deployment process
make deploy-dry-run

# Test the deployment script locally
./scripts/deploy.sh local
```

## Troubleshooting

### Docker Compose Build Failures

- Ensure Docker and Docker Compose are installed and the Docker daemon is running.
- Check that the `Dockerfile` and `docker‑compose.yml` are in the correct location.
- Verify that there are no syntax errors in the `Dockerfile` or compose file.

### SSH Deployment Issues

- Confirm that the SSH private key secret (`SSH_PRIVATE_KEY`) is correctly set in GitHub repository secrets.
- Ensure the deployment server has Docker and Docker Compose installed.
- The remote server must have the repository cloned at the path used in the workflow (`/opt/village` by default). Adjust the `cd` command in `.github/workflows/deploy.yml` if needed.

### Health Check Failures

- Wait a few seconds after startup for services to become ready.
- Check container logs: `docker‑compose logs [service]`.
- Verify that the backend can connect to the PostgreSQL container (check network configuration in `docker‑compose.yml`).

## Migration from Registry‑Based Deployment

If you previously used the GitHub Container Registry (ghcr.io) deployment:

1. Remove any references to `ghcr.io` from your environment variables, scripts, or configuration files.
2. Delete the `build‑and‑push` job from `.github/workflows/deploy.yml` (already removed in this version).
3. Update any external documentation that still mentions pushing/pulling images.

## Further Reading

- [README.md](../README.md) – project overview and quick start
- [Dockerfile](../Dockerfile) – multi‑stage build configuration
- [docker‑compose.yml](../docker-compose.yml) – service definitions and networking
- [Makefile](../Makefile) – convenient development and deployment targets
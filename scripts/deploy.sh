#!/bin/bash
set -e

# Deployment script for Village Economy System
# Usage: ./scripts/deploy.sh [environment]
# environment: local (default), production, ci

ENVIRONMENT=${1:-local}
COMPOSE_FILE="docker-compose.yml"
DOCKER_COMPOSE="docker-compose"

echo "Deploying Village Economy System ($ENVIRONMENT)"

# Check prerequisites
command -v docker >/dev/null 2>&1 || { echo "docker is required but not installed. Aborting." >&2; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "docker-compose is required but not installed. Aborting." >&2; exit 1; }

case "$ENVIRONMENT" in
  local)
    echo "Starting local deployment..."
    # Build images
    $DOCKER_COMPOSE -f $COMPOSE_FILE build
    # Start services
    $DOCKER_COMPOSE -f $COMPOSE_FILE up -d
    ;;
  production)
    echo "Starting production deployment..."
    # Ensure we have production environment variables
    if [ ! -f .env.production ]; then
       echo "Warning: .env.production not found. Using default environment variables."
    fi
    # Pull latest images (if using registry)
    # $DOCKER_COMPOSE -f $COMPOSE_FILE pull
    # Build images locally
    $DOCKER_COMPOSE -f $COMPOSE_FILE build --no-cache
    # Stop existing services
    $DOCKER_COMPOSE -f $COMPOSE_FILE down
    # Start services
    $DOCKER_COMPOSE -f $COMPOSE_FILE up -d
    ;;
  ci)
    echo "CI/CD environment detected"
    # Validate configuration
    $DOCKER_COMPOSE -f $COMPOSE_FILE config
    # Build images (no cache)
    $DOCKER_COMPOSE -f $COMPOSE_FILE build --no-cache
    ;;
  *)
    echo "Unknown environment: $ENVIRONMENT"
    echo "Usage: $0 [local|production|ci]"
    exit 1
    ;;
esac

echo "Waiting for services to become healthy..."
sleep 10

# Health checks
echo "Performing health checks..."
if curl -f http://localhost:8080/health; then
  echo "Backend is healthy"
else
  echo "Backend health check failed"
  exit 1
fi

if curl -f http://localhost:80/health; then
  echo "Frontend is healthy"
else
  echo "Frontend health check failed"
  exit 1
fi

echo "Deployment completed successfully!"
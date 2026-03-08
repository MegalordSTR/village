#!/bin/bash
set -e

# Deployment script for Village Economy System
# Usage: ./scripts/deploy.sh [environment]
# environment: local (default), production, ci

ENVIRONMENT=${1:-local}
COMPOSE_FILE="docker-compose.yml"
DOCKER_COMPOSE="docker-compose"
SKIP_HEALTH_CHECKS=0

echo "Deploying Village Economy System ($ENVIRONMENT)"

# Check prerequisites
command -v docker >/dev/null 2>&1 || { echo "docker is required but not installed. Aborting." >&2; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "docker-compose is required but not installed. Aborting." >&2; exit 1; }

case "$ENVIRONMENT" in
  local)
    echo "Starting local deployment..."
    # Build images and start services
    $DOCKER_COMPOSE -f $COMPOSE_FILE up -d --build
    ;;
  production)
    echo "Starting production deployment..."
    # Require production environment variables
    if [ ! -f .env.production ]; then
       echo "Error: .env.production not found. Production deployment requires environment variables."
       exit 1
    fi
    # Build images locally
    $DOCKER_COMPOSE -f $COMPOSE_FILE build --no-cache
    # Stop existing services and start new ones
    $DOCKER_COMPOSE -f $COMPOSE_FILE up -d --build
    ;;
  ci)
    echo "CI/CD environment detected"
    SKIP_HEALTH_CHECKS=1
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

if [ "$SKIP_HEALTH_CHECKS" -eq 0 ]; then
  echo "Waiting for services to become healthy..."
  sleep 10

  # Health check with retry
  health_check() {
    local url=$1
    local service=$2
    local max_attempts=5
    local attempt=1
    echo "Checking $service health at $url"
    while [ $attempt -le $max_attempts ]; do
      if curl -f -s $url >/dev/null; then
        echo "$service is healthy"
        return 0
      fi
      echo "Attempt $attempt/$max_attempts failed, retrying in 5 seconds..."
      sleep 5
      attempt=$((attempt + 1))
    done
    echo "$service health check failed after $max_attempts attempts"
    return 1
  }

  # Perform health checks
  if ! health_check "http://localhost:8080/health" "Backend"; then
    echo "Backend health check failed"
    exit 1
  fi

  if ! health_check "http://localhost:80/health" "Frontend"; then
    echo "Frontend health check failed"
    exit 1
  fi

  echo "Deployment completed successfully!"
else
  echo "CI validation completed"
fi
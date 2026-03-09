# SRE Review - Feature F005 (Deployment Simplification) - Round 2

## SCOPE
Files modified/added:
- `.github/workflows/deploy.yml` (modified)
- `scripts/deploy.sh` (modified)
- `README.md` (modified)
- `docs/deployment.md` (new)
- `Makefile` (unchanged but in scope)
- `Dockerfile` (modified)
- `docker-compose.yml` (modified)

## RISK MAP
Top 3 reliability risk areas:
1. **Health Check Reliability**: Frontend health check uses curl not installed in image; inconsistent retry logic between deployment methods.
2. **Production Configuration Loading**: Deployment script checks for .env.production but does not load it, causing production deployments to use wrong environment variables.
3. **Rollback & Recovery**: No rollback mechanism on failed deployment; no image versioning makes recovery difficult.

## EVIDENCE

### Previous P1 Findings Status
- ✓ SSH host key verification added (deploy.yml lines 55-59)
- ✓ Backend container now runs as non-root user (Dockerfile line 17)
- ✓ Health check updated to use curl (docker-compose.yml line 35)
- ~ Default PostgreSQL credentials now use environment variables with defaults (still weak)
- ~ Build integrity reduction (design choice accepted)

### New Reliability Findings

**Finding 1: Frontend health check uses curl not installed**
- **File**: `docker-compose.yml:51`
- **Issue**: Health check `CMD ["curl", "-f", "http://localhost:80/health"]` requires curl in nginx:alpine image, but frontend Dockerfile does not install curl.
- **Impact**: Health checks will always fail, potentially causing container restart loops and false unhealthy status.

**Finding 2: Production environment file not loaded**
- **File**: `scripts/deploy.sh:28-31`
- **Issue**: Script checks for `.env.production` existence but does not export variables or pass to docker-compose. Docker-compose uses `.env` by default, not `.env.production`.
- **Impact**: Production deployments may use wrong environment variables (e.g., default credentials), leading to security and configuration issues.

**Finding 3: No rollback mechanism on failed deployment**
- **File**: `scripts/deploy.sh:35`, `.github/workflows/deploy.yml:73`
- **Issue**: If health check fails after `docker-compose up --build -d`, script exits with error but containers remain in broken state. No automatic reversion to previous version.
- **Impact**: Extended downtime and manual recovery required.

**Finding 4: Inconsistent health check retry logic**
- **File**: `.github/workflows/deploy.yml:76` vs `scripts/deploy.sh:57-74`
- **Issue**: GitHub Actions deployment has simple 10s sleep + curl with no retry; deploy.sh has 5-attempt retry logic. Inconsistent behavior can cause false failures in CI/CD.
- **Impact**: Unreliable deployments depending on method used.

**Finding 5: Default PostgreSQL credentials still weak**
- **File**: `docker-compose.yml:5-7`
- **Issue**: Environment variable defaults `village:village` remain; production deployment requires `.env.production` but nothing prevents copying defaults.
- **Impact**: Risk of weak credentials in production if not properly overridden.

**Finding 6: Potential brief downtime during deployment**
- **Issue**: `docker-compose up --build -d` stops old containers before starting new ones, causing brief downtime.
- **Impact**: Acceptable for non-critical services but violates zero-downtime expectations.

**Finding 7: No Docker image versioning**
- **Issue**: Images built locally without version tags; previous images become dangling `<none>` tags.
- **Impact**: Difficult rollback and debugging of which image version is running.

## SEVERITY
1. Frontend health check uses curl not installed: **P1** (breaks deployment health monitoring)
2. Production environment file not loaded: **P1** (breaks production configuration)
3. No rollback mechanism: **P2** (maintenance debt, extended downtime)
4. Inconsistent health check retry logic: **P3** (style)
5. Default PostgreSQL credentials still weak: **P2** (security risk)
6. Potential brief downtime: **P3** (style)
7. No Docker image versioning: **P3** (maintenance debt)

## VERDICT
**FAIL** – Maximum severity is P1 (frontend health check and production configuration loading). These issues must be addressed before feature can be considered reliable for production use.

## FINDINGS_CREATED
village-00d, village-hqs, village-jty, village-tii, village-0hv, village-c9q, village-2nf

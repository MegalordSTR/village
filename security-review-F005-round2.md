# Security Review Round 2 - Feature F005 (Deployment Simplification)

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
Top 3 security risk areas:
1. **Credential Management**: Default PostgreSQL credentials still present as fallback values.
2. **SSH Deployment Security**: SSH host key verification via ssh-keyscan still vulnerable to MITM on first connection.
3. **Container Health Monitoring**: Frontend health check requires curl but nginx:alpine image does not include curl; backend health endpoint may not exist.

## EVIDENCE

### Finding 1: Default PostgreSQL credentials still present as defaults
- **File**: `docker-compose.yml:5-7`
- **Issue**: Environment variables have default values (`village:village`) that remain active if not overridden.
- **Impact**: If someone deploys using `docker-compose up` directly without setting environment variables, the database will use weak credentials.
- **Mitigation**: Production deployment script requires `.env.production`, but defaults still exist for local development.

### Finding 2: SSH host key verification still vulnerable to MITM
- **File**: `.github/workflows/deploy.yml:55-60`
- **Issue**: `ssh-keyscan` is used to add the host key to `known_hosts` without prior verification. An attacker performing MITM during the first workflow run could inject their own host key.
- **Impact**: Potential man‑in‑the‑middle attack during SSH connection.
- **Mitigation**: Consider storing the expected host key as a secret and writing it directly to `known_hosts`.

### Finding 3: Frontend health check uses curl but image lacks curl
- **File**: `docker-compose.yml:50-54`
- **Issue**: Frontend health check uses `CMD ["curl", "-f", "http://localhost:80/health"]`. The base image `nginx:alpine` does not include `curl`.
- **Impact**: Health check will always fail, causing false deployment failures or (if health check is ignored) masking real service issues.
- **Mitigation**: Install `curl` in the frontend Docker stage or use a different health check method (e.g., `wget` or `nc`).

### Finding 4: Backend health endpoint may not exist
- **File**: `cmd/village/main.go` (no HTTP server)
- **Issue**: The backend container currently does not start an HTTP server, yet the health check expects `http://localhost:8080/health`. This will cause the health check to fail.
- **Impact**: Deployment will fail until a proper health endpoint is implemented. While this is not a direct security flaw, it indicates a mismatch between deployment expectations and actual service capabilities.
- **Mitigation**: Ensure the backend implements a `/health` endpoint or adjust the health check to match the actual service.

## SEVERITY
- **Maximum severity**: **P2** (maintenance debt)
- All remaining findings are P2 (should be addressed but do not constitute an immediate security threat).

## VERDICT
**PASS** – The P1 findings from the previous review have been adequately addressed. Remaining issues are classified as P2 (maintenance debt) and do not block the feature.

## FINDINGS_CREATED
village-j7a village-9oh village-2wq village-dw8
# Security Review - Feature F005 (Deployment Simplification)

## SCOPE
Files modified/added:
- `.github/workflows/deploy.yml` (modified)
- `scripts/deploy.sh` (modified)
- `README.md` (modified)
- `docs/deployment.md` (new)
- `Makefile` (unchanged but in scope)
- `Dockerfile` (unchanged)
- `docker-compose.yml` (unchanged)

## RISK MAP
Top 3 security risk areas:
1. **SSH Key Management & Host Verification**: SSH deployment introduces private key handling and potential MITM risks.
2. **Secrets Exposure & Credential Handling**: Default database credentials and secret leakage risks.
3. **Container Security & Default Configurations**: Containers running as root, unreliable health checks.

## EVIDENCE

### Finding 1: Backend Docker container runs as root
- **File**: `Dockerfile:11-16`
- **Issue**: No `USER` directive specified; container runs as root.
- **Impact**: Increased attack surface if container is compromised.

### Finding 2: Default PostgreSQL credentials in docker-compose.yml
- **File**: `docker-compose.yml:5-7`
- **Issue**: Hardcoded default credentials (`village:village`).
- **Impact**: If deployed to production without overriding, database is exposed with weak credentials.

### Finding 3: SSH deployment lacks host key verification
- **File**: `.github/workflows/deploy.yml:49-63`
- **Issue**: `appleboy/ssh-action` used without `StrictHostKeyChecking` or known_hosts.
- **Impact**: Potential MITM attack during SSH connection.

### Finding 4: Unreliable backend health check
- **File**: `docker-compose.yml:35-38`
- **Issue**: Health check uses `ps aux | grep village | grep -v grep` which is unreliable.
- **Impact**: May not detect unhealthy backend states.

### Finding 5: Removal of Docker registry reduces image integrity guarantees
- **Issue**: Elimination of registry dependency removes image signing and integrity checks.
- **Impact**: Increased risk of supply chain attacks; no verification of built images.

## SEVERITY
1. Default PostgreSQL credentials: **P1** (breaks on edge case – deployment without overriding credentials)
2. SSH host key verification missing: **P1** (breaks on edge case – MITM in network)
3. Backend container runs as root: **P2** (maintenance debt)
4. Build integrity reduction: **P2** (maintenance debt)
5. Unreliable health check: **P3** (style)

## VERDICT
**FAIL** – Maximum severity is P1 (exploitable in edge cases).

## FINDINGS_CREATED
village-bna village-sno village-akd village-s7e village-gyo
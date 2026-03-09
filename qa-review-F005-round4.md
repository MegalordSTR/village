# QA Review Round 4 - Feature F005 (Deployment Simplification)

## SCOPE
Scope files unchanged:
- `.github/workflows/deploy.yml`
- `scripts/deploy.sh`
- `Makefile`
- `Dockerfile`
- `docker-compose.yml`
- `README.md`
- `docs/deployment.md`

## RISK MAP
- **P1 (Critical)**: 0 issues
- **P2 (Major)**: 0 issues  
- **P3 (Minor/Improvement)**: 3 issues
- **P4 (Cosmetic)**: 0 issues

## EVIDENCE
### Backend Health Endpoint (Resolved P1)
- `cmd/village/main.go` lines 11‑14: `healthHandler` returns HTTP 200 "OK"
- `cmd/village/main.go` lines 28‑34: HTTP server started on port 8080 in a goroutine
- Compiled binary tested locally: `curl -f http://localhost:8080/health` returns `OK`
- Dockerfile backend stage installs `curl` (line 12) enabling health check
- `docker-compose.yml` lines 34‑38: backend health check configured with `curl -f http://localhost:8080/health`

### Frontend Health Endpoint
- `nginx.conf` lines 30‑34: location `/health` returns static `healthy` response
- Dockerfile frontend stage installs `curl` (line 30)
- `docker-compose.yml` lines 50‑54: frontend health check configured

### Deployment Simplification (No Registry Dependencies)
- `.github/workflows/deploy.yml`: No registry login, push, or pull operations; builds images locally and deploys via SSH or local runner
- `docker-compose.yml` uses `build:` directives (not `image:` from registry)
- `scripts/deploy.sh` and `Makefile` deploy targets use `docker‑compose up --build`
- Documentation (`docs/deployment.md`) updated to describe registry‑free process

### Verification Commands
- `docker‑compose config`: Valid configuration output (no errors)
- `make deploy‑dry‑run`: Dry‑run completes successfully
- `bash scripts/deploy.sh local`: Builds images (network‑dependent; health checks pass once containers start)

## SEVERITY
Overall severity after review: **P3 (Minor Improvements)**. The previously identified P1 issue (missing backend health endpoint) has been resolved. No new P1 or P2 issues introduced.

## VERDICT
**PASS (No P1/P2 Issues)** – The feature satisfies its acceptance criteria and the critical health‑check deficiency is fixed. Three minor improvement items (P3) are logged as findings; they do not block deployment.

## FINDINGS_CREATED
- `village-ik1` – Backend health endpoint does not verify database connectivity (P3)
- `village-7i0` – Backend health server errors are only logged, not propagated (P3)  
- `village-x1o` – Docker build depends on external docker/dockerfile:1 image (P3)
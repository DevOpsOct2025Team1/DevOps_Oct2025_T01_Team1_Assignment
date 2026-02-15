# Devops Assignment

[![CodeQL SAST Scan](https://github.com/DevOpsOct2025Team1/DevOps_Oct2025_T01_Team1_Assignment/actions/workflows/codeql.yml/badge.svg)](https://github.com/DevOpsOct2025Team1/DevOps_Oct2025_T01_Team1_Assignment/actions/workflows/codeql.yml)
[![API Gateway](https://codecov.io/gh/DevOpsOct2025Team1/DevOps_Oct2025_T01_Team1_Assignment/branch/main/graph/badge.svg?flag=api-gateway)](https://codecov.io/gh/DevOpsOct2025Team1/DevOps_Oct2025_T01_Team1_Assignment)
[![Auth Service](https://codecov.io/gh/DevOpsOct2025Team1/DevOps_Oct2025_T01_Team1_Assignment/branch/main/graph/badge.svg?flag=auth-service)](https://codecov.io/gh/DevOpsOct2025Team1/DevOps_Oct2025_T01_Team1_Assignment)
[![User Service](https://codecov.io/gh/DevOpsOct2025Team1/DevOps_Oct2025_T01_Team1_Assignment/branch/main/graph/badge.svg?flag=user-service)](https://codecov.io/gh/DevOpsOct2025Team1/DevOps_Oct2025_T01_Team1_Assignment)
[![File Service](https://codecov.io/gh/DevOpsOct2025Team1/DevOps_Oct2025_T01_Team1_Assignment/branch/main/graph/badge.svg?flag=file-service)](https://codecov.io/gh/DevOpsOct2025Team1/DevOps_Oct2025_T01_Team1_Assignment)
[![Frontend](https://codecov.io/gh/DevOpsOct2025Team1/DevOps_Oct2025_T01_Team1_Assignment/branch/main/graph/badge.svg?flag=frontend)](https://codecov.io/gh/DevOpsOct2025Team1/DevOps_Oct2025_T01_Team1_Assignment)

This assignment project uses a microservices architecture.
## Services

- **User Service** (Port 8080): Manages user data with MongoDB
- **Auth Service** (Port 8081): Handles authentication with JWT
- **API Gateway** (Port 3000): HTTP/REST interface for clients
- **File Service** (Port 50054): Handles files uploads, downloads.
- **Frontend** (Port 5173): React application for user interface

## Prerequisites

### For Development (Docker Compose - Recommended)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) or Docker + Docker Compose

### Serve all services
```bash
make serve
# or
./nx run-many --target=serve --all
```

This starts all services via NX in development mode.

### Docker Compose
To run with Docker Compose instead:

Ensure you have created your .env file, then run:
```bash
docker compose up
```

This starts:
- API Gateway: http://localhost:3000

**Stop everything:**
```bash
docker compose down
```

**View logs:**
```bash
docker compose logs -f
```

### Run services individually:
```bash
./nx serve user-service

./nx serve auth-service

./nx serve api-gateway
,
./nx serve frontend
```

## Frontend Development

### Running the Frontend Locally

The frontend is a React 19 application built with React Router v7 and Tailwind CSS v4.
 
**Configure the API endpoint (optional):**
   
   By default, the frontend connects to the API Gateway at `http://localhost:3001`. To change this:
   
   - Create a `.env` file in `apps/frontend/`:
     ```
     VITE_API_BASE_URL=http://localhost:3001
     ```

## CI/CD Pipeline

### Pipeline Overview

The pipeline is built on GitHub Actions with NX-based affected project detection, multi-layer security scanning, and automated deployments via SSH over Cloudflare Tunnel.

```
PR / Push to main
  ├── Unit Tests (unit-test.yml)
  ├── Integration Tests (integration-test.yml)
  ├── CodeQL SAST Scan (codeql.yml)
  ├── Gitleaks Secret Scan (gitleaks.yml)
  └── DAST via OWASP ZAP (dast.yml)
         │
         ▼ (on success, main branch only)
  Build & Push Docker Images (container-images.yml)
  ├── Push to GHCR with :dev tag
  ├── Trivy vulnerability scan
  └── Auto-deploy to dev environment (deploy.yml)
         │
         ▼ (on git tag: {service}/v{version}[-staging])
  Release Build (container-images-releases.yml)
  ├── -staging suffix → deploy to staging
  └── No suffix → deploy to production + GitHub Release
```

### GitHub Secrets Setup

To get the pipeline working on a new repository, configure these secrets under **Settings > Secrets and variables > Actions**.

#### Email Notifications (all workflows)

| Secret              | Description                                 |
|---------------------|---------------------------------------------|
| `SMTP_HOST`         | SMTP server hostname                        |
| `SMTP_PORT`         | SMTP server port                            |
| `SMTP_USER`         | SMTP authentication username                |
| `SMTP_PASS`         | SMTP authentication password                |
| `NOTIFY_EMAIL_FROM` | Sender email address                        |
| `DEV_EMAIL_TO`      | Developer team email (build failure alerts) |
| `QA_EMAIL_TO`       | QA team email (security scan alerts)        |

#### Code Coverage

| Secret          | Description                                                                |
|-----------------|----------------------------------------------------------------------------|
| `CODECOV_TOKEN` | Token from [codecov.io](https://codecov.io) for uploading coverage reports |

#### Deployment

| Secret            | Description                                         |
|-------------------|-----------------------------------------------------|
| `SSH_PRIVATE_KEY` | Private key for SSH access to the deployment server |
| `SSH_HOSTNAME`    | Deployment server hostname                          |
| `SSH_USER`        | SSH username on the deployment server               |

#### GitHub Environments

Create the following environments under **Settings > Environments**:

- **dev** — set `DEPLOY_URL` to the dev server URL
- **staging** — set `DEPLOY_URL` to the staging server URL
- **production** — set `DEPLOY_URL` to the production server URL

### Workflow Files Reference

| Workflow          | File                            | Trigger                       | Purpose                                              |
|-------------------|---------------------------------|-------------------------------|------------------------------------------------------|
| Unit Tests        | `unit-test.yml`                 | PR, push to main              | Run unit tests on affected services, upload coverage |
| Integration Tests | `integration-test.yml`          | PR, push to main              | Run integration tests on affected services           |
| CodeQL            | `codeql.yml`                    | PR to main, push to main      | SAST scan for Go, JS, Python                         |
| Gitleaks          | `gitleaks.yml`                  | PR, push to main              | Scan for hardcoded secrets                           |
| DAST              | `dast.yml`                      | PR to main, push to main      | OWASP ZAP scan against running services              |
| Container Images  | `container-images.yml`          | After unit tests pass on main | Build, push to GHCR, Trivy scan, deploy to dev       |
| Release Images    | `container-images-releases.yml` | Git tag `*/v*`                | Build release images, deploy to staging/production   |
| Deploy            | `deploy.yml`                    | Called by other workflows     | SSH deploy via Cloudflare Tunnel                     |
| Detect Affected   | `detect-affected.yml`           | Called by other workflows     | NX-based affected project detection                  |

### Deployment Server Setup

The deployment server must have:

1. **Docker and Docker Compose** installed
2. **SSH access** configured with the public key matching `SSH_PRIVATE_KEY`
3. **Cloudflare Tunnel** (`cloudflared`) routing traffic to the server
4. **Directory structure:**
   ```
   ~/apps/
   ├── dev/
   │   └── compose.dev.yml       # or symlink to compose.yml
   ├── staging/
   │   └── compose.staging.yml
   └── production/
       └── compose.production.yml
   ```
5. **Environment variables** set on the server (or in `.env` files next to compose files):
   - `MONGODB_URI`, `MONGODB_DATABASE`
   - `JWT_SECRET`, `JWT_EXPIRY`
   - `S3_ENDPOINT`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`, `S3_BUCKET_NAME`
   - `AXIOM_API_TOKEN`, `AXIOM_ENDPOINT`, `AXIOM_DATASET`
   - `USER_SERVICE_DEFAULT_ADMIN_USERNAME`, `USER_SERVICE_DEFAULT_ADMIN_PASSWORD`

### Container Registry

All images are pushed to **GitHub Container Registry (GHCR)** at:

```
ghcr.io/devopsoct2025team1/devops_oct2025_t01_team1_assignment/{service}
```

Tag conventions:
- `:dev` — latest build from main
- `:main-{sha}` — pinned to a specific commit
- `:staging` — current staging release
- `:latest` — current production release
- `:{version}` — specific version (e.g., `v1.2.0`)

### Releasing a New Version

To deploy a service to **staging**:
```bash
git tag api-gateway/v1.2.0-staging
git push origin api-gateway/v1.2.0-staging
```

To deploy a service to **production**:
```bash
git tag api-gateway/v1.2.0
git push origin api-gateway/v1.2.0
```

This triggers `container-images-releases.yml`, which builds the image, pushes it, and deploys to the appropriate environment. Production releases also create a GitHub Release with auto-generated notes.

### Dependency Management

Dependabot is configured in `.github/dependabot.yml` to create weekly PRs for:

| Ecosystem   | Services                                        |
|-------------|-------------------------------------------------|
| Go modules  | api-gateway, auth-service, common, user-service |
| Bun         | frontend                                        |
| uv (Python) | file-service                                    |

### Security Scanning

The pipeline includes four layers of security scanning:

1. **CodeQL (SAST)** — Static analysis for Go, JavaScript, and Python. Results appear in the GitHub Security tab.
2. **Gitleaks** — Scans git history for hardcoded secrets and credentials.
3. **Trivy** — Scans built Docker images for known vulnerabilities. Results uploaded to GitHub Security tab.
4. **OWASP ZAP (DAST)** — Dynamic scan against the running application using the OpenAPI spec.

### Code Coverage

Coverage is tracked via [Codecov](https://codecov.io) with per-service flags configured in `codecov.yml`. Each service's coverage is reported independently using carryforward flags so partial CI runs don't reset other services' coverage.

### Adding a New Service

To add a new service to the pipeline:

1. Create the service under `apps/{service-name}/` with a `Dockerfile` and `project.json`
2. Add test targets (`test`, `integration-coverage`) to the service's `project.json`
3. The `detect-affected.yml` workflow will automatically pick it up if it has a Dockerfile
4. Add a Codecov flag in `codecov.yml`
5. Add a Dependabot entry in `.github/dependabot.yml`
6. Add the service to all compose files (`compose.yml`, `compose.staging.yml`, `compose.production.yml`)
# GitHub Actions CI/CD

## Workflow Architecture

Titan uses GitHub Actions for continuous integration and deployment. Workflows are organized into reusable composites to reduce duplication across 15+ engine repositories.

## CI Workflows

- **Build & Lint**: Triggered on every PR. Runs `golangci-lint`, `ruff`, `clippy`, `eslint` depending on language.
- **Unit Tests**: Runs with coverage thresholds enforced per engine (minimum 80% line coverage).
- **Integration Tests**: Executes against ephemeral environments with testcontainers.
- **Security Scan**: Trivy on images, `gitleaks` on secrets, `semgrep` on SAST rules.
- **SBOM Generation**: `syft` produces SPDX-formatted SBOM attached as build artifact.

## CD Workflows

- **Build & Push**: On merge to `main` — builds multi-arch images, signs with Cosign, pushes to ECR.
- **Deploy Staging**: ArgoCD syncs the staging cluster automatically; smoke tests validate.
- **Deploy Production**: Manual approval gate. Runs blue-green deployment with automated rollback criteria.
- **Canary Analysis**: After deploy, runs metric comparison (p99 latency, error rate, CPU) vs. baseline.

## Self-Hosted Runners

GPU-intensive workflows (model training, chaos tests) run on self-hosted runners in AWS ECG (G5 instances). Standard workflows use GitHub-hosted runners.

## Secrets & Caching

All secrets stored in GitHub Actions secrets scoped to environment. OIDC token exchange authenticates to AWS without long-lived credentials. Docker layer caching via registry cache. Go module, pip, and npm caches restored on `setup-*` actions with key derived from lockfile hash.
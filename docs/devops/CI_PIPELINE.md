# CI Pipeline

## Trigger

The CI pipeline runs automatically on every pull request targeting `main`, `staging`, or `release/*` branches. It can also be triggered manually via workflow dispatch.

## Stages

### 1. Lint & Format
Static analysis: `golangci-lint` (Go), `clippy` (Rust), `ruff` / `mypy` (Python), `eslint` (TypeScript). Format checking: `gofmt`, `rustfmt`, `black`, `prettier`. Style enforcement via `editorconfig-checker`.

### 2. Build
Compile all binaries with `-race` flag for data race detection. Build Docker images (development variant with debug symbols). Cache build artifacts for downstream stages.

### 3. Unit Tests
Run per-engine with `go test -count=1 -race -coverprofile=coverage.out`. Coverage must meet minimum thresholds. Results published as PR comment via `test-reporter`.

### 4. Security Scan
SAST via Semgrep with custom Titan rules and OWASP Top 10. Secrets scanning via Gitleaks. Dependency scanning via Dependabot alerts. License compliance via FOSSA.

### 5. Integration Tests
Launches ephemeral test environment via Docker Compose or Kind. Runs API contract tests (Pact), database migration tests, and cross-engine integration suites. Tears down environment regardless of pass/fail.

### 6. Artifact Publishing (Post-merge only)
Builds production images (distroless, stripped, multi-arch). Signs images and attachments with Cosign. Generates and uploads SBOM to Dependency-Track. Tags commit with `build-{short-sha}`.
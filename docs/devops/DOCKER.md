# Docker Containers

## Philosophy

Every Titan service is packaged as a distroless or Alpine-based Docker image. Images are minimal, signed, and scanned for vulnerabilities before promotion to any environment.

## Base Images

- **Go services**: `golang:1.22-alpine` (build) → `gcr.io/distroless/base` (runtime)
- **Rust services**: `rust:1.77-slim` (build) → `gcr.io/distroless/cc` (runtime)
- **Python/ML services**: custom `titan/python-base` with pinned system deps
- **Node.js services**: `node:24-alpine` with production-only dependencies

## Image Requirements

No shell, no package manager in production images. USER directive switches to non-root (UID 1001). HEALTHCHECK instruction defined for every image. Labels include `org.opencontainers.image.source`, `org.opencontainers.image.revision`, and `org.opencontainers.image.created`.

## Build Process

Multi-stage builds ensure separation of build and runtime environments. Layer caching is optimized: system dependencies first, vendored libraries second, application code last. Images are built via `docker buildx` for multi-architecture support (amd64 + arm64).

## Registry & Signing

Images are pushed to Amazon ECR with tag immutability enabled. Every image is signed using Cosign, and admission controllers enforce that only signed images from the Titan registry can be deployed.

## Vulnerability Scanning

Trivy runs on every PR to gate images with critical or high CVEs. Amazon Inspector continuously scans the registry and raises findings. Weekly full-rescan with Grype and Snyk.
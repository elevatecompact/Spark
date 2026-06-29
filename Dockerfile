# ============================================================================
# Spark - Multi-stage Dockerfile for Go microservices
# ============================================================================
# Build with: docker build --build-arg SERVICE=identity-service -t spark/identity-service .
# ============================================================================

# -- Stage 1: Build ----------------------------------------------------------
FROM golang:1.25-alpine AS builder

ARG SERVICE
ENV SERVICE=${SERVICE}
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOPROXY=https://proxy.golang.org,direct

WORKDIR /workspace

# Install system dependencies
RUN apk add --no-cache \
    ca-certificates \
    git \
    make \
    protobuf-dev \
    protoc \
    curl \
    upx \
    && update-ca-certificates

# Copy go.mod and go.sum first for layer caching
COPY services/${SERVICE}/go.mod services/${SERVICE}/go.sum ./services/${SERVICE}/
COPY api/go.mod api/go.sum ./api/

# Download dependencies
RUN cd services/${SERVICE} && go mod download && go mod verify

# Copy source code
COPY services/${SERVICE}/ ./services/${SERVICE}/
COPY api/ ./api/
COPY packages/ ./packages/

# Build the binary
RUN cd services/${SERVICE} && \
    go build \
        -ldflags="-s -w \
            -X github.com/spark-platform/services/${SERVICE}/internal/version.Version=$(git describe --tags --always 2>/dev/null || echo dev) \
            -X github.com/spark-platform/services/${SERVICE}/internal/version.Commit=$(git rev-parse --short HEAD 2>/dev/null || echo unknown) \
            -X github.com/spark-platform/services/${SERVICE}/internal/version.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')" \
        -o /app/service \
        ./cmd/...

# Strip and compress binary
RUN upx --best --lzma /app/service 2>/dev/null || true

# -- Stage 2: Runtime ---------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

ARG SERVICE
LABEL org.opencontainers.image.title="Spark ${SERVICE}"
LABEL org.opencontainers.image.description="Spark ${SERVICE} microservice"
LABEL org.opencontainers.image.source="https://github.com/spark-platform/spark"
LABEL org.opencontainers.image.licenses="MIT"

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/service /app/service

# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Use non-root user
USER 65532:65532

EXPOSE 8080
EXPOSE 9090

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ["/app/service", "healthcheck"]

ENTRYPOINT ["/app/service"]


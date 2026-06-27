# Nexus Engine — Media Storage

## Purpose

Nexus is Titan's media storage engine. It provides an S3-compatible object storage layer for all platform media assets — uploaded videos, transcoded renditions, thumbnails, profile images, and archival data. Nexus abstracts multi-region replication and storage tiering behind a unified API.

## Architecture

Nexus wraps MinIO for S3-compatible storage with a Titan-specific metadata layer that tracks content hash, replication status, and lifecycle policies. Storage is organized in a multi-tier hierarchy: hot (SSD-backed), warm (S3 Standard), cold (S3 Glacier).

## Tech Stack

- **Language**: Go
- **Storage Core**: MinIO (S3-compatible API)
- **Backend**: AWS S3 (primary), with GCS and Azure Blob as replication targets
- **Metadata DB**: PostgreSQL for content catalog and replication state
- **CDN**: Velocity engine integration for origin pull
- **Encryption**: Server-side encryption with AWS KMS (AES-256)

## Key Features

- **S3-compatible API**: Standard S3 API for all uploads and downloads — works with any S3 SDK
- **Multi-region replication**: Automatic content replication across 3+ AWS regions
- **Lifecycle management**: Automatic tiering based on access patterns (hot to warm to cold)
- **Deduplication**: Content-addressed storage — identical uploads are stored once with reference counting
- **Presigned URLs**: Time-limited URLs for direct uploads and downloads
- **Access control**: IAM-style bucket policies and user-level access controls
- **Bandwidth limiting**: Per-user and per-bucket bandwidth throttling
- **Integrity checks**: Automatic checksum verification on upload and periodic scrub of stored objects

## Performance Targets

| Metric | Target |
|--------|--------|
| Upload throughput (per connection) | 100 MB/s |
| Download throughput (per connection) | 200 MB/s |
| P99 read latency (hot tier) | < 10ms |
| Replication lag (cross-region) | < 5 seconds |
| Durability | 99.999999999% (11 9's) |
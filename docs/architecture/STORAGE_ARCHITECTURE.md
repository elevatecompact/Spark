# Storage Architecture

Spark employs a multi-model storage architecture optimized for different data types and access patterns. The platform uses S3-compatible object storage powered by the Nexus engine for media assets, PostgreSQL for relational data, and purpose-built stores for specific workloads.

## Object Storage: Nexus Engine

The Nexus engine is Spark's high-performance storage layer for video segments, thumbnails, and archival content. It implements the S3 API for compatibility with existing tools and libraries.

### Nexus Architecture
`
Client → Nexus Proxy → Storage Nodes (NVMe) → Cold Tier (HDD/Cloud)
`

- **Hot Tier**: NVMe flash on edge nodes for frequently accessed content
- **Warm Tier**: SSD-based storage in regional data centers
- **Cold Tier**: S3-compatible cloud storage for archival content

### Key Features
- **Erasure Coding**: Data is striped across storage nodes with 12+4 erasure coding for durability
- **Geo-Replication**: Content is automatically replicated to at least two regions
- **Prefix-Based Routing**: Storage paths are hash-partitioned for even distribution
- **Consistent Hashing**: Node additions and removals minimize data movement

## Relational Storage

PostgreSQL clusters are deployed per microservice with logical replication for cross-region availability.

### Sharding Strategy
User data is sharded by user ID hash across 64 logical shards per region. Each shard runs as a separate PostgreSQL instance with its own replica set.

### Time-Series Data
Stream metrics, viewer engagement, and quality-of-service data are stored in TimescaleDB for efficient time-based queries and automatic data retention policies.

## Media Asset Storage

Videos are stored as segmented files in Nexus object storage:
`
/streams/{streamId}/{segmentIndex}.ts
/streams/{streamId}/manifest.m3u8
/thumbnails/{streamId}/{timestamp}.jpg
`

Metadata about each asset is stored in PostgreSQL for indexing and search. Content lifecycle policies automatically transition cold streams to cheaper storage tiers.

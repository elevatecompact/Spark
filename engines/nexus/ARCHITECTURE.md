# Architecture

Nexus uses content-addressable storage. Uploaded files are SHA-256 hashed and stored at hash-prefix/hash/filename in S3. Metadata indexed in PostgreSQL with Redis caching. Transformation pipeline runs asynchronously - thumbnails, transcoded copies enqueued in Redis. Lifecycle manager enforces retention policies and migrates between hot (SSD S3), warm (standard S3), and cold (Glacier) tiers based on access recency.

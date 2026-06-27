# viewer-service — Troubleshooting
## Watch events not recording: Check Kafka producer, DB pool saturation, event batching buffer. Verify: ./viewer check kafka.
## Bookmarks missing: Cache eviction failed, content deleted. Verify DB existence, rebuild cache: ./viewer cache warm --type bookmarks --id {viewerId}.
## High write latency: WAL growth on PG, missing partition. Run VACUUM ANALYZE, verify monthly partitioning active, increase max_wal_size.

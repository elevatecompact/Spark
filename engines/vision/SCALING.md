# Scaling

Vision scales across multiple ClickHouse nodes using distributed tables sharded by tenant ID. Kafka partitions match shards for data locality. Consumer instances run on each ClickHouse node. Query service nodes route queries to correct shard. Multi-region uses per-region clusters with federation. Cold data auto-migrated to S3 using ClickHouse tiered storage with zero-copy replication.

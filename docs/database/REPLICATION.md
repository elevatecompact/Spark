# Replication

SPARK implements replication across all database systems to ensure high availability, disaster recovery, and read scaling.

## PostgreSQL Replication

PostgreSQL uses streaming replication with one primary and two synchronous replicas in different availability zones. The synchronous_standby_names configuration ensures at least one replica acknowledges every write before the transaction commits, providing zero data loss. Replication slots prevent WAL from being removed before replicas consume it.

Read replicas serve analytics queries, reporting, and read-heavy API endpoints. Applications use a read/write splitter that routes SELECT queries to replicas and all other statements to the primary. Replication lag is monitored through the pg_stat_replication view with PagerDuty alerts if lag exceeds 10 seconds.

## ClickHouse Replication

ClickHouse uses ReplicatedMergeTree engines with ZooKeeper-based coordination. Each shard has two replicas across availability zones. Inserts are asynchronous and eventually consistent, with most replicas converging within one second.

## OpenSearch Replication

Each OpenSearch index has two replica shards distributed across nodes and availability zones. The shard allocation awareness ensures replicas never reside on the same node or zone as their primary. Rolling restarts update nodes one at a time without cluster downtime.

## Redis Replication

Redis Sentinel manages automatic failover. Replication is asynchronous with 1-second AOF fsync for durability. Client connections use the Sentinel-discovered master address and retry on failover.

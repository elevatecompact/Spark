# ADR-0004: Apache Kafka for Event-Driven Messaging

## Status

Accepted

## Context

Spark's platform generates hundreds of discrete event types including content uploads, livestream state changes, viewer interactions, moderation actions, payout triggers, and recommendation signals. These events must be reliably delivered to multiple consumers with different latency and throughput requirements. Some consumers need replay capability for state reconstruction, while others require exactly-once semantics for financial transactions. The evaluation compared Apache Kafka, RabbitMQ, Amazon SQS/SNS, and Apache Pulsar. RabbitMQ excelled at complex routing but struggled with throughput at Spark's scale. SQS/SNS simplified operations but imposed vendor lock-in and limited retention. Pulsar offered strong features but a smaller ecosystem. Kafka provided the best combination of throughput, durability, replayability, and ecosystem maturity.

## Decision

Adopt Apache Kafka as the central event backbone. Deploy a multi-cluster topology with a primary aggregation cluster and regional clusters for data locality. Topics are partitioned by event key with compaction enabled for keyed event streams. Schema Registry enforces Avro schema evolution with full backward and forward compatibility. Kafka Connect handles ingestion from PostgreSQL (Debezium CDC) and sinks to OpenSearch, ClickHouse, and object storage. Exactly-once semantics are enabled for payment and subscription topics. Consumer groups are monitored through Kafka Lag Exporter with alerts on threshold violations.

## Consequences

### Positive
- High throughput with horizontal partitioning scales to millions of events per second
- Durable log with configurable retention supports replay and state reconstruction
- Rich connector ecosystem simplifies integration with databases and search engines
- Strong community and enterprise support

### Negative
- Operational complexity of managing Kafka clusters, especially across regions
- Client library versioning requires careful coordination across services
- Schema registry adds a dependency that must be highly available
- Storage costs grow with retention duration and replication factor

# ADR-0009: OpenSearch for Search

## Status

Accepted

## Context

Spark requires a full-text search capability spanning content titles, descriptions, user profiles, livestream metadata, and transcriptions. Search must support fuzzy matching, multi-field ranking, faceted filtering, and typo-tolerant autocomplete. The platform also requires log analytics across microservice logs for operational debugging. Elasticsearch and OpenSearch were the primary candidates. When this decision was made, Elasticsearch had moved to the SSPL license, creating uncertainty for third-party consumption. OpenSearch, as an Apache 2.0-licensed fork, provided equivalent functionality with a fully open-source governance model under the OpenSearch Foundation. Amazon OpenSearch Service offered a managed option, while self-hosted OpenSearch provided deployment flexibility.

## Decision

Adopt OpenSearch (forked from Elasticsearch 7.10) as the search and log analytics platform. Deploy a multi-node cluster with dedicated master, data, and coordinating node roles. Indices use index-lifecycle management with hot-warm-cold phases: hot nodes on SSD for recent data, warm nodes on HDD with replica compression for medium-term data, and cold/delete policies for older data. Search indices are populated via Kafka Connect sinks consuming domain events. Log analytics indices receive structured logs via Filebeat and Fluentd. All indices are alias-based for zero-downtime reindexing. Cluster monitoring uses OpenSearch Dashboards with anomaly detection plugins.

## Consequences

### Positive
- Full Apache 2.0 license with open governance avoids licensing uncertainty
- Feature parity with Elasticsearch for search, analytics, and observability use cases
- Rich query DSL supports relevance scoring, faceted search, and aggregation
- Index lifecycle management automates data tiering and retention

### Negative
- Cluster maintenance requires expertise in shard sizing, index tuning, and cluster health
- Write-heavy workloads can trigger segment merges that impact query performance
- OpenSearch Dashboards lags behind Kibana in some observability features
- Migration path from managed Elasticsearch requires careful planning

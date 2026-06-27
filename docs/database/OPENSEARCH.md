# OpenSearch

OpenSearch powers SPARK's full-text search, content discovery, and log analytics. It provides fast, scalable search across the platform's content catalog, user profiles, and operational logs.

## Cluster Architecture

The OpenSearch cluster runs across six nodes in a three-tier topology: dedicated master nodes for cluster management, dedicated data nodes for indexing and querying, and dedicated coordinating nodes for request routing. Index templates define shard counts (5 primary, 2 replica) and analyzer configurations for each document type.

## Index Design

The primary content index includes fields for title, description, tags, creator name, transcript text, and content category. Custom analyzers handle multilingual stemming, synonym expansion, and n-gram partial matching for autocomplete. Field-level mappings specify English text analysis for content fields and keyword analysis for filterable fields like category and status.

## Search Features

The search pipeline supports multi-match queries across title, description, and tags with tiered boosting. Function score queries incorporate popularity signals, recency, and user preferences for personalized results. Highlighting returns matching fragments in search results. Suggesters power the autocomplete feature with completion suggestions based on the top 10,000 most searched terms.

## Monitoring

Cluster health, search latency (p99 < 100ms), indexing rate, and merge throughput are monitored through Prometheus. Slow logs capture queries exceeding 500ms for performance optimization.

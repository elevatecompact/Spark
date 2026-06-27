# Search Architecture

Spark uses OpenSearch for full-text search, content discovery, and real-time analytics across the platform. The search infrastructure is designed to index streaming content metadata, user profiles, chat history, and moderation records at scale.

## Cluster Topology

The OpenSearch deployment spans three regions with a primary-active configuration:

`
Region A (Primary)          Region B (Active)          Region C (Active)
  ┌─────────────┐           ┌─────────────┐            ┌─────────────┐
  │ Master +    │◄─Cross──►│ Master +    │◄──Cross──►│ Master +    │
  │ Data Nodes  │   Cluster │ Data Nodes  │   Cluster │ Data Nodes  │
  └─────────────┘  Replication└─────────────┘  Replication└─────────────┘
`

Each cluster runs three dedicated master nodes for cluster management and a configurable number of data nodes scaled by throughput requirements.

## Index Strategy

### Stream Index
Indexes stream metadata including title, description, category, tags, broadcaster name, and language. Optimized for full-text search with custom analyzers per language.

### User Index
Indexes user profiles for discovery. Uses edge n-gram tokenizer for autocomplete functionality and phonetic analyzers for name matching.

### Chat Index
Indexes chat messages for moderation review and historical search. Time-based indices with rollover every 6 hours for high-volume streams.

### Analytics Index
Stores aggregated viewership metrics for dashboard visualization. Pre-aggregated at minute and hour granularities using OpenSearch's composite aggregation framework.

## Search Features

- **Multi-language support**: Custom analyzers for English, Spanish, Mandarin, Arabic, and Hindi
- **Faceted search**: Category, language, viewer count range, and content rating filters
- **Typo tolerance**: Fuzzy matching up to 2 edits for user-facing search
- **Personalized ranking**: Boost results based on user viewing history and location
- **Real-time indexing**: Documents are indexed within 500ms of domain event publication

## Scaling Strategy

Data nodes scale horizontally based on indexing throughput and query concurrency. Hot-warm-cold architecture optimizes storage costs: hot nodes use NVMe, warm uses SSD, and cold indices are stored in S3 with a snapshot lifecycle policy.

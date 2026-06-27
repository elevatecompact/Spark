# Atlas Engine — Content Discovery

## Purpose

Atlas is Titan's content discovery engine. It enables users to explore, search, and navigate the platform's content catalog through faceted search, topic clustering, and personalized browse surfaces.

## Architecture

Atlas indexes all content metadata into Elasticsearch with custom analyzers for multilingual text, tags, and audio transcripts. The browse layer generates curated surfaces using a mix of editorial rules and machine-learned topic models.

## Tech Stack

- **Language**: Go
- **Search Index**: Elasticsearch (multi-node cluster with hot-warm architecture)
- **Cache**: Redis for browse surface caching and trending computation
- **Queue**: RabbitMQ for async indexing jobs
- **Natural Language**: spaCy for entity extraction, topic modeling via BERT

## Key Features

- **Full-text search**: Multi-language search with typo tolerance, synonym expansion, and phrase matching
- **Faceted navigation**: Filter by category, language, duration, upload date, popularity, and custom tags
- **Trending detection**: Real-time trending computation based on engagement velocity (views, shares, comments in sliding window)
- **Topic clusters**: Automated topic grouping of related content using content embeddings
- **Personalized browse**: Reranked browse surfaces using Oracle's recommendation scores
- **Autocomplete**: Prefix-based search suggestions with popularity weighting

## Performance Targets

| Metric | Target |
|--------|--------|
| Search p99 latency | < 100ms |
| Indexing throughput | 10,000 docs/s |
| Browse surface generation | < 200ms |
| Trending freshness | < 1 minute |
| Catalog coverage | 100% of published content indexed within 60 seconds of publish |
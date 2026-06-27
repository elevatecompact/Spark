# search-service — Database Schema
## Elasticsearch — Primary search index
### Index: creators — Mapping: name(text+keyword), description(text), category(keyword), follower_count(long), verification_status(keyword). Analyzer: english with edge_ngram for autocomplete
### Index: streams — title(text), description(text), creator_name(text), category(keyword), tags(text), viewer_count(long), status(keyword), started_at(date)
### Index: ecordings — title(text), description(text), creator(text), duration(long), view_count(long), upload_date(date)
### Index: clips — title(text), creator(text), stream_title(text), view_count(long), created_at(date)
### Synonym sets stored in Elasticsearch synonyms API
## PostgreSQL — Index metadata, search analytics logs, synonym dictionary
## Redis — Query result cache (TTL 60s), autocomplete cache (TTL 30s), trending queries

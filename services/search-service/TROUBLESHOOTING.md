# search-service — Troubleshooting
## No results for valid query: Index missing, mapping changed, analyzer not producing tokens. Check ES index exist, verify mapping PUT, test analyzer: POST /_analyze.
## Autocomplete slow: Edge n-gram index too large, too many shards, cache cold. Optimize index, increase shards, warm cache.
## Relevance poor: Boosts misconfigured, missing field data, synonym set wrong. Review mapping boosts, verify field data, check synonym file.
## Index lag: Kafka index-worker consumer lag, ES bulk rejection. Scale index-worker pods, increase ES bulk queue size.

# search-service — Configuration
SEARCH_PORT=4016, SEARCH_ES_URLS=https://es1:9200,https://es2:9200, SEARCH_REDIS_URL, SEARCH_KAFKA_BROKERS, ES_INDEX_REPLICAS=1, ES_INDEX_SHARDS=3, SEARCH_RESULTS_PER_PAGE=20, AUTOCOMPLETE_MAX=10, SEARCH_CACHE_TTL=60, REINDEX_BATCH_SIZE=1000
FF: fulltext_search=true, semantic_search=false, autocomplete=true, personalized_ranking=true, synonym_expansion=true
Rate limits: 60 searches/min per user, 300 searches/min per IP. Admin: 1000/min.

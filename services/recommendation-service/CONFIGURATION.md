# recommendation-service — Configuration
REC_PORT=4015, REC_FEATURE_STORE_URL, REC_REDIS_URL, REC_KAFKA_BROKERS, REC_MODEL_PATH=/models/current, REC_EMBEDDING_DIM=128, REC_FEED_SIZE=50 (home), REC_UP_NEXT_SIZE=10, REC_TRENDING_SIZE=100, REC_CACHE_TTL_SECONDS=300, REC_FEATURE_REFRESH_INTERVAL=60
FF: personalization_enabled=true, trending_enabled=true, cold_start_enabled=true, explications_enabled=true, ab_testing_enabled=false
Rate limits: 100 feed requests/min per user, 1000 feedback events/min per user

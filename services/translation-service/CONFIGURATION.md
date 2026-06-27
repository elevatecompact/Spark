# translation-service — Configuration
TRANSLATION_PORT=4017, TRANSLATION_DB_URL, TRANSLATION_REDIS_URL, TRANSLATION_KAFKA_BROKERS, DEEPL_API_KEY, DEEPL_API_URL=https://api.deepl.com/v2, GOOGLE_TRANSLATE_API_KEY, PRIMARY_PROVIDER=deepl, FALLBACK_PROVIDER=google, MAX_TEXT_LENGTH=5000 (chars), BATCH_MAX_SIZE=100, CACHE_TTL_HOURS=24, STREAMING_BUFFER_SIZE=50
FF: auto_translate_enabled=true, realtime_chat_translate=true, translation_memory=true, human_review=false, provider_failover=true
Rate limits: 1000 translations/min (total), 100 translations/min per user, 5000 chars/request

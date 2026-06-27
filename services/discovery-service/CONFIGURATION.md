# discovery-service — Configuration
DISCOVERY_PORT=4026, DISCOVERY_DB_URL, DISCOVERY_REDIS_URL, DISCOVERY_KAFKA_BROKERS, FEED_SIZE_HOME=50, FEED_SIZE_TRENDING=100, FEED_SIZE_CATEGORY=50, TRENDING_WINDOW_HOURS=24 (freshness decay), TRENDING_VELOCITY_WINDOW_MINUTES=15, CACHE_TTL_FEED_SECONDS=120, CACHE_TTL_CATEGORY_SECONDS=60, TRENDING_REFRESH_INTERVAL_SECONDS=60, CATEGORY_TREE_CACHE_TTL_SECONDS=300
FF: personalized_feeds=true, trending_enabled=true, category_browsing=true, editorial_picks=true, collections_enabled=true, holiday_campaigns=false
Rate limits: 200 feed requests/min per user, 500 category browse/min per user, 30 collection views/min per user

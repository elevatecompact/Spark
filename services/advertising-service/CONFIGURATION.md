# advertising-service — Configuration
ADVERTISING_PORT=4022, ADVERTISING_DB_URL, ADVERTISING_REDIS_URL, ADVERTISING_KAFKA_BROKERS, ADVERTISING_CLICKHOUSE_URL, DEFAULT_BID_CPM_CENTS=500, MIN_BID_CPM_CENTS=50, MAX_DAILY_BUDGET_CENTS=10000000, IMPRESSION_FRAUD_THRESHOLD=0.05 (5% suspicious), CREATOR_REVENUE_SHARE=0.70 (70% to creator), AD_SERVE_TIMEOUT_MS=100
FF: advertising_enabled=true, programmatic_ads=true, direct_ads=true, ad_fraud_detection=true, creator_revenue_sharing=true
Rate limits: 100 ad requests/s per placement (capped), 10 campaign creations/day per advertiser

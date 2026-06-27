# trust-service — Configuration
TRUST_PORT=4027, TRUST_DB_URL, TRUST_REDIS_URL, TRUST_KAFKA_BROKERS, REPUTATION_DECAY_FACTOR=0.95 (daily decay), REPUTATION_RECALCULATION_CRON="0 3 * * *" (daily), TRUST_LEVEL_THRESHOLDS={"low":300,"medium":500,"high":700,"verified":900}, FRAUD_PAYMENT_AMOUNT_THRESHOLD_CENTS=1000000 (), MAX_FAILED_LOGINS_BEFORE_RISK=5, IP_REPUTATION_CHECK_ENABLED=true, DEVICE_FINGERPRINT_ENABLED=true, SIGNAL_RETENTION_DAYS=730
FF: reputation_scoring_enabled=true, realtime_risk_assessment=true, fraud_detection_enabled=true, ip_reputation=true, device_fingerprinting=false, behavioral_analytics=false
Rate limits: 1000 risk assessments/min (total), 100 reputation lookups/s, 50 fraud report/h per user

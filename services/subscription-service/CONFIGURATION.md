# subscription-service — Configuration
SUBSCRIPTION_PORT=4010, SUBSCRIPTION_DB_URL, SUBSCRIPTION_REDIS_URL, BILLING_CYCLE_CRON="0 2 * * *", GRACE_PERIOD_DAYS=3, MAX_ACTIVE_SUBSCRIPTIONS=50, TRIAL_PERIOD_DAYS=0, RETRY_MAX_ATTEMPTS=3, RETRY_INTERVAL_HOURS=24
FF: billing_enabled=true, free_trials=false, proration_enabled=true, grace_period_enabled=true, annual_billing=true
Stripe/PayPal keys from Vault. Proration uses exact day counting.

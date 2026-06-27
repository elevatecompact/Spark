# wallet-service — Configuration
WALLET_PORT=4009, WALLET_DB_URL, WALLET_REDIS_URL (idempotency keys), STRIPE_SECRET_KEY, STRIPE_WEBHOOK_SECRET, PAYPAL_CLIENT_ID/SECRET, PAYOUT_MINIMUM_CENTS=5000, MAX_WALLET_BALANCE_CENTS=100000000, RECONCILIATION_CRON="0 3 * * *"
FF: deposits_enabled=true, withdrawals_enabled=true, crypto_enabled=false, auto_payout=true
All mutations require Idempotency-Key. Keys in Redis with 24h TTL.

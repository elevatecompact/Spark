# payment-service — Configuration
PAYMENT_PORT=4012, PAYMENT_DB_URL, PAYMENT_REDIS_URL (idempotency), STRIPE_SECRET_KEY, STRIPE_WEBHOOK_SECRET, STRIPE_PUBLISHABLE_KEY, PAYPAL_CLIENT_ID/SECRET, PAYPAL_WEBHOOK_ID, DEFAULT_CURRENCY=USD
FF: stripe_enabled=true, paypal_enabled=false, saving_payment_methods=true, refunds_enabled=true, disputes_enabled=true
Idempotency keys expire after 24h in Redis.

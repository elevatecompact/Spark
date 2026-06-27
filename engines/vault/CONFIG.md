# Configuration

| Key | Default | Description |
|-----|---------|-------------|
| stripe.secret_key | - | Stripe API key (from Vault) |
| stripe.webhook_secret | - | Stripe webhook secret |
| paypal.client_id | - | PayPal client ID |
| billing.default_currency | USD | Default currency |
| billing.trial_days | 7 | Default trial period |
| billing.dunning.max_attempts | 3 | Max retry attempts |
| billing.dunning.interval_hours | 72 | Hours between retries |
| reconciliation.schedule | 0 * * * * | Cron schedule |
| idempotency.ttl | 86400 | Idempotency key TTL |

# subscription-service — Troubleshooting
## Not activating: Payment webhook not processed, Kafka delay, grace period still active. Check invoice status, webhook logs, subscription.grace_period_end.
## Billing not running: CronJob pod failed, DB pool exhausted. Check CronJob logs, manually trigger: ./subscription billing run.
## Double billing: Missing idempotency check, retry creating duplicate. Void duplicate invoice, check idempotency keys, issue refund.

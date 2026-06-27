# wallet-service — Troubleshooting
## Transaction stuck pending: Idempotency key collision, Kafka consumer lag, webhook not received. Check status in DB, idempotency key in Redis, consumer group lag. Manual settle via admin.
## Balance discrepancy: Missed event, race condition on concurrent update. Run reconcile, check optimistic lock versions. Escalate to Financial Systems.
## Withdrawal fails: External provider rejection, insufficient balance after pending holds. Check provider response, sum pending holds, verify payment method.

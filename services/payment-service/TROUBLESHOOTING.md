# payment-service — Troubleshooting
## Payment declined: Insufficient funds, bank fraud block, expired method, processor API error. Check intent metadata for processor response, suggest different method.
## Webhook not processing: Endpoint URL changed, signature verification failed, secret rotated. Check webhook_events table, verify STRIPE_WEBHOOK_SECRET match, retry.
## Processor downtime: Stripe/PayPal outage, network issue. Check status page, switch to backup processor via FF, queue for retry.

# Troubleshooting

## Payment failures
1. Check Stripe dashboard for decline codes.
2. Verify webhook signature validation.
3. Check idempotency key collisions.
4. Review payment.failed event.

## Subscription billing missed
1. Check scheduler leader lock.
2. Verify customer payment method valid.
3. Check dunning queue.
4. Manual trigger: POST /v1/invoice/:id/retry.

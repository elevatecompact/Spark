# payment-service — Runbook
## Alerts: PaymentSuccessRate < 95%, StripeAPIErrorRate > 2%, WebhookProcessingLag > 5m, DisputeRate > 1% (fraud pattern)
## Retry webhook: POST /v1/admin/webhooks/{id}/retry
## Manual refund: POST /v1/admin/refund {paymentIntentId, amountCents}
## Switch processor: Set stripe_enabled=false, paypal_enabled=true (emergency)
## Force settle: POST /v1/admin/payment-intents/{id}/settle {status: "succeeded"}

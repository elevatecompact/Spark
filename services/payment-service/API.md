# payment-service — API Contract
## Intents: POST/GET/POST{id}/confirm/POST{id}/cancel /v1/payment-intents
## Methods: POST/GET/DELETE{id}/PATCH{id} /v1/payment-methods
## Refunds: POST /v1/payment-intents/{id}/refund, GET /v1/refunds
## Payouts: POST /v1/payouts, GET /v1/payouts/{id}
## Webhooks: POST /v1/webhooks/stripe, POST /v1/webhooks/paypal
## Admin: GET /v1/admin/processors/status, POST /v1/admin/webhooks/retry/{id}
All amounts in smallest unit. Idempotency-Key header required.

# API

## Payment
- POST /v1/payment/intent - Create payment intent.
- POST /v1/payment/intent/:id/confirm - Confirm payment.
- GET /v1/payment/intent/:id - Retrieve status.
- POST /v1/payment/refund - Issue refund.

## Subscription
- POST /v1/subscription - Create subscription.
- GET /v1/subscription/:id - Retrieve details.
- POST /v1/subscription/:id/cancel - Cancel at period end.
- POST /v1/subscription/:id/update - Change plan.

## Billing
- GET /v1/invoice/:id, POST /v1/invoice/:id/retry.

## Admin
- POST /v1/admin/plan - Create plans.
- GET /v1/admin/reconciliation - Reports.

# subscription-service — API Contract
## Plans: POST/GET/GET{id}/PATCH/DELETE /v1/plans
## User Subs: POST /v1/subscriptions, GET /v1/subscriptions/me, GET /v1/subscriptions/{id}, PUT /v1/subscriptions/{id}/plan (change), POST /v1/subscriptions/{id}/cancel, POST /v1/subscriptions/{id}/reactivate
## Billing: GET /v1/billing/invoices, GET /v1/billing/invoices/{id}, POST /v1/billing/payment-methods, DELETE /v1/billing/payment-methods/{id}
## Benefits: GET /v1/benefits/{subscriptionId}, POST /v1/benefits/verify
## Admin: POST /v1/admin/subscriptions/{id}/override, GET /v1/admin/revenue

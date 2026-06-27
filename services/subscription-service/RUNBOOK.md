# subscription-service — Runbook
## Alerts: BillingSuccessRate < 95%, ChurnRateSpike > 2x, SubscriptionActivationLag > 5m, BillingCronFailure
## Manual billing: ./subscription billing run --date 2026-06-27
## Force cancel: POST /v1/admin/subscriptions/{id}/force-cancel
## Grant free sub: POST /v1/admin/grants {userId, planId, durationDays}
## Refund: POST /v1/admin/invoices/{id}/refund

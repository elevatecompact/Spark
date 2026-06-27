# notification-service — Runbook
## Alerts: PushDeliveryRate < 95%, EmailBounceRate > 3%, SMSSendFailure > 5%, NotificationWorkerQueueDepth > 50000, DigestJobFailure
## Resend: POST /v1/admin/notifications/{id}/resend
## Test push: POST /v1/admin/test-push {userId} — sends test notification to all user devices
## Test email: POST /v1/admin/test-email {email} — sends template test
## Bounce handling: Check SendGrid suppression list, remove invalid emails from push tokens

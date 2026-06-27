# moderation-service — Runbook
## Alerts: ScanLatencyP99 > 2s, AutoActionAccuracy < 90%, ReviewQueueDepth > 1000, MLModelHealthFailure, StreamMonitoringLag > 30s
## Bypass moderation: POST /v1/admin/actions/bypass {contentId} — for legitimate content caught by false positive.
## Override rule: PATCH /v1/admin/rules/{id} {is_active: false} — disable problematic rule temporarily.
## Re-review: POST /v1/admin/queue/{id}/reopen — send item back to queue for re-review.
## Check model: GET /v1/admin/models/status — health and version of deployed ML models.

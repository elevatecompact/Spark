# advertising-service — Runbook
## Alerts: AdServerLatencyP99 > 100ms, ImpressionRecordingLag > 60s, CampaignBudgetOverspend > 1%, FraudRate > 5%, InventoryFillRate < 50%
## Pause campaign: POST /v1/admin/campaigns/{id}/pause — emergency stop for problematic ads.
## Force budget reset: POST /v1/admin/campaigns/{id}/budget-reset {spent_cents: 0}
## Flush impression buffer: POST /v1/admin/impressions/flush — forces pending impressions to ClickHouse.
## Fraud review: GET /v1/admin/fraud/suspicious-impressions — review flagged impressions.

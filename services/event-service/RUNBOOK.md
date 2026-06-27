# event-service — Runbook
## Alerts: TicketOversellDetected (immediate, critical), EventStartFailure > 2%, ReminderJobFailure, TicketPurchaseErrorRate > 3%
## Cancel event: POST /v1/admin/events/{id}/cancel — triggers refunds for all ticket holders.
## Refund all: POST /v1/admin/events/{id}/refund-all — for cancelled events.
## Force start: POST /v1/admin/events/{id}/force-start — if auto-start fails.
## Inventory fix: POST /v1/admin/events/{id}/inventory-rebuild — recalculates ticket counts.

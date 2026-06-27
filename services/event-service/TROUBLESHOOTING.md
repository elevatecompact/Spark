# event-service — Troubleshooting
## Tickets oversold: Race condition on concurrent purchases, inventory check without lock. Check ticket_tiers.quantity_sold vs quantity_total, issue refunds for oversold, fix with optimistic locking.
## Event not starting at scheduled time: Stream not associated, auto-start cron missed, stream-service issue. Verify stream_id on event, manually trigger start, check stream-service health.
## Reminders not sending: CronJob failed, notification-service unavailable, template missing. Check CronJob logs, verify notification-service health, validate reminder template.
## Series not generating next occurrence: Recurrence calculation error, series deactivated. Check series.is_active, test recurrence formula, manually trigger next occurrence.

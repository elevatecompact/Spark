# Events

## Published Events
- analytics.event.ingested - Payload: { eventType, source, size }.
- analytics.query.executed - Payload: { queryHash, durationMs, userId }.
- alert.triggered - Payload: { alertId, condition, currentValue, threshold }.
- alert.resolved, alert.silenced.
- data.latency.high - Kafka consumer lag exceeds threshold.
- storage.capacity.warning - Disk usage > threshold.

## Subscribed Events
- analytics.schema.updated, analytics.retention.updated.

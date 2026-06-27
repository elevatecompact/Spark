# analytics-service — Event Contracts
## Published: analytics.report.ready, analytics.alert.triggered (threshold crossed), analytics.anomaly.detected (statistical anomaly), analytics.dashboard.snapshot
## Consumed (ingested from ALL services): iam.*, creator.*, viewer.*, stream.*, chat.*, wallet.*, subscription.*, gift.*, payment.* — all platform events
## Schema: GenericTrackedEvent {event, properties{}, userId, anonymousId, timestamp, context{ip, userAgent}}

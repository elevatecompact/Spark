# advertising-service — Event Contracts
## Published: advertising.campaign.created, advertising.campaign.activated, advertising.campaign.ended (budget exhausted), advertising.campaign.paused, advertising.impression.recorded, advertising.click.recorded, advertising.inventory.updated
## Consumed: stream.session.started (open ad inventory), viewer.watch.started (targeting signal), analytics.anomaly.detected (fraud alert), creator.channel.created (new inventory source)
## Schema: AdImpressionRecordedEvent {impressionId, campaignId, adUnitId, placementId, userId, costCents, timestamp, servedLatencyMs}

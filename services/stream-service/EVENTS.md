# stream-service — Event Contracts
## Published: stream.session.created, stream.session.started, stream.session.ended, stream.transcoding.completed, stream.recording.ready, stream.health.degraded
## Consumed: moderation.content.flagged (halt stream), identity.user.deleted (end all streams), subscription.tier.activated (allow sub-only stream)
## Schema: StreamSessionStartedEvent {streamId, creatorId, ingestNode, startedAt, initialViewers}

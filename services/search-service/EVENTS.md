# search-service — Event Contracts
## Published: search.query.executed (query, resultIds[], latency, clickedResultId), search.index.updated (doc indexed), search.suggestions.clicked
## Consumed: creator.channel.created (index new), viewer.rating.submitted (update signals), media.content.uploaded (index media), moderation.content.flagged (remove from index), identity.user.deleted (remove user content)
## Schema: SearchQueryExecutedEvent {queryId, query, filters{}, resultCount, topResultIds[], clickedResultId nullable, latencyMs, userId, timestamp}

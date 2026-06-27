# discovery-service — Event Contracts
## Published: discovery.feed.served, discovery.trending.updated, discovery.category.contents.changed, discovery.collection.updated, discovery.editorial.pick.changed
## Consumed: viewer.watch.completed (signal for related content), recommendation.feed.served (incorporate recs), creator.channel.created (new category assignment), stream.session.started (trending velocity signal), analytics.anomaly.detected (trending spike detection)
## Schema: DiscoveryFeedServedEvent {sessionId, userId, feedType, contentIds[], source("algorithmic","editorial","mixed"), latencyMs, servedAt}

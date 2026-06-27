# recommendation-service — Event Contracts
## Published: recommendation.feed.served (recs shown), recommendation.click.recorded, recommendation.feedback.recorded, recommendation.model.deployed (new version), recommendation.features.updated
## Consumed: viewer.watch.completed (training signal), viewer.rating.submitted (preference signal), viewer.reaction.added (engagement signal), subscription.activated (interest signal), creator.channel.created (cold start trigger), search.query.executed (intent signal)
## Schema: RecommendationFeedServedEvent {sessionId, userId, recommendations[{contentId, score, reason}], feedType, servedAt, latencyMs}

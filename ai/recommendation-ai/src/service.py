import datetime
import random

from src.events import EventProducer, NoopProducer
from src.models import (
    BatchInteractionsResponse,
    GetRecommendationsRequest,
    GetRecommendationsResponse,
    GetSimilarRequest,
    GetSimilarResponse,
    GetTrendingRequest,
    GetTrendingResponse,
    PersonalizeFeedRequest,
    PersonalizeFeedResponse,
    RecommendedItem,
    TrackInteractionResponse,
    UserInteraction,
)
from src.repository import RecommendationRepository

MOCK_REASONS = [
    "Because you watched {item}",
    "Trending in your area",
    "Popular in category",
    "Recommended for you",
    "Others also liked this",
    "New and noteworthy",
    "Based on your interests",
]


class RecommendationService:
    def __init__(
        self,
        repository: RecommendationRepository | None = None,
        producer: EventProducer | None = None,
    ) -> None:
        self._repository = repository or RecommendationRepository()
        self._producer = producer or NoopProducer()

    async def get_recommendations(self, request: GetRecommendationsRequest) -> GetRecommendationsResponse:
        items = [
            RecommendedItem(
                item_id=f"rec-{request.user_id}-{request.offset + i}",
                score=round(random.random(), 4),
                reason=random.choice(MOCK_REASONS).format(item=f"item-{random.randint(1, 100)}"),
                position=i,
                metadata={"context": request.context},
            )
            for i in range(request.count)
        ]
        items.sort(key=lambda x: x.score, reverse=True)
        for pos, item in enumerate(items):
            item.position = pos
        total = request.count * 3
        await self._producer.publish(
            "recommendation.recommendations.completed",
            request.user_id,
            {"user_id": request.user_id, "context": request.context, "count": request.count, "total": total},
        )
        return GetRecommendationsResponse(items=items, total=total, context=request.context)

    async def get_similar(self, request: GetSimilarRequest) -> GetSimilarResponse:
        items = [
            RecommendedItem(
                item_id=f"sim-{request.item_id}-{i}",
                score=round(random.random(), 4),
                reason=f"Similar to {request.item_id}",
                position=i,
            )
            for i in range(request.count)
        ]
        items.sort(key=lambda x: x.score, reverse=True)
        for pos, item in enumerate(items):
            item.position = pos
        await self._producer.publish(
            "recommendation.similar.completed",
            request.item_id,
            {"item_id": request.item_id, "count": request.count},
        )
        return GetSimilarResponse(items=items, source_item_id=request.item_id)

    async def get_trending(self, request: GetTrendingRequest) -> GetTrendingResponse:
        items = [
            RecommendedItem(
                item_id=f"trend-{request.timeframe}-{i}",
                score=round(random.random(), 4),
                reason="Trending in " + (request.category or request.timeframe),
                position=i,
                metadata={"timeframe": request.timeframe, "category": request.category},
            )
            for i in range(request.count)
        ]
        items.sort(key=lambda x: x.score, reverse=True)
        for pos, item in enumerate(items):
            item.position = pos
        await self._producer.publish(
            "recommendation.trending.completed",
            request.timeframe,
            {"timeframe": request.timeframe, "category": request.category, "count": request.count},
        )
        return GetTrendingResponse(items=items, timeframe=request.timeframe, category=request.category)

    async def personalize_feed(self, request: PersonalizeFeedRequest) -> PersonalizeFeedResponse:
        profile = self._repository.get_user_profile(request.user_id)
        items = [
            RecommendedItem(
                item_id=f"feed-{request.user_id}-{request.offset + i}",
                score=round(random.random(), 4),
                reason=random.choice(MOCK_REASONS),
                position=i,
                metadata={"feed_type": request.feed_type},
            )
            for i in range(request.count)
        ]
        items.sort(key=lambda x: x.score, reverse=True)
        for pos, item in enumerate(items):
            item.position = pos
        total = request.count * 3
        personalization_score = round(random.uniform(0.5, 1.0), 4)
        await self._producer.publish(
            "recommendation.feed.completed",
            request.user_id,
            {
                "user_id": request.user_id,
                "feed_type": request.feed_type,
                "count": request.count,
                "total": total,
                "personalization_score": personalization_score,
                "preferred_categories": profile.get("preferred_categories", []),
            },
        )
        return PersonalizeFeedResponse(
            items=items,
            total=total,
            feed_type=request.feed_type,
            personalization_score=personalization_score,
        )

    async def track_interaction(self, interaction: UserInteraction) -> TrackInteractionResponse:
        self._repository.store_interaction(interaction)
        await self._producer.publish(
            "recommendation.interaction.tracked",
            interaction.user_id,
            {
                "user_id": interaction.user_id,
                "item_id": interaction.item_id,
                "interaction_type": interaction.interaction_type,
                "weight": interaction.weight,
                "timestamp": interaction.timestamp or datetime.datetime.now(datetime.timezone.utc).isoformat(),
            },
        )
        return TrackInteractionResponse(status="success", recorded=True)

    async def batch_track(self, interactions: list[UserInteraction]) -> BatchInteractionsResponse:
        for interaction in interactions:
            self._repository.store_interaction(interaction)
        recorded_count = len(interactions)
        await self._producer.publish(
            "recommendation.interactions.batch.completed",
            "batch",
            {"recorded_count": recorded_count},
        )
        return BatchInteractionsResponse(recorded_count=recorded_count)
